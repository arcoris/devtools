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
	"errors"
	"strings"
	"testing"
)

// TestArtifactDigestAlgorithm verifies algorithm helpers.
func TestArtifactDigestAlgorithm(t *testing.T) {
	t.Parallel()

	if got, want := ArtifactDigestSHA256.HexLength(), 64; got != want {
		t.Fatalf("SHA256 HexLength() = %d, want %d", got, want)
	}

	if got, want := ArtifactDigestSHA512.HexLength(), 128; got != want {
		t.Fatalf("SHA512 HexLength() = %d, want %d", got, want)
	}

	if !ArtifactDigestSHA256.IsKnown() {
		t.Fatalf("SHA256 IsKnown() = false, want true")
	}

	if ArtifactDigestAlgorithm("md5").IsKnown() {
		t.Fatalf("md5 IsKnown() = true, want false")
	}
}

// TestArtifactDigestValidation verifies digest construction and validation.
func TestArtifactDigestValidation(t *testing.T) {
	t.Parallel()

	digest, err := NewArtifactDigest(ArtifactDigestSpec{
		Algorithm: ArtifactDigestSHA256,
		Value:     strings.Repeat("a", 64),
	})
	if err != nil {
		t.Fatalf("NewArtifactDigest() returned unexpected error: %v", err)
	}

	if got, want := digest.String(), "sha256:"+strings.Repeat("a", 64); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	tests := []struct {
		name string
		spec ArtifactDigestSpec
	}{
		{
			name: "empty algorithm",
			spec: ArtifactDigestSpec{
				Value: strings.Repeat("a", 64),
			},
		},
		{
			name: "unknown algorithm",
			spec: ArtifactDigestSpec{
				Algorithm: ArtifactDigestAlgorithm("md5"),
				Value:     strings.Repeat("a", 64),
			},
		},
		{
			name: "wrong length",
			spec: ArtifactDigestSpec{
				Algorithm: ArtifactDigestSHA256,
				Value:     "abc",
			},
		},
		{
			name: "uppercase hex",
			spec: ArtifactDigestSpec{
				Algorithm: ArtifactDigestSHA256,
				Value:     strings.Repeat("A", 64),
			},
		},
		{
			name: "non hex",
			spec: ArtifactDigestSpec{
				Algorithm: ArtifactDigestSHA256,
				Value:     strings.Repeat("g", 64),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewArtifactDigest(tt.spec)
			if err == nil {
				t.Fatalf("NewArtifactDigest() returned nil error")
			}

			if !errors.Is(err, ErrInvalidArtifactDigest) {
				t.Fatalf("NewArtifactDigest() error = %v, want ErrInvalidArtifactDigest", err)
			}
		})
	}
}
