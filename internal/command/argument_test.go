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

// TestMustArgumentPanicsForInvalidArgument verifies fail-fast construction.
func TestMustArgumentPanicsForInvalidArgument(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustArgument did not panic")
		}
	}()

	_ = MustArgument(ArgumentSpec{
		Name: "Bad",
		Kind: OptionKindString,
	})
}

// TestNewArgumentAcceptsEnumArgument verifies enum allowed-value validation.
func TestNewArgumentAcceptsEnumArgument(t *testing.T) {
	t.Parallel()

	argument, err := NewArgument(ArgumentSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json", "markdown"},
	})
	if err != nil {
		t.Fatalf("NewArgument() returned unexpected error: %v", err)
	}

	if !argument.HasAllowedValues() {
		t.Fatalf("HasAllowedValues() = false, want true")
	}

	if !argument.AllowsValue("json") {
		t.Fatalf("AllowsValue(json) = false, want true")
	}

	if argument.AllowsValue("xml") {
		t.Fatalf("AllowsValue(xml) = true, want false")
	}
}

// TestNewArgumentAcceptsValidArgument verifies full argument declaration
// construction.
func TestNewArgumentAcceptsValidArgument(t *testing.T) {
	t.Parallel()

	argument, err := NewArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Requirement: ArgumentRequirementOptional,
		Cardinality: ArgumentCardinalitySingle,
		Metavar:     "PATTERN",
		DefaultValues: []string{
			"./...",
		},
		Documentation: MustSummaryDocumentation("Package pattern to test."),
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewArgument() returned unexpected error: %v", err)
	}

	if got, want := argument.Name(), MustArgumentName("package"); got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if got, want := argument.Kind(), OptionKindString; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := argument.Metavar(), "PATTERN"; got != want {
		t.Fatalf("Metavar() = %q, want %q", got, want)
	}

	if !argument.IsOptional() {
		t.Fatalf("IsOptional() = false, want true")
	}

	if argument.IsVariadic() {
		t.Fatalf("IsVariadic() = true, want false")
	}

	if got, ok := argument.DefaultValue(); !ok || got != "./..." {
		t.Fatalf("DefaultValue() = %q, %v; want ./..., true", got, ok)
	}
}

// TestNewArgumentAcceptsVariadicArgument verifies variadic declaration
// construction.
func TestNewArgumentAcceptsVariadicArgument(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Requirement: ArgumentRequirementOptional,
		Cardinality: ArgumentCardinalityVariadic,
	})

	if !argument.IsVariadic() {
		t.Fatalf("IsVariadic() = false, want true")
	}

	if !argument.AcceptsCount(0) {
		t.Fatalf("optional variadic should accept zero values")
	}

	if !argument.AcceptsCount(10) {
		t.Fatalf("optional variadic should accept many values")
	}
}

// TestNewArgumentDefaults verifies default requirement, cardinality, empty-value
// policy, visibility, and metavar.
func TestNewArgumentDefaults(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name: "bench-time",
		Kind: OptionKindDuration,
	})

	if !argument.IsRequired() {
		t.Fatalf("default argument IsRequired() = false, want true")
	}

	if !argument.IsSingle() {
		t.Fatalf("default argument IsSingle() = false, want true")
	}

	if argument.EmptyValue().AllowsEmpty() {
		t.Fatalf("default argument AllowsEmpty() = true, want false")
	}

	if got, want := argument.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if got, want := argument.Metavar(), "BENCH_TIME"; got != want {
		t.Fatalf("Metavar() = %q, want %q", got, want)
	}
}
