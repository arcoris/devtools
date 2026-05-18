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

// Resolve resolves adapter-facing command segments by walking child Use and
// alias segments.
//
// Resolve is useful for CLI adapters. It accepts aliases. Use FindByPath or
// ResolvePath when canonical-only lookup is required.
func (registry Registry) Resolve(segments ...string) (Node, bool) {
	return registry.tree.Resolve(segments...)
}

// ResolvePath resolves a canonical Path by child Use segments only.
func (registry Registry) ResolvePath(path Path) (Node, bool) {
	return registry.tree.ResolvePath(path)
}
