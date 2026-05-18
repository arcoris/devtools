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

// TestRuntimeExecutePropagatesPanicWhenRecoveryDisabled verifies recovery option.
func TestRuntimeExecutePropagatesPanicWhenRecoveryDisabled(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Options: MustRuntimeOptions(RuntimeOptionsSpec{
			RecoverPanics: boolPointer(false),
		}),
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			panic("boom")
		}),
	})

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("Execute() did not propagate panic")
		}
	}()

	_, _ = runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
}

// TestRuntimeExecuteRecoversPanic verifies panic recovery.
func TestRuntimeExecuteRecoversPanic(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Options: DefaultRuntimeOptions(),
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			panic("boom")
		}),
	})

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimePanic) {
		t.Fatalf("Execute() error = %v, want ErrRuntimePanic", err)
	}

	if !result.IsFailed() {
		t.Fatalf("result IsFailed() = false, want true")
	}
}
