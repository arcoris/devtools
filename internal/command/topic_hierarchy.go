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

// Parent returns the parent topic and whether such parent exists.
//
// Examples:
//
//   - "diagnostics.perf" returns "diagnostics", true;
//   - "profiling" returns "", false;
//   - the zero topic returns "", false.
//
// Parent does not validate the topic; call Validate first for untrusted input.
func (topic Topic) Parent() (Topic, bool) {
	if topic == "" {
		return "", false
	}

	raw := string(topic)
	index := strings.LastIndex(raw, TopicSeparator)
	if index < 0 {
		return "", false
	}

	return Topic(raw[:index]), true
}

// HasParent reports whether the topic has a hierarchical parent.
func (topic Topic) HasParent() bool {
	_, ok := topic.Parent()

	return ok
}

// HasPrefix reports whether topic is equal to prefix or belongs to prefix's
// hierarchical subtree.
//
// An empty prefix always returns false because topic prefixes are expected to be
// explicit.
func (topic Topic) HasPrefix(prefix Topic) bool {
	if topic == "" || prefix == "" {
		return false
	}

	if topic == prefix {
		return true
	}

	return strings.HasPrefix(string(topic), string(prefix)+TopicSeparator)
}

// TrimPrefix removes prefix from topic and returns the remaining relative topic.
//
// The second return value is false when prefix is empty or not a prefix of
// topic. Trimming an identical prefix returns the zero topic.
func (topic Topic) TrimPrefix(prefix Topic) (Topic, bool) {
	if !topic.HasPrefix(prefix) {
		return "", false
	}

	if topic == prefix {
		return "", true
	}

	trimmed := strings.TrimPrefix(string(topic), string(prefix)+TopicSeparator)

	return Topic(trimmed), true
}

// Append returns a child topic by appending one validated segment to topic.
//
// Append accepts exactly one segment, not a dotted suffix. This prevents
// accidental multi-level hierarchy changes in one call.
func (topic Topic) Append(segment string) (Topic, error) {
	if strings.Contains(segment, TopicSeparator) {
		return "", fmt.Errorf(
			"%w: appended segment must not contain %q",
			ErrInvalidTopic,
			TopicSeparator,
		)
	}

	if err := validateTopicSegment(0, segment); err != nil {
		return "", err
	}

	if topic.IsZero() {
		return Topic(segment), nil
	}

	if err := topic.Validate(); err != nil {
		return "", err
	}

	return Topic(string(topic) + TopicSeparator + segment), nil
}

// MustAppend returns a child topic by appending segment to topic.
//
// MustAppend panics on invalid input. It is intended for static command
// definitions and tests.
func (topic Topic) MustAppend(segment string) Topic {
	child, err := topic.Append(segment)
	if err != nil {
		panic(err)
	}

	return child
}

// Join returns a topic formed by appending another topic's segments.
func (topic Topic) Join(other Topic) (Topic, error) {
	switch {
	case topic.IsZero():
		return NewTopic(other.String())
	case other.IsZero():
		return NewTopic(topic.String())
	}

	if err := topic.Validate(); err != nil {
		return "", err
	}

	if err := other.Validate(); err != nil {
		return "", err
	}

	return NewTopic(string(topic) + TopicSeparator + string(other))
}

// MustJoin returns a topic formed by appending another topic's segments.
func (topic Topic) MustJoin(other Topic) Topic {
	joined, err := topic.Join(other)
	if err != nil {
		panic(err)
	}

	return joined
}
