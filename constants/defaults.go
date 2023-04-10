package constants

const (
	DomainName            = "borealisdb.io"
	AppName               = "borealis"
	RepositoryName        = "borealis"
	InfrastructuresHost   = "borealis-service"
	BackupHost            = "borealis-backup-service"
	OperatorHost          = "borealis-operator-service"
	ProxyPort             = "8080"
	GrpcProxyPort         = "8086"
	VictoriaMetricsPort   = "8428"
	BackupSystemPort      = "8333"
	MonitoringAPIGRPCPort = "8081"
	MonitoringAPIPort     = "8082"
	AuthorizationAPIPort  = "8083"
	OperatorAPIServerPort = "8085"

	PasswordLength = 64

	ClusterSecretsBackupEncryptionKey = "backupEncryptionKey"

	ImageRepository          = "borealis"
	BackupSystemImageName    = "chrislusf/seaweedfs"
	BackupSystemImageVersion = "latest"

	MonitoringDatabaseImageName    = "clickhouse/clickhouse-server"
	MonitoringDatabaseImageVersion = "latest"

	BorealisInfrastructuresName = "borealis-infrastructures"
)

// Users naming
const (
	AdminUsername       = "postgres"
	ReplicationUsername = "standby"
	MonitoringUsername  = "monitoring"

	Migrator                         = "migrator"
	Application                      = "application"
	Developer                        = "developer"
	Analyst                          = "analyst"
	PostgresClusterSecretPasswordKey = "password"
	PostgresClusterSecretUsernameKey = "user"
)

const (
	RootCaCertName = "root.crt"
	ServerCertName = "tls.crt"
	ServerKeyName  = "tls.key"
)
