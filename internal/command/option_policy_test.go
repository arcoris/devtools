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
	"testing"
)

// TestNewOptionPolicyAcceptsExplicitPolicy verifies full policy construction.
func TestNewOptionPolicyAcceptsExplicitPolicy(t *testing.T) {
	t.Parallel()

	policy, err := NewOptionPolicy(OptionPolicySpec{
		Requirement: OptionRequirementRequired,
		Scope:       OptionScopeGlobal,
		Occurrence:  OptionOccurrenceMultiple,
		EmptyValue:  OptionEmptyValueReject,
		AllowedSources: []OptionSource{
			OptionSourceCommandLine,
			OptionSourceEnvironment,
		},
	})
	if err != nil {
		t.Fatalf("NewOptionPolicy() returned unexpected error: %v", err)
	}

	if !policy.IsRequired() {
		t.Fatalf("IsRequired() = false, want true")
	}

	if !policy.Scope().IsGlobal() {
		t.Fatalf("Scope().IsGlobal() = false, want true")
	}

	if !policy.IsRepeatable() {
		t.Fatalf("IsRepeatable() = false, want true")
	}

	if !policy.AllowsSource(OptionSourceCommandLine) {
		t.Fatalf("AllowsSource(command-line) = false, want true")
	}

	if policy.AllowsSource(OptionSourceDefault) {
		t.Fatalf("AllowsSource(default) = true, want false")
	}

	if !policy.AllowsOnlyExplicitSources() {
		t.Fatalf("AllowsOnlyExplicitSources() = false, want true")
	}
}

// TestNewOptionPolicyDefaults verifies ordinary default policy construction.
func TestNewOptionPolicyDefaults(t *testing.T) {
	t.Parallel()

	policy, err := NewOptionPolicy(OptionPolicySpec{})
	if err != nil {
		t.Fatalf("NewOptionPolicy() returned unexpected error: %v", err)
	}

	if !policy.IsOptional() {
		t.Fatalf("default policy IsOptional() = false, want true")
	}

	if !policy.IsLocal() {
		t.Fatalf("default policy IsLocal() = false, want true")
	}

	if policy.IsRepeatable() {
		t.Fatalf("default policy IsRepeatable() = true, want false")
	}

	if policy.AllowsEmptyValue() {
		t.Fatalf("default policy AllowsEmptyValue() = true, want false")
	}

	if got, want := len(policy.AllowedSources()), len(KnownOptionSources()); got != want {
		t.Fatalf("len(AllowedSources()) = %d, want %d", got, want)
	}
}

// TestNewOptionPolicyForKindDefaultsOccurrence verifies kind-aware occurrence defaults.
func TestNewOptionPolicyForKindDefaultsOccurrence(t *testing.T) {
	t.Parallel()

	scalar := MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{})
	if got, want := scalar.Occurrence(), OptionOccurrenceSingle; got != want {
		t.Fatalf("scalar occurrence = %q, want %q", got, want)
	}

	list := MustOptionPolicyForKind(OptionKindStringList, OptionPolicySpec{})
	if got, want := list.Occurrence(), OptionOccurrenceMultiple; got != want {
		t.Fatalf("list occurrence = %q, want %q", got, want)
	}
}

// TestNewOptionPolicyNormalizesAllowedSources verifies source order normalization.
func TestNewOptionPolicyNormalizesAllowedSources(t *testing.T) {
	t.Parallel()

	policy := MustOptionPolicy(OptionPolicySpec{
		AllowedSources: []OptionSource{
			OptionSourceCommandLine,
			OptionSourceDefault,
			OptionSourceEnvironment,
		},
	})

	got := policy.AllowedSources()
	want := []OptionSource{
		OptionSourceDefault,
		OptionSourceEnvironment,
		OptionSourceCommandLine,
	}

	if len(got) != len(want) {
		t.Fatalf("len(AllowedSources()) = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("AllowedSources()[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}
