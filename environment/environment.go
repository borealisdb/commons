package environment

import "os"

const (
	Kubernetes = "kubernetes"
	VM         = "VM"
	Mock       = "Mock"
)

var allowedEnvironments = map[string]bool{Kubernetes: true, VM: true, Mock: true}

func DetermineEnvironment(environment string) (string, error) {
	// This force the environment
	if environment != "" {
		return environment, nil
	}
	if _, isRunningInKubernetes := os.LookupEnv("KUBERNETES_SERVICE_HOST"); isRunningInKubernetes {
		return Kubernetes, nil
	}

	return VM, nil
}
