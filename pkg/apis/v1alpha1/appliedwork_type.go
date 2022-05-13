/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// AppliedWorkSpec represents the desired configuration of AppliedWork
type AppliedWorkSpec struct {
	// ManifestWorkName represents the name of the related manifestwork on the hub.
	// +kubebuilder:validation:Required
	// +required
	ManifestWorkName string `json:"manifestWorkName"`
}

// AppliedtWorkStatus represents the current status of AppliedWork
type AppliedtWorkStatus struct {
	// AppliedResources represents a list of resources defined within the manifestwork that are applied.
	// Only resources with valid GroupVersionResource, namespace, and name are suitable.
	// An item in this slice is deleted when there is no mapped manifest in manifestwork.Spec or by finalizer.
	// The resource relating to the item will also be removed from managed cluster.
	// The deleted resource may still be present until the finalizers for that resource are finished.
	// However, the resource will not be undeleted, so it can be removed from this list and eventual consistency is preserved.
	// +optional
	AppliedResources []AppliedManifestResourceMeta `json:"appliedResources,omitempty"`
}

// AppliedManifestResourceMeta represents the group, version, resource, name and namespace of a resource.
// Since these resources have been created, they must have valid group, version, resource, namespace, and name.
type AppliedManifestResourceMeta struct {
	// Group is the API Group of the Kubernetes resource.
	// +required
	Group string `json:"group"`

	// Version is the version of the Kubernetes resource.
	// +required
	Version string `json:"version"`

	// Resource is the resource name of the Kubernetes resource.
	// +required
	Resource string `json:"resource"`

	// Name is the name of the Kubernetes resource.
	// +required
	Name string `json:"name"`

	// UID is set on successful deletion of the Kubernetes resource by controller. The
	// resource might be still visible on the managed cluster after this field is set.
	// It is not directly settable by a client.
	// +optional
	UID types.UID `json:"uid,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={fleet}
// +kubebuilder:object:root=true

// AppliedWork represents an applied work on managed cluster that is placed
// on a managed cluster. An appliedwork links to a work on a hub recording resources
// deployed in the managed cluster.
// When the agent is removed from managed cluster, cluster-admin on managed cluster
// can delete appliedmanifestwork to remove resources deployed by the agent.
// The name of the appliedwork must be the same as {manifestwork name}
// The namespace of the appliedwork should be the same as the resource applied on
// the managed cluster.
type AppliedWork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents the desired configuration of AppliedManifestWork.
	// +kubebuilder:validation:Required
	// +required
	Spec AppliedWorkSpec `json:"spec"`

	// Status represents the current status of AppliedManifestWork.
	// +optional
	Status AppliedtWorkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppliedWorkList contains a list of AppliedWork
type AppliedWorkList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of works.
	// +listType=set
	Items []AppliedWork `json:"items"`
}
