package credentials

import (
	"context"
	borealisdbv1 "github.com/borealisdb/commons/borealisdb.io/v1"
)

type Credentials interface {
	Init() error
	GetPostgresCredentials(
		ctx context.Context,
		clusterName string,
		username string,
		options Options,
	) (GetPostgresCredentialsResponse, error)
	GetPostgresSSLRootCert(ctx context.Context, clusterName string, options Options) (GetPostgresSSLRootCertResponse, error)
	GetClusterCredentials(ctx context.Context, clusterName string, args Options) (GetClusterCredentialsResponse, error)
	GetClusterInfo(ctx context.Context, clusterName string) (borealisdbv1.Postgresql, error)
}

type GetPostgresCredentialsResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type GetClusterCredentialsResponse struct {
	AwsAccessKeyId      string `json:"awsAccessKeyId"`
	AwsSecretAccessKey  string `json:"awsSecretAccessKey"`
	BackupEncryptionKey string `json:"backupEncryptionKey"`
}

type GetPostgresSSLRootCertResponse struct {
	RootCertBytes []byte `json:"rootCertBytes"`
}

type Options struct {
	Namespace               string `json:"namespace"`
	KubernetesSecretFromEnv bool
}
