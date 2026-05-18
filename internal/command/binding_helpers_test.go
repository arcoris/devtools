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

// TestCloneBindingHelpers verifies detached helper behavior.
func TestCloneBindingHelpers(t *testing.T) {
	t.Parallel()

	stringsValue := []string{"a"}
	clonedStrings := cloneBindingStrings(stringsValue)
	clonedStrings[0] = "b"

	if got, want := stringsValue[0], "a"; got != want {
		t.Fatalf("cloneBindingStrings mutated source: got %q, want %q", got, want)
	}

	options := []Option{bindingTestStringOption("output")}
	clonedOptions := cloneBindingOptions(options)
	clonedOptions[0] = bindingTestStringOption("changed")

	if options[0].Name() != MustOptionName("output") {
		t.Fatalf("cloneBindingOptions mutated source")
	}

	arguments := []Argument{bindingTestStringArgument("package")}
	clonedArguments := cloneBindingArguments(arguments)
	clonedArguments[0] = bindingTestStringArgument("changed")

	if arguments[0].Name() != MustArgumentName("package") {
		t.Fatalf("cloneBindingArguments mutated source")
	}
}

func bindingTestFormatOption() Option {
	return MustOption(OptionSpec{
		Name:      "format",
		Aliases:   []string{"fmt"},
		Shorthand: "f",
		Kind:      OptionKindEnum,
		AllowedValues: []string{
			"text",
			"json",
		},
		DefaultValues: []string{"text"},
	})
}

func bindingTestPackageOption() Option {
	return MustOption(OptionSpec{
		Name: "package",
		Kind: OptionKindStringList,
		DefaultValues: []string{
			"./...",
		},
	})
}

func bindingTestStringOption(name string) Option {
	return MustOption(OptionSpec{
		Name: name,
		Kind: OptionKindString,
	})
}

func bindingTestSuiteArgument() Argument {
	return MustArgument(ArgumentSpec{
		Name: "suite",
		Kind: OptionKindEnum,
		AllowedValues: []string{
			"smoke",
			"stable",
		},
	})
}

func bindingTestPackageArgument() Argument {
	return MustArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Requirement: ArgumentRequirementOptional,
	})
}

func bindingTestStringArgument(name string) Argument {
	return MustArgument(ArgumentSpec{
		Name: name,
		Kind: OptionKindString,
	})
}
