package plugins

import (
	v12 "github.com/borealisdb/commons/borealisdb.io/v1"
	"github.com/borealisdb/commons/constants"
)

func SetBackupDefaults(backup v12.Backup, clusterName string) v12.Backup {
	if backup.BackupEndpoint == "" {
		backup.BackupEndpoint = constants.GetDefaultBackupEndpoint()
	}
	if backup.S3BucketName == "" {
		backup.S3BucketName = clusterName
	}

	return backup
}
