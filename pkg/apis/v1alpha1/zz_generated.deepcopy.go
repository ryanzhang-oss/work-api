//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestResourceMeta) DeepCopyInto(out *AppliedManifestResourceMeta) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestResourceMeta.
func (in *AppliedManifestResourceMeta) DeepCopy() *AppliedManifestResourceMeta {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestResourceMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedWork) DeepCopyInto(out *AppliedWork) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedWork.
func (in *AppliedWork) DeepCopy() *AppliedWork {
	if in == nil {
		return nil
	}
	out := new(AppliedWork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppliedWork) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedWorkList) DeepCopyInto(out *AppliedWorkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppliedWork, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedWorkList.
func (in *AppliedWorkList) DeepCopy() *AppliedWorkList {
	if in == nil {
		return nil
	}
	out := new(AppliedWorkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppliedWorkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedWorkSpec) DeepCopyInto(out *AppliedWorkSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedWorkSpec.
func (in *AppliedWorkSpec) DeepCopy() *AppliedWorkSpec {
	if in == nil {
		return nil
	}
	out := new(AppliedWorkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedtWorkStatus) DeepCopyInto(out *AppliedtWorkStatus) {
	*out = *in
	if in.AppliedResources != nil {
		in, out := &in.AppliedResources, &out.AppliedResources
		*out = make([]AppliedManifestResourceMeta, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedtWorkStatus.
func (in *AppliedtWorkStatus) DeepCopy() *AppliedtWorkStatus {
	if in == nil {
		return nil
	}
	out := new(AppliedtWorkStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Manifest) DeepCopyInto(out *Manifest) {
	*out = *in
	in.RawExtension.DeepCopyInto(&out.RawExtension)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Manifest.
func (in *Manifest) DeepCopy() *Manifest {
	if in == nil {
		return nil
	}
	out := new(Manifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestCondition) DeepCopyInto(out *ManifestCondition) {
	*out = *in
	out.Identifier = in.Identifier
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestCondition.
func (in *ManifestCondition) DeepCopy() *ManifestCondition {
	if in == nil {
		return nil
	}
	out := new(ManifestCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceIdentifier) DeepCopyInto(out *ResourceIdentifier) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceIdentifier.
func (in *ResourceIdentifier) DeepCopy() *ResourceIdentifier {
	if in == nil {
		return nil
	}
	out := new(ResourceIdentifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Work) DeepCopyInto(out *Work) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Work.
func (in *Work) DeepCopy() *Work {
	if in == nil {
		return nil
	}
	out := new(Work)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Work) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkList) DeepCopyInto(out *WorkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Work, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkList.
func (in *WorkList) DeepCopy() *WorkList {
	if in == nil {
		return nil
	}
	out := new(WorkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkSpec) DeepCopyInto(out *WorkSpec) {
	*out = *in
	in.Workload.DeepCopyInto(&out.Workload)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkSpec.
func (in *WorkSpec) DeepCopy() *WorkSpec {
	if in == nil {
		return nil
	}
	out := new(WorkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkStatus) DeepCopyInto(out *WorkStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ManifestConditions != nil {
		in, out := &in.ManifestConditions, &out.ManifestConditions
		*out = make([]ManifestCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkStatus.
func (in *WorkStatus) DeepCopy() *WorkStatus {
	if in == nil {
		return nil
	}
	out := new(WorkStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadTemplate) DeepCopyInto(out *WorkloadTemplate) {
	*out = *in
	if in.Manifests != nil {
		in, out := &in.Manifests, &out.Manifests
		*out = make([]Manifest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadTemplate.
func (in *WorkloadTemplate) DeepCopy() *WorkloadTemplate {
	if in == nil {
		return nil
	}
	out := new(WorkloadTemplate)
	in.DeepCopyInto(out)
	return out
}
