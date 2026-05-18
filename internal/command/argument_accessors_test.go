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

// TestArgumentCopySemantics verifies detached default and allowed value slices.
func TestArgumentCopySemantics(t *testing.T) {
	t.Parallel()

	defaults := []string{"text"}
	allowed := []string{"text", "json"}

	argument := MustArgument(ArgumentSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		Requirement:   ArgumentRequirementOptional,
		DefaultValues: defaults,
		AllowedValues: allowed,
	})

	defaults[0] = "changed"
	allowed[0] = "changed"

	if got, want := argument.DefaultValues()[0], "text"; got != want {
		t.Fatalf("default changed through input slice: got %q, want %q", got, want)
	}

	if got, want := argument.AllowedValues()[0], "text"; got != want {
		t.Fatalf("allowed changed through input slice: got %q, want %q", got, want)
	}

	outDefaults := argument.DefaultValues()
	outDefaults[0] = "changed"

	if got, want := argument.DefaultValues()[0], "text"; got != want {
		t.Fatalf("default changed through output slice: got %q, want %q", got, want)
	}

	outAllowed := argument.AllowedValues()
	outAllowed[0] = "changed"

	if got, want := argument.AllowedValues()[0], "text"; got != want {
		t.Fatalf("allowed changed through output slice: got %q, want %q", got, want)
	}
}

// TestArgumentCounts verifies MinValues, MaxValues, and AcceptsCount.
func TestArgumentCounts(t *testing.T) {
	t.Parallel()

	requiredSingle := MustArgument(ArgumentSpec{
		Name: "package",
		Kind: OptionKindString,
	})

	if got, want := requiredSingle.MinValues(), 1; got != want {
		t.Fatalf("required single MinValues() = %d, want %d", got, want)
	}

	if maxValue, bounded := requiredSingle.MaxValues(); !bounded || maxValue != 1 {
		t.Fatalf("required single MaxValues() = %d, %v; want 1, true", maxValue, bounded)
	}

	if requiredSingle.AcceptsCount(0) {
		t.Fatalf("required single should not accept zero values")
	}

	if !requiredSingle.AcceptsCount(1) {
		t.Fatalf("required single should accept one value")
	}

	if requiredSingle.AcceptsCount(2) {
		t.Fatalf("required single should not accept two values")
	}

	optionalVariadic := MustArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Requirement: ArgumentRequirementOptional,
		Cardinality: ArgumentCardinalityVariadic,
	})

	if got, want := optionalVariadic.MinValues(), 0; got != want {
		t.Fatalf("optional variadic MinValues() = %d, want %d", got, want)
	}

	if _, bounded := optionalVariadic.MaxValues(); bounded {
		t.Fatalf("optional variadic MaxValues() bounded = true, want false")
	}

	if !optionalVariadic.AcceptsCount(0) {
		t.Fatalf("optional variadic should accept zero values")
	}

	if !optionalVariadic.AcceptsCount(5) {
		t.Fatalf("optional variadic should accept five values")
	}
}
