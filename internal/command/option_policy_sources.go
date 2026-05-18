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

// normalizeOptionSources removes duplicate sources and returns them in known
// default precedence order.
//
// Unknown sources are preserved after known sources so validation can report a
// precise invalid-source error later.
func normalizeOptionSources(sources []OptionSource) []OptionSource {
	if len(sources) == 0 {
		return nil
	}

	seen := make(map[OptionSource]struct{}, len(sources))
	for _, source := range sources {
		seen[source] = struct{}{}
	}

	out := make([]OptionSource, 0, len(seen))
	for _, source := range KnownOptionSources() {
		if _, ok := seen[source]; ok {
			out = append(out, source)
			delete(seen, source)
		}
	}

	for source := range seen {
		out = append(out, source)
	}

	return out
}

// isZeroOptionPolicy reports whether policy has not been initialized.
func isZeroOptionPolicy(policy OptionPolicy) bool {
	return policy.requirement == "" &&
		policy.scope == "" &&
		policy.occurrence == "" &&
		policy.emptyValue == "" &&
		len(policy.allowedSources) == 0
}
