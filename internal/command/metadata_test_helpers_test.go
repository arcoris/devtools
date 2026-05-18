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

func mustTestMetadata(t *testing.T) Metadata {
	t.Helper()

	return MustMetadata(MetadataSpec{
		Owner: "devtools",
		Area:  "command.registry",
		Since: "v0.1.0",
		Deprecation: &DeprecationSpec{
			Since:       "v0.2.0",
			Message:     "Use bench run instead.",
			Replacement: MustPath("bench", "run"),
		},
		Annotations: map[string]string{
			"docs.page": "commands/bench.md",
			"ci.mode":   "smoke",
		},
	})
}

func mustTestAnnotation(t *testing.T, metadata Metadata, key string) string {
	t.Helper()

	value, ok := metadata.Annotation(key)
	if !ok {
		t.Fatalf("annotation %q not found", key)
	}

	return value
}
