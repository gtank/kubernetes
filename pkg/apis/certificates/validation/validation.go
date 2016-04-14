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

package validation

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	apivalidation "k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/apis/certificates"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// validates the signature of a PEM-encoded PKCS#10 certificate signing
// request. If this is invalid, we must not accept the CSR for further
// processing.
func validateSignature(obj *certificates.CertificateSigningRequest) error {
	// extract PEM from request object
	pemCert := []byte(obj.Spec.CertificateRequest)

	// decode PEM
	block, _ := pem.Decode(pemCert)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return errors.New("PEM block type must be CERTIFICATE REQUEST")
	}

	// parse request
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return err
	}

	// check that the signature is valid
	return csr.CheckSignature()
}

// We don't care what you call your certificate requests.
func ValidateCertificateRequestName(name string, prefix bool) (bool, string) {
	return true, ""
}

func ValidateCertificateSigningRequest(csr *certificates.CertificateSigningRequest) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&csr.ObjectMeta, true, ValidateCertificateRequestName, field.NewPath("metadata"))
	err := validateSignature(csr)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("request"), err, "request signature is invalid"))
	}
	return allErrs
}

func ValidateCertificateSigningRequestUpdate(newCSR, oldCSR *certificates.CertificateSigningRequest) field.ErrorList {
	return ValidateCertificateSigningRequest(newCSR)
}
