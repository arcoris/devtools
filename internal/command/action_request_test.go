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

	"arcoris.dev/devtools/internal/textvalidate"
)

func TestNewActionRequestAcceptsValidRequest(t *testing.T) {
	t.Parallel()

	request, err := NewActionRequest(
		actionTestCommandNode(),
		[]string{"arg1", "arg2"},
		map[string]string{"source": "test"},
	)
	if err != nil {
		t.Fatalf("NewActionRequest() returned unexpected error: %v", err)
	}

	if !request.HasArguments() {
		t.Fatalf("HasArguments() = false, want true")
	}

	if got, want := request.ArgCount(), 2; got != want {
		t.Fatalf("ArgCount() = %d, want %d", got, want)
	}

	if got, ok := request.Argument(1); !ok || got != "arg2" {
		t.Fatalf("Argument(1) = %q, %v; want arg2, true", got, ok)
	}

	if _, ok := request.Argument(-1); ok {
		t.Fatalf("Argument(-1) ok = true, want false")
	}

	if got, ok := request.Field("source"); !ok || got != "test" {
		t.Fatalf("Field(source) = %q, %v; want test, true", got, ok)
	}
}

func TestNewActionRequestRejectsNonCommandNode(t *testing.T) {
	t.Parallel()

	_, err := NewActionRequest(actionTestFamilyNode(), nil, nil)
	if err == nil {
		t.Fatalf("NewActionRequest() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("NewActionRequest() error = %v, want ErrInvalidActionRequest", err)
	}
}

func TestNewActionRequestRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "zero node",
			run: func() error {
				return ActionRequest{}.Validate()
			},
		},
		{
			name: "invalid argument control",
			run: func() error {
				_, err := NewActionRequest(actionTestCommandNode(), []string{"bad\x00value"}, nil)
				return err
			},
		},
		{
			name: "invalid argument too long",
			run: func() error {
				_, err := NewActionRequest(actionTestCommandNode(), []string{strings.Repeat("x", maxActionArgumentLength+1)}, nil)
				return err
			},
		},
		{
			name: "invalid field key",
			run: func() error {
				_, err := NewActionRequest(actionTestCommandNode(), nil, map[string]string{"Bad": "value"})
				return err
			},
		},
		{
			name: "invalid field value",
			run: func() error {
				_, err := NewActionRequest(actionTestCommandNode(), nil, map[string]string{"source": "bad\x00value"})
				return err
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.run()
			if err == nil {
				t.Fatalf("%s returned nil error", test.name)
			}

			if !errors.Is(err, ErrInvalidActionRequest) {
				t.Fatalf("%s error = %v, want ErrInvalidActionRequest", test.name, err)
			}
		})
	}
}

func TestNewActionRequestWrapsTextvalidateErrorsWithoutResultSentinel(t *testing.T) {
	t.Parallel()

	_, err := NewActionRequest(actionTestCommandNode(), nil, map[string]string{"Bad": "value"})
	if err == nil {
		t.Fatalf("NewActionRequest() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("NewActionRequest() error = %v, want ErrInvalidActionRequest", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("NewActionRequest() error = %v, want ErrInvalidDottedKebabKey", err)
	}

	if errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("NewActionRequest() error = %v, must not wrap ErrInvalidActionResult", err)
	}
}

func TestActionRequestCopySemantics(t *testing.T) {
	t.Parallel()

	args := []string{"one"}
	fields := map[string]string{"source": "test"}

	request := MustActionRequest(actionTestCommandNode(), args, fields)

	args[0] = "changed"
	fields["source"] = "changed"

	if got, want := request.Args()[0], "one"; got != want {
		t.Fatalf("argument changed through input slice: got %q, want %q", got, want)
	}

	if got, want := actionTestRequestField(t, request, "source"), "test"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	outArgs := request.Args()
	outArgs[0] = "changed"

	if got, want := request.Args()[0], "one"; got != want {
		t.Fatalf("argument changed through output slice: got %q, want %q", got, want)
	}

	outFields := request.FieldMap()
	outFields["source"] = "changed"

	if got, want := actionTestRequestField(t, request, "source"), "test"; got != want {
		t.Fatalf("field changed through FieldMap: got %q, want %q", got, want)
	}
}

func TestActionRequestFieldAccessors(t *testing.T) {
	t.Parallel()

	request := MustActionRequest(actionTestCommandNode(), nil, map[string]string{
		"z.key": "z",
		"a.key": "a",
	})

	if !request.HasFields() {
		t.Fatalf("HasFields() = false, want true")
	}

	if got, want := request.FieldCount(), 2; got != want {
		t.Fatalf("FieldCount() = %d, want %d", got, want)
	}

	if !request.HasField("a.key") {
		t.Fatalf("HasField(a.key) = false, want true")
	}

	assertStringSlicesEqual(t, request.FieldKeys(), []string{"a.key", "z.key"})
}

func TestActionRequestWithHelpers(t *testing.T) {
	t.Parallel()

	request := actionTestRequest(t).
		MustWithNode(actionTestAlternateCommandNode()).
		MustWithArguments("one", "two").
		MustWithField("source", "test").
		MustWithField("mode", "smoke")

	if got, want := request.Node.ID(), ID("check"); got != want {
		t.Fatalf("Node.ID() = %q, want %q", got, want)
	}

	if got, want := request.ArgCount(), 2; got != want {
		t.Fatalf("ArgCount() = %d, want %d", got, want)
	}

	assertStringSlicesEqual(t, request.FieldKeys(), []string{"mode", "source"})

	if request.WithoutField("mode").HasField("mode") {
		t.Fatalf("WithoutField() still has removed field")
	}

	if request.WithoutFields().HasFields() {
		t.Fatalf("WithoutFields() still has fields")
	}

	if request.WithoutArguments().HasArguments() {
		t.Fatalf("WithoutArguments() still has arguments")
	}
}

func TestActionRequestWithFieldsCopiesInput(t *testing.T) {
	t.Parallel()

	fields := map[string]string{"source": "test"}
	request := actionTestRequest(t).MustWithFields(fields)

	fields["source"] = "changed"

	if got, want := actionTestRequestField(t, request, "source"), "test"; got != want {
		t.Fatalf("field changed through WithFields input: got %q, want %q", got, want)
	}
}

func TestActionRequestWithHelpersRejectInvalidValues(t *testing.T) {
	t.Parallel()

	request := actionTestRequest(t)

	if _, err := request.WithNode(actionTestFamilyNode()); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("WithNode() error = %v, want ErrInvalidActionRequest", err)
	}

	if _, err := request.WithArguments("bad\x00value"); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("WithArguments() error = %v, want ErrInvalidActionRequest", err)
	}

	if _, err := request.WithFields(map[string]string{"Bad": "value"}); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("WithFields() error = %v, want ErrInvalidActionRequest", err)
	}

	if _, err := request.WithField("Bad", "value"); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("WithField(key) error = %v, want ErrInvalidActionRequest", err)
	}

	if _, err := request.WithField("source", "bad\x00value"); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("WithField(value) error = %v, want ErrInvalidActionRequest", err)
	}
}

func TestMustActionRequestPanicsForInvalidRequest(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustActionRequest(actionTestFamilyNode(), nil, nil)
	})
}
