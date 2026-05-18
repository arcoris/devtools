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

// TestBindingBindRejectsDisallowedOptionSource verifies option policy source checks.
func TestBindingBindRejectsDisallowedOptionSource(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
					AllowedSources: []OptionSource{OptionSourceCommandLine},
				}),
			}),
		},
	})

	_, err := binding.Bind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceEnvironment, "out.txt"),
		},
	})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}

// TestBindingBindRejectsDuplicateOptionValue verifies duplicate resolved values.
func TestBindingBindRejectsDuplicateOptionValue(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
		},
	})

	_, err := binding.Bind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "text"),
			MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}

// TestBindingBindRejectsMissingRequiredOption verifies required option handling.
func TestBindingBindRejectsMissingRequiredOption(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
					Requirement: OptionRequirementRequired,
				}),
			}),
		},
	})

	_, err := binding.Bind(BindingValueSpec{})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}

// TestBindingBindRejectsUnknownOption verifies unknown option handling.
func TestBindingBindRejectsUnknownOption(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
		},
	})

	_, err := binding.Bind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("missing", OptionKindString, OptionSourceCommandLine, "value"),
		},
	})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}

// TestBindingBindRejectsWrongOptionKind verifies option value kind matching.
func TestBindingBindRejectsWrongOptionKind(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestStringOption("output"),
		},
	})

	_, err := binding.Bind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}
