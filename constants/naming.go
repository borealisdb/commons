package constants

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func getNameTemplate(clusterName, mod string) string {
	return fmt.Sprintf("%v-%v", clusterName, mod)
}

func GetBackupPath(clusterName string) string {
	return fmt.Sprintf("s3://%v", clusterName)
}

func GetLoadBalancerName(clusterName string) string {
	return getNameTemplate(clusterName, "loadbalancer")
}

func GetCredentialSecretNameForCluster(username string, clusterName string) string {
	return fmt.Sprintf("%v-%v-credentials", clusterName, username)
}

func GetPasswordFromSecret(secret *v1.Secret) string {
	return string(secret.Data[PostgresClusterSecretPasswordKey])
}

func GetUsernameFromSecret(secret *v1.Secret) string {
	return string(secret.Data[PostgresClusterSecretUsernameKey])
}

func GetSecretsForInfrastructures() string {
	return fmt.Sprintf("%v-secrets", AppName)
}

func GetClusterSecrets(clusterName string) string {
	return fmt.Sprintf("%v-secrets", clusterName)
}

func GetTLSSecretName(clusterName string) string {
	return fmt.Sprintf("%v-tls", clusterName)
}

func GetDefaultBackupEndpoint() string {
	return fmt.Sprintf("http://%v:%v", BackupHost, BackupSystemPort)
}

func GetImageName(name, version string) string {
	return fmt.Sprintf("%v/%v:%v", RepositoryName, name, version)
}
