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
)

const (
	// PathSeparator is the canonical human-facing separator used by Path.String.
	//
	// A Path represents a logical command-tree path, not a stable command ID.
	// Therefore it is rendered like a command invocation path:
	//
	//   - Path{"bench", "run"} is rendered as "bench run";
	//   - Path{"profile", "cpu"} is rendered as "profile cpu".
	//
	// Stable machine-facing identifiers are represented by ID and use a
	// different separator.
	PathSeparator = " "

	// maxPathLength is the maximum allowed byte length of a canonical command
	// path string.
	//
	// Command paths are used in diagnostics, generated docs, and adapter lookup
	// tables. Keeping them compact prevents accidental free-form display text
	// from becoming structural command metadata.
	maxPathLength = 255

	// maxPathDepth is the maximum number of command-tree segments in a path.
	maxPathDepth = 32

	// maxPathSegmentLength is the maximum allowed byte length of one path
	// segment.
	maxPathSegmentLength = 63
)

var (
	// ErrInvalidPath reports that a command path violates the path grammar.
	ErrInvalidPath = errors.New("command path is invalid")
)

// Path is an adapter-neutral logical location in the command tree.
//
// Path is intentionally separate from ID:
//
//   - ID is a stable machine-facing key such as "bench.run";
//   - Path is a logical command-tree path such as "bench run".
//
// A CLI adapter such as Cobra may render a Path directly as command hierarchy,
// but Path itself is not Cobra-specific. It is used by registries, traversal,
// help generation, documentation generation, tests, diagnostics, and adapter
// layers.
//
// Path is a value object. Its internal segment slice is not exposed. All
// constructors and mutating-style methods copy their inputs and outputs.
//
// The root path is represented by the zero value:
//
//   - Path{} is root;
//   - RootPath() returns root;
//   - ParsePath("") returns root;
//   - ParsePath("   ") returns root.
type Path struct {
	segments []string
}

// RootPath returns the root command-tree path.
//
// The root path has no segments. It is a valid path and is treated as a prefix
// of every other path.
func RootPath() Path {
	return Path{}
}

// NewPath validates segments and returns a Path.
//
// Passing no segments returns the root path.
//
// Each segment uses the command-name segment grammar from validate.go:
// lowercase ASCII letter first, then lowercase ASCII letters, digits, or
// interior hyphens. Path also applies command-tree bounds for segment length,
// total canonical length, and depth.
func NewPath(segments ...string) (Path, error) {
	if err := validatePathSegments(segments); err != nil {
		return Path{}, err
	}

	return Path{segments: cloneStringSlice(segments)}, nil
}

// MustPath validates segments and returns a Path.
//
// MustPath panics on invalid input. It is intended for static command
// definitions and tests where invalid paths are programmer errors.
func MustPath(segments ...string) Path {
	path, err := NewPath(segments...)
	if err != nil {
		panic(err)
	}

	return path
}

// ParsePath parses a human-facing command path string.
//
// ParsePath uses strings.Fields, so any run of Unicode whitespace is treated as
// a separator. The canonical string form produced by String always uses a
// single ASCII space.
//
// Examples:
//
//   - ParsePath("") returns the root path;
//   - ParsePath("bench run") returns Path("bench", "run");
//   - ParsePath("  bench   run  ") returns Path("bench", "run").
func ParsePath(raw string) (Path, error) {
	segments := strings.Fields(raw)

	return NewPath(segments...)
}

// MustParsePath parses a human-facing command path string and panics on error.
//
// MustParsePath is intended for static declarations and tests.
func MustParsePath(raw string) Path {
	path, err := ParsePath(raw)
	if err != nil {
		panic(err)
	}

	return path
}

// PathFromID converts a stable command ID into its human-facing command path.
func PathFromID(id ID) (Path, error) {
	if err := id.Validate(); err != nil {
		return Path{}, err
	}

	return NewPath(id.Parts()...)
}

