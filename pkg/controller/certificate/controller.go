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

package certificate

import (
	"time"

	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_1"
	"k8s.io/kubernetes/pkg/client/record"
	unversioned_legacy "k8s.io/kubernetes/pkg/client/typed/generated/legacy/unversioned"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util"
	utilruntime "k8s.io/kubernetes/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	// Periodically check for new certificate requests
	CertificateResyncPeriod = 15 * time.Second
)

type CertificateController struct {
	kubeClient clientset.Interface

	// CSR framework and store
	csrController *framework.Controller
	csrStore      cache.StoreToCSRLister

	// To allow injection of updateCertificateRequestStatus for testing.
	updateHandler func(csr *extensions.CertificateSigningRequest) error
	syncHandler   func(csrKey string) error

	signer *local.Signer

	queue *workqueue.Type
}

func NewCertificateController(kubeClient clientset.Interface, syncPeriod time.Duration, caCertFile, caKeyFile string) (*CertificateController, error) {
	// Send events to the apiserver
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&unversioned_legacy.EventSinkImpl{kubeClient.Legacy().Events("")})

	// Configure cfssl signer
	// TODO: support non-default policy and remote/pkcs11 signing
	policy := &config.Signing{
		Default: config.DefaultConfig(),
	}
	ca, err := local.NewSignerFromFile(caCertFile, caKeyFile, policy)
	if err != nil {
		glog.Errorf("Unable to initialize signer: %v", err)
		return nil, err
	}

	cc := &CertificateController{
		kubeClient: kubeClient,
		queue:      workqueue.New(),
		signer:     ca,
	}

	// Manage the addition/update of certificate requests
	cc.csrStore.Store, cc.csrController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return cc.kubeClient.Extensions().CertificateSigningRequests(api.NamespaceAll).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return cc.kubeClient.Extensions().CertificateSigningRequests(api.NamespaceAll).Watch(options)
			},
		},
		&extensions.CertificateSigningRequest{},
		syncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				csr := obj.(*extensions.CertificateSigningRequest)
				glog.V(4).Infof("Adding certificate request %s", csr.Name)
				cc.enqueueCertificateRequest(obj)
			},
			UpdateFunc: func(old, new interface{}) {
				oldCSR := old.(*extensions.CertificateSigningRequest)
				glog.V(4).Infof("Updating certificate request %s", oldCSR.Name)
				cc.enqueueCertificateRequest(new)
			},
			DeleteFunc: func(obj interface{}) {
				csr := obj.(*extensions.CertificateSigningRequest)
				glog.V(4).Infof("Deleting certificate request %s", csr.Name)
				cc.enqueueCertificateRequest(obj)
			},
		},
	)
	cc.syncHandler = cc.maybeSignCertificate
	return cc, nil
}

// Run the main goroutine responsible for watching and syncing jobs.
func (cc *CertificateController) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	go cc.csrController.Run(stopCh)
	glog.Infof("Starting certificate controller manager")
	for i := 0; i < workers; i++ {
		go util.Until(cc.worker, 1*time.Second, stopCh)
	}
	<-stopCh
	glog.Infof("Shutting down certificate controller")
	cc.queue.ShutDown()
}

// worker runs a thread that dequeues CSRs, handles them, and marks them done.
func (cc *CertificateController) worker() {
	for {
		func() {
			key, quit := cc.queue.Get()
			if quit {
				return
			}
			defer cc.queue.Done(key)
			err := cc.syncHandler(key.(string))
			if err != nil {
				glog.Errorf("Error syncing CSR: %v", err)
			}
		}()
	}
}

func (cc *CertificateController) enqueueCertificateRequest(obj interface{}) {
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	cc.queue.Add(key)
}

// maybeSignCertificate will inspect the certificate request and, if it has
// been approved and meets policy expectations, generate an X509 cert using the
// cluster CA assets. If successful it will update the CSR approve subresource
// with the signed certificate.
func (cc *CertificateController) maybeSignCertificate(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing certificate request %q (%v)", key, time.Now().Sub(startTime))
	}()
	obj, exists, err := cc.csrStore.Store.GetByKey(key)
	if err != nil {
		glog.Errorf("Unable to retrieve csr %v from store: %v", key, err)
		cc.queue.Add(key)
		return err
	}
	if !exists {
		glog.V(3).Infof("csr has been deleted: %v", key)
		return nil
	}
	csr := obj.(*extensions.CertificateSigningRequest)

	// At this point, the controller needs to:
	// 1. Derive information from the CSR and update the Status subresource
	// 2. Check for the approve subresource. If CSR was approved, then
	// 3. Generate a signed certificate, add it to /approve
	// 4. Update the Status resource to indicate certificate is available

	req := signer.SignRequest{Request: csr.Spec.CertificateRequest}
	certBytes, err := cc.signer.Sign(req)
	if err != nil {
		glog.Errorf("Unable to sign csr %v: %v", key, err)
		return err
	}
	cli.PrintCert(nil, []byte(csr.Spec.CertificateRequest), certBytes)
	return nil
}
