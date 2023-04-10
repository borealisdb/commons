package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var parseTimeTests = []struct {
	about string
	in    string
	out   metav1.Time
	err   error
}{
	{"parse common time with minutes", "16:08", mustParseTime("16:08"), nil},
	{"parse time with zeroed minutes", "11:00", mustParseTime("11:00"), nil},
	{"parse corner case last minute of the day", "23:59", mustParseTime("23:59"), nil},

	{"expect error as hour is out of range", "26:09", metav1.Now(), errors.New(`parsing time "26:09": hour out of range`)},
	{"expect error as minute is out of range", "23:69", metav1.Now(), errors.New(`parsing time "23:69": minute out of range`)},
}

var parseWeekdayTests = []struct {
	about string
	in    string
	out   time.Weekday
	err   error
}{
	{"parse common weekday", "Wed", time.Wednesday, nil},
	{"expect error as weekday is invalid", "Sunday", time.Weekday(0), errors.New("incorrect weekday")},
	{"expect error as weekday is empty", "", time.Weekday(0), errors.New("incorrect weekday")},
}

var clusterNames = []struct {
	about       string
	in          string
	clusterName string
	err         error
}{
	{"common team and cluster name", "test", "test", nil},
	{"cluster name with hyphen", "my-name", "my-name", nil},
	{"cluster and team name with hyphen", "another-test", "another-test", nil},
	{"expect error as cluster name is just hyphens", "-----", "cluster",
		errors.New(`name must confirm to DNS-1035, regex used for validation is "^[a-z]([-a-z0-9]*[a-z0-9])?$"`)},
	{"expect error as cluster name is too long", "fooobar-fooobarfooobarfooobarfooobarfooobarfooobarfooobarfooobar", "",
		errors.New("name cannot be longer than 58 characters")},
	{"expect error as cluster name is empty", "-test", "", errors.New("name must confirm to DNS-1035, regex used for validation is \"^[a-z]([-a-z0-9]*[a-z0-9])?$\"")},
	{"expect error as cluster and team name are hyphens", "-", "", errors.New("name must confirm to DNS-1035, regex used for validation is \"^[a-z]([-a-z0-9]*[a-z0-9])?$\"")},
}

var cloneClusterDescriptions = []struct {
	about string
	in    *CloneDescription
	err   error
}{
	{"cluster name invalid but EndTimeSet is not empty", &CloneDescription{"foo+bar", "", "NotEmpty", "", "", "", "", nil}, nil},
	{"expect error as cluster name does not match DNS-1035", &CloneDescription{"foo+bar", "", "", "", "", "", "", nil},
		errors.New(`clone cluster name must confirm to DNS-1035, regex used for validation is "^[a-z]([-a-z0-9]*[a-z0-9])?$"`)},
	{"expect error as cluster name is too long", &CloneDescription{"foobar123456789012345678901234567890123456789012345678901234567890", "", "", "", "", "", "", nil},
		errors.New("clone cluster name must be no longer than 63 characters")},
	{"common cluster name", &CloneDescription{"foobar", "", "", "", "", "", "", nil}, nil},
}

var postgresStatus = []struct {
	about string
	in    []byte
	out   PostgresStatus
	err   error
}{
	{"cluster running", []byte(`{"PostgresClusterStatus":"Running"}`),
		PostgresStatus{PostgresClusterStatus: ClusterStatusRunning}, nil},
	{"cluster status undefined", []byte(`{"PostgresClusterStatus":""}`),
		PostgresStatus{PostgresClusterStatus: ClusterStatusUnknown}, nil},
	{"cluster running without full JSON format", []byte(`"Running"`),
		PostgresStatus{PostgresClusterStatus: ClusterStatusRunning}, nil},
	{"cluster status empty", []byte(`""`),
		PostgresStatus{PostgresClusterStatus: ClusterStatusUnknown}, nil}}

