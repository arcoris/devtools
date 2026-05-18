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

// Validate verifies that the topic satisfies the topic grammar.
func (topic Topic) Validate() error {
	raw := string(topic)

	if raw == "" {
		return ErrEmptyTopic
	}

	if len(raw) > maxTopicLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidTopic,
			len(raw),
			maxTopicLength,
		)
	}

	if strings.HasPrefix(raw, TopicSeparator) {
		return fmt.Errorf("%w: must not start with %q", ErrInvalidTopic, TopicSeparator)
	}

	if strings.HasSuffix(raw, TopicSeparator) {
		return fmt.Errorf("%w: must not end with %q", ErrInvalidTopic, TopicSeparator)
	}

	segments := strings.Split(raw, TopicSeparator)
	return validateTopicSegments(segments)
}

// validateTopicSegments validates all topic segments plus whole-topic limits.
func validateTopicSegments(segments []string) error {
	if len(segments) == 0 {
		return ErrEmptyTopic
	}

	if len(segments) > maxTopicDepth {
		return fmt.Errorf(
			"%w: depth %d exceeds maximum depth %d",
			ErrInvalidTopic,
			len(segments),
			maxTopicDepth,
		)
	}

	if length := canonicalTopicLength(segments); length > maxTopicLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidTopic,
			length,
			maxTopicLength,
		)
	}

	for index, segment := range segments {
		if err := validateTopicSegment(index, segment); err != nil {
			return err
		}
	}

	return nil
}

// validateTopicSegment validates one topic segment and wraps generic segment
// validation errors with topic-specific diagnostics.
func validateTopicSegment(index int, segment string) error {
	if len(segment) > maxTopicSegmentLength {
		return fmt.Errorf(
			"%w: segment %d length %d exceeds maximum length %d",
			ErrInvalidTopic,
			index,
			len(segment),
			maxTopicSegmentLength,
		)
	}

	if err := validateCommandNameSegment(segment); err != nil {
		return fmt.Errorf("%w: segment %d: %w", ErrInvalidTopic, index, err)
	}

	return nil
}

// joinTopicSegments validates segments and returns their canonical topic text.
func joinTopicSegments(segments []string) (string, error) {
	if err := validateTopicSegments(segments); err != nil {
		return "", err
	}

	return strings.Join(segments, TopicSeparator), nil
}

// canonicalTopicLength returns the byte length of a joined topic key without
// allocating the joined string.
func canonicalTopicLength(segments []string) int {
	if len(segments) == 0 {
		return 0
	}

	length := 0
	for index, segment := range segments {
		if index > 0 {
			length += len(TopicSeparator)
		}

		length += len(segment)
	}

	return length
}
