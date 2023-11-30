//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapper) DeepCopyInto(out *AppWrapper) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapper.
func (in *AppWrapper) DeepCopy() *AppWrapper {
	if in == nil {
		return nil
	}
	out := new(AppWrapper)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppWrapper) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperList) DeepCopyInto(out *AppWrapperList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppWrapper, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperList.
func (in *AppWrapperList) DeepCopy() *AppWrapperList {
	if in == nil {
		return nil
	}
	out := new(AppWrapperList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppWrapperList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperResources) DeepCopyInto(out *AppWrapperResources) {
	*out = *in
	if in.GenericItems != nil {
		in, out := &in.GenericItems, &out.GenericItems
		*out = make([]GenericItem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperResources.
func (in *AppWrapperResources) DeepCopy() *AppWrapperResources {
	if in == nil {
		return nil
	}
	out := new(AppWrapperResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperService) DeepCopyInto(out *AppWrapperService) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperService.
func (in *AppWrapperService) DeepCopy() *AppWrapperService {
	if in == nil {
		return nil
	}
	out := new(AppWrapperService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperSpec) DeepCopyInto(out *AppWrapperSpec) {
	*out = *in
	out.DoNotUsePrioritySlope = in.DoNotUsePrioritySlope.DeepCopy()
	in.DoNotUseService.DeepCopyInto(&out.DoNotUseService)
	in.Resources.DeepCopyInto(&out.Resources)
	if in.DoNotUseSelector != nil {
		in, out := &in.DoNotUseSelector, &out.DoNotUseSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	in.Scheduling.DeepCopyInto(&out.Scheduling)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperSpec.
func (in *AppWrapperSpec) DeepCopy() *AppWrapperSpec {
	if in == nil {
		return nil
	}
	out := new(AppWrapperSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperStatus) DeepCopyInto(out *AppWrapperStatus) {
	*out = *in
	in.DispatchTimestamp.DeepCopyInto(&out.DispatchTimestamp)
	in.RequeueTimestamp.DeepCopyInto(&out.RequeueTimestamp)
	if in.Transitions != nil {
		in, out := &in.Transitions, &out.Transitions
		*out = make([]AppWrapperTransition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperStatus.
func (in *AppWrapperStatus) DeepCopy() *AppWrapperStatus {
	if in == nil {
		return nil
	}
	out := new(AppWrapperStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppWrapperTransition) DeepCopyInto(out *AppWrapperTransition) {
	*out = *in
	in.Time.DeepCopyInto(&out.Time)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppWrapperTransition.
func (in *AppWrapperTransition) DeepCopy() *AppWrapperTransition {
	if in == nil {
		return nil
	}
	out := new(AppWrapperTransition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomPodResource) DeepCopyInto(out *CustomPodResource) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make(corev1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	if in.DoNotUseLimits != nil {
		in, out := &in.DoNotUseLimits, &out.DoNotUseLimits
		*out = make(corev1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomPodResource.
func (in *CustomPodResource) DeepCopy() *CustomPodResource {
	if in == nil {
		return nil
	}
	out := new(CustomPodResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoNotUseDispatchDurationSpec) DeepCopyInto(out *DoNotUseDispatchDurationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoNotUseDispatchDurationSpec.
func (in *DoNotUseDispatchDurationSpec) DeepCopy() *DoNotUseDispatchDurationSpec {
	if in == nil {
		return nil
	}
	out := new(DoNotUseDispatchDurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericItem) DeepCopyInto(out *GenericItem) {
	*out = *in
	if in.DoNotUseMinAvailable != nil {
		in, out := &in.DoNotUseMinAvailable, &out.DoNotUseMinAvailable
		*out = new(int32)
		**out = **in
	}
	out.DoNotUsePrioritySlope = in.DoNotUsePrioritySlope.DeepCopy()
	in.GenericTemplate.DeepCopyInto(&out.GenericTemplate)
	if in.CustomPodResources != nil {
		in, out := &in.CustomPodResources, &out.CustomPodResources
		*out = make([]CustomPodResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericItem.
func (in *GenericItem) DeepCopy() *GenericItem {
	if in == nil {
		return nil
	}
	out := new(GenericItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequeuingSpec) DeepCopyInto(out *RequeuingSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequeuingSpec.
func (in *RequeuingSpec) DeepCopy() *RequeuingSpec {
	if in == nil {
		return nil
	}
	out := new(RequeuingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchedulingSpec) DeepCopyInto(out *SchedulingSpec) {
	*out = *in
	if in.DoNotUseNodeSelector != nil {
		in, out := &in.DoNotUseNodeSelector, &out.DoNotUseNodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.Requeuing = in.Requeuing
	out.DoNotUseDispatchDuration = in.DoNotUseDispatchDuration
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchedulingSpec.
func (in *SchedulingSpec) DeepCopy() *SchedulingSpec {
	if in == nil {
		return nil
	}
	out := new(SchedulingSpec)
	in.DeepCopyInto(out)
	return out
}
