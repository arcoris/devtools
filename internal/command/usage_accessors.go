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

import "sort"

// String returns the primary usage syntax.
func (usage Usage) String() string {
	return usage.syntax.String()
}

// Spec returns a detached construction spec for usage.
func (usage Usage) Spec() UsageSpec {
	return usage.spec()
}

// Syntax returns the primary usage line.
func (usage Usage) Syntax() UsageLine {
	return usage.syntax
}

// Alternatives returns a detached copy of alternative usage lines.
func (usage Usage) Alternatives() []UsageLine {
	return cloneUsageLines(usage.alternatives)
}

// Alternative returns the alternative usage line at index.
//
// The second return value is false when index is out of range. Alternative
// never panics.
func (usage Usage) Alternative(index int) (UsageLine, bool) {
	if index < 0 || index >= len(usage.alternatives) {
		return "", false
	}

	return usage.alternatives[index], true
}

// Lines returns all usage lines in declaration order.
//
// The first line is always the primary syntax. Alternatives follow it.
func (usage Usage) Lines() []UsageLine {
	if usage.IsZero() {
		return nil
	}

	lines := make([]UsageLine, 0, 1+len(usage.alternatives))
	lines = append(lines, usage.syntax)
	lines = append(lines, usage.alternatives...)

	return lines
}

// Line returns one usage line by declaration-order index.
//
// Index 0 is the primary syntax. Alternative lines start at index 1.
func (usage Usage) Line(index int) (UsageLine, bool) {
	lines := usage.Lines()
	if index < 0 || index >= len(lines) {
		return "", false
	}

	return lines[index], true
}

// LineStrings returns all usage lines as strings in declaration order.
func (usage Usage) LineStrings() []string {
	lines := usage.Lines()
	out := make([]string, len(lines))

	for index, line := range lines {
		out[index] = line.String()
	}

	return out
}

// SortedLines returns all usage lines sorted lexically.
//
// This is useful for deterministic comparisons and generated indexes. It
// should not be used for help rendering when declaration order matters.
func (usage Usage) SortedLines() []UsageLine {
	lines := usage.Lines()
	sort.Slice(lines, func(i int, j int) bool {
		return lines[i].String() < lines[j].String()
	})

	return lines
}

// Len returns the number of usage lines.
//
// The primary syntax counts as one line.
func (usage Usage) Len() int {
	if usage.IsZero() {
		return 0
	}

	return 1 + len(usage.alternatives)
}

// LineCount returns the number of usage lines.
func (usage Usage) LineCount() int {
	return usage.Len()
}

// AlternativeCount returns the number of alternative syntax lines.
func (usage Usage) AlternativeCount() int {
	return len(usage.alternatives)
}

// HasAlternatives reports whether alternative usage lines are present.
func (usage Usage) HasAlternatives() bool {
	return len(usage.alternatives) > 0
}

// IsZero reports whether usage has not been set.
func (usage Usage) IsZero() bool {
	return usage.syntax == "" && len(usage.alternatives) == 0
}

// IsValid reports whether usage satisfies the usage grammar.
func (usage Usage) IsValid() bool {
	return usage.Validate() == nil
}

// Contains reports whether usage contains line as primary syntax or
// alternative syntax.
//
// The input is normalized before comparison.
func (usage Usage) Contains(line string) bool {
	normalized, err := NewUsageLine(line)
	if err != nil {
		return false
	}

	return usage.ContainsLine(normalized)
}

// ContainsLine reports whether usage contains line as primary syntax or
// alternative syntax.
func (usage Usage) ContainsLine(line UsageLine) bool {
	for _, candidate := range usage.Lines() {
		if candidate == line {
			return true
		}
	}

	return false
}

// spec returns a detached construction spec.
func (usage Usage) spec() UsageSpec {
	alternatives := make([]string, len(usage.alternatives))
	for index, alternative := range usage.alternatives {
		alternatives[index] = alternative.String()
	}

	return UsageSpec{
		Syntax:       usage.syntax.String(),
		Alternatives: alternatives,
	}
}
