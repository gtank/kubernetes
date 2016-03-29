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

// Provides utility functions for parsing and approving certificate requests.
package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

//TODO;
// construct a cfssl Signer from /var/lib/kubernetes/ca/
// signer.ParseCertificateRequest on DER bytes
//

// DecodeCertificateRequest decodes and validates a PEM-encoded Certificate
// Signing Request. It returns the decoded request bytes if both operations are
// successful.
func DecodeCertificateRequest(encoded []byte) ([]byte, error) {
	block, _ := pem.Decode(encoded)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("could not decode PEM block type %s", block.Type)
	}

	// check the signature
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}
	err = csr.CheckSignature()
	if err != nil {
		return nil, err
	}

	return block.Bytes, nil
}
