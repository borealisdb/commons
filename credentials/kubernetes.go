package credentials

import (
	"context"
	"fmt"
	borealisdbv1 "github.com/borealisdb/commons/borealisdb.io/v1"
	"github.com/borealisdb/commons/constants"
	"github.com/borealisdb/commons/k8sutil"
	"github.com/borealisdb/commons/plugins"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Kubernetes struct {
	kubeClient k8sutil.KubernetesClient
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
	if options.Namespace == "" {
		info, err := k.GetClusterInfo(ctx, clusterName)
		if err != nil {
			return GetPostgresCredentialsResponse{}, err
		}
		options.Namespace = info.Namespace
	}

	secretName := constants.GetCredentialSecretNameForCluster(username, clusterName)
	secret, err := k.kubeClient.Secrets(options.Namespace).Get(ctx, secretName, metav1.GetOptions{})
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
		Host:     clusterName,
	}, nil
}

func (k *Kubernetes) GetPostgresSSLRootCert(ctx context.Context, clusterName string, options Options) (GetPostgresSSLRootCertResponse, error) {
	if options.Namespace == "" {
		info, err := k.GetClusterInfo(ctx, clusterName)
		if err != nil {
			return GetPostgresSSLRootCertResponse{}, err
		}
		options.Namespace = info.Namespace
	}
	tlsSecretName := constants.GetTLSSecretName(clusterName)
	tlsSecret, err := k.kubeClient.Secrets(options.Namespace).Get(ctx, tlsSecretName, metav1.GetOptions{})
	if err != nil {
		return GetPostgresSSLRootCertResponse{}, fmt.Errorf("could not Get Secrets for %v: %v", tlsSecretName, err)
	}

	return GetPostgresSSLRootCertResponse{
		RootCertBytes: tlsSecret.Data[constants.RootCaCertName],
	}, err
}

func (k *Kubernetes) GetClusterCredentials(ctx context.Context, clusterName string, args Options) (GetClusterCredentialsResponse, error) {
	info, err := k.GetClusterInfo(ctx, clusterName)
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

func (k *Kubernetes) GetClusterInfo(ctx context.Context, clusterName string) (borealisdbv1.Postgresql, error) {
	list, err := k.kubeClient.Postgresqls("").List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%v", clusterName)})
	if err != nil {
		return borealisdbv1.Postgresql{}, err
	}
	if len(list.Items) == 0 {
		return borealisdbv1.Postgresql{}, fmt.Errorf("not cluster found with name %v", clusterName)
	}

	cluster := getDefaults(list)
	return cluster, nil
}

func getDefaults(list *borealisdbv1.PostgresqlList) borealisdbv1.Postgresql {
	cluster := list.Items[0]

	// Adding defaults if fields are empty
	backup := plugins.SetBackupDefaults(cluster.Spec.Backup, cluster.Name)
	cluster.Spec.Backup = backup

	if cluster.Spec.ClusterSecretsName == "" {
		cluster.Spec.ClusterSecretsName = constants.GetClusterSecrets(cluster.Name)
	}

	return cluster
}
