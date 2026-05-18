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

// TestMustArtifactPanicsForInvalidArtifact verifies fail-fast construction.
func TestMustArtifactPanicsForInvalidArtifact(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustArtifact did not panic")
		}
	}()

	_ = MustArtifact(ArtifactSpec{})
}

// TestNewArtifactAcceptsValidArtifact verifies full artifact construction.
func TestNewArtifactAcceptsValidArtifact(t *testing.T) {
	t.Parallel()

	size := int64(128)
	artifact, err := NewArtifact(ArtifactSpec{
		ID:          "bench.report",
		Kind:        ArtifactKindReport,
		Location:    "bench/reports/report.md",
		Format:      ArtifactFormatMarkdown,
		MediaType:   "text/markdown",
		Description: "Benchmark report.",
		SizeBytes:   &size,
		Digest: &ArtifactDigestSpec{
			Algorithm: ArtifactDigestSHA256,
			Value:     strings.Repeat("a", 64),
		},
		Labels: []string{"bench", "report.markdown"},
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewArtifact() returned unexpected error: %v", err)
	}

	if got, want := artifact.ID(), MustArtifactID("bench.report"); got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}

	if got, want := artifact.Kind(), ArtifactKindReport; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := artifact.Location(), "bench/reports/report.md"; got != want {
		t.Fatalf("Location() = %q, want %q", got, want)
	}

	if got, want := artifact.Format(), ArtifactFormatMarkdown; got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}

	if got, ok := artifact.SizeBytes(); !ok || got != 128 {
		t.Fatalf("SizeBytes() = %d, %v; want 128, true", got, ok)
	}

	if !artifact.HasDigest() {
		t.Fatalf("HasDigest() = false, want true")
	}

	if !artifact.HasLabel("report.markdown") {
		t.Fatalf("HasLabel(report.markdown) = false, want true")
	}

	if !artifact.IsRepositoryRelative() {
		t.Fatalf("IsRepositoryRelative() = false, want true")
	}
}

// TestNewArtifactDefaults verifies visibility and optional fields.
func TestNewArtifactDefaults(t *testing.T) {
	t.Parallel()

	artifact := MustArtifact(ArtifactSpec{
		ID:       "bench.raw",
		Kind:     ArtifactKindRaw,
		Location: "bench/raw/out.txt",
	})

	if got, want := artifact.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if artifact.HasFormat() {
		t.Fatalf("HasFormat() = true, want false")
	}

	if artifact.HasSize() {
		t.Fatalf("HasSize() = true, want false")
	}

	if artifact.HasDigest() {
		t.Fatalf("HasDigest() = true, want false")
	}
}

// TestNewArtifactNormalizesText verifies normalization of compact text fields.
func TestNewArtifactNormalizesText(t *testing.T) {
	t.Parallel()

	artifact := MustArtifact(ArtifactSpec{
		ID:          "bench.report",
		Kind:        ArtifactKindReport,
		Location:    "  bench/reports/report.md  ",
		MediaType:   "  text/markdown  ",
		Description: "  First line.  \n  Second line.  ",
	})

	if got, want := artifact.Location(), "bench/reports/report.md"; got != want {
		t.Fatalf("Location() = %q, want %q", got, want)
	}

	if got, want := artifact.MediaType(), "text/markdown"; got != want {
		t.Fatalf("MediaType() = %q, want %q", got, want)
	}

	if got, want := artifact.Description(), "First line.\nSecond line."; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
}
