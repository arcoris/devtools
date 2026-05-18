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

// ValidateEnvName validates a conventional ASCII environment variable name.
//
// Accepted names are non-empty, valid UTF-8, ASCII-only, start with an
// uppercase letter or underscore, and continue with uppercase letters, digits,
// or underscores. maxLength is a byte length limit.
func ValidateEnvName(name string, maxLength int) error {
	if name == "" {
		return ErrEmptyEnvName
	}

	if maxLength <= 0 {
		return fmt.Errorf("%w: maximum length must be positive", ErrInvalidEnvName)
	}

	if !utf8.ValidString(name) {
		return fmt.Errorf("%w: value is not valid UTF-8", ErrInvalidEnvName)
	}

	if len(name) > maxLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidEnvName,
			len(name),
			maxLength,
		)
	}

	for offset, r := range name {
		if !IsASCII(r) {
			return fmt.Errorf(
				"%w: contains non-ASCII rune %q at byte offset %d",
				ErrInvalidEnvName,
				r,
				offset,
			)
		}

		ch := byte(r)
		if offset == 0 {
			if IsEnvNameStart(ch) {
				continue
			}

			return fmt.Errorf("%w: must start with uppercase ASCII letter or underscore", ErrInvalidEnvName)
		}

		if IsEnvNameContinuation(ch) {
			continue
		}

		return fmt.Errorf(
			"%w: contains invalid character %q at byte offset %d",
			ErrInvalidEnvName,
			ch,
			offset,
		)
	}

	return nil
}

// IsEnvNameStart reports whether ch can start an environment variable name.
func IsEnvNameStart(ch byte) bool {
	return IsASCIIUpperLetter(ch) || ch == '_'
}

// IsEnvNameContinuation reports whether ch can appear after the first byte in
// an environment variable name.
func IsEnvNameContinuation(ch byte) bool {
	return IsEnvNameStart(ch) || IsASCIIDigit(ch)
}
