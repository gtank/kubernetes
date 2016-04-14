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

package certificates

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/certificates"
	"k8s.io/kubernetes/pkg/apis/certificates/validation"
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
	csr := obj.(*certificates.CertificateSigningRequest)

	// Users cannot create any derived information, but we expect
	// information about the requesting user to be injected by the registry
	// interface. Clear everything else.
	// TODO: check these ordering assumptions. better way to inject this info?
	oldStatus := csr.Status
	csr.Status = certificates.CertificateSigningRequestStatus{}
	csr.Status.Username = oldStatus.Username
	csr.Status.UID = oldStatus.UID
	csr.Status.Groups = oldStatus.Groups

	// Users cannot create pre-approved certificate requests.
	csr.Approve = certificates.CertificateSigningRequestApproval{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (csrStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newCSR := obj.(*certificates.CertificateSigningRequest)
	oldCSR := old.(*certificates.CertificateSigningRequest)

	// Certificate requests are immutable after creation.
	newCSR.Spec = oldCSR.Spec

	// The user cannot update these by updating the CSR.
	newCSR.Status = oldCSR.Status
	newCSR.Approve = oldCSR.Approve
}

// Validate validates a new CSR.
func (csrStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	csr := obj.(*certificates.CertificateSigningRequest)
	return validation.ValidateCertificateSigningRequest(csr)
}

// Canonicalize normalizes the object after validation.
func (csrStrategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (csrStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateCertificateSigningRequest(obj.(*certificates.CertificateSigningRequest))
	return append(errorList, validation.ValidateCertificateSigningRequestUpdate(obj.(*certificates.CertificateSigningRequest), old.(*certificates.CertificateSigningRequest))...)
}

// If AllowUnconditionalUpdate() is true and the object specified by
// the user does not have a resource version, then generic Update()
// populates it with the latest version. Else, it checks that the
// version specified by the user matches the version of latest etcd
// object.
func (csrStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (s csrStrategy) Export(obj runtime.Object, exact bool) error {
	csr, ok := obj.(*certificates.CertificateSigningRequest)
	if !ok {
		// unexpected programmer error
		return fmt.Errorf("unexpected object: %v", obj)
	}
	s.PrepareForCreate(obj)
	if exact {
		return nil
	}
	// CSRs allow direct subresource edits, we clear them without exact so the CSR value can be reused.
	csr.Status = certificates.CertificateSigningRequestStatus{}
	csr.Approve = certificates.CertificateSigningRequestApproval{}
	return nil
}

// Storage strategy for the Status subresource
type csrStatusStrategy struct {
	csrStrategy
}

var StatusStrategy = csrStatusStrategy{Strategy}

func (csrStatusStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newCSR = obj.(*certificates.CertificateSigningRequest)
	oldCSR = old.(*certificates.CertificateSigningRequest)

	// Updating the Status should only update the Status.
	newCSR.Spec = oldCSR.Spec
	newCSR.Approve = oldCSR.Approve
}

func (csrStatusStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateCertificateSigningRequestUpdate(obj.(*certificates.CertificateSigningRequest), old.(*certificates.CertificateSigningRequest))
}

// Canonicalize normalizes the object after validation.
func (csrStatusStrategy) Canonicalize(obj runtime.Object) {
	// TODO: derive information in Status here?
}

// Storage strategy for the Approve subresource
type csrApproveStrategy struct {
	csrStrategy
}

var ApproveStrategy = csrApproveStrategy{Strategy}

func (csrApproveStrategy) PrepareForUpdate(obj, old runtime.Object) {
	newCSR = obj.(*certificates.CertificateSigningRequest)
	oldCSR = old.(*certificates.CertificateSigningRequest)

	// Updating the Approve should only update the Approve.
	newCSR.Spec = oldCSR.Spec
	newCSR.Status = oldCSR.Status
}

func (csrApproveStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateCertificateSigningRequestUpdate(obj.(*certificates.CertificateSigningRequest), old.(*certificates.CertificateSigningRequest))
}

// Matcher returns a generic matcher for a given label and field selector.
func Matcher(label labels.Selector, field fields.Selector) generic.Matcher {
	return generic.MatcherFunc(func(obj runtime.Object) (bool, error) {
		sa, ok := obj.(*certificates.CertificateSigningRequest)
		if !ok {
			return false, fmt.Errorf("not a CertificateSigningRequest")
		}
		fields := SelectableFields(sa)
		return label.Matches(labels.Set(sa.Labels)) && field.Matches(fields), nil
	})
}

// SelectableFields returns a label set that can be used for filter selection
func SelectableFields(obj *certificates.CertificateSigningRequest) labels.Set {
	return labels.Set{}
}
