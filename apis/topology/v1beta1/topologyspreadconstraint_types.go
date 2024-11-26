/*
Copyright 2023.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"context"
)

// TopologySpreadConstraintSpec defines the desired state of TopologySpreadConstraint
type TopologySpreadConstraintSpec struct {
	// Name is a field used to identify the CR referenced by service operators
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	// APITopologySpreadConstraint exposes topologySpreadConstraint that are
	// applied to the StatefulSet
	TopologySpreadConstraint *[]corev1.TopologySpreadConstraint `json:"topologySpreadConstraint,omitempty"`
}

// TopologySpreadConstraintStatus defines the observed state of TopologySpreadConstraint
type TopologySpreadConstraintStatus struct {
	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// Hash of the topology configuration
	Hash string `json:"hash,omitempty"`

	// ObservedGeneration - the most recent generation observed for this
	// service. If the observed generation is less than the spec generation,
	// then the controller has not processed the latest changes injected by
	// the opentack-operator in the top-level CR (e.g. the ContainerImage)
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[0].status",description="Ready"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// TopologySpreadConstraint is the Schema for the topologyspreadconstraints API
type TopologySpreadConstraint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologySpreadConstraintSpec   `json:"spec,omitempty"`
	Status TopologySpreadConstraintStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TopologySpreadConstraintList contains a list of TopologySpreadConstraint
type TopologySpreadConstraintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopologySpreadConstraint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopologySpreadConstraint{}, &TopologySpreadConstraintList{})
}

// IsReady returns true if TopologySpreadConstraint reconciled successfully
func (instance TopologySpreadConstraint) IsReady() bool {
	return instance.Status.Conditions.IsTrue(condition.ReadyCondition)
}

// GetTopologyConstraintWithName - a function exposed to the service operators
// that need to retrieve the referenced topologyConstraint by name
func GetTopologyConstraintWithName(
	ctx context.Context,
	h *helper.Helper,
	name string,
	namespace string,
) (*TopologySpreadConstraint, string, error) {

	topologyConstraint := &TopologySpreadConstraint{}
	err := h.GetClient().Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, topologyConstraint)
	if err != nil {
		return topologyConstraint, "", err
	}
	return topologyConstraint, topologyConstraint.Status.Hash, nil
}


// AddFinalizer -
func (t *TopologySpreadConstraint) AddFinalizer(
		ctx context.Context,
		h *helper.Helper,
		caller string,
) error {
	if !hasFinalizer(t.GetFinalizers(), caller) {
		// Add finalizer to the resource because it is going to be consumed by GlanceAPI
		t.SetFinalizers(append(t.GetFinalizers(), caller))
		// Update the resource
		if err := h.GetClient().Update(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

// hasFinalizer -
// NOTE: This might be a utility function provided by lib-common
func hasFinalizer(
	finalizers []string,
	caller string,
) bool {
	for _, f := range finalizers {
		if f == caller {
			return true
		}
	}
	return false
}
