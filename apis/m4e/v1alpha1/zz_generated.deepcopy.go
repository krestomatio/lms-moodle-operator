//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Flavor) DeepCopyInto(out *Flavor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Flavor.
func (in *Flavor) DeepCopy() *Flavor {
	if in == nil {
		return nil
	}
	out := new(Flavor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Flavor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlavorList) DeepCopyInto(out *FlavorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Flavor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlavorList.
func (in *FlavorList) DeepCopy() *FlavorList {
	if in == nil {
		return nil
	}
	out := new(FlavorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlavorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlavorSpec) DeepCopyInto(out *FlavorSpec) {
	*out = *in
	in.MoodleSpec.DeepCopyInto(&out.MoodleSpec)
	in.PostgresSpec.DeepCopyInto(&out.PostgresSpec)
	in.NfsSpec.DeepCopyInto(&out.NfsSpec)
	in.KeydbSpec.DeepCopyInto(&out.KeydbSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlavorSpec.
func (in *FlavorSpec) DeepCopy() *FlavorSpec {
	if in == nil {
		return nil
	}
	out := new(FlavorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlavorStatus) DeepCopyInto(out *FlavorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlavorStatus.
func (in *FlavorStatus) DeepCopy() *FlavorStatus {
	if in == nil {
		return nil
	}
	out := new(FlavorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeydbSpec) DeepCopyInto(out *KeydbSpec) {
	*out = *in
	if in.KeydbTolerations != nil {
		in, out := &in.KeydbTolerations, &out.KeydbTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeydbSpec.
func (in *KeydbSpec) DeepCopy() *KeydbSpec {
	if in == nil {
		return nil
	}
	out := new(KeydbSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MoodleConfigProperty) DeepCopyInto(out *MoodleConfigProperty) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MoodleConfigProperty.
func (in *MoodleConfigProperty) DeepCopy() *MoodleConfigProperty {
	if in == nil {
		return nil
	}
	out := new(MoodleConfigProperty)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MoodleSpec) DeepCopyInto(out *MoodleSpec) {
	*out = *in
	if in.MoodleCronjobTolerations != nil {
		in, out := &in.MoodleCronjobTolerations, &out.MoodleCronjobTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MoodleUpdateJobTolerations != nil {
		in, out := &in.MoodleUpdateJobTolerations, &out.MoodleUpdateJobTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MoodleNewInstanceJobTolerations != nil {
		in, out := &in.MoodleNewInstanceJobTolerations, &out.MoodleNewInstanceJobTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.MoodleConfigAdditionalCfg = in.MoodleConfigAdditionalCfg
	if in.NginxTolerations != nil {
		in, out := &in.NginxTolerations, &out.NginxTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PhpFpmTolerations != nil {
		in, out := &in.PhpFpmTolerations, &out.PhpFpmTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.RoutineStatusCrNotify.DeepCopyInto(&out.RoutineStatusCrNotify)
	in.RoutineStatusCrNotifyTermination.DeepCopyInto(&out.RoutineStatusCrNotifyTermination)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MoodleSpec.
func (in *MoodleSpec) DeepCopy() *MoodleSpec {
	if in == nil {
		return nil
	}
	out := new(MoodleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NfsSpec) DeepCopyInto(out *NfsSpec) {
	*out = *in
	if in.GaneshaTolerations != nil {
		in, out := &in.GaneshaTolerations, &out.GaneshaTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NfsSpec.
func (in *NfsSpec) DeepCopy() *NfsSpec {
	if in == nil {
		return nil
	}
	out := new(NfsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresSpec) DeepCopyInto(out *PostgresSpec) {
	*out = *in
	if in.PostgresTolerations != nil {
		in, out := &in.PostgresTolerations, &out.PostgresTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PostgresReadreplicasTolerations != nil {
		in, out := &in.PostgresReadreplicasTolerations, &out.PostgresReadreplicasTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PgbouncerTolerations != nil {
		in, out := &in.PgbouncerTolerations, &out.PgbouncerTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PgbouncerReadonlyTolerations != nil {
		in, out := &in.PgbouncerReadonlyTolerations, &out.PgbouncerReadonlyTolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresSpec.
func (in *PostgresSpec) DeepCopy() *PostgresSpec {
	if in == nil {
		return nil
	}
	out := new(PostgresSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RoutineStatusCrNotify) DeepCopyInto(out *RoutineStatusCrNotify) {
	*out = *in
	if in.StatusCode != nil {
		in, out := &in.StatusCode, &out.StatusCode
		*out = make([]int8, len(*in))
		copy(*out, *in)
	}
	out.Headers = in.Headers
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RoutineStatusCrNotify.
func (in *RoutineStatusCrNotify) DeepCopy() *RoutineStatusCrNotify {
	if in == nil {
		return nil
	}
	out := new(RoutineStatusCrNotify)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RoutineStatusCrNotifyHeaders) DeepCopyInto(out *RoutineStatusCrNotifyHeaders) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RoutineStatusCrNotifyHeaders.
func (in *RoutineStatusCrNotifyHeaders) DeepCopy() *RoutineStatusCrNotifyHeaders {
	if in == nil {
		return nil
	}
	out := new(RoutineStatusCrNotifyHeaders)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Site) DeepCopyInto(out *Site) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Site.
func (in *Site) DeepCopy() *Site {
	if in == nil {
		return nil
	}
	out := new(Site)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Site) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SiteList) DeepCopyInto(out *SiteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Site, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SiteList.
func (in *SiteList) DeepCopy() *SiteList {
	if in == nil {
		return nil
	}
	out := new(SiteList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SiteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SiteSpec) DeepCopyInto(out *SiteSpec) {
	*out = *in
	in.FlavorSpec.DeepCopyInto(&out.FlavorSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SiteSpec.
func (in *SiteSpec) DeepCopy() *SiteSpec {
	if in == nil {
		return nil
	}
	out := new(SiteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SiteStatus) DeepCopyInto(out *SiteStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SiteStatus.
func (in *SiteStatus) DeepCopy() *SiteStatus {
	if in == nil {
		return nil
	}
	out := new(SiteStatus)
	in.DeepCopyInto(out)
	return out
}
