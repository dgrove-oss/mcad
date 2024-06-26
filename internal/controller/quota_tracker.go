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
	"strings"

	mcadv1beta1 "github.com/project-codeflare/mcad/api/v1beta1"
	v1 "k8s.io/api/core/v1"
)

// Should be defined in api/core/v1/types.go
const DefaultResourceLimitsPrefix = "limits."

// A tracker of allocated quota, mapped by namespace
type QuotaTracker struct {
	// state of quotas, used, and allocated amounts
	state map[string]*QuotaState

	// used amounts by dispatched AppWrappers partially unaccounted by the ResourceQuota,
	// as some pods may have not passed the ResourceQuota admission controller
	unAdmittedWeightsMap map[string]*WeightsPair
}

// Create a new QuotaTracker
func NewQuotaTracker() *QuotaTracker {
	return &QuotaTracker{
		state:                map[string]*QuotaState{},
		unAdmittedWeightsMap: map[string]*WeightsPair{},
	}
}

// State includes total quota, used quota, and currently allocated quota
type QuotaState struct {
	// quota enforced in the ResourceQuota object
	quota *WeightsPair
	// used amount in the status of the ResourceQuota object
	used *WeightsPair
	// allocated amount by dispatched AppWrappers in the current dispatching cycle
	allocated *WeightsPair
}

// Create a QuotaState from a ResourceQuota object
func NewQuotaStateFromResourceQuota(resourceQuota *v1.ResourceQuota) *QuotaState {
	quotaWeights, usedWeights := getQuotaAndUsedWeightsPairsForResourceQuota(resourceQuota)
	return &QuotaState{
		quota:     quotaWeights,
		used:      usedWeights,
		allocated: NewWeightsPair(Weights{}, Weights{}),
	}
}

// Account for all in-flight AppWrappers with their resource demand not yet reflected in
// the Used status of any ResourceQuota object in their corresponding namespace
func (tracker *QuotaTracker) Init(weightsPairMap map[string]*WeightsPair) {
	tracker.unAdmittedWeightsMap = weightsPairMap
}

// Check if the resource demand of an AppWrapper satisfies a ResourceQuota,
// without changing the current quota allocation, returning resource names with insufficient quota
func (tracker *QuotaTracker) Satisfies(appWrapperAskWeights *WeightsPair, resourceQuota *v1.ResourceQuota) (bool, []v1.ResourceName) {
	namespace := resourceQuota.GetNamespace()
	var quotaState *QuotaState
	var exists bool
	if quotaState, exists = tracker.state[namespace]; !exists {
		quotaState = NewQuotaStateFromResourceQuota(resourceQuota)
		tracker.state[namespace] = quotaState
	}
	// check if both appwrapper requests and limits fit available resource quota
	quotaWeights := quotaState.quota.Clone()
	quotaWeights.QuotaSub(quotaState.used)
	quotaWeights.QuotaSub(quotaState.allocated)
	var unAdmittedWeights *WeightsPair
	if unAdmittedWeights, exists = tracker.unAdmittedWeightsMap[namespace]; exists {
		quotaWeights.QuotaSub(unAdmittedWeights)
	}
	quotaFits, insufficientResources := appWrapperAskWeights.Fits(quotaWeights)

	// mcadLog.Info("QuotaTracker.Satisfies():", "namespace", namespace,
	// 	"QuotaWeights", quotaState.quota, "UsedWeights", quotaState.used,
	// 	"AllocatedWeights", quotaState.allocated, "unAdmittedWeights", unAdmittedWeights,
	// 	"AvailableWeights", quotaWeights, "appWrapperAskWeights", appWrapperAskWeights,
	// 	"quotaFits", quotaFits)
	return quotaFits, insufficientResources
}

// Update the QuotaState by the allocated weights of an AppWrapper in a namespace,
// fails if QuotaState does not exist in the QuotaTracker
func (tracker *QuotaTracker) Allocate(namespace string, appWrapperAskWeights *WeightsPair) bool {
	if state, exists := tracker.state[namespace]; exists && appWrapperAskWeights != nil {
		state.allocated.Add(appWrapperAskWeights)
		return true
	}
	return false
}

// Get requests and limits from AppWrapper specs
func getWeightsPairForAppWrapper(appWrapper *mcadv1beta1.AppWrapper) *WeightsPair {
	requests := aggregateRequests(appWrapper)
	limits := aggregateLimits(appWrapper)
	return NewWeightsPair(requests, limits)
}

// Get requests and limits for both quota and used from ResourceQuota object
func getQuotaAndUsedWeightsPairsForResourceQuota(resourceQuota *v1.ResourceQuota) (quotaWeights *WeightsPair,
	usedWeights *WeightsPair) {
	quotaWeights = getWeightsPairForResourceList(&resourceQuota.Status.Hard)
	usedWeights = getWeightsPairForResourceList(&resourceQuota.Status.Used)
	return quotaWeights, usedWeights
}

// Create a pair of Weights for requests and limits
// given in a ResourceList of a ResourceQuota
func getWeightsPairForResourceList(r *v1.ResourceList) *WeightsPair {
	requests := Weights{}
	limits := Weights{}
	for k, v := range *r {
		if strings.HasPrefix(k.String(), DefaultResourceLimitsPrefix) {
			trimmedName := strings.Replace(k.String(), DefaultResourceLimitsPrefix, "", 1)
			if _, exists := limits[v1.ResourceName(trimmedName)]; !exists {
				limits[v1.ResourceName(trimmedName)] = v.AsDec()
			}
			continue
		}
		if strings.HasPrefix(k.String(), v1.DefaultResourceRequestsPrefix) {
			trimmedName := strings.Replace(k.String(), v1.DefaultResourceRequestsPrefix, "", 1)
			k = v1.ResourceName(trimmedName)
		}
		// in case of two keys: requests.xxx and xxx, take the minimum quota of the two
		if value, exists := requests[k]; !exists || value.Cmp(v.AsDec()) > 0 {
			requests[k] = v.AsDec()
		}
	}
	return NewWeightsPair(requests, limits)
}
