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

package v1alpha1

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

// +genclient=true

// Describes a certificate signing request
type CertificateSigningRequest struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// The certificate request itself and any additonal information.
	Spec CertificateSigningRequestSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Derived information about the request.
	Status CertificateSigningRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// This information is immutable after the request is created. Only the Request
// and ExtraInfo fields can be set on creation, other fields are derived by
// Kubernetes and cannot be modified by users.
type CertificateSigningRequestSpec struct {
	// Base64-encoded PKCS#10 CSR data
	Request string `json:"request" protobuf:"bytes,1,opt,name=request"`

	// Any extra information the node wishes to send with the request.
	ExtraInfo []string `json:"extrainfo,omitempty" protobuf:"bytes,2,rep,name=extrainfo"`

	// Fingerprint of the public key in request
	Fingerprint string `json:"fingerprint,omitempty" protobuf:"bytes,3,opt,name=fingerprint"`

	// Subject fields from the request
	Subject Subject `json:"subject,omitempty" protobuf:"bytes,4,opt,name=subject"`

	// DNS SANs from the request
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,5,rep,name=hostnames"`

	// IP SANs from the request
	IPAddresses []string `json:"ipaddresses,omitempty" protobuf:"bytes,6,rep,name=ipaddresses"`

	// Information about the requesting user (if relevant)
	// See user.Info interface for details
	Username string   `json:"username,omitempty" protobuf:"bytes,7,opt,name=username"`
	UID      string   `json:"uid,omitempty" protobuf:"bytes,8,opt,name=uid"`
	Groups   []string `json:"groups,omitempty" protobuf:"bytes,9,rep,name=groups"`
}

type CertificateSigningRequestStatus struct {
	// Conditions applied to the request, such as approval or denial.
	Conditions []CertificateSigningRequestCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`

	// If request was approved, the controller will place the issued certificate here.
	Certificate []byte `json:"certificate,omitempty" protobuf:"bytes,2,opt,name=certificate"`
}

type RequestConditionType string

// These are the possible states for a certificate request.
const (
	CertificateApproved RequestConditionType = "Approved"
	CertificateDenied   RequestConditionType = "Denied"
)

type CertificateSigningRequestCondition struct {
	// request approval state, currently Approved or Denied.
	Type RequestConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=RequestConditionType"`
	// brief reason for the request state
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// human readable message with details about the request state
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// timestamp for the last update to this condition
	LastUpdateTime unversioned.Time `json:"lastupdatetime,omitempty" protobuf:"bytes,4,opt,name=lastupdatetime"`
}

type CertificateSigningRequestList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []CertificateSigningRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
