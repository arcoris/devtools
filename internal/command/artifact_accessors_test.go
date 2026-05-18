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

// TestArtifactCopySemantics verifies detached labels.
func TestArtifactCopySemantics(t *testing.T) {
	t.Parallel()

	labels := []string{"bench"}

	artifact := MustArtifact(ArtifactSpec{
		ID:       "report",
		Kind:     ArtifactKindReport,
		Location: "report.md",
		Labels:   labels,
	})

	labels[0] = "changed"

	if artifact.HasLabel("changed") {
		t.Fatalf("artifact changed through input label slice")
	}

	out := artifact.Labels()
	out[0] = "changed"

	if artifact.HasLabel("changed") {
		t.Fatalf("artifact changed through output label slice")
	}
}
