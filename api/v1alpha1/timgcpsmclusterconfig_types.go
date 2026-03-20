package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefaultClusterConfigName is the TimGcpSmClusterConfig resource name used when TimGcpSmSecret
// does not set spec.projectId, spec.gcpSmConfig, or spec.clusterConfig.
const DefaultClusterConfigName = "default"

// TimGcpSmClusterConfigSpec holds cluster-wide defaults (single GCP project for Secret Manager).
type TimGcpSmClusterConfigSpec struct {
	// ProjectID is the GCP project id used when TimGcpSmSecret resources do not specify projectId
	// and do not reference a namespaced TimGcpSmSecretConfig.
	// +kubebuilder:validation:Required
	ProjectID string `json:"projectId"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=tgcc

// TimGcpSmClusterConfig is a cluster-scoped default for GCP project (optional; use namespaced TimGcpSmSecretConfig when you need per-namespace overrides).
type TimGcpSmClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TimGcpSmClusterConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// TimGcpSmClusterConfigList contains a list of TimGcpSmClusterConfig
type TimGcpSmClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimGcpSmClusterConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TimGcpSmClusterConfig{}, &TimGcpSmClusterConfigList{})
}
