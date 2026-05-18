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
	"strings"
	"testing"
)

func mustTestPathFromID(t *testing.T, raw string) Path {
	t.Helper()

	path := MustParsePath(strings.ReplaceAll(raw, IDSeparator, PathSeparator))
	if path.IsRoot() {
		t.Fatalf("test path for %q is root", raw)
	}

	return path
}

func mustTestCommandNode(t *testing.T, raw string, aliases ...string) Node {
	t.Helper()

	path := mustTestPathFromID(t, raw)

	return MustNode(NodeSpec{
		Kind:    NodeCommand,
		ID:      MustID(raw),
		Path:    path,
		Use:     path.Leaf(),
		Aliases: aliases,
	})
}

func mustTestFamilyNode(t *testing.T, raw string, children ...Node) Node {
	t.Helper()

	path := mustTestPathFromID(t, raw)

	return MustFamilyNode(MustID(raw), path, path.Leaf(), children...)
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("function did not panic")
		}
	}()

	fn()
}
