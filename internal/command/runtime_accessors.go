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

// Name returns the runtime name.
func (runtime Runtime) Name() string {
	return runtime.name
}

// CommandID returns the optional command ID and whether it is set.
func (runtime Runtime) CommandID() (ID, bool) {
	if runtime.commandID == "" {
		return "", false
	}

	return runtime.commandID, true
}

// HasCommandID reports whether the runtime has a command ID.
func (runtime Runtime) HasCommandID() bool {
	return runtime.commandID != ""
}

// Binding returns the runtime input binding.
func (runtime Runtime) Binding() Binding {
	return runtime.binding
}

// Handler returns the runtime handler.
func (runtime Runtime) Handler() RuntimeHandler {
	return runtime.handler
}

// Options returns runtime options.
func (runtime Runtime) Options() RuntimeOptions {
	return runtime.options
}

// Metadata returns runtime metadata.
func (runtime Runtime) Metadata() Metadata {
	return runtime.metadata
}

// Visibility returns runtime visibility.
func (runtime Runtime) Visibility() Visibility {
	return runtime.visibility
}

// IsVisibleByDefault reports whether default reports/logs/discovery should
// expose the runtime.
func (runtime Runtime) IsVisibleByDefault() bool {
	return runtime.visibility.IsDiscoverableByDefault()
}
