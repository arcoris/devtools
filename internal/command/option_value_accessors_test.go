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

// TestOptionValueAccessorsAndSpec verifies detached snapshots and zero state.
func TestOptionValueAccessorsAndSpec(t *testing.T) {
	t.Parallel()

	value := MustListOptionValue("package", OptionKindStringList, OptionSourceEnvironment, "./...", "./internal/...")

	if value.IsZero() {
		t.Fatalf("IsZero() = true, want false")
	}

	if !value.IsValid() {
		t.Fatalf("IsValid() = false, want true")
	}

	if value.IsEmpty() {
		t.Fatalf("IsEmpty() = true, want false")
	}

	if got, ok := value.Value(); !ok || got != "./..." {
		t.Fatalf("Value() = %q, %v; want ./..., true", got, ok)
	}

	if got, ok := value.ValueAt(1); !ok || got != "./internal/..." {
		t.Fatalf("ValueAt(1) = %q, %v; want ./internal/..., true", got, ok)
	}

	if got, ok := value.ValueAt(-1); ok || got != "" {
		t.Fatalf("ValueAt(-1) = %q, %v; want empty, false", got, ok)
	}

	if got, ok := value.ValueAt(2); ok || got != "" {
		t.Fatalf("ValueAt(2) = %q, %v; want empty, false", got, ok)
	}

	spec := value.Spec()
	if got, want := spec.Name, "package"; got != want {
		t.Fatalf("Spec().Name = %q, want %q", got, want)
	}

	if got, want := spec.Kind, OptionKindStringList; got != want {
		t.Fatalf("Spec().Kind = %q, want %q", got, want)
	}

	if got, want := spec.Source, OptionSourceEnvironment; got != want {
		t.Fatalf("Spec().Source = %q, want %q", got, want)
	}

	spec.Values[0] = "changed"
	if got, want := value.Values()[0], "./..."; got != want {
		t.Fatalf("Spec() leaked values slice: got %q, want %q", got, want)
	}

	var zero OptionValue
	if !zero.IsZero() {
		t.Fatalf("zero IsZero() = false, want true")
	}

	if zero.IsValid() {
		t.Fatalf("zero IsValid() = true, want false")
	}

	if !zero.IsEmpty() {
		t.Fatalf("zero IsEmpty() = false, want true")
	}

	if got, ok := zero.Value(); ok || got != "" {
		t.Fatalf("zero Value() = %q, %v; want empty, false", got, ok)
	}

	assertOptionValuePanics(t, func() {
		_ = zero.MustValue()
	})
}

// TestOptionValueCopySemantics verifies detached value slices.
func TestOptionValueCopySemantics(t *testing.T) {
	t.Parallel()

	values := []string{"./..."}

	optionValue := MustOptionValue(OptionValueSpec{
		Name:   "package",
		Kind:   OptionKindStringList,
		Source: OptionSourceCommandLine,
		Values: values,
	})

	values[0] = "changed"

	if got, want := optionValue.Values()[0], "./..."; got != want {
		t.Fatalf("value changed through input slice: got %q, want %q", got, want)
	}

	out := optionValue.Values()
	out[0] = "changed"

	if got, want := optionValue.Values()[0], "./..."; got != want {
		t.Fatalf("value changed through output slice: got %q, want %q", got, want)
	}
}
