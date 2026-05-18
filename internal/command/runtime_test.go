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
	"context"
	"errors"
	"testing"
)

// TestMustRuntimePanicsForInvalidRuntime verifies fail-fast construction.
func TestMustRuntimePanicsForInvalidRuntime(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustRuntime did not panic")
		}
	}()

	_ = MustRuntime(RuntimeSpec{})
}

// TestNewRuntimeAcceptsValidRuntime verifies full runtime construction.
func TestNewRuntimeAcceptsValidRuntime(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}

	runtime, err := NewRuntime(RuntimeSpec{
		Name:      "test-runtime",
		CommandID: MustID("bench.run"),
		Binding:   runtimeTestBinding(),
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			return OKResult("ok"), nil
		}),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Options:   DefaultRuntimeOptions(),
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewRuntime() returned unexpected error: %v", err)
	}

	if got, want := runtime.Name(), "test-runtime"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if commandID, ok := runtime.CommandID(); !ok || commandID != MustID("bench.run") {
		t.Fatalf("CommandID() = %q, %v; want bench.run, true", commandID, ok)
	}

	if got, want := runtime.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}
}

// TestNewRuntimeDefaults verifies default name, clock, options, and visibility.
func TestNewRuntimeDefaults(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Handler: runtimeTestOKHandler(),
	})

	if got, want := runtime.Name(), "runtime"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if !runtime.Options().RecoverPanics() {
		t.Fatalf("default options RecoverPanics() = false, want true")
	}

	if got, want := runtime.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}
}

// TestNewRuntimeRejectsInvalidRuntime verifies runtime validation.
func TestNewRuntimeRejectsInvalidRuntime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec RuntimeSpec
	}{
		{
			name: "invalid name",
			spec: RuntimeSpec{
				Name:    "bad\nname",
				Handler: runtimeTestOKHandler(),
			},
		},
		{
			name: "invalid command id",
			spec: RuntimeSpec{
				CommandID: ID("Bad.Command"),
				Handler:   runtimeTestOKHandler(),
			},
		},
		{
			name: "invalid binding",
			spec: RuntimeSpec{
				Binding: Binding{
					options: []Option{
						runtimeTestStringOption("output"),
						runtimeTestStringOption("output"),
					},
				},
				Handler: runtimeTestOKHandler(),
			},
		},
		{
			name: "nil handler",
			spec: RuntimeSpec{},
		},
		{
			name: "invalid options",
			spec: RuntimeSpec{
				Handler: runtimeTestOKHandler(),
				Options: RuntimeOptions{
					recoverPanics:     false,
					includePanicStack: true,
				},
			},
		},
		{
			name: "invalid metadata",
			spec: RuntimeSpec{
				Handler: runtimeTestOKHandler(),
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
		},
		{
			name: "invalid visibility",
			spec: RuntimeSpec{
				Handler:    runtimeTestOKHandler(),
				Visibility: Visibility("private"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewRuntime(tt.spec)
			if err == nil {
				t.Fatalf("NewRuntime() returned nil error")
			}

			if !errors.Is(err, ErrInvalidRuntime) && !errors.Is(err, ErrInvalidBinding) {
				t.Fatalf("NewRuntime() error = %v, want ErrInvalidRuntime or ErrInvalidBinding", err)
			}
		})
	}
}
