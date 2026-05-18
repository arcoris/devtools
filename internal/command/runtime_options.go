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

import "fmt"

// RuntimeOptions controls runtime behavior.
type RuntimeOptions struct {
	configured        bool
	recoverPanics     bool
	includePanicStack bool
	suppressEvents    bool
}

// RuntimeOptionsSpec describes runtime options before validation.
type RuntimeOptionsSpec struct {
	// RecoverPanics controls whether handler panics are converted into failed
	// Results.
	RecoverPanics bool

	// IncludePanicStack controls whether recovered panic messages include
	// debug.Stack output.
	IncludePanicStack bool

	// SuppressEvents disables lifecycle event emission even when EventSink is
	// configured.
	SuppressEvents bool
}

// NewRuntimeOptions validates spec and returns RuntimeOptions.
func NewRuntimeOptions(spec RuntimeOptionsSpec) (RuntimeOptions, error) {
	options := RuntimeOptions{
		configured:        true,
		recoverPanics:     spec.RecoverPanics,
		includePanicStack: spec.IncludePanicStack,
		suppressEvents:    spec.SuppressEvents,
	}

	if err := options.Validate(); err != nil {
		return RuntimeOptions{}, err
	}

	return options, nil
}

// MustRuntimeOptions validates spec and returns RuntimeOptions.
//
// MustRuntimeOptions panics on invalid input.
func MustRuntimeOptions(spec RuntimeOptionsSpec) RuntimeOptions {
	options, err := NewRuntimeOptions(spec)
	if err != nil {
		panic(err)
	}

	return options
}

// DefaultRuntimeOptions returns ordinary runtime behavior.
func DefaultRuntimeOptions() RuntimeOptions {
	return RuntimeOptions{
		configured:        true,
		recoverPanics:     true,
		includePanicStack: false,
		suppressEvents:    false,
	}
}

// OrDefault returns DefaultRuntimeOptions when options is zero.
func (options RuntimeOptions) OrDefault() RuntimeOptions {
	if !options.configured && options == (RuntimeOptions{}) {
		return DefaultRuntimeOptions()
	}

	return options
}

// RecoverPanics reports whether panics are recovered.
func (options RuntimeOptions) RecoverPanics() bool {
	return options.recoverPanics
}

// IncludePanicStack reports whether recovered panic errors include stack text.
func (options RuntimeOptions) IncludePanicStack() bool {
	return options.includePanicStack
}

// SuppressEvents reports whether lifecycle events are disabled.
func (options RuntimeOptions) SuppressEvents() bool {
	return options.suppressEvents
}

// Validate verifies runtime option consistency.
func (options RuntimeOptions) Validate() error {
	if options.includePanicStack && !options.recoverPanics {
		return fmt.Errorf("%w: panic stack requires panic recovery", ErrInvalidRuntime)
	}

	return nil
}

// WithRecoverPanics returns a validated copy with RecoverPanics replaced.
func (options RuntimeOptions) WithRecoverPanics(enabled bool) (RuntimeOptions, error) {
	next := options
	next.configured = true
	next.recoverPanics = enabled

	if !enabled {
		next.includePanicStack = false
	}

	if err := next.Validate(); err != nil {
		return RuntimeOptions{}, err
	}

	return next, nil
}

// MustWithRecoverPanics returns a validated copy with RecoverPanics replaced
// and panics on invalid input.
func (options RuntimeOptions) MustWithRecoverPanics(enabled bool) RuntimeOptions {
	next, err := options.WithRecoverPanics(enabled)
	if err != nil {
		panic(err)
	}

	return next
}

// WithIncludePanicStack returns a validated copy with IncludePanicStack
// replaced.
func (options RuntimeOptions) WithIncludePanicStack(enabled bool) (RuntimeOptions, error) {
	next := options
	next.configured = true
	next.includePanicStack = enabled

	if err := next.Validate(); err != nil {
		return RuntimeOptions{}, err
	}

	return next, nil
}

// MustWithIncludePanicStack returns a validated copy with IncludePanicStack
// replaced and panics on invalid input.
func (options RuntimeOptions) MustWithIncludePanicStack(enabled bool) RuntimeOptions {
	next, err := options.WithIncludePanicStack(enabled)
	if err != nil {
		panic(err)
	}

	return next
}

// WithSuppressEvents returns a validated copy with SuppressEvents replaced.
func (options RuntimeOptions) WithSuppressEvents(enabled bool) (RuntimeOptions, error) {
	next := options
	next.configured = true
	next.suppressEvents = enabled

	if err := next.Validate(); err != nil {
		return RuntimeOptions{}, err
	}

	return next, nil
}

// MustWithSuppressEvents returns a validated copy with SuppressEvents replaced
// and panics on invalid input.
func (options RuntimeOptions) MustWithSuppressEvents(enabled bool) RuntimeOptions {
	next, err := options.WithSuppressEvents(enabled)
	if err != nil {
		panic(err)
	}

	return next
}
