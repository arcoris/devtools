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

// TestBoundInputCopySemantics verifies detached bound output slices.
func TestBoundInputCopySemantics(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
		},
		Arguments: []Argument{
			bindingTestSuiteArgument(),
		},
	})

	bound := binding.MustBind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})

	options := bound.Options()
	options[0] = MustScalarOptionValue("other", OptionKindString, OptionSourceCommandLine, "value")

	if _, ok := bound.Option(MustOptionName("other")); ok {
		t.Fatalf("bound input changed through Options() slice")
	}

	arguments := bound.Arguments()
	arguments[0] = BoundArgument{}

	if suite, ok := bound.Argument(MustArgumentName("suite")); !ok || suite.MustValue() != "stable" {
		t.Fatalf("bound input changed through Arguments() slice")
	}
}

// TestBoundInputValidateRejectsDuplicates verifies bound input structural validation.
func TestBoundInputValidateRejectsDuplicates(t *testing.T) {
	t.Parallel()

	input := BoundInput{
		options: []OptionValue{
			MustScalarOptionValue("format", OptionKindString, OptionSourceCommandLine, "text"),
			MustScalarOptionValue("format", OptionKindString, OptionSourceCommandLine, "json"),
		},
	}

	err := input.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidBindingValue) {
		t.Fatalf("Validate() error = %v, want ErrInvalidBindingValue", err)
	}
}
