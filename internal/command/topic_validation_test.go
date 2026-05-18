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
	"errors"
	"strings"
	"testing"
)

func TestNewTopicRejectsInvalidTopics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{name: "empty", raw: "", err: ErrEmptyTopic},
		{name: "leading separator", raw: ".profiling", err: ErrInvalidTopic},
		{name: "trailing separator", raw: "profiling.", err: ErrInvalidTopic},
		{name: "consecutive separators", raw: "diagnostics..perf", err: ErrInvalidTopic},
		{name: "invalid segment", raw: "diagnostics.Perf", err: ErrInvalidTopic},
		{name: "too long complete topic", raw: "a." + strings.Repeat("b", maxTopicLength), err: ErrInvalidTopic},
		{name: "too deep", raw: strings.Repeat("a.", maxTopicDepth) + "a", err: ErrInvalidTopic},
		{name: "too long segment", raw: "a" + strings.Repeat("b", maxTopicSegmentLength), err: ErrInvalidTopic},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			topic, err := NewTopic(test.raw)
			if err == nil {
				t.Fatalf("NewTopic(%q) returned nil error and topic %q", test.raw, topic)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("NewTopic(%q) error = %v, want errors.Is(..., %v)", test.raw, err, test.err)
			}

			if Topic(test.raw).IsValid() {
				t.Fatalf("Topic(%q).IsValid() = true, want false", test.raw)
			}
		})
	}
}

func TestValidateTopicSegmentWrapsGenericValidation(t *testing.T) {
	t.Parallel()

	err := validateTopicSegment(0, "Invalid")
	if err == nil {
		t.Fatalf("validateTopicSegment() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("validateTopicSegment() error = %v, want ErrInvalidTopic", err)
	}

	if !errors.Is(err, ErrInvalidCommandNameSegment) {
		t.Fatalf("validateTopicSegment() error = %v, want ErrInvalidCommandNameSegment", err)
	}
}

func TestValidateTopicSegmentWrapsEmptySegment(t *testing.T) {
	t.Parallel()

	err := validateTopicSegment(1, "")
	if err == nil {
		t.Fatalf("validateTopicSegment() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("validateTopicSegment() error = %v, want ErrInvalidTopic", err)
	}

	if !errors.Is(err, ErrEmptyCommandNameSegment) {
		t.Fatalf("validateTopicSegment() error = %v, want ErrEmptyCommandNameSegment", err)
	}
}

func TestNewTopicPartsRejectsNoSegments(t *testing.T) {
	t.Parallel()

	_, err := NewTopicParts()
	if err == nil {
		t.Fatalf("NewTopicParts() returned nil error")
	}

	if !errors.Is(err, ErrEmptyTopic) {
		t.Fatalf("NewTopicParts() error = %v, want ErrEmptyTopic", err)
	}
}

func TestCanonicalTopicLength(t *testing.T) {
	t.Parallel()

	if got, want := canonicalTopicLength([]string{"diagnostics", "perf"}), len("diagnostics.perf"); got != want {
		t.Fatalf("canonicalTopicLength() = %d, want %d", got, want)
	}
}
