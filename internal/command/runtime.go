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

import "errors"

const (
	// maxRuntimeNameLength is the maximum byte length of a runtime name.
	maxRuntimeNameLength = 255

	// maxRuntimePanicMessageLength is the maximum byte length kept from a panic
	// value when RuntimeOptions.RecoverPanics is enabled.
	maxRuntimePanicMessageLength = 4096
)

var (
	// ErrInvalidRuntime reports that a runtime declaration or runtime request is
	// malformed.
	ErrInvalidRuntime = errors.New("command runtime is invalid")

	// ErrRuntimeExecution reports that a runtime execution failed.
	ErrRuntimeExecution = errors.New("command runtime execution failed")

	// ErrRuntimeCanceled reports that a runtime execution was canceled through
	// context cancellation.
	ErrRuntimeCanceled = errors.New("command runtime canceled")

	// ErrRuntimePanic reports that a runtime handler panicked and the runtime
	// converted that panic into a structured failure.
	ErrRuntimePanic = errors.New("command runtime panic")
)

// RuntimeSpec describes command runtime wiring before validation.
//
// RuntimeSpec is a construction DTO. Runtime stores values in immutable-style
// form. Runtime does not own CLI parsing, Cobra command construction,
// configuration loading, environment loading, file writing, or process exit.
// It coordinates the already declared command binding with an executable
// handler and optional lifecycle event sinks.
type RuntimeSpec struct {
	// Name is an optional runtime name used in diagnostics and events.
	//
	// Empty defaults to "runtime".
	Name string

	// CommandID optionally identifies the command executed by this runtime.
	CommandID ID

	// Binding describes accepted command inputs.
	Binding Binding

	// Handler executes a command after input binding.
	Handler RuntimeHandler

	// Clock provides timestamps for result timing and lifecycle events.
	//
	// Nil defaults to SystemRuntimeClock.
	Clock RuntimeClock

	// EventSink receives lifecycle events.
	//
	// Nil means lifecycle events are created internally but not exported.
	// RecordEvent failures currently fail runtime execution. This keeps audit
	// and test observers strict by default; ordinary logging adapters should
	// wrap their sink if they want best-effort behavior.
	EventSink RuntimeEventSink

	// Options controls runtime behavior.
	Options RuntimeOptions

	// Metadata contains optional machine-facing runtime metadata.
	Metadata Metadata

	// Visibility controls default exposure in reports, logs, docs, and
	// discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Runtime is an adapter-neutral command execution orchestrator.
//
// Runtime coordinates the command lifecycle:
//
//   - validates runtime wiring;
//   - binds already parsed input values using Binding;
//   - emits lifecycle Events through RuntimeEventSink;
//   - calls RuntimeHandler with a RuntimeRequest;
//   - normalizes Result timing, status, cancellation, panic recovery, and
//     handler errors.
//
// Runtime deliberately does not parse os.Args, read environment variables,
// load config, use Cobra, write files, or call os.Exit. CLI adapters and
// higher-level application packages should perform those responsibilities
// around this kernel.
type Runtime struct {
	name       string
	commandID  ID
	binding    Binding
	handler    RuntimeHandler
	clock      RuntimeClock
	eventSink  RuntimeEventSink
	options    RuntimeOptions
	metadata   Metadata
	visibility Visibility
}
