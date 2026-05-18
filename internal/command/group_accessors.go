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

import "strings"

// String returns the canonical string representation of the group.
func (group Group) String() string {
	return string(group)
}

// Key returns the stable map key for group.
func (group Group) Key() string {
	return string(group)
}

// IsZero reports whether the group has not been set.
func (group Group) IsZero() bool {
	return group == ""
}

// IsValid reports whether the group satisfies the group grammar.
func (group Group) IsValid() bool {
	return group.Validate() == nil
}

// Equal reports whether two groups are exactly the same key.
func (group Group) Equal(other Group) bool {
	return group == other
}

// Parts returns a detached copy of the group segments.
//
// Parts performs a lexical split and does not validate the group. Call Validate
// first when the value may be untrusted. The returned slice can be safely
// modified by the caller.
func (group Group) Parts() []string {
	if group == "" {
		return nil
	}

	parts := strings.Split(string(group), GroupSeparator)
	out := make([]string, len(parts))
	copy(out, parts)

	return out
}

// Len returns the number of hierarchical segments in the group.
func (group Group) Len() int {
	return group.Depth()
}

// Depth returns the number of hierarchical segments in the group.
func (group Group) Depth() int {
	if group == "" {
		return 0
	}

	return strings.Count(string(group), GroupSeparator) + 1
}

// At returns the segment at index.
//
// The second return value is false when index is out of range. At never panics.
func (group Group) At(index int) (string, bool) {
	if index < 0 || index >= group.Depth() {
		return "", false
	}

	start := 0
	raw := string(group)
	for current := 0; current < index; current++ {
		next := strings.Index(raw[start:], GroupSeparator)
		if next < 0 {
			return "", false
		}

		start += next + len(GroupSeparator)
	}

	end := strings.Index(raw[start:], GroupSeparator)
	if end < 0 {
		return raw[start:], true
	}

	return raw[start : start+end], true
}

// Leaf returns the final segment of the group.
//
// For "diagnostics.perf", Leaf returns "perf". For the zero group, Leaf
// returns an empty string. Leaf does not validate the group; call Validate first
// for untrusted input.
func (group Group) Leaf() string {
	if group == "" {
		return ""
	}

	raw := string(group)
	index := strings.LastIndex(raw, GroupSeparator)
	if index < 0 {
		return raw
	}

	return raw[index+len(GroupSeparator):]
}
