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

import (
	"errors"
	"testing"
)

// TestOptionValueWithHelpers verifies immutable-style updates.
func TestOptionValueWithHelpers(t *testing.T) {
	t.Parallel()

	value := MustScalarOptionValue("timeout", OptionKindDuration, OptionSourceCommandLine, "10s")

	next := value.
		MustWithName("delay").
		MustWithSource(OptionSourceEnvironment).
		MustWithValues("20s")

	if got, want := value.Source(), OptionSourceCommandLine; got != want {
		t.Fatalf("WithSource mutated original source: got %q, want %q", got, want)
	}

	if got, want := next.Name(), MustOptionName("delay"); got != want {
		t.Fatalf("next Name() = %q, want %q", got, want)
	}

	if got, want := next.Source(), OptionSourceEnvironment; got != want {
		t.Fatalf("next Source() = %q, want %q", got, want)
	}

	if got, want := next.MustValue(), "20s"; got != want {
		t.Fatalf("next MustValue() = %q, want %q", got, want)
	}

	if _, err := value.WithName("Delay"); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("WithName(invalid) error = %v, want ErrInvalidOptionValue", err)
	}

	if _, err := value.WithKind(OptionKindInt); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("WithKind(invalid) error = %v, want ErrInvalidOptionValue", err)
	}

	if _, err := value.WithValues("soon"); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("WithValues(invalid) error = %v, want ErrInvalidOptionValue", err)
	}

	assertOptionValuePanics(t, func() {
		_ = value.MustWithValues("soon")
	})
}
