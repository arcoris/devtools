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

// normalizeResultHints normalizes hints and drops blank entries.
func normalizeResultHints(hints []string) []string {
	if hints == nil {
		return nil
	}

	out := make([]string, 0, len(hints))
	for _, hint := range hints {
		normalized := normalizeResultBlock(hint)
		if normalized == "" {
			continue
		}

		out = append(out, normalized)
	}

	return out
}

// cloneResultArtifacts returns a detached artifact slice.
func cloneResultArtifacts(values []Artifact) []Artifact {
	if values == nil {
		return nil
	}

	out := make([]Artifact, len(values))
	copy(out, values)

	return out
}

// cloneResultWarnings returns a detached warning slice.
func cloneResultWarnings(values []ResultWarning) []ResultWarning {
	if values == nil {
		return nil
	}

	out := make([]ResultWarning, len(values))
	for index, warning := range values {
		out[index] = ResultWarning{
			kind:    warning.kind,
			message: warning.message,
			hints:   cloneResultStrings(warning.hints),
		}
	}

	return out
}

// cloneResultStringMap returns a detached string map.
func cloneResultStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}

	out := make(map[string]string, len(values))
	for key, value := range values {
		out[key] = value
	}

	return out
}

// cloneResultStrings returns a detached string slice.
func cloneResultStrings(values []string) []string {
	if values == nil {
		return nil
	}

	out := make([]string, len(values))
	copy(out, values)

	return out
}
