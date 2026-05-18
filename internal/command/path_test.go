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
)

// TestNewPathAcceptsValidPaths verifies construction from valid path segments.
func TestNewPathAcceptsValidPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		segments []string
		want     string
	}{
		{
			name:     "root",
			segments: nil,
			want:     "",
		},
		{
			name:     "single",
			segments: []string{"check"},
			want:     "check",
		},
		{
			name:     "nested",
			segments: []string{"bench", "run"},
			want:     "bench run",
		},
		{
			name:     "kebab",
			segments: []string{"release-notes", "generate2"},
			want:     "release-notes generate2",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path, err := NewPath(tt.segments...)
			if err != nil {
				t.Fatalf("NewPath(%v) returned unexpected error: %v", tt.segments, err)
			}

			if got := path.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}

			if err := path.Validate(); err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}
		})
	}
}

// TestNewPathRejectsInvalidSegments verifies that Path uses generic
// command-name segment validation and wraps errors as ErrInvalidPath.
func TestNewPathRejectsInvalidSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		segments []string
	}{
		{
			name:     "empty segment",
			segments: []string{"bench", ""},
		},
		{
			name:     "uppercase",
			segments: []string{"bench", "Run"},
		},
		{
			name:     "underscore",
			segments: []string{"bench_run"},
		},
		{
			name:     "dot",
			segments: []string{"bench.run"},
		},
		{
			name:     "slash",
			segments: []string{"bench/run"},
		},
		{
			name:     "space",
			segments: []string{"bench run"},
		},
		{
			name:     "unicode",
			segments: []string{"ран"},
		},
		{
			name:     "starts with digit",
			segments: []string{"1bench"},
		},
		{
			name:     "starts with hyphen",
			segments: []string{"-bench"},
		},
		{
			name:     "trailing hyphen",
			segments: []string{"bench-"},
		},
		{
			name:     "repeated hyphen",
			segments: []string{"bench--run"},
		},
		{
			name:     "too long segment",
			segments: []string{strings.Repeat("a", maxPathSegmentLength+1)},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path, err := NewPath(tt.segments...)
			if err == nil {
				t.Fatalf("NewPath(%v) returned nil error and path %q", tt.segments, path)
			}

			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("NewPath(%v) error = %v, want ErrInvalidPath", tt.segments, err)
			}
		})
	}
}

// TestNewPathRejectsWholePathLimitViolations verifies path-level bounds that
// are independent from individual segment grammar.
func TestNewPathRejectsWholePathLimitViolations(t *testing.T) {
	t.Parallel()

	t.Run("too deep", func(t *testing.T) {
		t.Parallel()

		segments := make([]string, maxPathDepth+1)
		for index := range segments {
			segments[index] = "a"
		}

		path, err := NewPath(segments...)
		if err == nil {
			t.Fatalf("NewPath returned nil error and path %q", path)
		}

		if !errors.Is(err, ErrInvalidPath) {
			t.Fatalf("NewPath error = %v, want ErrInvalidPath", err)
		}
	})

	t.Run("too long", func(t *testing.T) {
		t.Parallel()

		segment := strings.Repeat("a", maxPathSegmentLength)
		segments := []string{segment, segment, segment, segment, segment}

		path, err := NewPath(segments...)
		if err == nil {
			t.Fatalf("NewPath returned nil error and path %q", path)
		}

		if !errors.Is(err, ErrInvalidPath) {
			t.Fatalf("NewPath error = %v, want ErrInvalidPath", err)
		}
	})
}

// TestNewPathPreservesSegmentValidationSentinels verifies that path validation
// keeps both path-level and command-name validation sentinels available.
func TestNewPathPreservesSegmentValidationSentinels(t *testing.T) {
	t.Parallel()

	_, err := NewPath("bench", "Run")
	if err == nil {
		t.Fatalf("NewPath returned nil error")
	}

	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("NewPath error = %v, want ErrInvalidPath", err)
	}

	if !errors.Is(err, ErrInvalidCommandNameSegment) {
		t.Fatalf("NewPath error = %v, want ErrInvalidCommandNameSegment", err)
	}
}

