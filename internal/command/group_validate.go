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

// Validate verifies that the group satisfies the group grammar.
func (group Group) Validate() error {
	raw := string(group)

	if raw == "" {
		return ErrEmptyGroup
	}

	if len(raw) > maxGroupLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidGroup,
			len(raw),
			maxGroupLength,
		)
	}

	if strings.HasPrefix(raw, GroupSeparator) {
		return fmt.Errorf("%w: must not start with %q", ErrInvalidGroup, GroupSeparator)
	}

	if strings.HasSuffix(raw, GroupSeparator) {
		return fmt.Errorf("%w: must not end with %q", ErrInvalidGroup, GroupSeparator)
	}

	segments := strings.Split(raw, GroupSeparator)
	return validateGroupSegments(segments)
}

// validateGroupSegments validates all group segments plus whole-group limits.
func validateGroupSegments(segments []string) error {
	if len(segments) == 0 {
		return ErrEmptyGroup
	}

	if len(segments) > maxGroupDepth {
		return fmt.Errorf(
			"%w: depth %d exceeds maximum depth %d",
			ErrInvalidGroup,
			len(segments),
			maxGroupDepth,
		)
	}

	if length := canonicalGroupLength(segments); length > maxGroupLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidGroup,
			length,
			maxGroupLength,
		)
	}

	for index, segment := range segments {
		if err := validateGroupSegment(index, segment); err != nil {
			return err
		}
	}

	return nil
}

// validateGroupSegment validates one group segment and wraps generic segment
// validation errors with group-specific diagnostics.
func validateGroupSegment(index int, segment string) error {
	if len(segment) > maxGroupSegmentLength {
		return fmt.Errorf(
			"%w: segment %d length %d exceeds maximum length %d",
			ErrInvalidGroup,
			index,
			len(segment),
			maxGroupSegmentLength,
		)
	}

	if err := validateCommandNameSegment(segment); err != nil {
		return fmt.Errorf("%w: segment %d: %w", ErrInvalidGroup, index, err)
	}

	return nil
}

// joinGroupSegments validates segments and returns their canonical group text.
func joinGroupSegments(segments []string) (string, error) {
	if err := validateGroupSegments(segments); err != nil {
		return "", err
	}

	return strings.Join(segments, GroupSeparator), nil
}

// canonicalGroupLength returns the byte length of a joined group key without
// allocating the joined string.
func canonicalGroupLength(segments []string) int {
	if len(segments) == 0 {
		return 0
	}

	length := 0
	for index, segment := range segments {
		if index > 0 {
			length += len(GroupSeparator)
		}

		length += len(segment)
	}

	return length
}
