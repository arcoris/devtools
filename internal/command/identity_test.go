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
	"fmt"
	"strings"
	"testing"

	"arcoris.dev/devtools/internal/textvalidate"
)

// TestNewIDAcceptsValidIDs verifies that structurally valid command IDs are
// accepted by the ID value object.
func TestNewIDAcceptsValidIDs(t *testing.T) {
	t.Parallel()

	tests := []string{
		"check",
		"test",
		"bench.run",
		"bench.compare",
		"profile.cpu",
		"trace.pprof",
		"config.validate",
		"release-notes.generate",
		"a",
		"a1",
		"a-1",
		"a.b.c",
	}

	for _, raw := range tests {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			id, err := NewID(raw)
			if err != nil {
				t.Fatalf("NewID(%q) returned unexpected error: %v", raw, err)
			}

			if got := id.String(); got != raw {
				t.Fatalf("String() = %q, want %q", got, raw)
			}

			if !id.IsValid() {
				t.Fatalf("IsValid() = false, want true")
			}
		})
	}
}

// TestNewIDRejectsInvalidIDs verifies only ID-level validation: empty IDs,
// invalid separators, excessive complete ID length, and invalid ID segments.
func TestNewIDRejectsInvalidIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{
			name: "empty",
			raw:  "",
			err:  ErrEmptyID,
		},
		{
			name: "leading separator",
			raw:  ".bench",
			err:  ErrInvalidID,
		},
		{
			name: "trailing separator",
			raw:  "bench.",
			err:  ErrInvalidID,
		},
		{
			name: "consecutive separators",
			raw:  "bench..run",
			err:  ErrInvalidID,
		},
		{
			name: "invalid segment",
			raw:  "bench.Run",
			err:  ErrInvalidID,
		},
		{
			name: "trailing hyphen segment",
			raw:  "bench-.run",
			err:  ErrInvalidID,
		},
		{
			name: "repeated hyphen segment",
			raw:  "bench--run",
			err:  ErrInvalidID,
		},
		{
			name: "too long complete id",
			raw:  "a." + strings.Repeat("b", maxIDLength),
			err:  ErrInvalidID,
		},
		{
			name: "too long segment",
			raw:  "a" + strings.Repeat("b", maxIDSegmentLength),
			err:  ErrInvalidID,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewID(tt.raw)
			if err == nil {
				t.Fatalf("NewID(%q) returned nil error and ID %q", tt.raw, id)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewID(%q) error = %v, want errors.Is(..., %v)", tt.raw, err, tt.err)
			}

			if ID(tt.raw).IsValid() {
				t.Fatalf("ID(%q).IsValid() = true, want false", tt.raw)
			}
		})
	}
}

// TestNewIDPreservesSegmentValidationSentinels verifies that ID-level errors
// still expose the underlying command-name and shared text-validation errors.
func TestNewIDPreservesSegmentValidationSentinels(t *testing.T) {
	t.Parallel()

	_, err := NewID("bench.Run")
	if err == nil {
		t.Fatalf("NewID returned nil error")
	}

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("NewID error = %v, want ErrInvalidID", err)
	}

	if !errors.Is(err, ErrInvalidCommandNameSegment) {
		t.Fatalf("NewID error = %v, want ErrInvalidCommandNameSegment", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidKebabSegment) {
		t.Fatalf("NewID error = %v, want ErrInvalidKebabSegment", err)
	}
}

