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

// cloneOptionNames returns a detached copy of option names.
func cloneOptionNames(values []OptionName) []OptionName {
	if values == nil {
		return nil
	}

	out := make([]OptionName, len(values))
	copy(out, values)

	return out
}

// cloneOptionKinds returns a detached copy of option kinds.
func cloneOptionKinds(values []OptionKind) []OptionKind {
	if values == nil {
		return nil
	}

	out := make([]OptionKind, len(values))
	copy(out, values)

	return out
}

// cloneOptionSources returns a detached copy of option sources.
func cloneOptionSources(values []OptionSource) []OptionSource {
	if values == nil {
		return nil
	}

	out := make([]OptionSource, len(values))
	copy(out, values)

	return out
}

// containsString reports whether values contains target.
func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
