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
	"strings"
)

// ValidateDottedKebabKey validates a dot-separated stable machine-facing key.
//
// A valid key has one or more kebab-case segments separated by ".". maxLength
// is a byte length limit for the whole key. Passing maxDepth <= 0 disables the
// depth limit.
func ValidateDottedKebabKey(key string, maxLength int, maxDepth int) error {
	if key == "" {
		return ErrEmptyDottedKebabKey
	}

	if maxLength <= 0 {
		return fmt.Errorf("%w: maximum length must be positive", ErrInvalidDottedKebabKey)
	}

	if len(key) > maxLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidDottedKebabKey,
			len(key),
			maxLength,
		)
	}

	if strings.HasPrefix(key, DottedKebabKeySeparator) {
		return fmt.Errorf("%w: must not start with %q", ErrInvalidDottedKebabKey, DottedKebabKeySeparator)
	}

	if strings.HasSuffix(key, DottedKebabKeySeparator) {
		return fmt.Errorf("%w: must not end with %q", ErrInvalidDottedKebabKey, DottedKebabKeySeparator)
	}

	segments := strings.Split(key, DottedKebabKeySeparator)
	if maxDepth > 0 && len(segments) > maxDepth {
		return fmt.Errorf(
			"%w: depth %d exceeds maximum depth %d",
			ErrInvalidDottedKebabKey,
			len(segments),
			maxDepth,
		)
	}

	for index, segment := range segments {
		if err := ValidateKebabSegment(segment); err != nil {
			return fmt.Errorf("%w: segment %d: %w", ErrInvalidDottedKebabKey, index, err)
		}
	}

	return nil
}