// TestParseIDIsAliasForNewID verifies that ParseID preserves NewID behavior.
func TestParseIDIsAliasForNewID(t *testing.T) {
	t.Parallel()

	const raw = "bench.run"

	fromNew, err := NewID(raw)
	if err != nil {
		t.Fatalf("NewID(%q) returned unexpected error: %v", raw, err)
	}

	fromParse, err := ParseID(raw)
	if err != nil {
		t.Fatalf("ParseID(%q) returned unexpected error: %v", raw, err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseID(%q) = %q, want %q", raw, fromParse, fromNew)
	}
}

// TestMustIDReturnsValidID verifies the happy path for static ID declarations.
func TestMustIDReturnsValidID(t *testing.T) {
	t.Parallel()

	id := MustID("profile.cpu")

	if got, want := id.String(), "profile.cpu"; got != want {
		t.Fatalf("MustID returned %q, want %q", got, want)
	}
}

// TestMustIDPanicsForInvalidID verifies fail-fast behavior for invalid static IDs.
func TestMustIDPanicsForInvalidID(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustID did not panic")
		}
	}()

	_ = MustID("Profile.CPU")
}

// TestIDPartsReturnsDetachedCopy verifies that callers cannot mutate internal state.
func TestIDPartsReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	id := MustID("bench.run")

	parts := id.Parts()
	if got, want := fmt.Sprint(parts), "[bench run]"; got != want {
		t.Fatalf("Parts() = %s, want %s", got, want)
	}

	parts[0] = "changed"

	again := id.Parts()
	if got, want := fmt.Sprint(again), "[bench run]"; got != want {
		t.Fatalf("Parts() after modifying previous result = %s, want %s", got, want)
	}
}

// TestIDPartsForZeroID verifies zero-value behavior.
func TestIDPartsForZeroID(t *testing.T) {
	t.Parallel()

	var id ID

	if parts := id.Parts(); parts != nil {
		t.Fatalf("zero ID Parts() = %#v, want nil", parts)
	}
}

// TestIDDepth verifies hierarchical depth calculation.
func TestIDDepth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   ID
		want int
	}{
		{
			name: "zero",
			id:   "",
			want: 0,
		},
		{
			name: "root",
			id:   MustID("check"),
			want: 1,
		},
		{
			name: "child",
			id:   MustID("bench.run"),
			want: 2,
		},
		{
			name: "deep",
			id:   MustID("a.b.c"),
			want: 3,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.Depth(); got != tt.want {
				t.Fatalf("Depth() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestIDLeaf verifies leaf-segment extraction.
func TestIDLeaf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "zero",
			id:   "",
			want: "",
		},
		{
			name: "root",
			id:   MustID("check"),
			want: "check",
		},
		{
			name: "child",
			id:   MustID("bench.run"),
			want: "run",
		},
		{
			name: "deep",
			id:   MustID("a.b.c"),
			want: "c",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.Leaf(); got != tt.want {
				t.Fatalf("Leaf() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestIDParent verifies parent ID extraction.
func TestIDParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        ID
		want      ID
		wantFound bool
	}{
		{
			name:      "zero",
			id:        "",
			want:      "",
			wantFound: false,
		},
		{
			name:      "root",
			id:        MustID("check"),
			want:      "",
			wantFound: false,
		},
		{
			name:      "child",
			id:        MustID("bench.run"),
			want:      MustID("bench"),
			wantFound: true,
		},
		{
			name:      "deep",
			id:        MustID("a.b.c"),
			want:      MustID("a.b"),
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := tt.id.Parent()

			if found != tt.wantFound {
				t.Fatalf("Parent() found = %v, want %v", found, tt.wantFound)
			}

			if got != tt.want {
				t.Fatalf("Parent() ID = %q, want %q", got, tt.want)
			}

			if tt.id.HasParent() != tt.wantFound {
				t.Fatalf("HasParent() = %v, want %v", tt.id.HasParent(), tt.wantFound)
			}
		})
	}
}

// TestIDHasPrefix verifies hierarchical prefix checks.
func TestIDHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		id     ID
		prefix ID
		want   bool
	}{
		{
			name:   "same id",
			id:     MustID("bench.run"),
			prefix: MustID("bench.run"),
			want:   true,
		},
		{
			name:   "parent prefix",
			id:     MustID("bench.run"),
			prefix: MustID("bench"),
			want:   true,
		},
		{
			name:   "grandparent prefix",
			id:     MustID("a.b.c"),
			prefix: MustID("a"),
			want:   true,
		},
		{
			name:   "similar text is not hierarchy",
			id:     MustID("benchmark.run"),
			prefix: MustID("bench"),
			want:   false,
		},
		{
			name:   "sibling is false",
			id:     MustID("bench.run"),
			prefix: MustID("bench.compare"),
			want:   false,
		},
		{
			name:   "empty id",
			id:     "",
			prefix: MustID("bench"),
			want:   false,
		},
		{
			name:   "empty prefix",
			id:     MustID("bench.run"),
			prefix: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.HasPrefix(tt.prefix); got != tt.want {
				t.Fatalf("%q.HasPrefix(%q) = %v, want %v", tt.id, tt.prefix, got, tt.want)
			}
		})
	}
}

