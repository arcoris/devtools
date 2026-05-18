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

// normalizeDocumentationSingleLine returns canonical one-line documentation
// text.
func normalizeDocumentationSingleLine(raw string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
}

// normalizeDocumentationBlock returns canonical block documentation text.
//
// It trims leading and trailing whitespace from the whole block and from each
// line, removes leading and trailing empty lines, and preserves paragraph
// boundaries.
func normalizeDocumentationBlock(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return ""
	}

	lines := strings.Split(raw, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimSpace(line)
	}

	return strings.Join(lines, "\n")
}

// normalizeDocumentationNotes normalizes notes and drops empty note entries.
func normalizeDocumentationNotes(notes []string) []string {
	if notes == nil {
		return nil
	}

	out := make([]string, 0, len(notes))
	for _, note := range notes {
		normalized := normalizeDocumentationBlock(note)
		if normalized == "" {
			continue
		}

		out = append(out, normalized)
	}

	return out
}

// normalizeDocumentationReferenceKey normalizes reference keys.
func normalizeDocumentationReferenceKey(raw string) string {
	return strings.TrimSpace(raw)
}