// String returns the canonical human-facing command path.
//
// The root path is rendered as an empty string.
func (path Path) String() string {
	return strings.Join(path.segments, PathSeparator)
}

// Key returns a stable map key for the path.
//
// Key currently matches String. The separate method makes map-key usage
// explicit and keeps call sites semantically clear.
func (path Path) Key() string {
	return path.String()
}

// MarshalText returns the canonical text form of path.
//
// MarshalText validates package-internal or zero-value constructed paths before
// emitting text, so malformed values do not silently enter config or report
// output.
func (path Path) MarshalText() ([]byte, error) {
	if err := path.Validate(); err != nil {
		return nil, err
	}

	return []byte(path.String()), nil
}

// UnmarshalText parses text into path.
func (path *Path) UnmarshalText(text []byte) error {
	if path == nil {
		return fmt.Errorf("%w: cannot unmarshal into nil *Path", ErrInvalidPath)
	}

	parsed, err := ParsePath(string(text))
	if err != nil {
		return err
	}

	*path = parsed

	return nil
}

// ID converts path into the equivalent stable command ID.
//
// The root path has no command ID and returns ErrEmptyID.
func (path Path) ID() (ID, error) {
	if path.IsRoot() {
		return "", ErrEmptyID
	}

	if err := path.Validate(); err != nil {
		return "", err
	}

	return NewID(strings.Join(path.segments, IDSeparator))
}

// IsRoot reports whether the path is the command-tree root.
func (path Path) IsRoot() bool {
	return len(path.segments) == 0
}

// IsZero reports whether the path is the zero value.
//
// For Path, the zero value is intentionally the root path.
func (path Path) IsZero() bool {
	return path.IsRoot()
}

// Len returns the number of path segments.
func (path Path) Len() int {
	return len(path.segments)
}

// Depth returns the number of path segments.
//
// Depth is an alias for Len and is provided for symmetry with ID.
func (path Path) Depth() int {
	return path.Len()
}

// Parts returns a detached copy of the path segments.
//
// The returned slice can be safely modified by the caller.
func (path Path) Parts() []string {
	return cloneStringSlice(path.segments)
}

// At returns the segment at index.
//
// The second return value is false when index is out of range. At never panics.
func (path Path) At(index int) (string, bool) {
	if index < 0 || index >= len(path.segments) {
		return "", false
	}

	return path.segments[index], true
}

// Leaf returns the final path segment.
//
// For the root path, Leaf returns an empty string.
func (path Path) Leaf() string {
	if path.IsRoot() {
		return ""
	}

	return path.segments[len(path.segments)-1]
}

// Parent returns the parent path and whether such parent exists.
//
// Examples:
//
//   - root returns root, false;
//   - "check" returns root, true;
//   - "bench run" returns "bench", true.
func (path Path) Parent() (Path, bool) {
	if path.IsRoot() {
		return RootPath(), false
	}

	parent := Path{
		segments: cloneStringSlice(path.segments[:len(path.segments)-1]),
	}

	return parent, true
}

// HasParent reports whether the path has a parent path.
func (path Path) HasParent() bool {
	_, ok := path.Parent()
	return ok
}

// Append returns a new path with segment appended.
//
// Append validates the resulting path and never modifies the receiver.
func (path Path) Append(segment string) (Path, error) {
	next := make([]string, 0, len(path.segments)+1)
	next = append(next, path.segments...)
	next = append(next, segment)

	if err := validatePathSegments(next); err != nil {
		return Path{}, err
	}

	return Path{segments: next}, nil
}

// MustAppend returns a new path with segment appended and panics on error.
//
// MustAppend is intended for static command definitions and tests.
func (path Path) MustAppend(segment string) Path {
	next, err := path.Append(segment)
	if err != nil {
		panic(err)
	}

	return next
}

