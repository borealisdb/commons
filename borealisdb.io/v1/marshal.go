package v1

import (
	"encoding/json"
	"fmt"
	"time"
)

type postgresqlCopy Postgresql
type postgresStatusCopy PostgresStatus

// UnmarshalJSON converts a JSON to the status subresource definition.
func (ps *PostgresStatus) UnmarshalJSON(data []byte) error {
	var (
		tmp    postgresStatusCopy
		status string
	)

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		metaErr := json.Unmarshal(data, &status)
		if metaErr != nil {
			return fmt.Errorf("could not parse status: %v; err %v", string(data), metaErr)
		}
		tmp.PostgresClusterStatus = status
	}
	*ps = PostgresStatus(tmp)

	return nil
}

// UnmarshalJSON converts a JSON into the PostgreSQL object.
func (p *Postgresql) UnmarshalJSON(data []byte) error {
	var tmp postgresqlCopy

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		metaErr := json.Unmarshal(data, &tmp.ObjectMeta)
		if metaErr != nil {
			return err
		}

		tmp.Error = err.Error()
		tmp.Status.PostgresClusterStatus = ClusterStatusInvalid

		*p = Postgresql(tmp)

		return nil
	}
	tmp2 := Postgresql(tmp)

	if clusterName, err := extractClusterName(tmp2.ObjectMeta.Name); err != nil {
		tmp2.Error = err.Error()
		tmp2.Status = PostgresStatus{PostgresClusterStatus: ClusterStatusInvalid}
	} else if err := validateCloneClusterDescription(tmp2.Spec.Clone); err != nil {

		tmp2.Error = err.Error()
		tmp2.Status.PostgresClusterStatus = ClusterStatusInvalid
	} else {
		tmp2.Spec.ClusterName = clusterName
	}

	*p = tmp2

	return nil
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var (
		v   interface{}
		err error
	)
	if err = json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		t, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		*d = Duration(t)
		return nil
	case float64:
		t := time.Duration(val)
		*d = Duration(t)
		return nil
	default:
		return fmt.Errorf("could not recognize type %T as a valid type to unmarshal to Duration", val)
	}
}
