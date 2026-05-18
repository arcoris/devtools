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

// TestMustOptionPolicyPanicsForInvalidPolicy verifies fail-fast construction.
func TestMustOptionPolicyPanicsForInvalidPolicy(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustOptionPolicy did not panic")
		}
	}()

	_ = MustOptionPolicy(OptionPolicySpec{
		Requirement: OptionRequirement("mandatory"),
	})
}

// TestNewOptionPolicyForKindRejectsKindInconsistency verifies kind-aware checks.
func TestNewOptionPolicyForKindRejectsKindInconsistency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind OptionKind
		spec OptionPolicySpec
	}{
		{
			name: "invalid kind",
			kind: OptionKind("path"),
			spec: OptionPolicySpec{},
		},
		{
			name: "list kind with single occurrence",
			kind: OptionKindStringList,
			spec: OptionPolicySpec{
				Occurrence: OptionOccurrenceSingle,
			},
		},
		{
			name: "bool allows empty",
			kind: OptionKindBool,
			spec: OptionPolicySpec{
				EmptyValue: OptionEmptyValueAllow,
			},
		},
		{
			name: "int allows empty",
			kind: OptionKindInt,
			spec: OptionPolicySpec{
				EmptyValue: OptionEmptyValueAllow,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionPolicyForKind(tt.kind, tt.spec)
			if err == nil {
				t.Fatalf("NewOptionPolicyForKind() returned nil error")
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionPolicyForKind() error = %v, want ErrInvalidOptionPolicy", err)
			}
		})
	}
}

// TestNewOptionPolicyRejectsInvalidPolicy verifies policy validation.
func TestNewOptionPolicyRejectsInvalidPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec OptionPolicySpec
	}{
		{
			name: "invalid requirement",
			spec: OptionPolicySpec{
				Requirement: OptionRequirement("mandatory"),
			},
		},
		{
			name: "invalid scope",
			spec: OptionPolicySpec{
				Scope: OptionScope("tree"),
			},
		},
		{
			name: "invalid occurrence",
			spec: OptionPolicySpec{
				Occurrence: OptionOccurrence("many"),
			},
		},
		{
			name: "invalid empty policy",
			spec: OptionPolicySpec{
				EmptyValue: OptionEmptyValuePolicy("allow"),
			},
		},
		{
			name: "invalid source",
			spec: OptionPolicySpec{
				AllowedSources: []OptionSource{OptionSource("file")},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionPolicy(tt.spec)
			if err == nil {
				t.Fatalf("NewOptionPolicy() returned nil error")
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionPolicy() error = %v, want ErrInvalidOptionPolicy", err)
			}
		})
	}
}
