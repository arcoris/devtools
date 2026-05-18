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

// NewRuntime validates spec and returns Runtime.
func NewRuntime(spec RuntimeSpec) (Runtime, error) {
	name := strings.TrimSpace(spec.Name)
	if name == "" {
		name = "runtime"
	}

	clock := spec.Clock
	if clock == nil {
		clock = SystemRuntimeClock{}
	}

	options := spec.Options.OrDefault()
	visibility := spec.Visibility.OrDefault()

	runtime := Runtime{
		name:       name,
		commandID:  spec.CommandID,
		binding:    spec.Binding,
		handler:    spec.Handler,
		clock:      clock,
		eventSink:  spec.EventSink,
		options:    options,
		metadata:   spec.Metadata,
		visibility: visibility,
	}

	if err := runtime.Validate(); err != nil {
		return Runtime{}, err
	}

	return runtime, nil
}

// MustRuntime validates spec and returns Runtime.
//
// MustRuntime panics on invalid input. It is intended for tests and controlled
// static wiring.
func MustRuntime(spec RuntimeSpec) Runtime {
	runtime, err := NewRuntime(spec)
	if err != nil {
		panic(err)
	}

	return runtime
}
