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

package fake

import (
	api "k8s.io/kubernetes/pkg/api"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	core "k8s.io/kubernetes/pkg/client/testing/core"
	labels "k8s.io/kubernetes/pkg/labels"
	watch "k8s.io/kubernetes/pkg/watch"
)

// FakeCertificateSigningRequests implements CertificateSigningRequestInterface
type FakeCertificateSigningRequests struct {
	Fake *FakeExtensions
	ns   string
}

func (c *FakeCertificateSigningRequests) Create(certificateSigningRequest *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewCreateAction("certificatesigningrequests", c.ns, certificateSigningRequest), &extensions.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*extensions.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) Update(certificateSigningRequest *extensions.CertificateSigningRequest) (result *extensions.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateAction("certificatesigningrequests", c.ns, certificateSigningRequest), &extensions.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*extensions.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) Delete(name string, options *api.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(core.NewDeleteAction("certificatesigningrequests", c.ns, name), &extensions.CertificateSigningRequest{})

	return err
}

func (c *FakeCertificateSigningRequests) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	action := core.NewDeleteCollectionAction("events", c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &extensions.CertificateSigningRequestList{})
	return err
}

func (c *FakeCertificateSigningRequests) Get(name string) (result *extensions.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewGetAction("certificatesigningrequests", c.ns, name), &extensions.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*extensions.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) List(opts api.ListOptions) (result *extensions.CertificateSigningRequestList, err error) {
	obj, err := c.Fake.
		Invokes(core.NewListAction("certificatesigningrequests", c.ns, opts), &extensions.CertificateSigningRequestList{})

	if obj == nil {
		return nil, err
	}

	label := opts.LabelSelector
	if label == nil {
		label = labels.Everything()
	}
	list := &extensions.CertificateSigningRequestList{}
	for _, item := range obj.(*extensions.CertificateSigningRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificateSigningRequests.
func (c *FakeCertificateSigningRequests) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(core.NewWatchAction("certificatesigningrequests", c.ns, opts))

}
