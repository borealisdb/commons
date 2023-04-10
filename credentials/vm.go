package credentials

// Temporary solution

import (
	"context"
	"fmt"
	borealisdbv1 "github.com/borealisdb/commons/borealisdb.io/v1"
	"os"
)

type VM struct{}

func (m VM) Init() error {
	return nil
}

func (m VM) GetPostgresCredentials(
	ctx context.Context,
	clusterName string,
	username string,
	options Options,
) (GetPostgresCredentialsResponse, error) {
	return GetPostgresCredentialsResponse{
		Username: os.Getenv(fmt.Sprintf("%v_%v_CLUSTER_USERNAME", clusterName, username)),
		Password: os.Getenv(fmt.Sprintf("%v_%v_CLUSTER_PASSWORD", clusterName, username)),
		Host:     os.Getenv(fmt.Sprintf("%v_CLUSTER_HOST", clusterName)),
	}, nil
}

func (m VM) GetClusterCredentials(ctx context.Context, clusterName string, args Options) (GetClusterCredentialsResponse, error) {
	return GetClusterCredentialsResponse{}, nil
}

func (m VM) GetClusterInfo(ctx context.Context, clusterName string) (borealisdbv1.Postgresql, error) {
	return borealisdbv1.Postgresql{}, nil
}

func (m VM) GetPostgresSSLRootCert(ctx context.Context, clusterName string, options Options) (GetPostgresSSLRootCertResponse, error) {
	return GetPostgresSSLRootCertResponse{}, nil
}
