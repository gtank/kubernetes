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

import "crypto/x509/pkix"

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

func NewInternalSubject(name pkix.Name) Subject {
	subject := Subject{}
	subject.Country = name.Country
	subject.Organization = name.Organization
	subject.OrganizationalUnit = name.OrganizationalUnit
	subject.Locality = name.Locality
	subject.Province = name.Province
	subject.StreetAddress = name.StreetAddress
	subject.PostalCode = name.PostalCode
	subject.SerialNumber = name.SerialNumber
	subject.CommonName = name.CommonName
	for _, name := range name.Names {
		subject.Names = append(subject.Names, name.Value.(string))
	}
	for _, extraName := range name.ExtraNames {
		subject.ExtraNames = append(subject.ExtraNames, extraName.Value.(string))
	}
	return subject
}
