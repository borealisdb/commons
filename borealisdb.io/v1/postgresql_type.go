package v1

// Postgres CRD definition, please use CamelCase for field names.

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

// Postgresql defines PostgreSQL Custom Resource Definition Object.
type Postgresql struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresSpec   `json:"spec"`
	Status PostgresStatus `json:"status,omitempty"`
	Error  string         `json:"-"`
}

// PostgresSpec defines the specification for the PostgreSQL TPR.
type PostgresSpec struct {
	Patroni   `json:"patroni,omitempty"`
	Resources `json:"resources,omitempty"`

	ClusterSecretsName string   `json:"clusterSecretsName,omitempty"`
	Databases          []string `json:"databases"`

	// Plugins
	TLS            TLS            `json:"tls,omitempty"`
	Authentication Authentication `json:"authentication,omitempty"`
	Monitoring     Monitoring     `json:"monitoring,omitempty"`
	Backup         Backup         `json:"backup,omitempty"`
	LoadBalancer   LoadBalancer   `json:"loadBalancer,omitempty"`

	ClusterParameters map[string]string `json:"clusterParameters,omitempty"`

	EngineVersion       string `json:"engineVersion"`
	EngineMode          string `json:"engineMode,omitempty"`
	MaxAllocatedStorage string `json:"maxAllocatedStorage"`
	DeleteProtection    bool   `json:"deleteProtection,omitempty"`
	NumberOfInstances   int32  `json:"numberOfInstances"`
	DockerImage         string `json:"dockerImage,omitempty"`

	Clone          *CloneDescription   `json:"clone,omitempty"`
	ClusterName    string              `json:"-"`
	StandbyCluster *StandbyDescription `json:"standby,omitempty"`

	Advanced Advanced `json:"advanced,omitempty"`

	// load balancers' source ranges are the same for master and replica services
	AllowedSourceRanges []string `json:"allowedSourceRanges,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresqlList defines a list of PostgreSQL clusters.
type PostgresqlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Postgresql `json:"items"`
}

type Authentication struct {
	PluginName  string `json:"pluginName,omitempty"`
	Host        string `json:"host"`
	RootUrlPath string `json:"rootUrlPath,omitempty"`
	LogLevel    string `json:"logLevel,omitempty"`
}

type Monitoring struct {
	PluginName           string `json:"pluginName,omitempty"`
	PgVersion            string `json:"pgVersion,omitempty"`
	PgData               string `json:"pgData,omitempty"`
	PgPasswordSecretName string `json:"pgPasswordSecretName,omitempty"`
	PgUsername           string `json:"pgUsername,omitempty"`
	PgDataVolumeName     string `json:"pgDataVolumeName,omitempty"`
	InfrastructureHost   string `json:"infrastructureHost,omitempty"`
	GrpcCollectorPort    string `json:"grpcCollectorPort,omitempty"`
	LogLevel             string `json:"logLevel,omitempty"`
	VictoriaMetricsPort  string `json:"victoriaMetricsPort,omitempty"`
	SidecarImage         string `json:"sidecarImage,omitempty"`
}

type Backup struct {
	PluginName            string `json:"pluginName,omitempty"`
	BackupEndpoint        string `json:"backupEndpoint,omitempty"`
	S3BucketName          string `json:"s3BucketName,omitempty"`
	PreferredBackupWindow string `json:"preferredBackupWindow,omitempty"`
	BackupRetentionPeriod string `json:"backupRetentionPeriod,omitempty"`
	BackupRetentionNumber string `json:"backupRetentionNumber,omitempty"`
	EnableEncryption      string `json:"enableEncryption,omitempty"`
	OwnEncryptionKey      string `json:"ownEncryptionKey,omitempty"`
	DeletePolicy          string `json:"deletePolicy,omitempty" defaults:"Snapshot"` // Delete, Retain, Snapshot
	ClusterSecretsName    string `json:"clusterSecretsName,omitempty"`

	RestoreConfig Restore `json:"restoreConfig,omitempty"`
}

type Restore struct {
	UID              string `json:"uid,omitempty"`
	EndTimestamp     string `json:"timestamp,omitempty"`
	S3WalPath        string `json:"s3WalPath,omitempty"`
	S3ForcePathStyle *bool  `json:"s3ForcePathStyle,omitempty" defaults:"false"`
}

type TLS struct {
	PluginName string `json:"pluginName,omitempty"`
	LogLevel   string `json:"logLevel,omitempty"`
}

// Volume describes a single volume in the manifest.
type Volume struct {
	Selector     *metav1.LabelSelector `json:"selector,omitempty"`
	Size         string                `json:"-"` // Not used anymore, we can keep it for consistency
	StorageClass string                `json:"storageClass,omitempty"`
	SubPath      string                `json:"subPath,omitempty"`
	Iops         *int64                `json:"iops,omitempty"`
	Throughput   *int64                `json:"throughput,omitempty"`
	VolumeType   string                `json:"type,omitempty"`
}

type AdditionalVolume struct {
	Name             string          `json:"name"`
	MountPath        string          `json:"mountPath"`
	SubPath          string          `json:"subPath,omitempty"`
	TargetContainers []string        `json:"targetContainers"`
	VolumeSource     v1.VolumeSource `json:"volumeSource,omitempty"`
}

// ResourceDescription describes CPU and memory resources defined for a cluster.
type ResourceDescription struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// Resources describes requests and limits for the cluster resources.
type Resources struct {
	ResourceRequests ResourceDescription `json:"requests,omitempty"`
	ResourceLimits   ResourceDescription `json:"limits,omitempty"`
}

// Patroni contains Patroni-specific configuration
type Patroni struct {
	InitDB                map[string]string            `json:"initdb,omitempty"`
	PgHba                 []string                     `json:"pg_hba,omitempty"`
	TTL                   uint32                       `json:"ttl,omitempty"`
	LoopWait              uint32                       `json:"loop_wait,omitempty"`
	RetryTimeout          uint32                       `json:"retry_timeout,omitempty"`
	MaximumLagOnFailover  float32                      `json:"maximum_lag_on_failover,omitempty"` // float32 because https://github.com/kubernetes/kubernetes/issues/30213
	Slots                 map[string]map[string]string `json:"slots,omitempty"`
	SynchronousMode       bool                         `json:"synchronous_mode,omitempty"`
	SynchronousModeStrict bool                         `json:"synchronous_mode_strict,omitempty"`
}

// StandbyDescription contains s3 wal path
type StandbyDescription struct {
	S3WalPath string `json:"s3_wal_path,omitempty"`
	GSWalPath string `json:"gs_wal_path,omitempty"`
}

// TLSDescription specs TLS properties
type TLSDescription struct {
	SecretName      string `json:"secretName,omitempty"`
	CertificateFile string `json:"certificateFile,omitempty"`
	PrivateKeyFile  string `json:"privateKeyFile,omitempty"`
	CAFile          string `json:"caFile,omitempty"`
	CASecretName    string `json:"caSecretName,omitempty"`
}

// CloneDescription describes which cluster the new should clone and up to which point in time
type CloneDescription struct {
	ClusterName       string `json:"cluster,omitempty"`
	UID               string `json:"uid,omitempty"`
	EndTimestamp      string `json:"timestamp,omitempty"`
	S3WalPath         string `json:"s3_wal_path,omitempty"`
	S3Endpoint        string `json:"s3_endpoint,omitempty"`
	S3AccessKeyId     string `json:"s3_access_key_id,omitempty"`
	S3SecretAccessKey string `json:"s3_secret_access_key,omitempty"`
	S3ForcePathStyle  *bool  `json:"s3_force_path_style,omitempty" defaults:"false"`
}

// Sidecar defines a container to be run in the same pod as the Postgres container.
type Sidecar struct {
	Resources   `json:"resources,omitempty"`
	Name        string             `json:"name,omitempty"`
	DockerImage string             `json:"image,omitempty"`
	Ports       []v1.ContainerPort `json:"ports,omitempty"`
	Env         []v1.EnvVar        `json:"env,omitempty"`
}

// PostgresStatus contains status of the PostgreSQL cluster (running, creation failed etc.)
type PostgresStatus struct {
	PostgresClusterStatus string `json:"PostgresClusterStatus"`
}

type LoadBalancer struct {
	PluginName string `json:"pluginName,omitempty"`
	Image      string `json:"image,omitempty"`
	Disabled   bool   `json:"disabled,omitempty" default:"false"`

	NumberOfInstances *int32 `json:"numberOfInstances,omitempty"`
	Schema            string `json:"schema,omitempty"`
	User              string `json:"user,omitempty"`
	Mode              string `json:"mode,omitempty"`
	MaxDBConnections  *int32 `json:"maxDBConnections,omitempty"`
	PgPort            int32  `json:"pgPort,omitempty" default:"5432"`

	Resources `json:"resources,omitempty"`
}

type Advanced struct {
	Patroni Patroni `json:"patroni,omitempty"`

	Volume               Volume             `json:"volume,omitempty"`
	Sidecars             []Sidecar          `json:"sidecars,omitempty"`
	InitContainers       []v1.Container     `json:"initContainers,omitempty"`
	NodeAffinity         *v1.NodeAffinity   `json:"nodeAffinity,omitempty"`
	Tolerations          []v1.Toleration    `json:"tolerations,omitempty"`
	PodPriorityClassName string             `json:"podPriorityClassName,omitempty"`
	PodAnnotations       map[string]string  `json:"podAnnotations,omitempty"`
	ServiceAnnotations   map[string]string  `json:"serviceAnnotations,omitempty"`
	AdditionalVolumes    []AdditionalVolume `json:"additionalVolumes,omitempty"`
	ShmVolume            *bool              `json:"enableShmVolume,omitempty"`

	SpiloRunAsUser  *int64 `json:"spiloRunAsUser,omitempty"`
	SpiloRunAsGroup *int64 `json:"spiloRunAsGroup,omitempty"`
	SpiloFSGroup    *int64 `json:"spiloFSGroup,omitempty"`

	SchedulerName *string `json:"schedulerName,omitempty"`
}