// Join returns a new path formed by appending another path's segments.
//
// Join never modifies either receiver or argument. Joining with root returns a
// detached copy of the non-root side.
func (path Path) Join(other Path) Path {
	if path.IsRoot() {
		return Path{segments: cloneStringSlice(other.segments)}
	}

	if other.IsRoot() {
		return Path{segments: cloneStringSlice(path.segments)}
	}

	joined := make([]string, 0, len(path.segments)+len(other.segments))
	joined = append(joined, path.segments...)
	joined = append(joined, other.segments...)

	return Path{segments: joined}
}

// Equal reports whether two paths contain exactly the same segments.
func (path Path) Equal(other Path) bool {
	if len(path.segments) != len(other.segments) {
		return false
	}

	for index := range path.segments {
		if path.segments[index] != other.segments[index] {
			return false
		}
	}

	return true
}

// HasPrefix reports whether path is equal to prefix or belongs to prefix's
// subtree.
//
// The root path is a prefix of every path, including itself.
func (path Path) HasPrefix(prefix Path) bool {
	if len(prefix.segments) > len(path.segments) {
		return false
	}

	for index := range prefix.segments {
		if path.segments[index] != prefix.segments[index] {
			return false
		}
	}

	return true
}

// TrimPrefix removes prefix from path and returns the remaining relative path.
//
// The second return value is false when prefix is not a prefix of path.
//
// Examples:
//
//   - "bench run".TrimPrefix("bench") returns "run", true;
//   - "bench run".TrimPrefix(root) returns "bench run", true;
//   - "bench run".TrimPrefix("profile") returns root, false.
func (path Path) TrimPrefix(prefix Path) (Path, bool) {
	if !path.HasPrefix(prefix) {
		return RootPath(), false
	}

	rest := path.segments[len(prefix.segments):]

	return Path{segments: cloneStringSlice(rest)}, true
}

// clone returns a detached copy of path.
func (path Path) clone() Path {
	return Path{segments: cloneStringSlice(path.segments)}
}

// Validate verifies that the path satisfies the command path grammar.
//
// The root path is valid.
func (path Path) Validate() error {
	return validatePathSegments(path.segments)
}

// validatePathSegments validates all path segments plus whole-path limits.
func validatePathSegments(segments []string) error {
	if len(segments) > maxPathDepth {
		return fmt.Errorf(
			"%w: depth %d exceeds maximum depth %d",
			ErrInvalidPath,
			len(segments),
			maxPathDepth,
		)
	}

	if length := canonicalPathLength(segments); length > maxPathLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidPath,
			length,
			maxPathLength,
		)
	}

	for index, segment := range segments {
		if err := validatePathSegment(index, segment); err != nil {
			return err
		}
	}

	return nil
}

// validatePathSegment validates one path segment and wraps command-name
// validation errors with path-specific diagnostics.
func validatePathSegment(index int, segment string) error {
	if len(segment) > maxPathSegmentLength {
		return fmt.Errorf(
			"%w: segment %d length %d exceeds maximum length %d",
			ErrInvalidPath,
			index,
			len(segment),
			maxPathSegmentLength,
		)
	}

	if err := validateCommandNameSegment(segment); err != nil {
		return fmt.Errorf("%w: segment %d: %w", ErrInvalidPath, index, err)
	}

	return nil
}

// canonicalPathLength returns the byte length of Path.String for segments
// without allocating the joined string.
func canonicalPathLength(segments []string) int {
	if len(segments) == 0 {
		return 0
	}

	length := 0
	for index, segment := range segments {
		if index > 0 {
			length += len(PathSeparator)
		}

		length += len(segment)
	}

	return length
}

// cloneStringSlice returns a detached copy of values.
//
// A nil input stays nil. A non-nil empty input returns a non-nil empty slice,
// which is acceptable for value-object internals and tests.
func cloneStringSlice(values []string) []string {
	if values == nil {
		return nil
	}

	out := make([]string, len(values))
	copy(out, values)

	return out
}

// clonePaths returns a detached copy of path values and their segment slices.
func clonePaths(paths []Path) []Path {
	if paths == nil {
		return nil
	}

	out := make([]Path, len(paths))
	for index, path := range paths {
		out[index] = path.clone()
	}

	return out
}
