/*
Copyright 2023 IBM Corporation.

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

package controller

import (
	"bytes"
	"fmt"

	"gopkg.in/inf.v0"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Weights represent a set of resource requests or available resources
// Quantities are encoded as *inf.Dec to maintain precision and make arithmetic easy
type Weights map[v1.ResourceName]*inf.Dec

// Converts a ResourceList to Weights
func NewWeights(r v1.ResourceList) Weights {
	w := Weights{}
	for k, v := range r {
		w[k] = v.AsDec() // should be lossless
	}
	return w
}

// Converts effective resources of a pod to Weights
func NewWeightsForPod(pod *v1.Pod) Weights {
	podRequest := Weights{}
	// add up resources of all containers
	for _, container := range pod.Spec.Containers {
		podRequest.Add(NewWeights(container.Resources.Requests))
	}
	// take max(sum_pod, any_init_container)
	for _, initContainer := range pod.Spec.InitContainers {
		podRequest.Max(NewWeights(initContainer.Resources.Requests))
	}
	// add any pod overhead
	if pod.Spec.Overhead != nil {
		podRequest.Add(NewWeights(pod.Spec.Overhead))
	}
	return podRequest
}

// Add weights to receiver
func (w Weights) Add(r Weights) {
	for k, v := range r {
		if w[k] == nil {
			w[k] = &inf.Dec{} // fresh zero
		}
		w[k].Add(w[k], v)
	}
}

// Subtract weights from receiver
func (w Weights) Sub(r Weights) {
	for k, v := range r {
		if w[k] == nil {
			w[k] = &inf.Dec{} // fresh zero
		}
		w[k].Sub(w[k], v)
	}
}

// Subtract use from quota
func (w Weights) QuotaSub(r Weights) {
	for k, v := range r {
		if w[k] == nil {
			continue // ignore undefined quota
		}
		w[k].Sub(w[k], v)
	}
}

// Add coefficient * weights to receiver
func (w Weights) AddProd(coefficient int32, r Weights) {
	for k, v := range r {
		if w[k] == nil {
			w[k] = &inf.Dec{} // fresh zero
		}
		tmp := inf.NewDec(int64(coefficient), 0)
		tmp.Mul(tmp, v)
		w[k].Add(w[k], tmp)
	}
}

// Update receiver to max of receiver and argument in each dimension
func (w Weights) Max(r Weights) {
	for k, v := range r {
		if w[k] == nil {
			w[k] = &inf.Dec{} // fresh zero
		}
		if w[k].Cmp(v) == -1 {
			w[k].Set(v) // w[k] = v would not be correct due to aliasing
		}
	}
}

// Compare receiver to argument
// True if receiver is less than or equal to argument in every dimension
func (w Weights) Fits(r Weights) (bool, []v1.ResourceName) {
	insufficient := []v1.ResourceName{}
	zero := &inf.Dec{}    // shared zero, never mutated
	for k, v := range w { // range over receiver not argument
		// ignore 0 requests in case r does not contain k
		if v.Cmp(zero) <= 0 {
			continue
		}
		// v > 0 so r[k] must be defined and no less than v
		if r[k] == nil || v.Cmp(r[k]) == 1 {
			insufficient = append(insufficient, k)
		}
	}
	if len(insufficient) == 0 {
		return true, nil
	} else {
		return false, insufficient
	}
}

// Compare receiver to argument
// True if receiver is less than or equal to argument in every dimension where argument is defined
func (w Weights) QuotaFits(r Weights) (bool, []v1.ResourceName) {
	insufficient := []v1.ResourceName{}
	zero := &inf.Dec{}    // shared zero, never mutated
	for k, v := range w { // range over receiver not argument
		// ignore 0 requests or no quota
		if v.Cmp(zero) <= 0 || r[k] == nil {
			continue
		}
		// v > 0 so r[k] must be no less than v
		if v.Cmp(r[k]) == 1 {
			insufficient = append(insufficient, k)
		}
	}
	if len(insufficient) == 0 {
		return true, nil
	} else {
		return false, insufficient
	}
}

// Converts Weights to a ResourceList
func (w Weights) AsResources() v1.ResourceList {
	resources := v1.ResourceList{}
	for k, v := range w {
		resources[k] = *resource.NewDecimalQuantity(*v, resource.DecimalSI)
	}
	return resources
}

// Creates a (deep) copy of Weights
// (Easier to add to zero than creating all fields)
func (w Weights) Clone() Weights {
	r := Weights{}
	r.Add(w)
	return r
}

// A pair of Weights used to represent requests and limits for an object
type WeightsPair struct {
	requests Weights
	limits   Weights
}

// Creates a new pair of weights
func NewWeightsPair(requests Weights, limits Weights) *WeightsPair {
	return &WeightsPair{
		requests: requests,
		limits:   limits,
	}
}

// Add pair of weights to receiver
func (w *WeightsPair) Add(r *WeightsPair) {
	w.requests.Add(r.requests)
	w.limits.Add(r.limits)
}

// Subtract pair of weights from receiver
func (w *WeightsPair) QuotaSub(r *WeightsPair) {
	w.requests.QuotaSub(r.requests)
	w.limits.QuotaSub(r.limits)
}

// Max of two pairs of weights
func (w *WeightsPair) Max(r *WeightsPair) {
	w.requests.Max(r.requests)
	w.limits.Max(r.limits)
}

// Clone a pair of weights
func (w *WeightsPair) Clone() *WeightsPair {
	requestsClone := w.requests.Clone()
	limitsClone := w.limits.Clone()
	return NewWeightsPair(requestsClone, limitsClone)
}

// Compare receiver to argument -
// True if both weights of receiver fit corresponding weights of argument.
// If False, return list of insufficient resource names
func (w *WeightsPair) Fits(r *WeightsPair) (bool, []v1.ResourceName) {
	insufficient := []v1.ResourceName{}
	requestsFits, requestsInsufficient := w.requests.QuotaFits(r.requests)
	limitsFits, limitsInsufficient := w.limits.QuotaFits(r.limits)
	if requestsFits && limitsFits {
		return true, insufficient
	}
	insufficient = append(insufficient, requestsInsufficient...)
	insufficient = append(insufficient, limitsInsufficient...)
	return false, RemoveDuplicateResources(insufficient)
}

func (w *WeightsPair) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{Requests: %v}; ", w.requests)
	fmt.Fprintf(&b, "{Limits: %v}", w.limits)
	return b.String()
}

// Create new WeightsPair for a pod
func NewWeightsPairForPod(pod *v1.Pod) *WeightsPair {
	return NewWeightsPair(NewWeightsForPod(pod), NewLimitsWeightsForPod(pod))
}

// Converts resource limits of a pod to Weights
func NewLimitsWeightsForPod(pod *v1.Pod) Weights {
	podLimit := Weights{}
	// add up limits of resources of all containers
	for _, container := range pod.Spec.Containers {
		podLimit.Add(NewWeights(container.Resources.Limits))
	}
	// take max(sum_pod, any_init_container)
	for _, initContainer := range pod.Spec.InitContainers {
		podLimit.Max(NewWeights(initContainer.Resources.Limits))
	}
	return podLimit
}

// Remove duplicate resource names in a slice
func RemoveDuplicateResources(slice []v1.ResourceName) []v1.ResourceName {
	unique := make(map[v1.ResourceName]bool)
	result := []v1.ResourceName{}
	for _, val := range slice {
		if _, ok := unique[val]; !ok {
			unique[val] = true
			result = append(result, val)
		}
	}
	return result
}
