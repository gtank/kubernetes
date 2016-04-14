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

package testclient

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/certificates"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

// NewSimpleFakeCertificate returns a client that will respond with the provided objects
func NewSimpleFakeCertificates(objects ...runtime.Object) *FakeCertificates {
	return &FakeCertificates{Fake: NewSimpleFake(objects...)}
}

// FakeCertificates implements CertificatesInterface. Meant to be
// embedded into a struct to get a default implementation. This makes faking
// out just the method you want to test easier.
type FakeCertificates struct {
	*Fake
}

func (c *FakeCertificates) CertificateSigningRequests(namespace string) unversioned.CertificateSigningRequestInterface {
	return &FakeCertificateSigningRequest{Fake: c, Namespace: namespace}
}

// FakeCertificateSigningRequest implements CertificateSigningRequestInterface
type FakeCertificateSigningRequest struct {
	Fake      *FakeCertificates
	Namespace string
}

func (c *FakeCertificateSigningRequest) Get(name string) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.Invokes(NewGetAction("certificatesigningrequests", c.Namespace, name), &certificates.CertificateSigningRequest{})
	if obj == nil {
		return nil, err
	}

	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequest) List(opts api.ListOptions) (*certificates.CertificateSigningRequestList, error) {
	obj, err := c.Fake.Invokes(NewListAction("certificatesigningrequests", c.Namespace, opts), &certificates.CertificateSigningRequestList{})
	if obj == nil {
		return nil, err
	}

	return obj.(*certificates.CertificateSigningRequestList), err
}

func (c *FakeCertificateSigningRequest) Create(csr *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.Invokes(NewCreateAction("certificatesigningrequests", c.Namespace, csr), csr)
	if obj == nil {
		return nil, err
	}

	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequest) Update(csr *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.Invokes(NewUpdateAction("certificatesigningrequests", c.Namespace, csr), csr)
	if obj == nil {
		return nil, err
	}

	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequest) UpdateStatus(csr *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.Invokes(NewUpdateSubresourceAction("certificatesigningrequests", "status", c.Namespace, csr), csr)
	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequest) UpdateApprove(csr *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.Invokes(NewUpdateSubresourceAction("certificatesigningrequests", "approve", c.Namespace, csr), csr)
	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequest) Delete(name string, opts *api.DeleteOptions) error {
	_, err := c.Fake.Invokes(NewDeleteAction("certificatesigningrequests", c.Namespace, name), &certificates.CertificateSigningRequest{})
	return err
}

func (c *FakeCertificateSigningRequest) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(NewWatchAction("certificatesigningrequests", c.Namespace, opts))
}
