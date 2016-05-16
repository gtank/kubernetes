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
	unversioned "k8s.io/kubernetes/pkg/api/unversioned"
	certificates "k8s.io/kubernetes/pkg/apis/certificates"
	core "k8s.io/kubernetes/pkg/client/testing/core"
	labels "k8s.io/kubernetes/pkg/labels"
	watch "k8s.io/kubernetes/pkg/watch"
)

// FakeCertificateSigningRequests implements CertificateSigningRequestInterface
type FakeCertificateSigningRequests struct {
	Fake *FakeCertificates
	ns   string
}

var certificatesigningrequestsResource = unversioned.GroupVersionResource{Group: "certificates", Version: "", Resource: "certificatesigningrequests"}

func (c *FakeCertificateSigningRequests) Create(certificateSigningRequest *certificates.CertificateSigningRequest) (result *certificates.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewCreateAction(certificatesigningrequestsResource, c.ns, certificateSigningRequest), &certificates.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) Update(certificateSigningRequest *certificates.CertificateSigningRequest) (result *certificates.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateAction(certificatesigningrequestsResource, c.ns, certificateSigningRequest), &certificates.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) UpdateStatus(certificateSigningRequest *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	obj, err := c.Fake.
		Invokes(core.NewUpdateSubresourceAction(certificatesigningrequestsResource, "status", c.ns, certificateSigningRequest), &certificates.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) Delete(name string, options *api.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(core.NewDeleteAction(certificatesigningrequestsResource, c.ns, name), &certificates.CertificateSigningRequest{})

	return err
}

func (c *FakeCertificateSigningRequests) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	action := core.NewDeleteCollectionAction(certificatesigningrequestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &certificates.CertificateSigningRequestList{})
	return err
}

func (c *FakeCertificateSigningRequests) Get(name string) (result *certificates.CertificateSigningRequest, err error) {
	obj, err := c.Fake.
		Invokes(core.NewGetAction(certificatesigningrequestsResource, c.ns, name), &certificates.CertificateSigningRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*certificates.CertificateSigningRequest), err
}

func (c *FakeCertificateSigningRequests) List(opts api.ListOptions) (result *certificates.CertificateSigningRequestList, err error) {
	obj, err := c.Fake.
		Invokes(core.NewListAction(certificatesigningrequestsResource, c.ns, opts), &certificates.CertificateSigningRequestList{})

	if obj == nil {
		return nil, err
	}

	label := opts.LabelSelector
	if label == nil {
		label = labels.Everything()
	}
	list := &certificates.CertificateSigningRequestList{}
	for _, item := range obj.(*certificates.CertificateSigningRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificateSigningRequests.
func (c *FakeCertificateSigningRequests) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(core.NewWatchAction(certificatesigningrequestsResource, c.ns, opts))

}
