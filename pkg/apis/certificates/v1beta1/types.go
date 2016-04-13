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

package v1beta1

import (
	"crypto/x509/pkix"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

// +genclient=true

// Describes a certificate signing request
type CertificateSigningRequest struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

	// The certificate request itself and any additonal information.
	Spec CertificateSigningRequestSpec `json:"spec,omitempty"`

	// Derived information about the request.
	Status CertificateSigningRequestStatus `json:"status,omitempty"`

	// The current approval state of the request.
	Approve CertificateSigningRequestApproval `json:"approve,omitempty"`
}

// This information is immutable after the request is created.
type CertificateSigningRequestSpec struct {
	// Raw PKCS#10 CSR data
	CertificateRequest string `json:"request"`

	// Any extra information the node wishes to send with the request.
	ExtraInfo []string `json:"extra,omitempty"`
}

// This information is derived from the request by Kubernetes and cannot be
// modified by users. All information is optional since it might not be
// available in the underlying request. This is intented to aid approval
// decisions.
type CertificateSigningRequestStatus struct {
	// Information about the requesting user (if relevant)
	// See user.Info interface for details
	Username string   `json:"username,omitempty"`
	UID      string   `json:"uid,omitempty"`
	Groups   []string `json:"groups,omitempty"`

	// Fingerprint of the public key in request
	Fingerprint string `json:"fingerprint,omitempty"`

	// Subject fields from the request
	Subject pkix.Name `json:"subject,omitempty"`

	// DNS SANs from the request
	Hostnames []string `json:"dns,omitempty"`

	// IP SANs from the request
	IPAddresses []string `json:"ip,omitempty"`
}

type CertificateSigningRequestApproval struct {
	// CSR approval state, one of Submitted, Approved, or Denied
	State CertificateRequestState `json:"state"`

	// brief reason for the request state
	Reason string `json:"reason,omitempty"`
	// human readable message with details about the request state
	Message string `json:"message,omitempty"`

	// If request was approved, the controller will place the issued certificate here.
	Certificate []byte `json:"certificate,omitempty"`
}

type CertificateRequestState string

// These are the possible states for a certificate request.
const (
	RequestSubmitted CertificateRequestState = "Submitted"
	RequestApproved  CertificateRequestState = "Approved"
	RequestDenied    CertificateRequestState = "Denied"
)

type CertificateSigningRequestList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`

	Items []CertificateSigningRequest `json:"items,omitempty"`
}
