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

// NewOptionKind validates raw and returns it as an OptionKind.
func NewOptionKind(raw string) (OptionKind, error) {
	kind := OptionKind(raw)
	if err := kind.Validate(); err != nil {
		return "", err
	}

	return kind, nil
}

// ParseOptionKind is an alias for NewOptionKind.
func ParseOptionKind(raw string) (OptionKind, error) {
	return NewOptionKind(raw)
}

// MustOptionKind validates raw and returns it as an OptionKind.
//
// MustOptionKind panics on invalid input. It is intended for static option
// declarations and tests where invalid option kinds are programmer errors.
func MustOptionKind(raw string) OptionKind {
	kind, err := NewOptionKind(raw)
	if err != nil {
		panic(err)
	}

	return kind
}