// TestMustPathReturnsValidPath verifies fail-fast construction for static paths.
func TestMustPathReturnsValidPath(t *testing.T) {
	t.Parallel()

	path := MustPath("profile", "cpu")

	if got, want := path.String(), "profile cpu"; got != want {
		t.Fatalf("MustPath returned %q, want %q", got, want)
	}
}

// TestMustPathPanicsForInvalidPath verifies fail-fast behavior.
func TestMustPathPanicsForInvalidPath(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustPath did not panic")
		}
	}()

	_ = MustPath("Profile", "CPU")
}

// TestParsePathParsesHumanFacingPaths verifies whitespace-based parsing.
func TestParsePathParsesHumanFacingPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want Path
	}{
		{
			name: "empty",
			raw:  "",
			want: RootPath(),
		},
		{
			name: "spaces only",
			raw:  "   ",
			want: RootPath(),
		},
		{
			name: "single",
			raw:  "check",
			want: MustPath("check"),
		},
		{
			name: "nested",
			raw:  "bench run",
			want: MustPath("bench", "run"),
		},
		{
			name: "trimmed and normalized",
			raw:  "  bench   run  ",
			want: MustPath("bench", "run"),
		},
		{
			name: "tabs",
			raw:  "bench\trun",
			want: MustPath("bench", "run"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParsePath(tt.raw)
			if err != nil {
				t.Fatalf("ParsePath(%q) returned unexpected error: %v", tt.raw, err)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("ParsePath(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

// TestParsePathRejectsInvalidPath verifies parsing plus segment validation.
func TestParsePathRejectsInvalidPath(t *testing.T) {
	t.Parallel()

	tests := []string{
		"Bench run",
		"bench run_now",
		"bench run.now",
		"bench /run",
		"bench ран",
	}

	for _, raw := range tests {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			path, err := ParsePath(raw)
			if err == nil {
				t.Fatalf("ParsePath(%q) returned nil error and path %q", raw, path)
			}

			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("ParsePath(%q) error = %v, want ErrInvalidPath", raw, err)
			}
		})
	}
}

// TestMustParsePathReturnsValidPath verifies static parse construction.
func TestMustParsePathReturnsValidPath(t *testing.T) {
	t.Parallel()

	path := MustParsePath("bench run")

	if got, want := path, MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("MustParsePath() = %q, want %q", got, want)
	}
}

// TestMustParsePathPanicsForInvalidPath verifies fail-fast parsing.
func TestMustParsePathPanicsForInvalidPath(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustParsePath did not panic")
		}
	}()

	_ = MustParsePath("Bench run")
}

// TestPathFromIDConvertsIDParts verifies conversion from stable ID to logical
// command path.
func TestPathFromIDConvertsIDParts(t *testing.T) {
	t.Parallel()

	path, err := PathFromID(MustID("bench.run"))
	if err != nil {
		t.Fatalf("PathFromID returned unexpected error: %v", err)
	}

	if got, want := path, MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("PathFromID() = %q, want %q", got, want)
	}
}

// TestPathFromIDRejectsInvalidID verifies that ID validation is preserved.
func TestPathFromIDRejectsInvalidID(t *testing.T) {
	t.Parallel()

	path, err := PathFromID(ID("Bench.Run"))
	if err == nil {
		t.Fatalf("PathFromID returned nil error and path %q", path)
	}

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("PathFromID error = %v, want ErrInvalidID", err)
	}
}

