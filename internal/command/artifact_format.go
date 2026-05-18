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

// ArtifactFormat describes an optional machine-facing artifact format.
//
// ArtifactFormat is open but validated. Projects may define additional formats.
type ArtifactFormat string

const (
	ArtifactFormatText     ArtifactFormat = "text"
	ArtifactFormatJSON     ArtifactFormat = "json"
	ArtifactFormatMarkdown ArtifactFormat = "markdown"
	ArtifactFormatHTML     ArtifactFormat = "html"
	ArtifactFormatCSV      ArtifactFormat = "csv"
	ArtifactFormatSVG      ArtifactFormat = "svg"
	ArtifactFormatPPROF    ArtifactFormat = "pprof"
	ArtifactFormatTrace    ArtifactFormat = "trace"
	ArtifactFormatBinary   ArtifactFormat = "binary"
)

// NewArtifactFormat validates raw and returns it as an ArtifactFormat.
func NewArtifactFormat(raw string) (ArtifactFormat, error) {
	format := ArtifactFormat(raw)
	if err := format.Validate(); err != nil {
		return "", err
	}

	return format, nil
}

// MustArtifactFormat validates raw and returns it as an ArtifactFormat.
//
// MustArtifactFormat panics on invalid input.
func MustArtifactFormat(raw string) ArtifactFormat {
	format, err := NewArtifactFormat(raw)
	if err != nil {
		panic(err)
	}

	return format
}

// String returns the canonical format string.
func (format ArtifactFormat) String() string {
	return string(format)
}

// IsZero reports whether format has not been set.
func (format ArtifactFormat) IsZero() bool {
	return format == ""
}

// Validate verifies artifact format structural rules.
func (format ArtifactFormat) Validate() error {
	raw := string(format)
	if raw == "" {
		return nil
	}

	if len(raw) > maxArtifactFormatLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidArtifactFormat,
			len(raw),
			maxArtifactFormatLength,
		)
	}

	if err := validateCommandNameSegment(raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArtifactFormat, err)
	}

	return nil
}
