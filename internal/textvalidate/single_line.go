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

package textvalidate

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// ValidateSingleLineText validates bounded UTF-8 text that must stay on one
// logical line.
//
// Single-line text may contain ordinary spaces and punctuation. It rejects tab,
// newline, carriage return, and other C0 control runes.
func ValidateSingleLineText(text string, maxLength int) error {
	if maxLength < 0 {
		return fmt.Errorf("%w: maximum length must not be negative", ErrInvalidSingleLineText)
	}

	if !utf8.ValidString(text) {
		return fmt.Errorf("%w: value is not valid UTF-8", ErrInvalidSingleLineText)
	}

	if len(text) > maxLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidSingleLineText,
			len(text),
			maxLength,
		)
	}

	for offset, r := range text {
		if IsDisallowedSingleLineTextRune(r) {
			return fmt.Errorf(
				"%w: contains disallowed rune U+%04X at byte offset %d",
				ErrInvalidSingleLineText,
				r,
				offset,
			)
		}
	}

	return nil
}

// ValidateTokenText validates bounded UTF-8 text that must be one token.
//
// Token text rejects every Unicode whitespace rune and all disallowed
// single-line runes.
func ValidateTokenText(text string, maxLength int) error {
	if text == "" {
		return ErrEmptyTokenText
	}

	if maxLength <= 0 {
		return fmt.Errorf("%w: maximum length must be positive", ErrInvalidTokenText)
	}

	if !utf8.ValidString(text) {
		return fmt.Errorf("%w: value is not valid UTF-8", ErrInvalidTokenText)
	}

	if len(text) > maxLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidTokenText,
			len(text),
			maxLength,
		)
	}

	for offset, r := range text {
		if IsDisallowedSingleLineTextRune(r) {
			return fmt.Errorf(
				"%w: contains disallowed rune U+%04X at byte offset %d",
				ErrInvalidTokenText,
				r,
				offset,
			)
		}

		if unicode.IsSpace(r) {
			return fmt.Errorf(
				"%w: contains whitespace rune U+%04X at byte offset %d",
				ErrInvalidTokenText,
				r,
				offset,
			)
		}
	}

	return nil
}

// IsDisallowedSingleLineTextRune reports whether r must be rejected from
// single-line text.
func IsDisallowedSingleLineTextRune(r rune) bool {
	return r >= 0 && r < 0x20
}