var tmp postgresqlCopy
var unmarshalCluster = []struct {
	about   string
	in      []byte
	out     Postgresql
	marshal []byte
	err     error
}{
	{
		about: "example with simple status field",
		in: []byte(`{
	  "kind": "Postgresql","apiVersion": "borealisdb.io/v1",
	  "metadata": {"name": "borealisdb-testcluster1"}, "spec": {}}`),
		out: Postgresql{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Postgresql",
				APIVersion: "borealisdb.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "borealisdb-testcluster1",
			},
			Spec: PostgresSpec{ClusterName: "borealisdb-testcluster1"},
			// This error message can vary between Go versions, so compute it for the current version.
		},
		marshal: []byte(`{"kind":"Postgresql","apiVersion":"borealisdb.io/v1","metadata":{"name":"borealisdb-testcluster1","creationTimestamp":null},"spec":{"postgresql":{"version":"","parameters":null},"volume":{"size":"","storageClass":""},"patroni":{"initdb":null,"pg_hba":null,"ttl":0,"loop_wait":0,"retry_timeout":0,"maximum_lag_on_failover":0,"slots":null},"resources":{"requests":{"cpu":"","memory":""},"limits":{"cpu":"","memory":""}},"allowedSourceRanges":null,"numberOfInstances":0,"clone":null}}`),
		err:     nil},
	{
		about: "example with /status subresource",
		in: []byte(`{
	  "kind": "Postgresql","apiVersion": "borealisdb.io/v1",
	  "metadata": {"name": "borealisdb-testcluster1"}, "spec": {}}`),
		out: Postgresql{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Postgresql",
				APIVersion: "borealisdb.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "borealisdb-testcluster1",
			},
			Spec: PostgresSpec{ClusterName: "borealisdb-testcluster1"},
			// This error message can vary between Go versions, so compute it for the current version.
		},
		marshal: []byte(`{"kind":"Postgresql","apiVersion":"borealisdb.io/v1","metadata":{"name":"borealisdb-testcluster1","creationTimestamp":null},"spec":{"postgresql":{"version":"","parameters":null},"volume":{"size":"","storageClass":""},"patroni":{"initdb":null,"pg_hba":null,"ttl":0,"loop_wait":0,"retry_timeout":0,"maximum_lag_on_failover":0,"slots":null},"resources":{"requests":{"cpu":"","memory":""},"limits":{"cpu":"","memory":""}},"allowedSourceRanges":null,"numberOfInstances":0,"clone":null},"status":{"PostgresClusterStatus":""}}`),
		err:     nil},
	{
		about: "example with clone",
		in:    []byte(`{"kind": "Postgresql","apiVersion": "borealisdb.io/v1","metadata": {"name": "borealis-testcluster1"}, "spec": {"clone": {"cluster": "batman"}}}`),
		out: Postgresql{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Postgresql",
				APIVersion: "borealisdb.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "borealis-testcluster1",
			},
			Spec: PostgresSpec{
				Clone: &CloneDescription{
					ClusterName: "batman",
				},
				ClusterName: "borealis-testcluster1",
			},
			Error: "",
		},
		marshal: []byte(`{"kind":"Postgresql","apiVersion":"borealisdb.io/v1","metadata":{"name":"borealis-testcluster1","creationTimestamp":null},"spec":{"postgresql":{"version":"","parameters":null},"volume":{"size":"","storageClass":""},"patroni":{"initdb":null,"pg_hba":null,"ttl":0,"loop_wait":0,"retry_timeout":0,"maximum_lag_on_failover":0,"slots":null},"resources":{"requests":{"cpu":"","memory":""},"limits":{"cpu":"","memory":""}},"allowedSourceRanges":null,"numberOfInstances":0,"users":null,"clone":{"cluster":"batman"}},"status":{"PostgresClusterStatus":""}}`),
		err:     nil},
	{
		about: "standby example",
		in:    []byte(`{"kind": "Postgresql","apiVersion": "borealisdb.io/v1","metadata": {"name": "testcluster1"}, "spec": {"standby": {"s3_wal_path": "s3://custom/path/to/bucket/"}}}`),
		out: Postgresql{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Postgresql",
				APIVersion: "borealisdb.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "testcluster1",
			},
			Spec: PostgresSpec{
				StandbyCluster: &StandbyDescription{
					S3WalPath: "s3://custom/path/to/bucket/",
				},
				ClusterName: "testcluster1",
			},
			Error: "",
		},
		marshal: []byte(`{"kind":"Postgresql","apiVersion":"borealisdb.io/v1","metadata":{"name":"testcluster1","creationTimestamp":null},"spec":{"postgresql":{"version":"","parameters":null},"volume":{"size":"","storageClass":""},"patroni":{"initdb":null,"pg_hba":null,"ttl":0,"loop_wait":0,"retry_timeout":0,"maximum_lag_on_failover":0,"slots":null},"resources":{"requests":{"cpu":"","memory":""},"limits":{"cpu":"","memory":""}},"allowedSourceRanges":null,"numberOfInstances":0,"users":null,"standby":{"s3_wal_path":"s3://custom/path/to/bucket/"}},"status":{"PostgresClusterStatus":""}}`),
		err:     nil},
	{
		about:   "expect error on malformatted JSON",
		in:      []byte(`{"kind": "Postgresql","apiVersion": "borealisdb.io/v1"`),
		out:     Postgresql{},
		marshal: []byte{},
		err:     errors.New("unexpected end of JSON input")},
	{
		about:   "expect error on JSON with field's value malformatted",
		in:      []byte(`{"kind":"Postgresql","apiVersion":"borealisdb.io/v1","metadata":{"name":"acid-testcluster","creationTimestamp":qaz},"spec":{"postgresql":{"version":"","parameters":null},"volume":{"size":"","storageClass":""},"patroni":{"initdb":null,"pg_hba":null,"ttl":0,"loop_wait":0,"retry_timeout":0,"maximum_lag_on_failover":0,"slots":null},"resources":{"requests":{"cpu":"","memory":""},"limits":{"cpu":"","memory":""}},"teamId":"acid","allowedSourceRanges":null,"numberOfInstances":0,"users":null,"clone":null},"status":{"PostgresClusterStatus":"Invalid"}}`),
		out:     Postgresql{},
		marshal: []byte{},
		err:     errors.New("invalid character 'q' looking for beginning of value"),
	},
}

