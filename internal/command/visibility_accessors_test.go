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

import "testing"

func TestVisibilityStringAndKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		visibility Visibility
		want       string
	}{
		{visibility: VisibilityPublic, want: "public"},
		{visibility: VisibilityHidden, want: "hidden"},
		{visibility: VisibilityInternal, want: "internal"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.want, func(t *testing.T) {
			t.Parallel()

			if got := test.visibility.String(); got != test.want {
				t.Fatalf("String() = %q, want %q", got, test.want)
			}

			if got := test.visibility.Key(); got != test.want {
				t.Fatalf("Key() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestVisibilityOrDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		visibility Visibility
		want       Visibility
	}{
		{name: "zero", visibility: "", want: VisibilityPublic},
		{name: "public", visibility: VisibilityPublic, want: VisibilityPublic},
		{name: "hidden", visibility: VisibilityHidden, want: VisibilityHidden},
		{name: "internal", visibility: VisibilityInternal, want: VisibilityInternal},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.visibility.OrDefault(); got != test.want {
				t.Fatalf("%q.OrDefault() = %q, want %q", test.visibility, got, test.want)
			}
		})
	}
}

func TestVisibilityPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		visibility Visibility
		public     bool
		hidden     bool
		internal   bool
		zero       bool
		known      bool
	}{
		{name: "zero", visibility: "", public: false, hidden: false, internal: false, zero: true, known: false},
		{name: "public", visibility: VisibilityPublic, public: true, hidden: false, internal: false, zero: false, known: true},
		{name: "hidden", visibility: VisibilityHidden, public: false, hidden: true, internal: false, zero: false, known: true},
		{name: "internal", visibility: VisibilityInternal, public: false, hidden: false, internal: true, zero: false, known: true},
		{name: "unknown", visibility: Visibility("private"), public: false, hidden: false, internal: false, zero: false, known: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.visibility.IsPublic(); got != test.public {
				t.Fatalf("IsPublic() = %v, want %v", got, test.public)
			}

			if got := test.visibility.IsHidden(); got != test.hidden {
				t.Fatalf("IsHidden() = %v, want %v", got, test.hidden)
			}

			if got := test.visibility.IsInternal(); got != test.internal {
				t.Fatalf("IsInternal() = %v, want %v", got, test.internal)
			}

			if got := test.visibility.IsZero(); got != test.zero {
				t.Fatalf("IsZero() = %v, want %v", got, test.zero)
			}

			if got := test.visibility.IsKnown(); got != test.known {
				t.Fatalf("IsKnown() = %v, want %v", got, test.known)
			}
		})
	}
}

func TestVisibilityEqual(t *testing.T) {
	t.Parallel()

	if !VisibilityPublic.Equal(VisibilityPublic) {
		t.Fatalf("Equal(same) = false, want true")
	}

	if VisibilityPublic.Equal(VisibilityHidden) {
		t.Fatalf("Equal(other) = true, want false")
	}
}