// TestIDAppend verifies child ID construction.
func TestIDAppend(t *testing.T) {
	t.Parallel()

	parent := MustID("bench")

	child, err := parent.Append("run")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustID("bench.run"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

// TestIDAppendToZeroID verifies that appending to a zero ID creates a root ID.
func TestIDAppendToZeroID(t *testing.T) {
	t.Parallel()

	var id ID

	child, err := id.Append("check")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustID("check"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

// TestIDAppendRejectsInvalidSegment verifies that Append accepts only one valid segment.
func TestIDAppendRejectsInvalidSegment(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"Run",
		"run.now",
		"run_now",
		"run-",
		"run--fast",
		"1run",
		"-run",
		"ран",
	}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			_, err := MustID("bench").Append(segment)
			if err == nil {
				t.Fatalf("Append(%q) returned nil error", segment)
			}

			if !errors.Is(err, ErrInvalidID) {
				t.Fatalf("Append(%q) error = %v, want ErrInvalidID", segment, err)
			}
		})
	}
}

// TestIDAppendRejectsTooLongChild verifies that Append validates the complete
// child ID, not only the parent and appended segment independently.
func TestIDAppendRejectsTooLongChild(t *testing.T) {
	t.Parallel()

	parent := MustID(strings.Join([]string{
		strings.Repeat("a", maxIDSegmentLength),
		strings.Repeat("b", maxIDSegmentLength),
		strings.Repeat("c", maxIDSegmentLength),
		strings.Repeat("d", maxIDSegmentLength),
	}, IDSeparator))

	_, err := parent.Append("e")
	if err == nil {
		t.Fatalf("Append returned nil error")
	}

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("Append error = %v, want ErrInvalidID", err)
	}
}

// TestIDAppendRejectsInvalidParent verifies that invalid parents are rejected.
func TestIDAppendRejectsInvalidParent(t *testing.T) {
	t.Parallel()

	_, err := ID("Bench").Append("run")
	if err == nil {
		t.Fatalf("Append returned nil error")
	}

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("Append error = %v, want ErrInvalidID", err)
	}
}

// TestIDMustAppendReturnsValidID verifies the happy path for static child IDs.
func TestIDMustAppendReturnsValidID(t *testing.T) {
	t.Parallel()

	id := MustID("bench").MustAppend("run")

	if got, want := id, MustID("bench.run"); got != want {
		t.Fatalf("MustAppend() = %q, want %q", got, want)
	}
}

// TestIDMustAppendPanicsForInvalidSegment verifies fail-fast behavior.
func TestIDMustAppendPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustAppend did not panic")
		}
	}()

	_ = MustID("bench").MustAppend("Run")
}

// TestIDIsZero verifies zero-value detection.
func TestIDIsZero(t *testing.T) {
	t.Parallel()

	var zero ID
	if !zero.IsZero() {
		t.Fatalf("zero ID IsZero() = false, want true")
	}

	if MustID("check").IsZero() {
		t.Fatalf("non-zero ID IsZero() = true, want false")
	}
}
