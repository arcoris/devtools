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

// NewGroup validates raw and returns it as a Group.
//
// Use NewGroup for group values loaded from configuration, generated metadata,
// tests, or other external sources where validation errors should be returned.
func NewGroup(raw string) (Group, error) {
	group := Group(raw)
	if err := group.Validate(); err != nil {
		return "", err
	}

	return group, nil
}

// ParseGroup is an alias for NewGroup.
//
// The name is useful at call sites where the value is parsed from an external
// string representation.
func ParseGroup(raw string) (Group, error) {
	return NewGroup(raw)
}

// NewGroupParts validates segments and returns their canonical group key.
func NewGroupParts(segments ...string) (Group, error) {
	raw, err := joinGroupSegments(segments)
	if err != nil {
		return "", err
	}

	return NewGroup(raw)
}

// MustGroup validates raw and returns it as a Group.
//
// MustGroup panics on invalid input. It is intended for static command
// definitions and tests where invalid group keys are programmer errors.
func MustGroup(raw string) Group {
	group, err := NewGroup(raw)
	if err != nil {
		panic(err)
	}

	return group
}

// MustGroupParts validates segments and returns their canonical group key.
func MustGroupParts(segments ...string) Group {
	group, err := NewGroupParts(segments...)
	if err != nil {
		panic(err)
	}

	return group
}
