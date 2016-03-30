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

package unversioned

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/watch"
)

type CertificateSigningRequestNamespacer interface {
	CertificateSigningRequests(namespace string) CertificateSigningRequestInterface
}

type CertificateSigningRequestInterface interface {
	List(opts api.ListOptions) (*extensions.CertificateSigningRequestList, error)
	Get(name string) (*extensions.CertificateSigningRequest, error)
	Create(ctrl *extensions.CertificateSigningRequest) (*extensions.CertificateSigningRequest, error)
	Update(ctrl *extensions.CertificateSigningRequest) (*extensions.CertificateSigningRequest, error)
	UpdateStatus(ctrl *extensions.CertificateSigningRequest) (*extensions.CertificateSigningRequest, error)
	Delete(name string) error
	Watch(opts api.ListOptions) (watch.Interface, error)
}

type certificateSigningRequests struct {
	client *ExtensionsClient
	ns     string
}

func newCertificateSigningRequests(c *ExtensionsClient, namespace string) *certificateSigningRequests {
	return &certificateSigningRequests{c, namespace}
}

// Ensure statically that certificateSigningRequests implements CertificateSigningRequestsInterface.
var _ CertificateSigningRequestInterface = &certificateSigningRequests{}

func (c *certificateSigningRequests) List(opts api.ListOptions) (result *extensions.CertificateSigningRequestList, err error) {
	result = &extensions.CertificateSigningRequestList{}
	err = c.client.Get().Namespace(c.ns).Resource("certificatesigningrequests").VersionedParams(&opts, api.ParameterCodec).Do().Into(result)
	return
}

// Get returns information about a particular certificate signing request.
func (c *certificateSigningRequests) Get(name string) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Get().Namespace(c.ns).Resource("certificatesigningrequests").Name(name).Do().Into(result)
	return
}

// Create creates a new certificate signing request.
func (c *certificateSigningRequests) Create(resource *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Post().Namespace(c.ns).Resource("certificatesigningrequests").Body(resource).Do().Into(result)
	return
}

// Update updates an existing certificate signing request.
func (c *certificateSigningRequests) Update(resource *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Put().Namespace(c.ns).Resource("certificatesigningrequests").Name(resource.Name).Body(resource).Do().Into(result)
	return
}

// UpdateStatus updates an existing certificate signing request status
func (c *certificateSigningRequests) UpdateStatus(resource *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Put().Namespace(c.ns).Resource("certificatesigningrequests").Name(resource.Name).SubResource("status").Body(resource).Do().Into(result)
	return
}

// Delete deletes an existing certificate signing request.
func (c *certificateSigningRequests) Delete(name string) error {
	return c.client.Delete().Namespace(c.ns).Resource("certificatesigningrequests").Name(name).Do().Error()
}

// Watch returns a watch.Interface that watches the requested certificate signing requests.
func (c *certificateSigningRequests) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
