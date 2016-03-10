/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package csr

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/extensions/validation"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// csrStrategy implements behavior for CSRs
type csrStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// csrStrategy is the default logic that applies when creating and updating
// CSR objects.
var Strategy = csrStrategy{api.Scheme, api.SimpleNameGenerator}

// NamespaceScoped is true for CSRs.
func (csrStrategy) NamespaceScoped() bool {
	return true
}

// AllowCreateOnUpdate is false for CSRs.
func (csrStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (csrStrategy) PrepareForCreate(obj runtime.Object) {
	_ = obj.(*extensions.CertificateSigningRequest)
	//TODO(gtank): Restrict derived fields & possibly derive them here.
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (csrStrategy) PrepareForUpdate(obj, old runtime.Object) {
	_ = obj.(*extensions.CertificateSigningRequest)
	_ = old.(*extensions.CertificateSigningRequest)
}

// Validate validates a new CSR.
func (csrStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	csr := obj.(*extensions.CertificateSigningRequest)
	return validation.ValidateCertificateSigningRequest(csr)
}

// Canonicalize normalizes the object after validation.
func (csrStrategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (csrStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateCertificateSigningRequest(obj.(*extensions.CertificateSigningRequest))
	return append(errorList, validation.ValidateCertificateSigningRequestUpdate(obj.(*extensions.CertificateSigningRequest), old.(*extensions.CertificateSigningRequest))...)
}

func (csrStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (s csrStrategy) Export(obj runtime.Object, exact bool) error {
	csr, ok := obj.(*extensions.CertificateSigningRequest)
	if !ok {
		// unexpected programmer error
		return fmt.Errorf("unexpected object: %v", obj)
	}
	s.PrepareForCreate(obj)
	if exact {
		return nil
	}
	// CSRs allow direct status edits, therefore we clear that without
	// exact so that the CSR value can be reused.
	csr.Status = extensions.CertificateSigningRequestStatus{}
	return nil
}

type csrStatusStrategy struct {
	csrStrategy
}

var StatusStrategy = csrStatusStrategy{Strategy}

func (csrStatusStrategy) PrepareForCreate(obj runtime.Object) {
	_ = obj.(*extensions.CertificateSigningRequest)
	// TODO: CSRs should not allow status to be set on create.
}

func (csrStatusStrategy) PrepareForUpdate(obj, old runtime.Object) {
	_ = obj.(*extensions.CertificateSigningRequest)
	_ = old.(*extensions.CertificateSigningRequest)
}

func (csrStatusStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateCertificateSigningRequestUpdate(obj.(*extensions.CertificateSigningRequest), old.(*extensions.CertificateSigningRequest))
}

// Canonicalize normalizes the object after validation.
func (csrStatusStrategy) Canonicalize(obj runtime.Object) {
}

// Matcher returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return generic.MatcherFunc(func(obj runtime.Object) (bool, error) {
		sa, ok := obj.(*extensions.CertificateSigningRequest)
		if !ok {
			return false, fmt.Errorf("not a CertificateSigningRequest")
		}
		fields := SelectableFields(sa)
		return label.Matches(labels.Set(sa.Labels)) && field.Matches(fields), nil
	})
}

// SelectableFields returns a label set that can be used for filter selection
func SelectableFields(obj *extensions.CertificateSigningRequest) labels.Set {
	return labels.Set{}
}
