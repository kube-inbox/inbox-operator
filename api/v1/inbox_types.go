/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InboxType defines the type of inbox
type InboxType string

const (
	// TextInbox represents a text-based inbox
	TextInbox InboxType = "text"
	// InteractiveInbox represents an interactive inbox
	InteractiveInbox InboxType = "interactive"
)

// InboxSpec defines the desired state of Inbox.
type InboxSpec struct {
	// InboxType specifies the type of inbox (text or interactive)
	// +kubebuilder:validation:Enum=text;interactive
	// +kubebuilder:default=text
	InboxType InboxType `json:"inboxType,omitempty"`
}

// InboxStatus defines the observed state of Inbox.
type InboxStatus struct {
	// Phase represents the current phase of the inbox (Pending, Ready, Failed,Deleting)
	// +kubebuilder:validation:Enum=Pending;Ready;Failed;Deleting
	// +kubebuilder:default=Pending
	Phase string `json:"phase,omitempty"`

	// LastUpdated represents the last time this inbox was updated
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=ibx
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.inboxType`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Last Updated",type=date,JSONPath=`.status.lastUpdated`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Inbox is the Schema for the inboxes API.
type Inbox struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InboxSpec   `json:"spec,omitempty"`
	Status InboxStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// InboxList contains a list of Inbox.
type InboxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Inbox `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Inbox{}, &InboxList{})
}
