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

func TestResolveOptionValuesCanonicalizesPrecedenceAndDefaults(t *testing.T) {
	t.Parallel()

	binding := optionResolverTestBinding()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		ConfigValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceConfig, "text"),
		},
		EnvironmentValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceEnvironment, "env.out"),
		},
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	if got, want := len(values), 3; got != want {
		t.Fatalf("resolved value count = %d, want %d", got, want)
	}

	assertResolvedOptionValue(t, values[0], "format", OptionSourceCommandLine, []string{"json"})
	assertResolvedOptionValue(t, values[1], "output", OptionSourceEnvironment, []string{"env.out"})
	assertResolvedOptionValue(t, values[2], "package", OptionSourceDefault, []string{"./..."})

	bound, err := binding.Bind(BindingValueSpec{
		OptionValues:     values,
		PositionalValues: []string{"smoke"},
	})
	if err != nil {
		t.Fatalf("Binding.Bind(resolved values) returned unexpected error: %v", err)
	}

	if got, ok := bound.OptionByName("format"); !ok || got.MustValue() != "json" {
		t.Fatalf("bound format = %v, %v; want json", got, ok)
	}
}

func TestResolveOptionValuesRejectsUnknownOptions(t *testing.T) {
	t.Parallel()

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("unknown", OptionKindString, OptionSourceCommandLine, "value"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesRejectsScalarDuplicateSourceValues(t *testing.T) {
	t.Parallel()

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "text"),
			MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesMergesListDuplicateSourceValues(t *testing.T) {
	t.Parallel()

	binding := optionResolverTestBinding()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "bench.out"),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./..."),
			MustListOptionValue("pkg", OptionKindStringList, OptionSourceCommandLine, "./internal/..."),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	assertResolvedOptionValue(t, values[0], "format", OptionSourceDefault, []string{"text"})
	assertResolvedOptionValue(t, values[1], "output", OptionSourceCommandLine, []string{"bench.out"})
	assertResolvedOptionValue(t, values[2], "package", OptionSourceCommandLine, []string{"./...", "./internal/..."})

	bound, err := binding.Bind(BindingValueSpec{
		OptionValues:     values,
		PositionalValues: []string{"smoke"},
	})
	if err != nil {
		t.Fatalf("Binding.Bind(resolved values) returned unexpected error: %v", err)
	}

	packages, ok := bound.OptionByName("package")
	if !ok {
		t.Fatalf("package option missing")
	}

	assertOptionValueStrings(t, packages.Values(), []string{"./...", "./internal/..."})
}

func TestResolveOptionValuesMergesMultipleListOccurrencesInSourceOrder(t *testing.T) {
	t.Parallel()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "bench.out"),
			MustListOptionValue("pkg", OptionKindStringList, OptionSourceCommandLine, "./one"),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./two", "./three"),
			MustListOptionValue("pkg", OptionKindStringList, OptionSourceCommandLine, "./four"),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	assertResolvedOptionValue(t, values[2], "package", OptionSourceCommandLine, []string{
		"./one",
		"./two",
		"./three",
		"./four",
	})
}

func TestResolveOptionValuesRejectsInvalidMergedListValue(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name:          "mode",
				Kind:          OptionKindEnumList,
				AllowedValues: []string{"fast", "stable"},
			}),
		},
	})

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		CommandLineValues: []OptionValue{
			MustListOptionValue("mode", OptionKindEnumList, OptionSourceCommandLine, "fast"),
			MustListOptionValue("mode", OptionKindEnumList, OptionSourceCommandLine, "slow"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionValue", err)
	}
}

