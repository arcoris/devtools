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
	"testing"
)

func TestTopicMarshalText(t *testing.T) {
	t.Parallel()

	text, err := MustTopic("diagnostics.perf").MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() returned unexpected error: %v", err)
	}

	if got, want := string(text), "diagnostics.perf"; got != want {
		t.Fatalf("MarshalText() = %q, want %q", got, want)
	}
}

func TestTopicMarshalTextRejectsInvalidTopic(t *testing.T) {
	t.Parallel()

	_, err := Topic("Diagnostics").MarshalText()
	if err == nil {
		t.Fatalf("MarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("MarshalText() error = %v, want ErrInvalidTopic", err)
	}
}

func TestTopicUnmarshalText(t *testing.T) {
	t.Parallel()

	var topic Topic

	if err := topic.UnmarshalText([]byte("diagnostics.perf")); err != nil {
		t.Fatalf("UnmarshalText() returned unexpected error: %v", err)
	}

	if got, want := topic, MustTopic("diagnostics.perf"); got != want {
		t.Fatalf("UnmarshalText() = %q, want %q", got, want)
	}
}

func TestTopicUnmarshalTextRejectsInvalidTopic(t *testing.T) {
	t.Parallel()

	var topic Topic

	err := topic.UnmarshalText([]byte("Diagnostics"))
	if err == nil {
		t.Fatalf("UnmarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("UnmarshalText() error = %v, want ErrInvalidTopic", err)
	}
}

func TestTopicUnmarshalTextRejectsNilReceiver(t *testing.T) {
	t.Parallel()

	var topic *Topic

	err := topic.UnmarshalText([]byte("diagnostics"))
	if err == nil {
		t.Fatalf("UnmarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("UnmarshalText() error = %v, want ErrInvalidTopic", err)
	}
}
