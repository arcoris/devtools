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

// TestEmptyBinding verifies empty binding behavior.
func TestEmptyBinding(t *testing.T) {
	t.Parallel()

	binding := EmptyBinding()

	if !binding.IsZero() {
		t.Fatalf("EmptyBinding().IsZero() = false, want true")
	}

	if binding.HasOptions() {
		t.Fatalf("EmptyBinding().HasOptions() = true, want false")
	}

	if binding.HasArguments() {
		t.Fatalf("EmptyBinding().HasArguments() = true, want false")
	}

	if err := binding.Validate(); err != nil {
		t.Fatalf("EmptyBinding().Validate() returned unexpected error: %v", err)
	}
}

// TestMustBindingPanicsForInvalidBinding verifies fail-fast construction.
func TestMustBindingPanicsForInvalidBinding(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustBinding did not panic")
		}
	}()

	_ = MustBinding(BindingSpec{
		Options: []Option{
			bindingTestStringOption("output"),
			bindingTestStringOption("output"),
		},
	})
}

// TestNewBindingAcceptsValidBinding verifies full binding declaration construction.
func TestNewBindingAcceptsValidBinding(t *testing.T) {
	t.Parallel()

	binding, err := NewBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
			bindingTestPackageOption(),
		},
		Arguments: []Argument{
			bindingTestSuiteArgument(),
			bindingTestPackageArgument(),
		},
	})
	if err != nil {
		t.Fatalf("NewBinding() returned unexpected error: %v", err)
	}

	if got, want := binding.OptionCount(), 2; got != want {
		t.Fatalf("OptionCount() = %d, want %d", got, want)
	}

	if got, want := binding.ArgumentCount(), 2; got != want {
		t.Fatalf("ArgumentCount() = %d, want %d", got, want)
	}

	if !binding.HasOption(MustOptionName("format")) {
		t.Fatalf("HasOption(format) = false, want true")
	}

	if !binding.HasOption(MustOptionName("fmt")) {
		t.Fatalf("HasOption(alias fmt) = false, want true")
	}

	if !binding.HasOptionShorthand("f") {
		t.Fatalf("HasOptionShorthand(f) = false, want true")
	}

	if !binding.HasArgument(MustArgumentName("suite")) {
		t.Fatalf("HasArgument(suite) = false, want true")
	}
}

// TestNewBindingRejectsInvalidBinding verifies binding-level declaration validation.
func TestNewBindingRejectsInvalidBinding(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec BindingSpec
	}{
		{
			name: "invalid option",
			spec: BindingSpec{
				Options: []Option{
					{
						name: OptionName("Bad"),
					},
				},
			},
		},
		{
			name: "duplicate option names",
			spec: BindingSpec{
				Options: []Option{
					bindingTestStringOption("output"),
					bindingTestStringOption("output"),
				},
			},
		},
		{
			name: "duplicate option alias",
			spec: BindingSpec{
				Options: []Option{
					MustOption(OptionSpec{
						Name:    "output",
						Aliases: []string{"out"},
						Kind:    OptionKindString,
					}),
					MustOption(OptionSpec{
						Name:    "report",
						Aliases: []string{"out"},
						Kind:    OptionKindString,
					}),
				},
			},
		},
		{
			name: "duplicate shorthand",
			spec: BindingSpec{
				Options: []Option{
					MustOption(OptionSpec{
						Name:      "output",
						Shorthand: "o",
						Kind:      OptionKindString,
					}),
					MustOption(OptionSpec{
						Name:      "overwrite",
						Shorthand: "o",
						Kind:      OptionKindBool,
					}),
				},
			},
		},
		{
			name: "invalid argument",
			spec: BindingSpec{
				Arguments: []Argument{
					{
						name: ArgumentName("Bad"),
					},
				},
			},
		},
		{
			name: "duplicate argument names",
			spec: BindingSpec{
				Arguments: []Argument{
					bindingTestStringArgument("package"),
					bindingTestStringArgument("package"),
				},
			},
		},
		{
			name: "required after optional",
			spec: BindingSpec{
				Arguments: []Argument{
					MustArgument(ArgumentSpec{
						Name:        "package",
						Kind:        OptionKindString,
						Requirement: ArgumentRequirementOptional,
					}),
					MustArgument(ArgumentSpec{
						Name: "suite",
						Kind: OptionKindString,
					}),
				},
			},
		},
		{
			name: "variadic is not last",
			spec: BindingSpec{
				Arguments: []Argument{
					MustArgument(ArgumentSpec{
						Name:        "package",
						Kind:        OptionKindString,
						Cardinality: ArgumentCardinalityVariadic,
					}),
					MustArgument(ArgumentSpec{
						Name: "suite",
						Kind: OptionKindString,
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewBinding(tt.spec)
			if err == nil {
				t.Fatalf("NewBinding() returned nil error")
			}

			if !errors.Is(err, ErrInvalidBinding) {
				t.Fatalf("NewBinding() error = %v, want ErrInvalidBinding", err)
			}
		})
	}
}