func TestResolveOptionValuesRejectsScalarAliasDuplicateSourceValues(t *testing.T) {
	t.Parallel()

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "text"),
			MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesRejectsRepeatableScalarDuplicateSourceValues(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name: "tag",
				Kind: OptionKindString,
				Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
					Occurrence: OptionOccurrenceMultiple,
				}),
			}),
		},
	})

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("tag", OptionKindString, OptionSourceCommandLine, "one"),
			MustScalarOptionValue("tag", OptionKindString, OptionSourceCommandLine, "two"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesOverridesLowerPrecedenceListValues(t *testing.T) {
	t.Parallel()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		ConfigValues: []OptionValue{
			MustListOptionValue("package", OptionKindStringList, OptionSourceConfig, "./from-config"),
		},
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "bench.out"),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./..."),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./internal/..."),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	assertResolvedOptionValue(t, values[2], "package", OptionSourceCommandLine, []string{"./...", "./internal/..."})
}

func TestResolveOptionValuesOverridesLowerPrecedenceScalarValues(t *testing.T) {
	t.Parallel()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		ConfigValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceConfig, "text"),
		},
		EnvironmentValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceEnvironment, "json"),
		},
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "bench.out"),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	assertResolvedOptionValue(t, values[0], "format", OptionSourceEnvironment, []string{"json"})
}

func TestResolveOptionValuesSourcePrecedenceOrder(t *testing.T) {
	t.Parallel()

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		InheritedValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceInherited, "inherited.out"),
		},
		ConfigValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceConfig, "config.out"),
		},
		EnvironmentValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceEnvironment, "env.out"),
		},
		RuntimeValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceRuntime, "runtime.out"),
		},
		InteractiveValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceInteractive, "interactive.out"),
		},
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "cli.out"),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	assertResolvedOptionValue(t, values[1], "output", OptionSourceCommandLine, []string{"cli.out"})
}

func TestResolveOptionValuesRespectsAllowedSources(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name: "token",
				Kind: OptionKindString,
				Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
					AllowedSources: []OptionSource{OptionSourceRuntime},
				}),
			}),
		},
	})

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("token", OptionKindString, OptionSourceCommandLine, "secret"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesEnforcesAllowedSourcesBeforeMerge(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name: "package",
				Kind: OptionKindStringList,
				Policy: MustOptionPolicyForKind(OptionKindStringList, OptionPolicySpec{
					AllowedSources: []OptionSource{OptionSourceConfig},
				}),
			}),
		},
	})

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		CommandLineValues: []OptionValue{
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./..."),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./internal/..."),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func TestResolveOptionValuesRejectsSourceMismatchBeforeMerge(t *testing.T) {
	t.Parallel()

	_, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: optionResolverTestBinding(),
		CommandLineValues: []OptionValue{
			MustListOptionValue("package", OptionKindStringList, OptionSourceConfig, "./from-config"),
			MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./from-cli"),
		},
	})
	if err == nil {
		t.Fatalf("ResolveOptionValues() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionResolution) {
		t.Fatalf("ResolveOptionValues() error = %v, want ErrInvalidOptionResolution", err)
	}
}

func optionResolverTestBinding() Binding {
	return MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name:          "format",
				Aliases:       []string{"fmt"},
				Kind:          OptionKindEnum,
				AllowedValues: []string{"text", "json"},
				DefaultValues: []string{"text"},
			}),
			MustOption(OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
					Requirement: OptionRequirementRequired,
				}),
			}),
			MustOption(OptionSpec{
				Name:          "package",
				Aliases:       []string{"pkg"},
				Kind:          OptionKindStringList,
				DefaultValues: []string{"./..."},
			}),
		},
		Arguments: []Argument{
			MustArgument(ArgumentSpec{
				Name:          "suite",
				Kind:          OptionKindEnum,
				AllowedValues: []string{"smoke", "stable"},
			}),
		},
	})
}

func assertResolvedOptionValue(
	t *testing.T,
	value OptionValue,
	name string,
	source OptionSource,
	values []string,
) {
	t.Helper()

	if got, want := value.Name(), MustOptionName(name); got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if got := value.Source(); got != source {
		t.Fatalf("Source() = %q, want %q", got, source)
	}

	assertOptionValueStrings(t, value.Values(), values)
}
