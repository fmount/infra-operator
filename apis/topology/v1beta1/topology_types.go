/*
Copyright 2024.

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
	"github.com/openstack-k8s-operators/lib-common/modules/common/affinity"
	"github.com/openstack-k8s-operators/lib-common/modules/common/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"context"
)

// TopologySpec defines the desired state of Topology
type TopologySpec struct {
	// +kubebuilder:validation:Optional
	// APITopologySpreadConstraint exposes topologySpreadConstraint that are
	// applied to the StatefulSet
	TopologySpreadConstraint *[]corev1.TopologySpreadConstraint `json:"topologySpreadConstraint,omitempty"`

	// APIAffinity exposes PodAffinity and PodAntiaffinity overrides that are applied
	// to the StatefulSet
	// +optional
	APIAffinity affinity.Overrides `json:"scheduling,omitempty"`
}

// TopologyStatus defines the observed state of Topology
type TopologyStatus struct {
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

// Topology is the Schema for the topologies API
type Topology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologySpec   `json:"spec,omitempty"`
	Status TopologyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TopologyList contains a list of Topology
type TopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Topology `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Topology{}, &TopologyList{})
}

// IsReady returns true if TopologySpreadConstraint reconciled successfully
func (instance Topology) IsReady() bool {
	return instance.Status.Conditions.IsTrue(condition.ReadyCondition)
}

// GetTopologyWithName - a function exposed to the service operators
// that need to retrieve the referenced topology by name
func GetTopologyWithName(
	ctx context.Context,
	h *helper.Helper,
	name string,
	namespace string,
) (*Topology, string, error) {

	topology := &Topology{}
	err := h.GetClient().Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, topology)
	if err != nil {
		return topology, "", err
	}
	return topology, topology.Status.Hash, nil
}


// AddFinalizer -
func (t *Topology) AddFinalizer(
		ctx context.Context,
		h *helper.Helper,
		caller string,
) error {
	if !hasFinalizer(t.GetFinalizers(), caller) {
		// Add finalizer to the resource because it is going to be consumed
		t.SetFinalizers(append(t.GetFinalizers(), caller))
		// Update the resource
		if err := h.GetClient().Update(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

// hasFinalizer -
// NOTE: Move this utility function in lib-common (util pkg)
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
