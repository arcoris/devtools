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

import "errors"

var (
	// ErrInvalidRegistry reports that a command registry definition is malformed.
	ErrInvalidRegistry = errors.New("command registry is invalid")
)

// Registry is an indexed, validated command-tree registry.
//
// Registry is the command model's lookup-oriented source of truth for CLI
// adapters, documentation generators, policy layers, diagnostics, tests, and
// command-discovery features.
//
// Tree owns command structure. Registry owns deterministic indexes built from
// that tree:
//
//   - by stable ID, including the zero ID for the root node;
//   - by canonical Path key;
//   - by pre-order ID and Path lists for stable iteration.
//
// Registry is framework-neutral. It must not depend on Cobra or any other CLI
// adapter.
type Registry struct {
	tree   Tree
	byID   map[ID]Node
	byPath map[string]Node
	ids    []ID
	paths  []Path
}
