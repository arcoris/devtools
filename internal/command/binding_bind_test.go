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

// TestBindingBindAppliesDefaultsAndCanonicalizesAliases verifies binding output.
func TestBindingBindAppliesDefaultsAndCanonicalizesAliases(t *testing.T) {
	t.Parallel()

	binding := MustBinding(BindingSpec{
		Options: []Option{
			bindingTestFormatOption(),
			bindingTestPackageOption(),
		},
		Arguments: []Argument{
			bindingTestSuiteArgument(),
			bindingTestPackageArgument(),
		},
	})

	inputOption := MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json")

	bound, err := binding.Bind(BindingValueSpec{
		OptionValues: []OptionValue{
			inputOption,
		},
		PositionalValues: []string{"stable", "./internal/..."},
	})
	if err != nil {
		t.Fatalf("Bind() returned unexpected error: %v", err)
	}

	if got, want := bound.OptionCount(), 2; got != want {
		t.Fatalf("OptionCount() = %d, want %d", got, want)
	}

	format, ok := bound.Option(MustOptionName("format"))
	if !ok {
		t.Fatalf("bound option format not found")
	}

	if got, want := format.Name(), MustOptionName("format"); got != want {
		t.Fatalf("canonical option name = %q, want %q", got, want)
	}

	if got, want := format.MustValue(), "json"; got != want {
		t.Fatalf("format value = %q, want %q", got, want)
	}

	pkg, ok := bound.Option(MustOptionName("package"))
	if !ok {
		t.Fatalf("bound option package not found")
	}

	if !pkg.IsDefault() {
		t.Fatalf("package option should be default-bound")
	}

	suite, ok := bound.Argument(MustArgumentName("suite"))
	if !ok {
		t.Fatalf("bound argument suite not found")
	}

	if got, want := suite.MustValue(), "stable"; got != want {
		t.Fatalf("suite value = %q, want %q", got, want)
	}
}