// TestPathRootSemantics verifies zero-value root behavior.
func TestPathRootSemantics(t *testing.T) {
	t.Parallel()

	var zero Path
	root := RootPath()

	if !zero.IsRoot() {
		t.Fatalf("zero path IsRoot() = false, want true")
	}

	if !zero.IsZero() {
		t.Fatalf("zero path IsZero() = false, want true")
	}

	if !zero.Equal(root) {
		t.Fatalf("zero path does not equal RootPath()")
	}

	if got := zero.String(); got != "" {
		t.Fatalf("zero path String() = %q, want empty string", got)
	}

	if got := zero.Key(); got != "" {
		t.Fatalf("zero path Key() = %q, want empty string", got)
	}
}

// TestPathTextEncoding verifies canonical text marshaling and parsing through
// encoding.TextMarshaler-style methods.
func TestPathTextEncoding(t *testing.T) {
	t.Parallel()

	path := MustPath("bench", "run")

	text, err := path.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText returned unexpected error: %v", err)
	}

	if got, want := string(text), "bench run"; got != want {
		t.Fatalf("MarshalText() = %q, want %q", got, want)
	}

	var decoded Path
	if err := decoded.UnmarshalText([]byte("  bench   run  ")); err != nil {
		t.Fatalf("UnmarshalText returned unexpected error: %v", err)
	}

	if !decoded.Equal(path) {
		t.Fatalf("UnmarshalText decoded %q, want %q", decoded, path)
	}
}

// TestPathTextEncodingRejectsInvalidValues verifies invalid manually
// constructed values and nil receivers are rejected.
func TestPathTextEncodingRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Path{segments: []string{"Bench"}}
	if _, err := invalid.MarshalText(); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("MarshalText invalid error = %v, want ErrInvalidPath", err)
	}

	var nilPath *Path
	if err := nilPath.UnmarshalText([]byte("bench")); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("UnmarshalText nil receiver error = %v, want ErrInvalidPath", err)
	}
}

// TestPathPartsReturnsDetachedCopy verifies value-object copy semantics.
func TestPathPartsReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	path := MustPath("bench", "run")

	parts := path.Parts()
	if got, want := fmt.Sprint(parts), "[bench run]"; got != want {
		t.Fatalf("Parts() = %s, want %s", got, want)
	}

	parts[0] = "changed"

	again := path.Parts()
	if got, want := fmt.Sprint(again), "[bench run]"; got != want {
		t.Fatalf("Parts() after mutation = %s, want %s", got, want)
	}
}

// TestPathLenAndDepth verifies segment count methods.
func TestPathLenAndDepth(t *testing.T) {
	t.Parallel()

	path := MustPath("a", "b", "c")

	if got, want := path.Len(), 3; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	if got, want := path.Depth(), 3; got != want {
		t.Fatalf("Depth() = %d, want %d", got, want)
	}
}

// TestPathAt verifies safe indexed access.
func TestPathAt(t *testing.T) {
	t.Parallel()

	path := MustPath("bench", "run")

	first, ok := path.At(0)
	if !ok || first != "bench" {
		t.Fatalf("At(0) = %q, %v; want bench, true", first, ok)
	}

	second, ok := path.At(1)
	if !ok || second != "run" {
		t.Fatalf("At(1) = %q, %v; want run, true", second, ok)
	}

	value, ok := path.At(2)
	if ok || value != "" {
		t.Fatalf("At(2) = %q, %v; want empty, false", value, ok)
	}

	value, ok = path.At(-1)
	if ok || value != "" {
		t.Fatalf("At(-1) = %q, %v; want empty, false", value, ok)
	}
}

