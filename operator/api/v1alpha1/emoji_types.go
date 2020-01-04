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

// EmojiSpec defines the desired state of Emoji
type EmojiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Description allows to add a description to an Emoji.
	Description string `json:"description,omitempty"`
}

// EmojiStatus defines the observed state of Emoji
type EmojiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Defines if the Emoji is part of the 100 registered Emojis
	// +optional
	Supported *bool `json:"supported,omitempty"`

	// Information when was the last time the Emoji was successfully published.
	// +optional
	LastPublishedTime *metav1.Time `json:"lastPublishedTime,omitempty"`

	// Information when the last time the Emoji was published what has been the output.
	// +optional
	LastPublishedOutput string `json:"lastPublishedOutput,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Emoji is the Schema for the emojis API
type Emoji struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmojiSpec   `json:"spec,omitempty"`
	Status EmojiStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmojiList contains a list of Emoji
type EmojiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Emoji `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Emoji{}, &EmojiList{})
}
