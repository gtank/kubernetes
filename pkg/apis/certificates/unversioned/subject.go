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
	"crypto/x509/pkix"
	"encoding/json"
)

// Subject is a wrapper around pkix.Name which supports correct marshaling to
// JSON. In particular, it marshals into strings, which can be used as map keys
// in json.
type Subject struct {
	pkix.Name
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (s *Subject) UnmarshalJSON(b []byte) error {
	var name pkix.Name
	err := json.Unmarshal(b, &name)
	if err != nil {
		return err
	}
	s.Name = name
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s *Subject) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Name)
}
