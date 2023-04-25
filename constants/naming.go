package constants

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

func getNameTemplate(clusterName, mod string) string {
	return fmt.Sprintf("%v-%v", clusterName, mod)
}

func GetClusterEndpoint(clusterName, namespace, role string) string {
	if namespace == "" {
		namespace = "default"
	}
	if role == RoleReplica {
		return fmt.Sprintf("%s-repl.%s.svc.%s", clusterName, namespace, "cluster.local")
	}
	return fmt.Sprintf("%s.%s.svc.%s", clusterName, namespace, "cluster.local")
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

func GetClusterSecrets(clusterName string) string {
	return fmt.Sprintf("%v-secrets", clusterName)
}

func GetTLSSecretName(clusterName string) string {
	return fmt.Sprintf("%v-tls", clusterName)
}

func GetDefaultBackupEndpoint(namespace string) string {
	return fmt.Sprintf("http://%v.%v.svc.cluster.local:%v", BackupHost, namespace, BackupSystemPort)
}

func GetImageName(name, version string) string {
	return fmt.Sprintf("%v/%v:%v", RepositoryName, name, version)
}
