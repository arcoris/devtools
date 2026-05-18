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

// TestArtifactWithHelpers verifies immutable-style update helpers.
func TestArtifactWithHelpers(t *testing.T) {
	t.Parallel()

	artifact := MustArtifact(ArtifactSpec{
		ID:       "report",
		Kind:     ArtifactKindReport,
		Location: "report.md",
	}).
		MustWithLocation("reports/report.md").
		MustWithFormat(ArtifactFormatMarkdown).
		MustWithDigest(ArtifactDigestSpec{
			Algorithm: ArtifactDigestSHA256,
			Value:     strings.Repeat("a", 64),
		}).
		MustWithLabel("bench").
		MustWithMetadata(MustMetadata(MetadataSpec{Owner: "devtools"})).
		MustWithVisibility(VisibilityHidden)

	if got, want := artifact.Location(), "reports/report.md"; got != want {
		t.Fatalf("Location() = %q, want %q", got, want)
	}

	if got, want := artifact.Format(), ArtifactFormatMarkdown; got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}

	if !artifact.HasDigest() {
		t.Fatalf("HasDigest() = false, want true")
	}

	if !artifact.HasLabel("bench") {
		t.Fatalf("HasLabel(bench) = false, want true")
	}

	if got, want := artifact.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if !artifact.Visibility().IsHidden() {
		t.Fatalf("Visibility().IsHidden() = false, want true")
	}

	withoutDigest := artifact.WithoutDigest()
	if withoutDigest.HasDigest() {
		t.Fatalf("WithoutDigest() still has digest")
	}

	withoutLabel := artifact.WithoutLabel("bench")
	if withoutLabel.HasLabel("bench") {
		t.Fatalf("WithoutLabel() still has label")
	}
}
