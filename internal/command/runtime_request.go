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

import "time"

// RuntimeExecutionSpec describes one runtime execution request before binding.
type RuntimeExecutionSpec struct {
	// OptionValues contains already resolved option values.
	OptionValues []OptionValue

	// PositionalValues contains runtime positional values in invocation order.
	PositionalValues []string

	// Metadata contains optional execution-local metadata visible to the handler.
	Metadata Metadata
}

// RuntimeRequest is passed to RuntimeHandler after binding succeeds.
type RuntimeRequest struct {
	runtimeName string
	commandID   ID
	binding     Binding
	input       BoundInput
	startedAt   time.Time
	metadata    Metadata
}

// RuntimeName returns the runtime name.
func (request RuntimeRequest) RuntimeName() string {
	return request.runtimeName
}

// CommandID returns the optional command ID and whether it is set.
func (request RuntimeRequest) CommandID() (ID, bool) {
	if request.commandID == "" {
		return "", false
	}

	return request.commandID, true
}

// Binding returns the runtime binding.
func (request RuntimeRequest) Binding() Binding {
	return request.binding
}

// Input returns the bound command input.
func (request RuntimeRequest) Input() BoundInput {
	return request.input
}

// StartedAt returns the execution start timestamp.
func (request RuntimeRequest) StartedAt() time.Time {
	return request.startedAt
}

// Metadata returns execution-local metadata.
func (request RuntimeRequest) Metadata() Metadata {
	return request.metadata
}

// Option returns a bound option value by canonical option name.
func (request RuntimeRequest) Option(name OptionName) (OptionValue, bool) {
	return request.input.Option(name)
}

// Argument returns a bound positional argument by name.
func (request RuntimeRequest) Argument(name ArgumentName) (BoundArgument, bool) {
	return request.input.Argument(name)
}
