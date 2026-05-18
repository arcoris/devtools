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

// cloneArtifactStrings returns a detached string slice.
func cloneArtifactStrings(values []string) []string {
	if values == nil {
		return nil
	}

	out := make([]string, len(values))
	copy(out, values)

	return out
}

// isArtifactLowerHex reports whether r is a lowercase hexadecimal rune.
func isArtifactLowerHex(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')
}
