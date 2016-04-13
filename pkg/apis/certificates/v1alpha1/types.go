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
	Request []byte `json:"request" protobuf:"bytes,1,opt,name=request"`

	// Any extra information the node wishes to send with the request.
	ExtraInfo []string `json:"extraInfo,omitempty" protobuf:"bytes,2,rep,name=extraInfo"`

	// Fingerprint of the public key in request
	Fingerprint string `json:"fingerprint,omitempty" protobuf:"bytes,3,opt,name=fingerprint"`

	// Subject fields from the request
	Subject Subject `json:"subject,omitempty" protobuf:"bytes,4,opt,name=subject"`

	// DNS SANs from the request
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,5,rep,name=hostnames"`

	// IP SANs from the request
	IPAddresses []string `json:"ipAddresses,omitempty" protobuf:"bytes,6,rep,name=ipAddresses"`

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

// These are the possible conditions for a certificate request.
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
	LastUpdateTime unversioned.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,4,opt,name=lastUpdateTime"`
}

type CertificateSigningRequestList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []CertificateSigningRequest `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

// Subject is a wrapper around pkix.Name which supports correct marshaling to
// JSON. In particular, it marshals into strings, which can be used as map keys
// in json.
type Subject struct {
	Country            []string `protobuf:"bytes,1,rep,name=country"`
	Organization       []string `protobuf:"bytes,2,rep,name=organization"`
	OrganizationalUnit []string `protobuf:"bytes,3,rep,name=organizationalUnit"`
	Locality           []string `protobuf:"bytes,4,rep,name=locality"`
	Province           []string `protobuf:"bytes,5,rep,name=province"`
	StreetAddress      []string `protobuf:"bytes,6,rep,name=streetAddress"`
	PostalCode         []string `protobuf:"bytes,7,rep,name=postalCode"`
	SerialNumber       string   `protobuf:"bytes,8,rep,name=serialNumber"`
	CommonName         string   `protobuf:"bytes,9,opt,name=commonName"`

	Names      []string `protobuf:"bytes,10,rep,name=names"`
	ExtraNames []string `protobuf:"bytes,11,rep,name=extraNames"`
}
