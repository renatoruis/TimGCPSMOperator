package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TimGcpSmSecretConfigSpec holds defaults for TimGcpSmSecret (GCP project). Auth uses Workload Identity / ADC on the operator.
type TimGcpSmSecretConfigSpec struct {
	// ProjectID is the GCP project id for Secret Manager
	// +kubebuilder:validation:Required
	ProjectID string `json:"projectId"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=tgsc

// TimGcpSmSecretConfig is the Schema for centralized GCP project configuration
type TimGcpSmSecretConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TimGcpSmSecretConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// TimGcpSmSecretConfigList contains a list of TimGcpSmSecretConfig
type TimGcpSmSecretConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimGcpSmSecretConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TimGcpSmSecretConfig{}, &TimGcpSmSecretConfigList{})
}
