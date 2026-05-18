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
	"unicode/utf8"
)

// ValidateCompactText validates bounded UTF-8 text without disallowed control
// runes.
//
// Compact text may contain spaces and punctuation. It may also contain tab,
// newline, and carriage return. Other C0 control runes are rejected.
func ValidateCompactText(text string, maxLength int) error {
	if maxLength < 0 {
		return fmt.Errorf("%w: maximum length must not be negative", ErrInvalidCompactText)
	}

	if !utf8.ValidString(text) {
		return fmt.Errorf("%w: value is not valid UTF-8", ErrInvalidCompactText)
	}

	if len(text) > maxLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidCompactText,
			len(text),
			maxLength,
		)
	}

	for offset, r := range text {
		if IsDisallowedCompactTextControlRune(r) {
			return fmt.Errorf(
				"%w: contains disallowed control rune U+%04X at byte offset %d",
				ErrInvalidCompactText,
				r,
				offset,
			)
		}
	}

	return nil
}

// IsDisallowedCompactTextControlRune reports whether r should be rejected from
// compact text values.
func IsDisallowedCompactTextControlRune(r rune) bool {
	if r == '\t' || r == '\n' || r == '\r' {
		return false
	}

	return r >= 0 && r < 0x20
}
