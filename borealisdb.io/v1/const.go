package v1

// ClusterStatusUnknown etc : status of a Postgres cluster known to the operator
const (
	ClusterStatusUnknown  = ""
	ClusterStatusCreating = "Creating"
	ClusterStatusUpdating = "Updating"
	ClusterStatusSyncing  = "Syncing"

	ClusterStatusUpdateFailed = "UpdateFailed"
	ClusterStatusSyncFailed   = "SyncFailed"
	ClusterStatusAddFailed    = "CreateFailed"
	ClusterStatusRunning      = "Running"
	ClusterStatusInvalid      = "Invalid"

	AccountsClusterNameLabel = "clusterName"

	PostgresCRDResourceKind = "postgresql"
)

const (
	serviceNameMaxLength   = 63
	clusterNameMaxLength   = serviceNameMaxLength - len("-repl")
	serviceNameRegexString = `^[a-z]([-a-z0-9]*[a-z0-9])?$`
)

var PostgresSupportedVersionImages = map[string]string{
	"14": "registry.opensource.zalan.do/acid/spilo-14:2.1-p3",
	"15": "ghcr.io/zalando/spilo-15:3.0-p1",
}

// GetPostgresSupportedVersionImages later we may need to fetch this map from a remote location to automated minor version upgrade.
// This method should be able to cache previous versions in case it will fail
func GetPostgresSupportedVersionImages() (map[string]string, error) {
	return PostgresSupportedVersionImages, nil
}
