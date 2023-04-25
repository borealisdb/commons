package credentials

import (
	"context"
	"fmt"
	borealisdbv1 "github.com/borealisdb/commons/borealisdb.io/v1"
	"github.com/borealisdb/commons/constants"
	"github.com/borealisdb/commons/k8sutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Kubernetes struct {
	kubeClient k8sutil.KubernetesClient
}

func (k *Kubernetes) GetClusterEndpoint(ctx context.Context, clusterName, role string) (GetClusterEndpointResponse, error) {
	info, err := k.getClusterInfo(ctx, clusterName)
	if err != nil {
		return GetClusterEndpointResponse{}, fmt.Errorf("could not getClusterInfo: %v", err)
	}

	return GetClusterEndpointResponse{
		Endpoint: constants.GetClusterEndpoint(clusterName, info.Namespace, role),
	}, nil
}

func (k *Kubernetes) Init() error {
	kubeClient, err := k8sutil.InitializeKubeClient()
	if err != nil {
		return err
	}

	k.kubeClient = kubeClient
	return nil
}

func (k *Kubernetes) SetKubeClient(client k8sutil.KubernetesClient) {
	k.kubeClient = client
}

func (k *Kubernetes) GetPostgresCredentials(
	ctx context.Context,
	clusterName string,
	username string,
	options Options,
) (GetPostgresCredentialsResponse, error) {
	info, err := k.getClusterInfo(ctx, clusterName)
	if err != nil {
		return GetPostgresCredentialsResponse{}, err
	}

	secretName := constants.GetCredentialSecretNameForCluster(username, clusterName)
	secret, err := k.kubeClient.Secrets(info.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return GetPostgresCredentialsResponse{}, err
	}
	pgPassword := constants.GetPasswordFromSecret(secret)
	if pgPassword == "" {
		return GetPostgresCredentialsResponse{}, fmt.Errorf("pgPassword for user %v not found in secret %v", username, secretName)
	}

	return GetPostgresCredentialsResponse{
		Username: username,
		Password: pgPassword,
		Host:     constants.GetClusterEndpoint(clusterName, info.Namespace, constants.RoleMaster),
	}, nil
}

func (k *Kubernetes) GetPostgresSSLRootCert(ctx context.Context, clusterName string, options Options) (GetPostgresSSLRootCertResponse, error) {
	info, err := k.getClusterInfo(ctx, clusterName)
	if err != nil {
		return GetPostgresSSLRootCertResponse{}, err
	}

	tlsSecretName := constants.GetTLSSecretName(clusterName)
	tlsSecret, err := k.kubeClient.Secrets(info.Namespace).Get(ctx, tlsSecretName, metav1.GetOptions{})
	if err != nil {
		return GetPostgresSSLRootCertResponse{}, fmt.Errorf("could not Get Secrets for %v: %v", tlsSecretName, err)
	}

	return GetPostgresSSLRootCertResponse{
		RootCertBytes: tlsSecret.Data[constants.RootCaCertName],
	}, err
}

func (k *Kubernetes) GetClusterCredentials(ctx context.Context, clusterName string, args Options) (GetClusterCredentialsResponse, error) {
	info, err := k.getClusterInfo(ctx, clusterName)
	if err != nil {
		return GetClusterCredentialsResponse{}, err
	}

	secret, err := k.kubeClient.Secrets(info.Namespace).Get(ctx, info.Spec.ClusterSecretsName, metav1.GetOptions{})
	if err != nil {
		return GetClusterCredentialsResponse{}, err
	}

	return GetClusterCredentialsResponse{
		AwsAccessKeyId:      string(secret.Data["awsAccessKeyId"]),
		AwsSecretAccessKey:  string(secret.Data["awsSecretAccessKey"]),
		BackupEncryptionKey: string(secret.Data["backupEncryptionKey"]),
	}, err
}

func (k *Kubernetes) getClusterInfo(ctx context.Context, clusterName string) (borealisdbv1.Postgresql, error) {
	list, err := k.kubeClient.Postgresqls("").List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%v", clusterName)})
	if err != nil {
		return borealisdbv1.Postgresql{}, err
	}
	if len(list.Items) == 0 {
		return borealisdbv1.Postgresql{}, fmt.Errorf("no cluster found with name %v", clusterName)
	}

	cluster := getDefaults(list)
	return cluster, nil
}

func getDefaults(list *borealisdbv1.PostgresqlList) borealisdbv1.Postgresql {
	cluster := list.Items[0]
	return cluster
}