// TestPathLeaf verifies final segment extraction.
func TestPathLeaf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path Path
		want string
	}{
		{
			name: "root",
			path: RootPath(),
			want: "",
		},
		{
			name: "single",
			path: MustPath("check"),
			want: "check",
		},
		{
			name: "nested",
			path: MustPath("bench", "run"),
			want: "run",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.path.Leaf(); got != tt.want {
				t.Fatalf("Leaf() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestPathParent verifies parent path calculation.
func TestPathParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      Path
		want      Path
		wantFound bool
	}{
		{
			name:      "root",
			path:      RootPath(),
			want:      RootPath(),
			wantFound: false,
		},
		{
			name:      "single",
			path:      MustPath("check"),
			want:      RootPath(),
			wantFound: true,
		},
		{
			name:      "nested",
			path:      MustPath("bench", "run"),
			want:      MustPath("bench"),
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := tt.path.Parent()

			if found != tt.wantFound {
				t.Fatalf("Parent() found = %v, want %v", found, tt.wantFound)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("Parent() = %q, want %q", got, tt.want)
			}

			if tt.path.HasParent() != tt.wantFound {
				t.Fatalf("HasParent() = %v, want %v", tt.path.HasParent(), tt.wantFound)
			}
		})
	}
}

// TestPathAppend verifies immutable-style child path construction.
func TestPathAppend(t *testing.T) {
	t.Parallel()

	parent := MustPath("bench")

	child, err := parent.Append("run")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("Append() = %q, want %q", got, want)
	}

	if got, want := parent, MustPath("bench"); !got.Equal(want) {
		t.Fatalf("Append mutated parent: got %q, want %q", got, want)
	}
}

// TestPathAppendRejectsInvalidSegment verifies append validation.
func TestPathAppendRejectsInvalidSegment(t *testing.T) {
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
		strings.Repeat("a", maxPathSegmentLength+1),
	}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			_, err := MustPath("bench").Append(segment)
			if err == nil {
				t.Fatalf("Append(%q) returned nil error", segment)
			}

			if !errors.Is(err, ErrInvalidPath) {
				t.Fatalf("Append(%q) error = %v, want ErrInvalidPath", segment, err)
			}
		})
	}
}

// TestPathMustAppendReturnsValidPath verifies static append construction.
func TestPathMustAppendReturnsValidPath(t *testing.T) {
	t.Parallel()

	got := MustPath("bench").MustAppend("run")
	want := MustPath("bench", "run")

	if !got.Equal(want) {
		t.Fatalf("MustAppend() = %q, want %q", got, want)
	}
}

// TestPathMustAppendPanicsForInvalidSegment verifies fail-fast append behavior.
func TestPathMustAppendPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustAppend did not panic")
		}
	}()

	_ = MustPath("bench").MustAppend("Run")
}

// TestPathAppendRejectsWholePathLimitViolations verifies that Append validates
// the resulting path, not only the appended segment.
func TestPathAppendRejectsWholePathLimitViolations(t *testing.T) {
	t.Parallel()

	segments := make([]string, maxPathDepth)
	for index := range segments {
		segments[index] = "a"
	}

	path := MustPath(segments...)

	child, err := path.Append("b")
	if err == nil {
		t.Fatalf("Append returned nil error and path %q", child)
	}

	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("Append error = %v, want ErrInvalidPath", err)
	}
}

