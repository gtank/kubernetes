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
	api "k8s.io/kubernetes/pkg/api"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	watch "k8s.io/kubernetes/pkg/watch"
)

// CertificateSigningRequestsGetter has a method to return a CertificateSigningRequestInterface.
// A group's client should implement this interface.
type CertificateSigningRequestsGetter interface {
	CertificateSigningRequests(namespace string) CertificateSigningRequestInterface
}

// CertificateSigningRequestInterface has methods to work with CertificateSigningRequest resources.
type CertificateSigningRequestInterface interface {
	Create(*extensions.CertificateSigningRequest) (*extensions.CertificateSigningRequest, error)
	Update(*extensions.CertificateSigningRequest) (*extensions.CertificateSigningRequest, error)
	Delete(name string, options *api.DeleteOptions) error
	DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error
	Get(name string) (*extensions.CertificateSigningRequest, error)
	List(opts api.ListOptions) (*extensions.CertificateSigningRequestList, error)
	Watch(opts api.ListOptions) (watch.Interface, error)
	CertificateSigningRequestExpansion
}

// certificateSigningRequests implements CertificateSigningRequestInterface
type certificateSigningRequests struct {
	client *ExtensionsClient
	ns     string
}

// newCertificateSigningRequests returns a CertificateSigningRequests
func newCertificateSigningRequests(c *ExtensionsClient, namespace string) *certificateSigningRequests {
	return &certificateSigningRequests{
		client: c,
		ns:     namespace,
	}
}

// Create takes the representation of a certificateSigningRequest and creates it.  Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *certificateSigningRequests) Create(certificateSigningRequest *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		Body(certificateSigningRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a certificateSigningRequest and updates it. Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *certificateSigningRequests) Update(certificateSigningRequest *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		Name(certificateSigningRequest.Name).
		Body(certificateSigningRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the certificateSigningRequest and deletes it. Returns an error if one occurs.
func (c *certificateSigningRequests) Delete(name string, options *api.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *certificateSigningRequests) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		VersionedParams(&listOptions, api.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the certificateSigningRequest, and returns the corresponding certificateSigningRequest object, and an error if there is any.
func (c *certificateSigningRequests) Get(name string) (result *extensions.CertificateSigningRequest, err error) {
	result = &extensions.CertificateSigningRequest{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CertificateSigningRequests that match those selectors.
func (c *certificateSigningRequests) List(opts api.ListOptions) (result *extensions.CertificateSigningRequestList, err error) {
	result = &extensions.CertificateSigningRequestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested certificateSigningRequests.
func (c *certificateSigningRequests) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("certificatesigningrequests").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}
