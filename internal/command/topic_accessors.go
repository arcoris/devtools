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

import "strings"

// String returns the canonical string representation of the topic.
func (topic Topic) String() string {
	return string(topic)
}

// Key returns the stable map key for topic.
func (topic Topic) Key() string {
	return string(topic)
}

// IsZero reports whether the topic has not been set.
func (topic Topic) IsZero() bool {
	return topic == ""
}

// IsValid reports whether the topic satisfies the topic grammar.
func (topic Topic) IsValid() bool {
	return topic.Validate() == nil
}

// Equal reports whether two topics are exactly the same key.
func (topic Topic) Equal(other Topic) bool {
	return topic == other
}

// Parts returns a detached copy of the topic segments.
//
// Parts performs a lexical split and does not validate the topic. Call Validate
// first when the value may be untrusted. The returned slice can be safely
// modified by the caller.
func (topic Topic) Parts() []string {
	if topic == "" {
		return nil
	}

	parts := strings.Split(string(topic), TopicSeparator)
	out := make([]string, len(parts))
	copy(out, parts)

	return out
}

// Len returns the number of hierarchical segments in the topic.
func (topic Topic) Len() int {
	return topic.Depth()
}

// Depth returns the number of hierarchical segments in the topic.
func (topic Topic) Depth() int {
	if topic == "" {
		return 0
	}

	return strings.Count(string(topic), TopicSeparator) + 1
}

// At returns the segment at index.
//
// The second return value is false when index is out of range. At never panics.
func (topic Topic) At(index int) (string, bool) {
	if index < 0 || index >= topic.Depth() {
		return "", false
	}

	start := 0
	raw := string(topic)
	for current := 0; current < index; current++ {
		next := strings.Index(raw[start:], TopicSeparator)
		if next < 0 {
			return "", false
		}

		start += next + len(TopicSeparator)
	}

	end := strings.Index(raw[start:], TopicSeparator)
	if end < 0 {
		return raw[start:], true
	}

	return raw[start : start+end], true
}

// Leaf returns the final segment of the topic.
//
// For "diagnostics.perf", Leaf returns "perf". For the zero topic, Leaf
// returns an empty string. Leaf does not validate the topic; call Validate first
// for untrusted input.
func (topic Topic) Leaf() string {
	if topic == "" {
		return ""
	}

	raw := string(topic)
	index := strings.LastIndex(raw, TopicSeparator)
	if index < 0 {
		return raw
	}

	return raw[index+len(TopicSeparator):]
}
