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

package etcd

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/certificates"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/cachesize"
	csrregistry "k8s.io/kubernetes/pkg/registry/certificates"
	"k8s.io/kubernetes/pkg/registry/generic"
	etcdgeneric "k8s.io/kubernetes/pkg/registry/generic/etcd"
	"k8s.io/kubernetes/pkg/runtime"
)

// REST implements a RESTStorage for CertificateSigningRequest against etcd
type REST struct {
	*etcdgeneric.Etcd
}

// NewREST returns a registry which will store CertificateSigningRequest in the given helper
func NewREST(opts generic.RESTOptions) (*REST, *StatusREST, *ApproveREST) {
	prefix := "/certificatesigningrequests"

	newListFunc := func() runtime.Object { return &certificates.CertificateSigningRequestList{} }
	storageInterface := opts.Decorator(opts.Storage, cachesize.GetWatchCacheSizeByResource(cachesize.CertificateSigningRequests), &certificates.CertificateSigningRequest{}, prefix, csrregistry.Strategy, newListFunc)

	store := &etcdgeneric.Etcd{
		NewFunc:     func() runtime.Object { return &certificates.CertificateSigningRequest{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx api.Context) string {
			return etcdgeneric.NamespaceKeyRootFunc(ctx, prefix)
		},
		KeyFunc: func(ctx api.Context, id string) (string, error) {
			return etcdgeneric.NamespaceKeyFunc(ctx, prefix, id)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*certificates.CertificateSigningRequest).Name, nil
		},
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return csrregistry.Matcher(label, field)
		},
		QualifiedResource:       certificates.Resource("certificatesigningrequests"),
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy: csrregistry.Strategy,
		UpdateStrategy: csrregistry.Strategy,

		Storage: storageInterface,
	}

	// Subresources use the same store and creation strategy, which only
	// allows empty subs. Updates to an existing subresource are handled by
	// dedicated strategies.
	statusStore := *store
	statusStore.UpdateStrategy = csrregistry.StatusStrategy

	approveStore := *store
	approveStore.UpdateStrategy = csrregistry.ApproveStrategy

	return &REST{store}, &StatusREST{store: &statusStore}, &ApproveREST{store: &approveStore}
}

// StatusREST implements the REST endpoint for changing the status of a CSR.
type StatusREST struct {
	store *etcdgeneric.Etcd
}

func (r *StatusREST) New() runtime.Object {
	return &certificates.CertificateSigningRequest{}
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx api.Context, obj runtime.Object) (runtime.Object, bool, error) {
	return r.store.Update(ctx, obj)
}

// ApproveREST implements the REST endpoint for changing the approval state of a CSR.
type ApproveREST struct {
	store *etcdgeneric.Etcd
}

func (r *ApproveREST) New() runtime.Object {
	return &certificates.CertificateSigningRequest{}
}

// Update alters the approve subset of an object.
func (r *ApproveREST) Update(ctx api.Context, obj runtime.Object) (runtime.Object, bool, error) {
	return r.store.Update(ctx, obj)
}
