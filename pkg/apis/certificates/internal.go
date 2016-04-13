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

import "crypto/x509/pkix"

// Subject is a wrapper around pkix.Name which supports correct marshaling to
// JSON. In particular, it marshals into strings, which can be used as map keys
// in json.
type Subject struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	Locality           []string
	Province           []string
	StreetAddress      []string
	PostalCode         []string
	SerialNumber       string
	CommonName         string

	Names      []string
	ExtraNames []string
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
