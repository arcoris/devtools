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

// TestUsageMutationHelpers verifies immutable-style update helpers.
func TestUsageMutationHelpers(t *testing.T) {
	t.Parallel()

	usage := MustSimpleUsage("bench run")

	next := usage.MustWithSyntax("bench smoke").
		MustWithAlternative("bench smoke --quick").
		MustWithAlternatives([]string{"bench smoke --quick", "bench smoke --full"})

	assertStringSlicesEqual(t, next.LineStrings(), []string{
		"bench smoke",
		"bench smoke --quick",
		"bench smoke --full",
	})

	if got, want := usage.String(), "bench run"; got != want {
		t.Fatalf("original usage changed: got %q, want %q", got, want)
	}

	withoutOne := next.WithoutAlternative("  bench   smoke   --quick  ")
	assertStringSlicesEqual(t, withoutOne.LineStrings(), []string{
		"bench smoke",
		"bench smoke --full",
	})

	withoutMissing := withoutOne.WithoutAlternative("missing")
	assertStringSlicesEqual(t, withoutMissing.LineStrings(), []string{
		"bench smoke",
		"bench smoke --full",
	})

	withoutAlternatives := next.WithoutAlternatives()
	assertStringSlicesEqual(t, withoutAlternatives.LineStrings(), []string{"bench smoke"})
}

// TestUsageMutationRejectsInvalidUpdates verifies update helper validation.
func TestUsageMutationRejectsInvalidUpdates(t *testing.T) {
	t.Parallel()

	usage := MustSimpleUsage("bench run")

	tests := []struct {
		name string
		fn   func() (Usage, error)
	}{
		{
			name: "invalid syntax",
			fn: func() (Usage, error) {
				return usage.WithSyntax("bench\nrun")
			},
		},
		{
			name: "invalid alternatives",
			fn: func() (Usage, error) {
				return usage.WithAlternatives([]string{"bench\nsmoke"})
			},
		},
		{
			name: "duplicate alternative",
			fn: func() (Usage, error) {
				return usage.WithAlternatives([]string{"bench smoke", "bench smoke"})
			},
		},
		{
			name: "invalid alternative append",
			fn: func() (Usage, error) {
				return usage.WithAlternative("bench\nsmoke")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.fn()
			if err == nil {
				t.Fatalf("%s returned nil error", tt.name)
			}

			if !errors.Is(err, ErrInvalidUsage) {
				t.Fatalf("%s error = %v, want ErrInvalidUsage", tt.name, err)
			}
		})
	}
}

// TestUsageMustMutationHelpersPanic verifies fail-fast update helpers.
func TestUsageMustMutationHelpersPanic(t *testing.T) {
	t.Parallel()

	usage := MustSimpleUsage("bench run")

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "syntax",
			fn: func() {
				_ = usage.MustWithSyntax("bench\nrun")
			},
		},
		{
			name: "alternatives",
			fn: func() {
				_ = usage.MustWithAlternatives([]string{"bench\nsmoke"})
			},
		},
		{
			name: "alternative",
			fn: func() {
				_ = usage.MustWithAlternative("bench\nsmoke")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assertPanics(t, tt.fn)
		})
	}
}
