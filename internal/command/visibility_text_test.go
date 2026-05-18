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

func TestVisibilityMarshalText(t *testing.T) {
	t.Parallel()

	text, err := VisibilityHidden.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() returned unexpected error: %v", err)
	}

	if got, want := string(text), "hidden"; got != want {
		t.Fatalf("MarshalText() = %q, want %q", got, want)
	}
}

func TestVisibilityMarshalTextRejectsInvalidVisibility(t *testing.T) {
	t.Parallel()

	_, err := Visibility("private").MarshalText()
	if err == nil {
		t.Fatalf("MarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("MarshalText() error = %v, want ErrInvalidVisibility", err)
	}
}

func TestVisibilityUnmarshalText(t *testing.T) {
	t.Parallel()

	var visibility Visibility

	if err := visibility.UnmarshalText([]byte("internal")); err != nil {
		t.Fatalf("UnmarshalText() returned unexpected error: %v", err)
	}

	if got, want := visibility, VisibilityInternal; got != want {
		t.Fatalf("UnmarshalText() = %q, want %q", got, want)
	}
}

func TestVisibilityUnmarshalTextRejectsInvalidVisibility(t *testing.T) {
	t.Parallel()

	var visibility Visibility

	err := visibility.UnmarshalText([]byte("private"))
	if err == nil {
		t.Fatalf("UnmarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("UnmarshalText() error = %v, want ErrInvalidVisibility", err)
	}
}

func TestVisibilityUnmarshalTextRejectsNilReceiver(t *testing.T) {
	t.Parallel()

	var visibility *Visibility

	err := visibility.UnmarshalText([]byte("public"))
	if err == nil {
		t.Fatalf("UnmarshalText() returned nil error")
	}

	if !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("UnmarshalText() error = %v, want ErrInvalidVisibility", err)
	}
}
