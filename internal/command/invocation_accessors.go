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

import "sort"

// Clone returns a detached copy of the invocation.
func (invocation Invocation) Clone() Invocation {
	return Invocation{
		arguments:    cloneStringSlice(invocation.arguments),
		rawArguments: cloneStringSlice(invocation.rawArguments),
		workingDir:   invocation.workingDir,
		env:          cloneStringMap(invocation.env),
		fields:       cloneStringMap(invocation.fields),
	}
}

// Spec returns a detached construction spec for invocation.
func (invocation Invocation) Spec() InvocationSpec {
	return invocation.spec()
}

// IsZero reports whether invocation has no metadata.
func (invocation Invocation) IsZero() bool {
	return len(invocation.arguments) == 0 &&
		len(invocation.rawArguments) == 0 &&
		invocation.workingDir == "" &&
		len(invocation.env) == 0 &&
		len(invocation.fields) == 0
}

// Arguments returns detached parsed positional arguments.
func (invocation Invocation) Arguments() []string {
	return cloneStringSlice(invocation.arguments)
}

// Argument returns the parsed positional argument at index.
//
// The second return value is false when index is out of range. Argument never
// panics.
func (invocation Invocation) Argument(index int) (string, bool) {
	if index < 0 || index >= len(invocation.arguments) {
		return "", false
	}

	return invocation.arguments[index], true
}

// ArgumentCount returns the number of parsed positional arguments.
func (invocation Invocation) ArgumentCount() int {
	return len(invocation.arguments)
}

// HasArguments reports whether parsed positional arguments are present.
func (invocation Invocation) HasArguments() bool {
	return len(invocation.arguments) > 0
}

// RawArguments returns detached raw adapter arguments.
func (invocation Invocation) RawArguments() []string {
	return cloneStringSlice(invocation.rawArguments)
}

// RawArgument returns the raw adapter argument at index.
//
// The second return value is false when index is out of range. RawArgument
// never panics.
func (invocation Invocation) RawArgument(index int) (string, bool) {
	if index < 0 || index >= len(invocation.rawArguments) {
		return "", false
	}

	return invocation.rawArguments[index], true
}

// RawArgumentCount returns the number of raw adapter arguments.
func (invocation Invocation) RawArgumentCount() int {
	return len(invocation.rawArguments)
}

// HasRawArguments reports whether raw adapter arguments are present.
func (invocation Invocation) HasRawArguments() bool {
	return len(invocation.rawArguments) > 0
}

// WorkingDir returns the invocation working directory.
func (invocation Invocation) WorkingDir() string {
	return invocation.workingDir
}

// HasWorkingDir reports whether WorkingDir is set.
func (invocation Invocation) HasWorkingDir() bool {
	return invocation.workingDir != ""
}

// Env returns a detached copy of selected environment metadata.
func (invocation Invocation) Env() map[string]string {
	return cloneStringMap(invocation.env)
}

// EnvValue returns an environment metadata value and whether it exists.
func (invocation Invocation) EnvValue(name string) (string, bool) {
	value, ok := invocation.env[name]

	return value, ok
}

// HasEnv reports whether invocation contains selected environment metadata.
func (invocation Invocation) HasEnv() bool {
	return len(invocation.env) > 0
}

// HasEnvValue reports whether an environment metadata name exists.
func (invocation Invocation) HasEnvValue(name string) bool {
	_, ok := invocation.EnvValue(name)

	return ok
}

// EnvCount returns the number of selected environment metadata values.
func (invocation Invocation) EnvCount() int {
	return len(invocation.env)
}

// EnvNames returns environment metadata names in deterministic lexical order.
func (invocation Invocation) EnvNames() []string {
	names := make([]string, 0, len(invocation.env))
	for name := range invocation.env {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Fields returns a detached copy of invocation metadata fields.
func (invocation Invocation) Fields() map[string]string {
	return cloneStringMap(invocation.fields)
}

// Field returns one invocation metadata field and whether it exists.
func (invocation Invocation) Field(key string) (string, bool) {
	value, ok := invocation.fields[key]

	return value, ok
}

// HasFields reports whether invocation metadata fields are present.
func (invocation Invocation) HasFields() bool {
	return len(invocation.fields) > 0
}

// HasField reports whether an invocation metadata field exists.
func (invocation Invocation) HasField(key string) bool {
	_, ok := invocation.Field(key)

	return ok
}

// FieldCount returns the number of invocation metadata fields.
func (invocation Invocation) FieldCount() int {
	return len(invocation.fields)
}

// FieldKeys returns invocation metadata keys in deterministic lexical order.
func (invocation Invocation) FieldKeys() []string {
	keys := make([]string, 0, len(invocation.fields))
	for key := range invocation.fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// spec returns a detached construction spec.
func (invocation Invocation) spec() InvocationSpec {
	return InvocationSpec{
		Arguments:    cloneStringSlice(invocation.arguments),
		RawArguments: cloneStringSlice(invocation.rawArguments),
		WorkingDir:   invocation.workingDir,
		Env:          cloneStringMap(invocation.env),
		Fields:       cloneStringMap(invocation.fields),
	}
}
