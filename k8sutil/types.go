package k8sutil

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"os"
	"strings"
)

// NamespacedName describes the namespace/name pairs used in Kubernetes names.
type NamespacedName types.NamespacedName

func (n NamespacedName) String() string {
	return types.NamespacedName(n).String()
}

// MarshalJSON defines marshaling rule for the namespaced name type.
func (n NamespacedName) MarshalJSON() ([]byte, error) {
	return []byte("\"" + n.String() + "\""), nil
}

// Decode converts a (possibly unqualified) string into the namespaced name object.
func (n *NamespacedName) Decode(value string) error {
	return n.DecodeWorker(value, GetOperatorNamespace())
}

// UnmarshalJSON converts a byte slice to NamespacedName
func (n *NamespacedName) UnmarshalJSON(data []byte) error {
	result := NamespacedName{}
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if err := result.Decode(tmp); err != nil {
		return err
	}
	*n = result
	return nil
}

// DecodeWorker separates the decode logic to (unit) test
// from obtaining the operator namespace that depends on k8s mounting files at runtime
func (n *NamespacedName) DecodeWorker(value, operatorNamespace string) error {
	var (
		name types.NamespacedName
	)

	result := strings.SplitN(value, string(types.Separator), 2)
	if len(result) < 2 {
		name.Name = result[0]
	} else {
		name.Name = strings.TrimLeft(result[1], string(types.Separator))
		name.Namespace = result[0]
	}
	if name.Name == "" {
		return fmt.Errorf("incorrect namespaced name: %v", value)
	}
	if name.Namespace == "" {
		name.Namespace = operatorNamespace
	}

	*n = NamespacedName(name)

	return nil
}

const fileWithNamespace = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// cached value for the GetOperatorNamespace
var operatorNamespace string

// GetOperatorNamespace assumes serviceaccount secret is mounted by kubernetes
// Placing this func here instead of pgk/util avoids circular import
func GetOperatorNamespace() string {
	if operatorNamespace == "" {
		if namespaceFromEnvironment := os.Getenv("OPERATOR_NAMESPACE"); namespaceFromEnvironment != "" {
			return namespaceFromEnvironment
		}
		operatorNamespaceBytes, err := os.ReadFile(fileWithNamespace)
		if err != nil {
			log.Fatalf("Unable to detect operator namespace from within its pod due to: %v", err)
		}
		operatorNamespace = string(operatorNamespaceBytes)
	}
	return operatorNamespace
}