var postgresqlList = []struct {
	about string
	in    []byte
	out   PostgresqlList
	err   error
}{
	{"expect success", []byte(`{"apiVersion":"v1","items":[{"apiVersion":"borealisdb.io/v1","kind":"Postgresql","metadata":{"labels":{},"name":"testcluster42","namespace":"default","resourceVersion":"30446957","uid":"857cd208-33dc-11e7-b20a-0699041e4b03"},"spec":{"allowedSourceRanges":["185.85.220.0/22"],"numberOfInstances":1,"postgresql":{"version":"9.6"},"volume":{"size":"10Gi"}},"status":{"PostgresClusterStatus":"Running"}}],"kind":"List","metadata":{},"resourceVersion":"","selfLink":""}`),
		PostgresqlList{
			TypeMeta: metav1.TypeMeta{
				Kind:       "List",
				APIVersion: "v1",
			},
			Items: []Postgresql{{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Postgresql",
					APIVersion: "borealisdb.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "borealis-testcluster42",
					Namespace:       "default",
					Labels:          map[string]string{"team": "borealis"},
					ResourceVersion: "30446957",
					UID:             "857cd208-33dc-11e7-b20a-0699041e4b03",
				},
				Spec: PostgresSpec{
					ClusterName:         "testcluster42",
					EngineVersion:       "14",
					AllowedSourceRanges: []string{"185.85.220.0/22"},
					NumberOfInstances:   1,
					MaxAllocatedStorage: "10Gi",
				},
				Status: PostgresStatus{
					PostgresClusterStatus: ClusterStatusRunning,
				},
				Error: "",
			}},
		},
		nil},
	{"expect error on malformatted JSON", []byte(`{"apiVersion":"v1","items":[{"apiVersion":"borealisdb.io/v1","kind":"Postgresql","metadata":{"labels":{"team":"borealis"},"name":"borealis-testcluster42","namespace"`),
		PostgresqlList{},
		errors.New("unexpected end of JSON input")}}

var podAnnotations = []struct {
	about       string
	in          []byte
	annotations map[string]string
	err         error
}{{
	about: "common annotations",
	in: []byte(`{
		"kind": "Postgresql",
		"apiVersion": "borealisdb.io/v1",
		"metadata": {
			"name": "acid-testcluster1"
		},
		"spec": {
			"podAnnotations": {
				"foo": "bar"
			},
			"teamId": "acid",
			"clone": {
				"cluster": "team-batman"
			}
		}
	}`),
	annotations: map[string]string{"foo": "bar"},
	err:         nil},
}