// TestPathJoin verifies immutable path concatenation.
func TestPathJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Path
		right Path
		want  Path
	}{
		{
			name:  "root root",
			left:  RootPath(),
			right: RootPath(),
			want:  RootPath(),
		},
		{
			name:  "root right",
			left:  RootPath(),
			right: MustPath("bench"),
			want:  MustPath("bench"),
		},
		{
			name:  "left root",
			left:  MustPath("bench"),
			right: RootPath(),
			want:  MustPath("bench"),
		},
		{
			name:  "nested",
			left:  MustPath("bench"),
			right: MustPath("run", "fast"),
			want:  MustPath("bench", "run", "fast"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.left.Join(tt.right)
			if !got.Equal(tt.want) {
				t.Fatalf("Join() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestPathHasPrefix verifies hierarchical prefix checks.
func TestPathHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		path   Path
		prefix Path
		want   bool
	}{
		{
			name:   "root prefixes root",
			path:   RootPath(),
			prefix: RootPath(),
			want:   true,
		},
		{
			name:   "root prefixes nested",
			path:   MustPath("bench", "run"),
			prefix: RootPath(),
			want:   true,
		},
		{
			name:   "same path",
			path:   MustPath("bench", "run"),
			prefix: MustPath("bench", "run"),
			want:   true,
		},
		{
			name:   "parent prefix",
			path:   MustPath("bench", "run"),
			prefix: MustPath("bench"),
			want:   true,
		},
		{
			name:   "sibling false",
			path:   MustPath("bench", "run"),
			prefix: MustPath("bench", "compare"),
			want:   false,
		},
		{
			name:   "similar text false",
			path:   MustPath("benchmark", "run"),
			prefix: MustPath("bench"),
			want:   false,
		},
		{
			name:   "longer prefix false",
			path:   MustPath("bench"),
			prefix: MustPath("bench", "run"),
			want:   false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.path.HasPrefix(tt.prefix); got != tt.want {
				t.Fatalf("%q.HasPrefix(%q) = %v, want %v", tt.path, tt.prefix, got, tt.want)
			}
		})
	}
}

// TestPathTrimPrefix verifies relative-path extraction.
func TestPathTrimPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      Path
		prefix    Path
		want      Path
		wantFound bool
	}{
		{
			name:      "trim root from nested",
			path:      MustPath("bench", "run"),
			prefix:    RootPath(),
			want:      MustPath("bench", "run"),
			wantFound: true,
		},
		{
			name:      "trim parent",
			path:      MustPath("bench", "run"),
			prefix:    MustPath("bench"),
			want:      MustPath("run"),
			wantFound: true,
		},
		{
			name:      "trim same",
			path:      MustPath("bench", "run"),
			prefix:    MustPath("bench", "run"),
			want:      RootPath(),
			wantFound: true,
		},
		{
			name:      "non-prefix",
			path:      MustPath("bench", "run"),
			prefix:    MustPath("profile"),
			want:      RootPath(),
			wantFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := tt.path.TrimPrefix(tt.prefix)

			if found != tt.wantFound {
				t.Fatalf("TrimPrefix() found = %v, want %v", found, tt.wantFound)
			}

			if !got.Equal(tt.want) {
				t.Fatalf("TrimPrefix() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestPathEqual verifies path equality semantics.
func TestPathEqual(t *testing.T) {
	t.Parallel()

	if !MustPath("bench", "run").Equal(MustPath("bench", "run")) {
		t.Fatalf("same paths are not equal")
	}

	if MustPath("bench", "run").Equal(MustPath("bench", "compare")) {
		t.Fatalf("different paths are equal")
	}

	if RootPath().Equal(MustPath("bench")) {
		t.Fatalf("root equals non-root")
	}
}

// TestPathValidate verifies validation of manually constructed package-internal paths.
func TestPathValidate(t *testing.T) {
	t.Parallel()

	valid := Path{segments: []string{"bench", "run"}}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid path Validate() returned unexpected error: %v", err)
	}

	invalid := Path{segments: []string{"bench", "Run"}}
	err := invalid.Validate()
	if err == nil {
		t.Fatalf("invalid path Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("invalid path Validate() error = %v, want ErrInvalidPath", err)
	}
}

// TestPathIDConvertsPathToStableID verifies conversion from logical command
// path to stable command ID.
func TestPathIDConvertsPathToStableID(t *testing.T) {
	t.Parallel()

	id, err := MustPath("bench", "run").ID()
	if err != nil {
		t.Fatalf("ID returned unexpected error: %v", err)
	}

	if got, want := id, MustID("bench.run"); got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}
}

// TestPathIDRejectsRootAndInvalidPath verifies conversion errors for paths
// that cannot produce a stable command ID.
func TestPathIDRejectsRootAndInvalidPath(t *testing.T) {
	t.Parallel()

	if _, err := RootPath().ID(); !errors.Is(err, ErrEmptyID) {
		t.Fatalf("root ID error = %v, want ErrEmptyID", err)
	}

	invalid := Path{segments: []string{"Bench"}}
	if _, err := invalid.ID(); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("invalid path ID error = %v, want ErrInvalidPath", err)
	}
}
