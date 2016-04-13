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
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	apivalidation "k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/apis/certificates"
	certutil "k8s.io/kubernetes/pkg/util/certificates"
	"k8s.io/kubernetes/pkg/util/validation/field"
)

// validateCSR validates the signature and formatting of a base64-wrapped,
// PEM-encoded PKCS#10 certificate signing request. If this is invalid, we must
// not accept the CSR for further processing.
func validateCSR(obj *certificates.CertificateSigningRequest) error {
	csr, err := certutil.ParseCertificateRequestObject(obj)
	if err != nil {
		return err
	}
	// check that the signature is valid
	err = csr.CheckSignature()
	if err != nil {
		return err
	}
	// make sure we can calculate a fingerprint
	_, err = calculateFingerprint(csr)
	if err != nil {
		return err
	}
	return nil
}

// calculateFingerprint returns the SHA256 hash of the PKIX encoded public key
// from a certificate request. This instance assumes that the request's
// signature has already been checked and should not be used in other contexts.
func calculateFingerprint(req *x509.CertificateRequest) (string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(req.PublicKey)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(derBytes)
	return hex.EncodeToString(digest[:]), nil
}

// We don't care what you call your certificate requests.
func ValidateCertificateRequestName(name string, prefix bool) []string {
	return nil
}

func ValidateCertificateSigningRequest(csr *certificates.CertificateSigningRequest) field.ErrorList {
	isNamespaced := false
	allErrs := apivalidation.ValidateObjectMeta(&csr.ObjectMeta, isNamespaced, ValidateCertificateRequestName, field.NewPath("metadata"))
	err := validateCSR(csr)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("request"), csr.Spec.Request, fmt.Sprintf("%v", err)))
	}
	return allErrs
}

func ValidateCertificateSigningRequestUpdate(newCSR, oldCSR *certificates.CertificateSigningRequest) field.ErrorList {
	validationErrorList := ValidateCertificateSigningRequest(newCSR)
	metaUpdateErrorList := apivalidation.ValidateObjectMetaUpdate(&newCSR.ObjectMeta, &oldCSR.ObjectMeta, field.NewPath("metadata"))
	return append(validationErrorList, metaUpdateErrorList...)
}
