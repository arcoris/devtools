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

var (
	// ErrEmptyOptionSource reports that an option value source was not provided.
	ErrEmptyOptionSource = errors.New("command option source is empty")

	// ErrInvalidOptionSource reports that an option value source is not one of
	// the supported option source states.
	ErrInvalidOptionSource = errors.New("command option source is invalid")
)

// OptionSource describes where an option value came from.
type OptionSource string

const (
	// OptionSourceDefault means the value came from the option declaration's
	// default value.
	OptionSourceDefault OptionSource = "default"

	// OptionSourceInherited means the value was inherited from a parent command,
	// global option scope, profile, or shared command context.
	OptionSourceInherited OptionSource = "inherited"

	// OptionSourceConfig means the value came from configuration.
	OptionSourceConfig OptionSource = "config"

	// OptionSourceEnvironment means the value came from an environment variable.
	OptionSourceEnvironment OptionSource = "environment"

	// OptionSourceRuntime means the value was supplied by runtime wiring.
	OptionSourceRuntime OptionSource = "runtime"

	// OptionSourceInteractive means the value was supplied through an explicit
	// interactive prompt.
	OptionSourceInteractive OptionSource = "interactive"

	// OptionSourceCommandLine means the value was explicitly supplied through
	// command-line syntax.
	OptionSourceCommandLine OptionSource = "command-line"
)

var knownOptionSources = []OptionSource{
	OptionSourceDefault,
	OptionSourceInherited,
	OptionSourceConfig,
	OptionSourceEnvironment,
	OptionSourceRuntime,
	OptionSourceInteractive,
	OptionSourceCommandLine,
}
