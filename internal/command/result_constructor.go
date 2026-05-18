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

// NewResult validates spec and returns a Result.
func NewResult(spec ResultSpec) (Result, error) {
	status := spec.Status.OrDefault()
	visibility := spec.Visibility.OrDefault()

	var exitCode int
	var hasExit bool
	if spec.ExitCode != nil {
		exitCode = *spec.ExitCode
		hasExit = true
	}

	result := Result{
		status:     status,
		message:    normalizeResultBlock(spec.Message),
		startedAt:  spec.StartedAt,
		finishedAt: spec.FinishedAt,
		exitCode:   exitCode,
		hasExit:    hasExit,
		artifacts:  cloneResultArtifacts(spec.Artifacts),
		warnings:   cloneResultWarnings(spec.Warnings),
		fields:     cloneResultStringMap(spec.Fields),
		metadata:   spec.Metadata,
		visibility: visibility,
	}

	if err := result.Validate(); err != nil {
		return Result{}, err
	}

	return result, nil
}

// MustResult validates spec and returns a Result.
//
// MustResult panics on invalid input. It is intended for tests and controlled
// static wiring.
func MustResult(spec ResultSpec) Result {
	result, err := NewResult(spec)
	if err != nil {
		panic(err)
	}

	return result
}

// OKResult returns a successful result with an optional message.
func OKResult(message string) Result {
	return MustResult(ResultSpec{
		Status:  ResultStatusOK,
		Message: message,
	})
}

// SkippedResult returns a skipped result with an optional message.
func SkippedResult(message string) Result {
	return MustResult(ResultSpec{
		Status:  ResultStatusSkipped,
		Message: message,
	})
}

// FailedResult returns a failed result with an optional message.
func FailedResult(message string) Result {
	return MustResult(ResultSpec{
		Status:  ResultStatusFailed,
		Message: message,
	})
}

// CanceledResult returns a canceled result with an optional message.
func CanceledResult(message string) Result {
	return MustResult(ResultSpec{
		Status:  ResultStatusCanceled,
		Message: message,
	})
}
