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

// TestRuntimeOptionsPreservesExplicitFalse verifies explicit false values do
// not collapse back to defaults when RuntimeSpec is normalized.
func TestRuntimeOptionsPreservesExplicitFalse(t *testing.T) {
	t.Parallel()

	options := MustRuntimeOptions(RuntimeOptionsSpec{
		RecoverPanics: false,
	})

	if options.OrDefault().RecoverPanics() {
		t.Fatalf("OrDefault() changed explicit RecoverPanics=false to true")
	}

	runtime := MustRuntime(RuntimeSpec{
		Binding: EmptyBinding(),
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			return OKResult("ok"), nil
		}),
		Options: options,
	})

	if runtime.Options().RecoverPanics() {
		t.Fatalf("NewRuntime() changed explicit RecoverPanics=false to true")
	}
}

// TestRuntimeOptionsRejectInvalid verifies invalid option combinations.
func TestRuntimeOptionsRejectInvalid(t *testing.T) {
	t.Parallel()

	_, err := NewRuntimeOptions(RuntimeOptionsSpec{
		RecoverPanics:     false,
		IncludePanicStack: true,
	})
	if err == nil {
		t.Fatalf("NewRuntimeOptions() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRuntime) {
		t.Fatalf("NewRuntimeOptions() error = %v, want ErrInvalidRuntime", err)
	}
}

// TestRuntimeOptions verifies option construction and helpers.
func TestRuntimeOptions(t *testing.T) {
	t.Parallel()

	options := DefaultRuntimeOptions()

	if !options.RecoverPanics() {
		t.Fatalf("RecoverPanics() = false, want true")
	}

	options = options.MustWithIncludePanicStack(true)

	if !options.IncludePanicStack() {
		t.Fatalf("IncludePanicStack() = false, want true")
	}

	options = options.MustWithSuppressEvents(true)

	if !options.SuppressEvents() {
		t.Fatalf("SuppressEvents() = false, want true")
	}

	options = options.MustWithRecoverPanics(false)

	if options.RecoverPanics() {
		t.Fatalf("RecoverPanics() = true, want false")
	}

	if options.IncludePanicStack() {
		t.Fatalf("IncludePanicStack() = true after disabling recovery, want false")
	}
}
