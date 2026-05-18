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
	"fmt"
	"strings"
)

// ArtifactDigestSpec describes artifact digest metadata before validation.
type ArtifactDigestSpec struct {
	// Algorithm is the digest algorithm.
	Algorithm ArtifactDigestAlgorithm

	// Value is the lowercase hexadecimal digest value.
	Value string
}

// ArtifactDigest is validated artifact digest metadata.
type ArtifactDigest struct {
	algorithm ArtifactDigestAlgorithm
	value     string
}

// NewArtifactDigest validates spec and returns ArtifactDigest.
func NewArtifactDigest(spec ArtifactDigestSpec) (ArtifactDigest, error) {
	digest := ArtifactDigest{
		algorithm: spec.Algorithm,
		value:     strings.TrimSpace(spec.Value),
	}

	if err := digest.Validate(); err != nil {
		return ArtifactDigest{}, err
	}

	return digest, nil
}

// MustArtifactDigest validates spec and returns ArtifactDigest.
//
// MustArtifactDigest panics on invalid input.
func MustArtifactDigest(spec ArtifactDigestSpec) ArtifactDigest {
	digest, err := NewArtifactDigest(spec)
	if err != nil {
		panic(err)
	}

	return digest
}

// Algorithm returns the digest algorithm.
func (digest ArtifactDigest) Algorithm() ArtifactDigestAlgorithm {
	return digest.algorithm
}

// Value returns the lowercase hexadecimal digest value.
func (digest ArtifactDigest) Value() string {
	return digest.value
}

// String returns the canonical digest string as "algorithm:value".
func (digest ArtifactDigest) String() string {
	if digest.algorithm == "" || digest.value == "" {
		return ""
	}

	return digest.algorithm.String() + ":" + digest.value
}

// Validate verifies digest structural rules.
func (digest ArtifactDigest) Validate() error {
	if err := digest.algorithm.Validate(); err != nil {
		return err
	}

	expectedLength := digest.algorithm.HexLength()
	if len(digest.value) != expectedLength {
		return fmt.Errorf(
			"%w: %s digest length %d, want %d",
			ErrInvalidArtifactDigest,
			digest.algorithm,
			len(digest.value),
			expectedLength,
		)
	}

	for offset, ch := range digest.value {
		if !isArtifactLowerHex(ch) {
			return fmt.Errorf(
				"%w: digest contains non-lowercase-hex rune %q at byte offset %d",
				ErrInvalidArtifactDigest,
				ch,
				offset,
			)
		}
	}

	return nil
}

// spec returns a detached construction spec.
func (digest ArtifactDigest) spec() ArtifactDigestSpec {
	return ArtifactDigestSpec{
		Algorithm: digest.algorithm,
		Value:     digest.value,
	}
}

// ArtifactDigestAlgorithm is a supported digest algorithm.
type ArtifactDigestAlgorithm string

const (
	ArtifactDigestSHA256 ArtifactDigestAlgorithm = "sha256"
	ArtifactDigestSHA512 ArtifactDigestAlgorithm = "sha512"
)

// String returns the canonical algorithm string.
func (algorithm ArtifactDigestAlgorithm) String() string {
	return string(algorithm)
}

// IsKnown reports whether algorithm is supported.
func (algorithm ArtifactDigestAlgorithm) IsKnown() bool {
	switch algorithm {
	case ArtifactDigestSHA256, ArtifactDigestSHA512:
		return true
	default:
		return false
	}
}

// Validate verifies digest algorithm structural rules.
func (algorithm ArtifactDigestAlgorithm) Validate() error {
	if algorithm == "" {
		return fmt.Errorf("%w: algorithm is empty", ErrInvalidArtifactDigest)
	}

	if algorithm.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported algorithm %q", ErrInvalidArtifactDigest, algorithm)
}

// HexLength returns expected lowercase hexadecimal digest length.
func (algorithm ArtifactDigestAlgorithm) HexLength() int {
	switch algorithm {
	case ArtifactDigestSHA256:
		return 64
	case ArtifactDigestSHA512:
		return 128
	default:
		return 0
	}
}
