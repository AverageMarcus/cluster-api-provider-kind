/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

// KindClusterPhase indicates the current phase of a clusters life
type KindClusterPhase string

var (
	// KindClusterPhasePending is the initial phase used before a KindCluster has been reconciled
	KindClusterPhasePending KindClusterPhase = ""
	// KindClusterPhaseCreating is the phase used during cluster creation
	KindClusterPhaseCreating KindClusterPhase = "Creating"
	// KindClusterPhaseReady is the phase used when the cluster is created and ready
	KindClusterPhaseReady KindClusterPhase = "Ready"
	// KindClusterPhaseDeleting is the phase used when deleting an existing cluster
	KindClusterPhaseDeleting KindClusterPhase = "Deleting"
)

// FailureReason contains machine-readable details of what error occurred
type FailureReason string

var (
	// FailureReasonCreateFailed indicates there was an error during cluster creation
	FailureReasonCreateFailed FailureReason = "CreateFailed"
	// FailureReasonKubeConfig indicates there was an error retrieving the KubeConfig for the cluster
	FailureReasonKubeConfig FailureReason = "KubeConfigNotFound"
	// FailureReasonEndpoint indicates there was an error retrieving and parsing the control plane endpoint details
	FailureReasonEndpoint FailureReason = "ControlPlaneEndpointInvalid"
	// FailureReasonClusterNotFound indicates there was an error trying to find the cluster in Kind
	FailureReasonClusterNotFound FailureReason = "ClusterNotFound"
	// FailureReasonDeleteFailed indicates there was an error while deleting a cluster
	FailureReasonDeleteFailed FailureReason = "DeleteFailed"
)

// KindClusterSpec defines the desired state of KindCluster
type KindClusterSpec struct {
	// Image is the node image used for the cluster nodes
	//
	// +kubebuilder:default=kindest/node
	Image string `json:"image,omitempty"`

	// Version is the Kubernetes version to use (e.g. v1.21.2)
	//
	// +kubebuilder:default=v1.21.2
	// +kubebuilder:validation:Pattern=^v\d\.\d+\.\d+$
	Version string `json:"version,omitempty"`

	// Replicas controls the number of control plane nodes to create
	//
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// FeatureGates enables or disabled Kubernetes feature gates
	//
	// See https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
	// for the available features.
	FeatureGates map[string]bool `json:"featureGates,omitempty"`

	// RuntimeConfig allows enabling or disabling built-in APIs.
	//
	// See https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/
	// for the available values.
	RuntimeConfig map[string]string `json:"runtimeConfig,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// KindClusterStatus defines the observed state of KindCluster
type KindClusterStatus struct {
	// Ready indicates if the cluster is ready to use or not
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// Phase contains details on the current phase of the cluster (e.g. creating, ready, deleting)
	// +optional
	Phase *KindClusterPhase `json:"phase"`

	// KubeConfig contains the KubeConfig to use to communicate with the cluster
	// +optional
	KubeConfig *string `json:"kubeConfig,omitempty"`

	// FailureReason indicates there is a fatal problem reconciling the infrastructure
	// suitable for programmatic interpretation
	// +optional
	FailureReason *FailureReason `json:"failureReason"`

	// FailureMessage indicates there is a fatal problem reconciling the infrastructure
	// descriptive interpretation
	// +optional
	FailureMessage *string `json:"failureMessage"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KindCluster is the Schema for the kindclusters API
type KindCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KindClusterSpec   `json:"spec,omitempty"`
	Status KindClusterStatus `json:"status,omitempty"`
}

// NamespacedName returns the KindCluster name prefixed with the namespace
func (kc *KindCluster) NamespacedName() string {
	return fmt.Sprintf("%s-%s", kc.Namespace, kc.Name)
}

//+kubebuilder:object:root=true

// KindClusterList contains a list of KindCluster
type KindClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KindCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KindCluster{}, &KindClusterList{})
}
