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

// TestNewUsageAcceptsValidUsage verifies full usage construction and
// normalization.
func TestNewUsageAcceptsValidUsage(t *testing.T) {
	t.Parallel()

	usage, err := NewUsage(UsageSpec{
		Syntax: "  bench   run   [flags]  ",
		Alternatives: []string{
			"bench run --suite <name>",
			"  bench   smoke   [flags]  ",
		},
	})
	if err != nil {
		t.Fatalf("NewUsage() returned unexpected error: %v", err)
	}

	if got, want := usage.String(), "bench run [flags]"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	assertStringSlicesEqual(t, usage.LineStrings(), []string{
		"bench run [flags]",
		"bench run --suite <name>",
		"bench smoke [flags]",
	})

	if !usage.IsValid() {
		t.Fatalf("IsValid() = false, want true")
	}
}

// TestNewUsageRejectsInvalidUsage verifies primary and alternative validation.
func TestNewUsageRejectsInvalidUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec UsageSpec
		want error
	}{
		{
			name: "empty syntax",
			spec: UsageSpec{},
			want: ErrEmptyUsage,
		},
		{
			name: "invalid syntax",
			spec: UsageSpec{
				Syntax: "bench\nrun",
			},
			want: ErrInvalidUsage,
		},
		{
			name: "invalid alternative",
			spec: UsageSpec{
				Syntax:       "bench run",
				Alternatives: []string{"bench\nsmoke"},
			},
			want: ErrInvalidUsage,
		},
		{
			name: "duplicate alternative",
			spec: UsageSpec{
				Syntax:       "bench run",
				Alternatives: []string{"bench smoke", "  bench   smoke  "},
			},
			want: ErrInvalidUsage,
		},
		{
			name: "alternative duplicates primary",
			spec: UsageSpec{
				Syntax:       "bench run",
				Alternatives: []string{"  bench   run  "},
			},
			want: ErrInvalidUsage,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewUsage(tt.spec)
			if err == nil {
				t.Fatalf("NewUsage() returned nil error")
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("NewUsage() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// TestUsageCopySemantics verifies that constructors and accessors detach
// mutable slices.
func TestUsageCopySemantics(t *testing.T) {
	t.Parallel()

	alternatives := []string{"bench smoke"}
	usage := MustUsage(UsageSpec{
		Syntax:       "bench run",
		Alternatives: alternatives,
	})

	alternatives[0] = "changed"
	if usage.Contains("changed") {
		t.Fatalf("usage changed through input alternatives slice")
	}

	out := usage.Alternatives()
	out[0] = MustUsageLine("changed")
	if usage.Contains("changed") {
		t.Fatalf("usage changed through Alternatives() slice")
	}

	spec := usage.Spec()
	spec.Alternatives[0] = "changed"
	if usage.Contains("changed") {
		t.Fatalf("usage changed through Spec() alternatives slice")
	}
}

// TestMustUsagePanicsForInvalidUsage verifies fail-fast static construction.
func TestMustUsagePanicsForInvalidUsage(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustUsage(UsageSpec{})
	})
}
