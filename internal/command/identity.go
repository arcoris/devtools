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

// Package command defines framework-neutral command metadata, identity,
// validation, and execution primitives for ARCORIS developer tooling.
package command

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// IDSeparator separates hierarchical command identity segments.
	//
	// The separator is part of the stable command identity grammar. It is not
	// derived from any specific CLI adapter. For example, Cobra may expose the
	// command path as "bench run", while the stable command ID remains
	// "bench.run".
	IDSeparator = "."

	// maxIDLength is the maximum allowed byte length of a complete command ID.
	//
	// Command IDs are expected to be compact, stable, human-readable keys used
	// in reports, registries, telemetry, documentation anchors, and policy
	// lookups. They are not free-form descriptions.
	maxIDLength = 255

	// maxIDSegmentLength is the maximum allowed byte length of one ID segment.
	//
	// Keeping individual segments short prevents unreadable report keys and
	// excessively long generated anchors.
	maxIDSegmentLength = 63
)

var (
	// ErrEmptyID reports that a command ID was not provided.
	ErrEmptyID = errors.New("command id is empty")

	// ErrInvalidID reports that a command ID violates the command ID grammar.
	ErrInvalidID = errors.New("command id is invalid")
)

// ID is a stable machine-facing identifier for a command or command family.
//
// ID is deliberately different from a CLI path. A CLI path is presentation
// state owned by an adapter such as Cobra and may include aliases, positional
// syntax, or adapter-specific formatting. ID is stable internal state used by
// registries, reports, generated documentation, policy lookup, audit records,
// telemetry, and tests.
//
// ID does not define the catalog of known commands. Known IDs MUST be declared
// by command-definition packages. This type only defines the value object and
// its structural grammar.
//
// The ID grammar is strict:
//
//   - an ID is one or more dot-separated segments;
//   - each segment uses the command-name segment grammar;
//   - separators MUST NOT appear at the beginning, at the end, or
//     consecutively;
//   - a complete ID MUST NOT exceed maxIDLength bytes;
//   - one segment MUST NOT exceed maxIDSegmentLength bytes.
//
// Valid shape examples:
//
//   - "check"
//   - "bench.run"
//   - "profile.cpu"
//   - "config.validate"
//   - "release-notes.generate"
type ID string

// NewID validates raw and returns it as an ID.
//
// Use NewID for IDs loaded from configuration, generated metadata, tests, or
// any other source where validation errors must be returned instead of causing
// a panic.
func NewID(raw string) (ID, error) {
	id := ID(raw)
	if err := id.Validate(); err != nil {
		return "", err
	}

	return id, nil
}

// ParseID is an alias for NewID.
//
// ParseID is useful when the call site is explicitly parsing an external string
// representation into an ID value.
func ParseID(raw string) (ID, error) {
	return NewID(raw)
}

// MustID validates raw and returns it as an ID.
//
// MustID panics on invalid input. It is intended for static command-definition
// declarations, where invalid IDs are programmer errors and should fail during
// tests or package initialization.
func MustID(raw string) ID {
	id, err := NewID(raw)
	if err != nil {
		panic(err)
	}

	return id
}

// String returns the canonical string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// IsZero reports whether the ID has not been set.
func (id ID) IsZero() bool {
	return id == ""
}

// IsValid reports whether the ID satisfies the command ID grammar.
func (id ID) IsValid() bool {
	return id.Validate() == nil
}

// Validate verifies that the ID satisfies the command ID grammar.
func (id ID) Validate() error {
	raw := string(id)

	if raw == "" {
		return ErrEmptyID
	}

	if len(raw) > maxIDLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidID,
			len(raw),
			maxIDLength,
		)
	}

	if strings.HasPrefix(raw, IDSeparator) {
		return fmt.Errorf("%w: must not start with %q", ErrInvalidID, IDSeparator)
	}

	if strings.HasSuffix(raw, IDSeparator) {
		return fmt.Errorf("%w: must not end with %q", ErrInvalidID, IDSeparator)
	}

	segments := strings.Split(raw, IDSeparator)
	for index, segment := range segments {
		if err := validateIDSegment(index, segment); err != nil {
			return err
		}
	}

	return nil
}