var serviceAnnotations = []struct {
	about       string
	in          []byte
	annotations map[string]string
	err         error
}{
	{
		about: "common single annotation",
		in: []byte(`{
			"kind": "Postgresql",
			"apiVersion": "borealisdb.io/v1",
			"metadata": {
				"name": "acid-testcluster1"
			},
			"spec": {
				"serviceAnnotations": {
					"foo": "bar"
				},
				"teamId": "acid",
				"clone": {
					"cluster": "team-batman"
				}
			}
		}`),
		annotations: map[string]string{"foo": "bar"},
		err:         nil,
	},
	{
		about: "common two annotations",
		in: []byte(`{
			"kind": "Postgresql",
			"apiVersion": "borealisdb.io/v1",
			"metadata": {
				"name": "acid-testcluster1"
			},
			"spec": {
				"serviceAnnotations": {
					"foo": "bar",
					"post": "gres"
				},
				"teamId": "acid",
				"clone": {
					"cluster": "team-batman"
				}
			}
		}`),
		annotations: map[string]string{"foo": "bar", "post": "gres"},
		err:         nil,
	},
}

func mustParseTime(s string) metav1.Time {
	v, err := time.Parse("15:04", s)
	if err != nil {
		panic(err)
	}

	return metav1.Time{Time: v.UTC()}
}

func TestParseTime(t *testing.T) {
	for _, tt := range parseTimeTests {
		t.Run(tt.about, func(t *testing.T) {
			aTime, err := parseTime(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("ParseTime expected error: %v, got: %v", tt.err, err)
				}
				return
			} else if tt.err != nil {
				t.Errorf("Expected error: %v", tt.err)
			}

			if aTime != tt.out {
				t.Errorf("Expected time: %v, got: %v", tt.out, aTime)
			}
		})
	}
}

func TestWeekdayTime(t *testing.T) {
	for _, tt := range parseWeekdayTests {
		t.Run(tt.about, func(t *testing.T) {
			aTime, err := parseWeekday(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("ParseWeekday expected error: %v, got: %v", tt.err, err)
				}
				return
			} else if tt.err != nil {
				t.Errorf("Expected error: %v", tt.err)
			}

			if aTime != tt.out {
				t.Errorf("Expected weekday: %v, got: %v", tt.out, aTime)
			}
		})
	}
}

func TestPodAnnotations(t *testing.T) {
	for _, tt := range podAnnotations {
		t.Run(tt.about, func(t *testing.T) {
			var cluster Postgresql
			err := cluster.UnmarshalJSON(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("Unable to marshal cluster with podAnnotations: expected %v got %v", tt.err, err)
				}
				return
			}
			for k, v := range cluster.Spec.Advanced.PodAnnotations {
				found, expected := v, tt.annotations[k]
				if found != expected {
					t.Errorf("Didn't find correct value for key %v in for podAnnotations: Expected %v found %v", k, expected, found)
				}
			}
		})
	}
}

func TestServiceAnnotations(t *testing.T) {
	for _, tt := range serviceAnnotations {
		t.Run(tt.about, func(t *testing.T) {
			var cluster Postgresql
			err := cluster.UnmarshalJSON(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("Unable to marshal cluster with serviceAnnotations: expected %v got %v", tt.err, err)
				}
				return
			}
			for k, v := range cluster.Spec.Advanced.ServiceAnnotations {
				found, expected := v, tt.annotations[k]
				if found != expected {
					t.Errorf("Didn't find correct value for key %v in for serviceAnnotations: Expected %v found %v", k, expected, found)
				}
			}
		})
	}
}

func TestClusterName(t *testing.T) {
	for _, tt := range clusterNames {
		t.Run(tt.about, func(t *testing.T) {
			name, err := extractClusterName(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("extractClusterName expected error: %v, got: %v", tt.err, err)
				}
				return
			} else if tt.err != nil {
				t.Errorf("Expected error: %v", tt.err)
			}
			if name != tt.clusterName {
				t.Errorf("Expected clusterName: %q, got: %q", tt.clusterName, name)
			}
		})
	}
}

