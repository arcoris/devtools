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

import "testing"

func actionTestCommandNode() Node {
	return MustCommandNode(
		MustID("bench.run"),
		MustPath("bench", "run"),
		"run",
	)
}

func actionTestAlternateCommandNode() Node {
	return MustCommandNode(
		MustID("check"),
		MustPath("check"),
		"check",
	)
}

func actionTestFamilyNode() Node {
	return MustFamilyNode(
		MustID("bench"),
		MustPath("bench"),
		"bench",
		actionTestCommandNode(),
	)
}

func actionTestRequest(t *testing.T) ActionRequest {
	t.Helper()

	return MustActionRequest(actionTestCommandNode(), nil, nil)
}

func actionTestRequestField(t *testing.T, request ActionRequest, key string) string {
	t.Helper()

	value, ok := request.Field(key)
	if !ok {
		t.Fatalf("request field %q not found", key)
	}

	return value
}

func actionTestResultField(t *testing.T, result ActionResult, key string) string {
	t.Helper()

	value, ok := result.Field(key)
	if !ok {
		t.Fatalf("result field %q not found", key)
	}

	return value
}
