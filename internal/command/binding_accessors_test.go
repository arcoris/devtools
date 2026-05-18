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

// TestBindingCopySemantics verifies detached option and argument slices.
func TestBindingCopySemantics(t *testing.T) {
	t.Parallel()

	options := []Option{bindingTestStringOption("output")}
	arguments := []Argument{bindingTestStringArgument("package")}

	binding := MustBinding(BindingSpec{
		Options:   options,
		Arguments: arguments,
	})

	options[0] = bindingTestStringOption("changed")
	arguments[0] = bindingTestStringArgument("changed")

	if binding.HasOption(MustOptionName("changed")) {
		t.Fatalf("binding changed through input option slice")
	}

	if binding.HasArgument(MustArgumentName("changed")) {
		t.Fatalf("binding changed through input argument slice")
	}

	outOptions := binding.Options()
	outOptions[0] = bindingTestStringOption("changed")

	if binding.HasOption(MustOptionName("changed")) {
		t.Fatalf("binding changed through output option slice")
	}

	outArguments := binding.Arguments()
	outArguments[0] = bindingTestStringArgument("changed")

	if binding.HasArgument(MustArgumentName("changed")) {
		t.Fatalf("binding changed through output argument slice")
	}
}

// TestBindingLookup verifies option and argument lookup helpers.
func TestBindingLookup(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
		},
		Arguments: []Argument{
			bindingTestSuiteArgument(),
		},
	})

	if _, ok := binding.OptionByName("format"); !ok {
		t.Fatalf("OptionByName(format) ok = false, want true")
	}

	if _, ok := binding.OptionByName("fmt"); !ok {
		t.Fatalf("OptionByName(fmt alias) ok = false, want true")
	}

	if _, ok := binding.OptionByShorthand("f"); !ok {
		t.Fatalf("OptionByShorthand(f) ok = false, want true")
	}

	if _, ok := binding.ArgumentByName("suite"); !ok {
		t.Fatalf("ArgumentByName(suite) ok = false, want true")
	}

	if _, ok := binding.OptionByName("missing"); ok {
		t.Fatalf("OptionByName(missing) ok = true, want false")
	}

	if _, ok := binding.ArgumentByName("missing"); ok {
		t.Fatalf("ArgumentByName(missing) ok = true, want false")
	}
}

// TestBindingSortedNames verifies deterministic name helpers.
func TestBindingSortedNames(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestStringOption("zeta"),
			bindingTestStringOption("alpha"),
			bindingTestStringOption("middle"),
		},
		Arguments: []Argument{
			bindingTestStringArgument("zeta"),
			bindingTestStringArgument("alpha"),
		},
	})

	optionNames := binding.SortedOptionNames()
	if got, want := optionNames[0], MustOptionName("alpha"); got != want {
		t.Fatalf("SortedOptionNames()[0] = %q, want %q", got, want)
	}

	argumentNames := binding.SortedArgumentNames()
	if got, want := argumentNames[0], MustArgumentName("alpha"); got != want {
		t.Fatalf("SortedArgumentNames()[0] = %q, want %q", got, want)
	}
}