func TestCloneClusterDescription(t *testing.T) {
	for _, tt := range cloneClusterDescriptions {
		t.Run(tt.about, func(t *testing.T) {
			if err := validateCloneClusterDescription(tt.in); err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("testCloneClusterDescription expected error: %v, got: %v", tt.err, err)
				}
			} else if tt.err != nil {
				t.Errorf("Expected error: %v", tt.err)
			}
		})
	}
}

func TestUnmarshalPostgresStatus(t *testing.T) {
	for _, tt := range postgresStatus {
		t.Run(tt.about, func(t *testing.T) {

			var ps PostgresStatus
			err := ps.UnmarshalJSON(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("CR status unmarshal expected error: %v, got %v", tt.err, err)
				}
				return
			}

			if !reflect.DeepEqual(ps, tt.out) {
				t.Errorf("Expected status: %#v, got: %#v", tt.out, ps)
			}
		})
	}
}

func TestPostgresUnmarshal(t *testing.T) {
	for _, tt := range unmarshalCluster {
		t.Run(tt.about, func(t *testing.T) {
			var cluster Postgresql
			err := cluster.UnmarshalJSON(tt.in)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("Unmarshal expected error: %v, got: %v", tt.err, err)
				}
				return
			} else if tt.err != nil {
				t.Errorf("Expected error: %v", tt.err)
			}

			if !reflect.DeepEqual(cluster, tt.out) {
				t.Errorf("Expected Postgresql: %#v, got %#v", tt.out, cluster)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	for _, tt := range unmarshalCluster {
		t.Run(tt.about, func(t *testing.T) {

			if tt.err != nil {
				return
			}

			// Unmarshal and marshal example to capture api changes
			var cluster Postgresql
			err := cluster.UnmarshalJSON(tt.marshal)
			if err != nil {
				if tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("Backwards compatibility unmarshal expected error: %v, got: %v", tt.err, err)
				}
				return
			}
			expected, err := json.Marshal(cluster)
			if err != nil {
				t.Errorf("Backwards compatibility marshal error: %v", err)
			}

			m, err := json.Marshal(tt.out)
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			if !bytes.Equal(m, expected) {
				t.Errorf("Marshal Postgresql \nexpected: %q, \ngot:      %q", string(expected), string(m))
			}
		})
	}
}

func TestPostgresMeta(t *testing.T) {
	for _, tt := range unmarshalCluster {
		t.Run(tt.about, func(t *testing.T) {

			if a := tt.out.GetObjectKind(); a != &tt.out.TypeMeta {
				t.Errorf("GetObjectKindMeta \nexpected: %v, \ngot:       %v", tt.out.TypeMeta, a)
			}

			if a := tt.out.GetObjectMeta(); reflect.DeepEqual(a, tt.out.ObjectMeta) {
				t.Errorf("GetObjectMeta \nexpected: %v, \ngot:       %v", tt.out.ObjectMeta, a)
			}
		})
	}
}

func TestPostgresListMeta(t *testing.T) {
	for _, tt := range postgresqlList {
		t.Run(tt.about, func(t *testing.T) {
			if tt.err != nil {
				return
			}

			if a := tt.out.GetObjectKind(); a != &tt.out.TypeMeta {
				t.Errorf("GetObjectKindMeta expected: %v, got: %v", tt.out.TypeMeta, a)
			}

			if a := tt.out.GetListMeta(); reflect.DeepEqual(a, tt.out.ListMeta) {
				t.Errorf("GetObjectMeta expected: %v, got: %v", tt.out.ListMeta, a)
			}

			return
		})
	}
}

func TestPostgresqlClone(t *testing.T) {
	for _, tt := range unmarshalCluster {
		t.Run(tt.about, func(t *testing.T) {
			cp := &tt.out
			cp.Error = ""
			clone := cp.Clone()
			if !reflect.DeepEqual(clone, cp) {
				t.Errorf("TestPostgresqlClone expected: \n%#v\n, got \n%#v", cp, clone)
			}
		})
	}
}
