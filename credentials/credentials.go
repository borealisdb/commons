package credentials

import (
	"context"
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
	GetClusterEndpoint(ctx context.Context, clusterName, role string) (GetClusterEndpointResponse, error)
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

type GetClusterEndpointResponse struct {
	Endpoint string `json:"endpoint"`
}

type Options struct{}
