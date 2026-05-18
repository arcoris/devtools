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

// NewTopic validates raw and returns it as a Topic.
//
// Use NewTopic for topic values loaded from configuration, generated metadata,
// tests, or other external sources where validation errors should be returned.
func NewTopic(raw string) (Topic, error) {
	topic := Topic(raw)
	if err := topic.Validate(); err != nil {
		return "", err
	}

	return topic, nil
}

// ParseTopic is an alias for NewTopic.
//
// The name is useful at call sites where the value is parsed from an external
// string representation.
func ParseTopic(raw string) (Topic, error) {
	return NewTopic(raw)
}

// NewTopicParts validates segments and returns their canonical topic key.
func NewTopicParts(segments ...string) (Topic, error) {
	raw, err := joinTopicSegments(segments)
	if err != nil {
		return "", err
	}

	return NewTopic(raw)
}

// MustTopic validates raw and returns it as a Topic.
//
// MustTopic panics on invalid input. It is intended for static command
// definitions and tests where invalid topic keys are programmer errors.
func MustTopic(raw string) Topic {
	topic, err := NewTopic(raw)
	if err != nil {
		panic(err)
	}

	return topic
}

// MustTopicParts validates segments and returns their canonical topic key.
func MustTopicParts(segments ...string) Topic {
	topic, err := NewTopicParts(segments...)
	if err != nil {
		panic(err)
	}

	return topic
}
