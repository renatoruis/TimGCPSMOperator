package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TimGcpSmSecretSpec defines the desired state: sync one GCP Secret Manager secret into a Kubernetes Secret.
type TimGcpSmSecretSpec struct {
	// GcpSmConfig is the name of TimGcpSmSecretConfig to use for default projectId
	// +optional
	GcpSmConfig string `json:"gcpSmConfig,omitempty"`

	// GcpSmConfigNamespace is the namespace of TimGcpSmSecretConfig (defaults to this resource's namespace)
	// +optional
	GcpSmConfigNamespace string `json:"gcpSmConfigNamespace,omitempty"`

	// ClusterConfig is the name of a cluster-scoped TimGcpSmClusterConfig (shared projectId). Ignored if projectId or gcpSmConfig is set.
	// +optional
	ClusterConfig string `json:"clusterConfig,omitempty"`

	// ProjectID is the GCP project ID. Required unless provided via GcpSmConfig, ClusterConfig, or the cluster default (see TimGcpSmClusterConfig).
	// +optional
	ProjectID string `json:"projectId,omitempty"`

	// SecretID is the Secret Manager secret id (short name, not the full resource name)
	// +kubebuilder:validation:Required
	SecretID string `json:"secretId"`

	// SecretVersion is the version to read: numeric id or "latest" (default)
	// +optional
	// +kubebuilder:default=latest
	SecretVersion string `json:"secretVersion,omitempty"`

	// SecretName is the Kubernetes Secret to create or update
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`

	// DeploymentName, if set, triggers a rollout when synced data changes (e.g. after console/API update in GSM)
	// +optional
	DeploymentName string `json:"deploymentName,omitempty"`

	// Namespace for the Kubernetes Secret and Deployment (defaults to this resource's namespace)
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// SyncInterval between polls of Secret Manager (default 5m, min 30s, max 1h)
	// +optional
	// +kubebuilder:default="5m"
	SyncInterval string `json:"syncInterval,omitempty"`

	// DecodeFormat: "text" stores the payload under SecretKey (default); "json" parses a JSON object into multiple keys
	// +optional
	// +kubebuilder:validation:Enum=text;json
	// +kubebuilder:default=text
	DecodeFormat string `json:"decodeFormat,omitempty"`

	// SecretKey used when DecodeFormat is text (default "value")
	// +optional
	SecretKey string `json:"secretKey,omitempty"`
}

// TimGcpSmSecretStatus defines the observed state
type TimGcpSmSecretStatus struct {
	// +optional
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
	// +optional
	SecretHash string `json:"secretHash,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// +optional
	RetryCount int `json:"retryCount,omitempty"`
	// +optional
	LastError string `json:"lastError,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=tgs

// TimGcpSmSecret is the Schema for the timgcpsmsecrets API
type TimGcpSmSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TimGcpSmSecretSpec   `json:"spec,omitempty"`
	Status TimGcpSmSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TimGcpSmSecretList contains a list of TimGcpSmSecret
type TimGcpSmSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimGcpSmSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TimGcpSmSecret{}, &TimGcpSmSecretList{})
}
