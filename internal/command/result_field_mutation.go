// Copyright 2026 The ARCORIS Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

// WithField returns a validated copy with one field added or replaced.
func (result Result) WithField(key string, value string) (Result, error) {
	spec := result.spec()
	if spec.Fields == nil {
		spec.Fields = make(map[string]string)
	}

	spec.Fields[key] = value

	return NewResult(spec)
}

// MustWithField returns a validated copy with one field added or replaced and
// panics on invalid input.
func (result Result) MustWithField(key string, value string) Result {
	next, err := result.WithField(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutField returns a validated copy without one field.
func (result Result) WithoutField(key string) Result {
	spec := result.spec()
	delete(spec.Fields, key)

	return MustResult(spec)
}
