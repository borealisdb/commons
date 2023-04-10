package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BorealisClusterAccount defines accounts Custom Resource Definition Object.
type BorealisClusterAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BorealisClusterAccountSpecs  `json:"spec"`
	Status BorealisClusterAccountStatus `json:"status,omitempty"`
	Error  string                       `json:"-"`
}

type Account struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type BorealisClusterAccountSpecs struct {
	Accounts []Account `json:"accounts"`
}

type BorealisClusterAccountStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BorealisClusterAccountList defines a list of Accounts clusters.
type BorealisClusterAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []BorealisClusterAccount `json:"items"`
}