// Parts returns a detached copy of the ID segments.
//
// Parts performs a lexical split and does not validate the ID. Call Validate
// first when the value may be untrusted. The returned slice is detached from
// internal state and can be safely modified by the caller.
func (id ID) Parts() []string {
	if id == "" {
		return nil
	}

	parts := strings.Split(string(id), IDSeparator)
	out := make([]string, len(parts))
	copy(out, parts)

	return out
}

// Depth returns the number of hierarchical segments in the ID.
//
// Examples:
//
//   - the zero ID has depth 0;
//   - "check" has depth 1;
//   - "bench.run" has depth 2.
func (id ID) Depth() int {
	if id == "" {
		return 0
	}

	return len(strings.Split(string(id), IDSeparator))
}

// Leaf returns the last segment of the ID.
//
// For "bench.run", Leaf returns "run". For the zero ID, Leaf returns an empty
// string. Leaf does not validate the ID; call Validate first for untrusted
// input.
func (id ID) Leaf() string {
	if id == "" {
		return ""
	}

	raw := string(id)
	index := strings.LastIndex(raw, IDSeparator)
	if index < 0 {
		return raw
	}

	return raw[index+len(IDSeparator):]
}

// Parent returns the parent ID and whether such parent exists.
//
// Examples:
//
//   - "bench.run" returns "bench", true;
//   - "bench" returns "", false;
//   - the zero ID returns "", false.
//
// Parent does not validate the ID; call Validate first for untrusted input.
func (id ID) Parent() (ID, bool) {
	if id == "" {
		return "", false
	}

	raw := string(id)
	index := strings.LastIndex(raw, IDSeparator)
	if index < 0 {
		return "", false
	}

	return ID(raw[:index]), true
}

// HasParent reports whether the ID has a hierarchical parent.
func (id ID) HasParent() bool {
	_, ok := id.Parent()
	return ok
}

// HasPrefix reports whether id is equal to prefix or belongs to prefix's
// hierarchical subtree.
//
// Examples:
//
//   - "bench.run".HasPrefix("bench") is true;
//   - "bench.run".HasPrefix("bench.run") is true;
//   - "benchmark.run".HasPrefix("bench") is false.
//
// An empty prefix always returns false because command ID prefixes are expected
// to be explicit.
func (id ID) HasPrefix(prefix ID) bool {
	if id == "" || prefix == "" {
		return false
	}

	if id == prefix {
		return true
	}

	return strings.HasPrefix(string(id), string(prefix)+IDSeparator)
}

// Append returns a child ID by appending one validated segment to id.
//
// Append accepts exactly one segment, not a dotted suffix. This prevents
// accidental multi-level hierarchy changes in a single call.
//
// Examples:
//
//   - ID("bench").Append("run") returns "bench.run";
//   - ID("").Append("check") returns "check";
//   - ID("bench").Append("run.fast") returns an error.
func (id ID) Append(segment string) (ID, error) {
	if strings.Contains(segment, IDSeparator) {
		return "", fmt.Errorf(
			"%w: appended segment must not contain %q",
			ErrInvalidID,
			IDSeparator,
		)
	}

	if err := validateIDSegment(0, segment); err != nil {
		return "", err
	}

	if id.IsZero() {
		return ID(segment), nil
	}

	if err := id.Validate(); err != nil {
		return "", err
	}

	child := ID(string(id) + IDSeparator + segment)
	if err := child.Validate(); err != nil {
		return "", err
	}

	return child, nil
}

// MustAppend returns a child ID by appending segment to id.
//
// MustAppend panics on invalid input. It is intended for static
// command-definition declarations.
func (id ID) MustAppend(segment string) ID {
	child, err := id.Append(segment)
	if err != nil {
		panic(err)
	}

	return child
}

// validateIDSegment validates one ID segment and wraps command-name validation
// errors with ID-specific diagnostics.
func validateIDSegment(index int, segment string) error {
	if len(segment) > maxIDSegmentLength {
		return fmt.Errorf(
			"%w: segment %d length %d exceeds maximum length %d",
			ErrInvalidID,
			index,
			len(segment),
			maxIDSegmentLength,
		)
	}

	if err := validateCommandNameSegment(segment); err != nil {
		return fmt.Errorf("%w: segment %d: %w", ErrInvalidID, index, err)
	}

	return nil
}
