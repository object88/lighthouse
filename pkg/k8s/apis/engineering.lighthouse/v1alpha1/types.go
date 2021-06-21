package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config describes a target configuration
// +genclient
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=config
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigSpec `json:"spec"`

	Status ConfigStatus `json:"status"`
}

// ConfigSpec is the spec for a Config
// +k8s:deepcopy-gen=true
type ConfigSpec struct {
	ReleaseName string `json:"releasename"`
}

// ConfigStatus is the status for a Config resource
// +k8s:deepcopy-gen=true
type ConfigStatus struct {
	DeployedAt metav1.Time `json:"deployedat"`
}

// ConfigList is a list of Config resources
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=config
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Config `json:"items"`
}
