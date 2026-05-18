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

// WithArguments returns a validated copy with parsed arguments replaced.
func (invocation Invocation) WithArguments(arguments ...string) (Invocation, error) {
	spec := invocation.spec()
	spec.Arguments = cloneStringSlice(arguments)

	return NewInvocation(spec)
}

// MustWithArguments returns a validated copy with parsed arguments replaced and
// panics on invalid input.
func (invocation Invocation) MustWithArguments(arguments ...string) Invocation {
	next, err := invocation.WithArguments(arguments...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutArguments returns a validated copy without parsed arguments.
func (invocation Invocation) WithoutArguments() Invocation {
	spec := invocation.spec()
	spec.Arguments = nil

	return MustInvocation(spec)
}

// WithRawArguments returns a validated copy with raw arguments replaced.
func (invocation Invocation) WithRawArguments(arguments ...string) (Invocation, error) {
	spec := invocation.spec()
	spec.RawArguments = cloneStringSlice(arguments)

	return NewInvocation(spec)
}

// MustWithRawArguments returns a validated copy with raw arguments replaced and
// panics on invalid input.
func (invocation Invocation) MustWithRawArguments(arguments ...string) Invocation {
	next, err := invocation.WithRawArguments(arguments...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutRawArguments returns a validated copy without raw adapter arguments.
func (invocation Invocation) WithoutRawArguments() Invocation {
	spec := invocation.spec()
	spec.RawArguments = nil

	return MustInvocation(spec)
}

// WithWorkingDir returns a validated copy with working directory replaced.
func (invocation Invocation) WithWorkingDir(workingDir string) (Invocation, error) {
	spec := invocation.spec()
	spec.WorkingDir = workingDir

	return NewInvocation(spec)
}

// MustWithWorkingDir returns a validated copy with working directory replaced
// and panics on invalid input.
func (invocation Invocation) MustWithWorkingDir(workingDir string) Invocation {
	next, err := invocation.WithWorkingDir(workingDir)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutWorkingDir returns a validated copy without working directory
// metadata.
func (invocation Invocation) WithoutWorkingDir() Invocation {
	spec := invocation.spec()
	spec.WorkingDir = ""

	return MustInvocation(spec)
}

// WithEnv returns a validated copy with environment metadata replaced.
func (invocation Invocation) WithEnv(env map[string]string) (Invocation, error) {
	spec := invocation.spec()
	spec.Env = cloneStringMap(env)

	return NewInvocation(spec)
}

// MustWithEnv returns a validated copy with environment metadata replaced and
// panics on invalid input.
func (invocation Invocation) MustWithEnv(env map[string]string) Invocation {
	next, err := invocation.WithEnv(env)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutEnv returns a validated copy without selected environment metadata.
func (invocation Invocation) WithoutEnv() Invocation {
	spec := invocation.spec()
	spec.Env = nil

	return MustInvocation(spec)
}

// WithEnvValue returns a validated copy with one environment metadata value
// added or replaced.
func (invocation Invocation) WithEnvValue(name string, value string) (Invocation, error) {
	spec := invocation.spec()
	if spec.Env == nil {
		spec.Env = make(map[string]string)
	}

	spec.Env[name] = value

	return NewInvocation(spec)
}

// MustWithEnvValue returns a validated copy with one environment metadata value
// added or replaced and panics on invalid input.
func (invocation Invocation) MustWithEnvValue(name string, value string) Invocation {
	next, err := invocation.WithEnvValue(name, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutEnvValue returns a validated copy without one environment metadata
// value.
func (invocation Invocation) WithoutEnvValue(name string) Invocation {
	spec := invocation.spec()
	delete(spec.Env, name)

	return MustInvocation(spec)
}

// WithFields returns a validated copy with invocation metadata fields replaced.
func (invocation Invocation) WithFields(fields map[string]string) (Invocation, error) {
	spec := invocation.spec()
	spec.Fields = cloneStringMap(fields)

	return NewInvocation(spec)
}

// MustWithFields returns a validated copy with invocation metadata fields
// replaced and panics on invalid input.
func (invocation Invocation) MustWithFields(fields map[string]string) Invocation {
	next, err := invocation.WithFields(fields)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutFields returns a validated copy without invocation metadata fields.
func (invocation Invocation) WithoutFields() Invocation {
	spec := invocation.spec()
	spec.Fields = nil

	return MustInvocation(spec)
}

// WithField returns a validated copy with one invocation metadata field added
// or replaced.
func (invocation Invocation) WithField(key string, value string) (Invocation, error) {
	spec := invocation.spec()
	if spec.Fields == nil {
		spec.Fields = make(map[string]string)
	}

	spec.Fields[key] = value

	return NewInvocation(spec)
}

// MustWithField returns a validated copy with one invocation metadata field
// added or replaced and panics on invalid input.
func (invocation Invocation) MustWithField(key string, value string) Invocation {
	next, err := invocation.WithField(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutField returns a validated copy without one invocation metadata field.
func (invocation Invocation) WithoutField(key string) Invocation {
	spec := invocation.spec()
	delete(spec.Fields, key)

	return MustInvocation(spec)
}
