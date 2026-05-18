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

// TestBindingBindHandlesVariadicArgument verifies variadic positional binding.
func TestBindingBindHandlesVariadicArgument(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Arguments: []Argument{
			MustArgument(ArgumentSpec{
				Name:        "package",
				Kind:        OptionKindString,
				Requirement: ArgumentRequirementOptional,
				Cardinality: ArgumentCardinalityVariadic,
			}),
		},
	})

	bound, err := binding.Bind(BindingValueSpec{
		PositionalValues: []string{"./...", "./internal/..."},
	})
	if err != nil {
		t.Fatalf("Bind() returned unexpected error: %v", err)
	}

	pkg, ok := bound.Argument(MustArgumentName("package"))
	if !ok {
		t.Fatalf("bound argument package not found")
	}

	if got, want := pkg.Len(), 2; got != want {
		t.Fatalf("package Len() = %d, want %d", got, want)
	}
}

// TestBindingBindRejectsMissingRequiredArgument verifies missing positional handling.
func TestBindingBindRejectsMissingRequiredArgument(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Arguments: []Argument{
			bindingTestSuiteArgument(),
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

// TestBindingBindRejectsTooManyPositionals verifies positional overflow.
func TestBindingBindRejectsTooManyPositionals(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Arguments: []Argument{
			bindingTestSuiteArgument(),
		},
	})

	_, err := binding.Bind(BindingValueSpec{
		PositionalValues: []string{"stable", "extra"},
	})
	if err == nil {
		t.Fatalf("Bind() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Bind() error = %v, want ErrInvalidBindingValue", err)
	}
}

// TestBindingBindUsesOptionalArgumentDefault verifies positional default binding.
func TestBindingBindUsesOptionalArgumentDefault(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Arguments: []Argument{
			MustArgument(ArgumentSpec{
				Name:          "suite",
				Kind:          OptionKindEnum,
				Requirement:   ArgumentRequirementOptional,
				AllowedValues: []string{"smoke", "stable"},
				DefaultValues: []string{"smoke"},
			}),
		},
	})

	bound, err := binding.Bind(BindingValueSpec{})
	if err != nil {
		t.Fatalf("Bind() returned unexpected error: %v", err)
	}

	suite, ok := bound.Argument(MustArgumentName("suite"))
	if !ok {
		t.Fatalf("bound argument suite not found")
	}

	if got, want := suite.MustValue(), "smoke"; got != want {
		t.Fatalf("suite value = %q, want %q", got, want)
	}
}
