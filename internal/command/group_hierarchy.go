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
	"fmt"
	"strings"
)

// Parent returns the parent group and whether such parent exists.
//
// Examples:
//
//   - "diagnostics.perf" returns "diagnostics", true;
//   - "benchmark" returns "", false;
//   - the zero group returns "", false.
//
// Parent does not validate the group; call Validate first for untrusted input.
func (group Group) Parent() (Group, bool) {
	if group == "" {
		return "", false
	}

	raw := string(group)
	index := strings.LastIndex(raw, GroupSeparator)
	if index < 0 {
		return "", false
	}

	return Group(raw[:index]), true
}

// HasParent reports whether the group has a hierarchical parent.
func (group Group) HasParent() bool {
	_, ok := group.Parent()

	return ok
}

// HasPrefix reports whether group is equal to prefix or belongs to prefix's
// hierarchical subtree.
//
// An empty prefix always returns false because group prefixes are expected to be
// explicit.
func (group Group) HasPrefix(prefix Group) bool {
	if group == "" || prefix == "" {
		return false
	}

	if group == prefix {
		return true
	}

	return strings.HasPrefix(string(group), string(prefix)+GroupSeparator)
}

// TrimPrefix removes prefix from group and returns the remaining relative group.
//
// The second return value is false when prefix is empty or not a prefix of
// group. Trimming an identical prefix returns the zero group.
func (group Group) TrimPrefix(prefix Group) (Group, bool) {
	if !group.HasPrefix(prefix) {
		return "", false
	}

	if group == prefix {
		return "", true
	}

	trimmed := strings.TrimPrefix(string(group), string(prefix)+GroupSeparator)

	return Group(trimmed), true
}

// Append returns a child group by appending one validated segment to group.
//
// Append accepts exactly one segment, not a dotted suffix. This prevents
// accidental multi-level hierarchy changes in one call.
func (group Group) Append(segment string) (Group, error) {
	if strings.Contains(segment, GroupSeparator) {
		return "", fmt.Errorf(
			"%w: appended segment must not contain %q",
			ErrInvalidGroup,
			GroupSeparator,
		)
	}

	if err := validateGroupSegment(0, segment); err != nil {
		return "", err
	}

	if group.IsZero() {
		return Group(segment), nil
	}

	if err := group.Validate(); err != nil {
		return "", err
	}

	return Group(string(group) + GroupSeparator + segment), nil
}

// MustAppend returns a child group by appending segment to group.
//
// MustAppend panics on invalid input. It is intended for static command
// definitions and tests.
func (group Group) MustAppend(segment string) Group {
	child, err := group.Append(segment)
	if err != nil {
		panic(err)
	}

	return child
}

// Join returns a group formed by appending another group's segments.
func (group Group) Join(other Group) (Group, error) {
	switch {
	case group.IsZero():
		return NewGroup(other.String())
	case other.IsZero():
		return NewGroup(group.String())
	}

	if err := group.Validate(); err != nil {
		return "", err
	}

	if err := other.Validate(); err != nil {
		return "", err
	}

	return NewGroup(string(group) + GroupSeparator + string(other))
}

// MustJoin returns a group formed by appending another group's segments.
func (group Group) MustJoin(other Group) Group {
	joined, err := group.Join(other)
	if err != nil {
		panic(err)
	}

	return joined
}
