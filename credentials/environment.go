package credentials

import (
	"context"
	"fmt"
	borealisdbv1 "github.com/borealisdb/commons/borealisdb.io/v1"
	"github.com/borealisdb/commons/constants"
	"os"
)

const EnvironmentProvider = "environment"

type Environment struct{}

func (m Environment) GetClusterEndpoint(ctx context.Context, clusterName string, role string) (GetClusterEndpointResponse, error) {
	port, ok := os.LookupEnv(fmt.Sprintf("%v_CLUSTER_PORT", clusterName))
	if !ok {
		port = constants.PostgresDefaultPort
	}
	return GetClusterEndpointResponse{
		Hostname: os.Getenv(fmt.Sprintf("%v_CLUSTER_HOSTNAME", clusterName)),
		Port:     port,
	}, nil
}

func (m Environment) Init() error {
	return nil
}

func (m Environment) GetPostgresCredentials(
	ctx context.Context,
	clusterName string,
	username string,
	options Options,
) (GetPostgresCredentialsResponse, error) {
	if username == "" {
		username = os.Getenv(fmt.Sprintf("%v_CLUSTER_USERNAME", clusterName))
	}
	return GetPostgresCredentialsResponse{
		Username: username,
		Password: os.Getenv(fmt.Sprintf("%v_%v_CLUSTER_PASSWORD", clusterName, username)),
	}, nil
}

func (m Environment) GetClusterCredentials(ctx context.Context, clusterName string, args Options) (GetClusterCredentialsResponse, error) {
	return GetClusterCredentialsResponse{}, nil
}

func (m Environment) GetClusterInfo(ctx context.Context, clusterName string) (borealisdbv1.Postgresql, error) {
	return borealisdbv1.Postgresql{}, nil
}

func (m Environment) GetPostgresSSLRootCert(ctx context.Context, clusterName string, options Options) (GetPostgresSSLRootCertResponse, error) {
	return GetPostgresSSLRootCertResponse{}, nil
}
