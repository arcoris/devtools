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

// ArtifactID is a stable machine-facing artifact identifier.
//
// Artifact IDs are dot-separated keys. They are suitable for manifests,
// generated indexes, test assertions, diagnostics, and cross-references.
type ArtifactID string

// NewArtifactID validates raw and returns it as an ArtifactID.
func NewArtifactID(raw string) (ArtifactID, error) {
	id := ArtifactID(raw)
	if err := id.Validate(); err != nil {
		return "", err
	}

	return id, nil
}

// MustArtifactID validates raw and returns it as an ArtifactID.
//
// MustArtifactID panics on invalid input. It is intended for static artifact
// declarations and tests.
func MustArtifactID(raw string) ArtifactID {
	id, err := NewArtifactID(raw)
	if err != nil {
		panic(err)
	}

	return id
}

// String returns the canonical artifact ID string.
func (id ArtifactID) String() string {
	return string(id)
}

// IsZero reports whether the ID has not been set.
func (id ArtifactID) IsZero() bool {
	return id == ""
}

// IsValid reports whether the ID satisfies the artifact ID grammar.
func (id ArtifactID) IsValid() bool {
	return id.Validate() == nil
}

// Validate verifies artifact ID structural rules.
func (id ArtifactID) Validate() error {
	raw := string(id)
	if raw == "" {
		return ErrEmptyArtifactID
	}

	if len(raw) > maxArtifactIDLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidArtifactID,
			len(raw),
			maxArtifactIDLength,
		)
	}

	if err := validateArtifactKey(raw, ErrInvalidArtifactID, "artifact id"); err != nil {
		return err
	}

	return nil
}
