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

// WithTiming returns a validated copy with StartedAt and FinishedAt replaced.
func (result Result) WithTiming(startedAt time.Time, finishedAt time.Time) (Result, error) {
	spec := result.spec()
	spec.StartedAt = startedAt
	spec.FinishedAt = finishedAt

	return NewResult(spec)
}

// MustWithTiming returns a validated copy with StartedAt and FinishedAt replaced
// and panics on invalid input.
func (result Result) MustWithTiming(startedAt time.Time, finishedAt time.Time) Result {
	next, err := result.WithTiming(startedAt, finishedAt)
	if err != nil {
		panic(err)
	}

	return next
}

// WithExitCode returns a validated copy with explicit ExitCode set.
func (result Result) WithExitCode(exitCode int) (Result, error) {
	spec := result.spec()
	spec.ExitCode = &exitCode

	return NewResult(spec)
}

// MustWithExitCode returns a validated copy with explicit ExitCode set and
// panics on invalid input.
func (result Result) MustWithExitCode(exitCode int) Result {
	next, err := result.WithExitCode(exitCode)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutExitCode returns a validated copy without explicit ExitCode.
func (result Result) WithoutExitCode() Result {
	spec := result.spec()
	spec.ExitCode = nil

	return MustResult(spec)
}
