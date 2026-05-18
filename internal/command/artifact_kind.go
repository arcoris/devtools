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

import "fmt"

// ArtifactKind describes an artifact category.
//
// ArtifactKind is open but validated. Projects may define additional artifact
// kinds while preserving a compact machine-facing key grammar.
type ArtifactKind string

const (
	ArtifactKindReport     ArtifactKind = "report"
	ArtifactKindRaw        ArtifactKind = "raw"
	ArtifactKindProfile    ArtifactKind = "profile"
	ArtifactKindTrace      ArtifactKind = "trace"
	ArtifactKindCoverage   ArtifactKind = "coverage"
	ArtifactKindBenchmark  ArtifactKind = "benchmark"
	ArtifactKindLog        ArtifactKind = "log"
	ArtifactKindSummary    ArtifactKind = "summary"
	ArtifactKindSnapshot   ArtifactKind = "snapshot"
	ArtifactKindDiagnostic ArtifactKind = "diagnostic"
	ArtifactKindMetadata   ArtifactKind = "metadata"
)

// NewArtifactKind validates raw and returns it as an ArtifactKind.
func NewArtifactKind(raw string) (ArtifactKind, error) {
	kind := ArtifactKind(raw)
	if err := kind.Validate(); err != nil {
		return "", err
	}

	return kind, nil
}

// MustArtifactKind validates raw and returns it as an ArtifactKind.
//
// MustArtifactKind panics on invalid input.
func MustArtifactKind(raw string) ArtifactKind {
	kind, err := NewArtifactKind(raw)
	if err != nil {
		panic(err)
	}

	return kind
}

// String returns the canonical kind string.
func (kind ArtifactKind) String() string {
	return string(kind)
}

// IsZero reports whether kind has not been set.
func (kind ArtifactKind) IsZero() bool {
	return kind == ""
}

// Validate verifies artifact kind structural rules.
func (kind ArtifactKind) Validate() error {
	raw := string(kind)
	if raw == "" {
		return ErrEmptyArtifactKind
	}

	if len(raw) > maxArtifactKindLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidArtifactKind,
			len(raw),
			maxArtifactKindLength,
		)
	}

	if err := validateCommandNameSegment(raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArtifactKind, err)
	}

	return nil
}
