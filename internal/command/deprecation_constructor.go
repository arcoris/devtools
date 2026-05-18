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

// NewDeprecation validates spec and returns Deprecation.
//
// A valid deprecation always has a non-blank Message. Since and Replacement are
// optional, with RootPath representing the absence of a replacement command.
func NewDeprecation(spec DeprecationSpec) (Deprecation, error) {
	deprecation := Deprecation{
		since:       spec.Since,
		message:     spec.Message,
		replacement: spec.Replacement,
	}

	if err := deprecation.Validate(); err != nil {
		return Deprecation{}, err
	}

	return deprecation, nil
}

// MustDeprecation validates spec and returns Deprecation.
//
// MustDeprecation panics on invalid input. It is intended for static command
// definitions and tests.
func MustDeprecation(spec DeprecationSpec) Deprecation {
	deprecation, err := NewDeprecation(spec)
	if err != nil {
		panic(err)
	}

	return deprecation
}
