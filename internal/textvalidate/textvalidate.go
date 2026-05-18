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

// Package textvalidate provides reusable validation helpers for compact
// machine-facing identifiers and bounded text values.
package textvalidate

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	// DottedKebabKeySeparator separates hierarchical machine-facing key segments.
	DottedKebabKeySeparator = "."
)

var (
	// ErrEmptyKebabSegment reports that a kebab-case segment is empty.
	ErrEmptyKebabSegment = errors.New("kebab segment is empty")

	// ErrInvalidKebabSegment reports that a kebab-case segment is malformed.
	ErrInvalidKebabSegment = errors.New("kebab segment is invalid")

	// ErrEmptyDottedKebabKey reports that a dotted kebab key is empty.
	ErrEmptyDottedKebabKey = errors.New("dotted kebab key is empty")

	// ErrInvalidDottedKebabKey reports that a dotted kebab key is malformed.
	ErrInvalidDottedKebabKey = errors.New("dotted kebab key is invalid")

	// ErrInvalidCompactText reports that a compact text value is malformed.
	ErrInvalidCompactText = errors.New("compact text is invalid")
)

// ValidateKebabSegment validates one ASCII kebab-case identifier segment.
//
// The accepted grammar is intentionally strict because these values are used
// as stable machine-facing keys, not display text:
//
//   - the segment must be non-empty;
//   - the segment must be valid UTF-8;
//   - the segment must start with a lowercase ASCII letter;
//   - the segment may contain lowercase ASCII letters, digits, or hyphens;
//   - hyphens must separate alphanumeric runs, so they cannot be trailing or
//     repeated;
//   - uppercase letters, underscores, spaces, slashes, dots, and non-ASCII
//     characters are rejected.
func ValidateKebabSegment(segment string) error {
	if segment == "" {
		return ErrEmptyKebabSegment
	}

	if !utf8.ValidString(segment) {
		return fmt.Errorf("%w: value is not valid UTF-8", ErrInvalidKebabSegment)
	}

	var previous byte
	for offset, r := range segment {
		if err := ValidateKebabSegmentRune(offset, r); err != nil {
			return err
		}

		ch := byte(r)
		if ch == '-' {
			if offset == len(segment)-1 {
				return fmt.Errorf("%w: must not end with hyphen", ErrInvalidKebabSegment)
			}

			if previous == '-' {
				return fmt.Errorf(
					"%w: contains repeated hyphen at byte offset %d",
					ErrInvalidKebabSegment,
					offset,
				)
			}
		}

		previous = ch
	}

	return nil
}

// ValidateKebabSegmentRune validates one rune at byteOffset inside an ASCII
// kebab-case segment.
//
// The offset is a byte offset because range over strings reports byte offsets.
// This helper only validates local character rules. Whole-segment checks such
// as repeated or trailing hyphens belong in ValidateKebabSegment.
func ValidateKebabSegmentRune(byteOffset int, r rune) error {
	if byteOffset < 0 {
		return fmt.Errorf("%w: byte offset %d is negative", ErrInvalidKebabSegment, byteOffset)
	}

	if !IsASCII(r) {
		return fmt.Errorf(
			"%w: contains non-ASCII rune %q at byte offset %d",
			ErrInvalidKebabSegment,
			r,
			byteOffset,
		)
	}

	ch := byte(r)
	if byteOffset == 0 {
		return ValidateKebabSegmentStart(ch)
	}

	return ValidateKebabSegmentContinuation(byteOffset, ch)
}

// ValidateKebabSegmentStart validates the first byte of a kebab-case segment.
func ValidateKebabSegmentStart(ch byte) error {
	if IsASCIILowerLetter(ch) {
		return nil
	}

	return fmt.Errorf("%w: must start with a lowercase ASCII letter", ErrInvalidKebabSegment)
}

// ValidateKebabSegmentContinuation validates a non-first byte of a kebab-case
// segment.
func ValidateKebabSegmentContinuation(byteOffset int, ch byte) error {
	if byteOffset <= 0 {
		return fmt.Errorf("%w: continuation byte offset must be positive", ErrInvalidKebabSegment)
	}

	if IsKebabSegmentContinuation(ch) {
		return nil
	}

	return fmt.Errorf(
		"%w: contains invalid character %q at byte offset %d",
		ErrInvalidKebabSegment,
		ch,
		byteOffset,
	)
}

// IsKebabSegmentContinuation reports whether ch is allowed after the first
// character of a kebab-case segment.
func IsKebabSegmentContinuation(ch byte) bool {
	return IsASCIILowerLetter(ch) || IsASCIIDigit(ch) || ch == '-'
}

// IsASCII reports whether r is an ASCII code point.
func IsASCII(r rune) bool {
	return r >= 0 && r <= 127
}

// IsASCIILowerLetter reports whether ch is an ASCII lowercase letter.
func IsASCIILowerLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

// IsASCIIDigit reports whether ch is an ASCII digit.
func IsASCIIDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
