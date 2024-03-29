/*

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TaskLauncherSpec defines the desired state of TaskLauncher
type TaskLauncherSpec struct {
	// Build resolves the image from a build resource.
	Build *Build `json:"build,omitempty"`
}

// TaskLauncherStatus defines the observed state of TaskLauncher
type TaskLauncherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type Build struct {
	// ApplicationRef references an application in this namespace.
	ApplicationRef string `json:"applicationRef,omitempty"`

	// ContainerRef references a container in this namespace.
	ContainerRef string `json:"containerRef,omitempty"`
}

// +kubebuilder:object:root=true

// TaskLauncher is the Schema for the tasklaunchers API
type TaskLauncher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskLauncherSpec   `json:"spec,omitempty"`
	Status TaskLauncherStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskLauncherList contains a list of TaskLauncher
type TaskLauncherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskLauncher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TaskLauncher{}, &TaskLauncherList{})
}
