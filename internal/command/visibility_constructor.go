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

// DefaultVisibility returns the default visibility for ordinary command nodes.
//
// The command kernel defaults to public exposure because most statically
// declared commands are intended to be visible unless they explicitly opt out.
func DefaultVisibility() Visibility {
	return VisibilityPublic
}

// KnownVisibilities returns all supported visibility states in stable order.
func KnownVisibilities() []Visibility {
	out := make([]Visibility, len(knownVisibilities))
	copy(out, knownVisibilities)

	return out
}

// NewVisibility validates raw and returns it as a Visibility.
//
// Use NewVisibility for values loaded from configuration, generated metadata,
// tests, or any other external source where invalid values should be reported
// as errors.
func NewVisibility(raw string) (Visibility, error) {
	visibility := Visibility(raw)
	if err := visibility.Validate(); err != nil {
		return "", err
	}

	return visibility, nil
}

// ParseVisibility is an alias for NewVisibility.
//
// The name is useful at call sites where the value is parsed from an external
// string representation.
func ParseVisibility(raw string) (Visibility, error) {
	return NewVisibility(raw)
}

// MustVisibility validates raw and returns it as a Visibility.
//
// MustVisibility panics on invalid input. It is intended for static command
// definitions and tests where invalid visibility values are programmer errors.
func MustVisibility(raw string) Visibility {
	visibility, err := NewVisibility(raw)
	if err != nil {
		panic(err)
	}

	return visibility
}

// VisibilityFromHidden converts the legacy Hidden boolean shape into the
// explicit visibility enum.
func VisibilityFromHidden(hidden bool) Visibility {
	if hidden {
		return VisibilityHidden
	}

	return DefaultVisibility()
}
