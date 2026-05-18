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

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

// executeHandler emits action lifecycle events and executes the handler.
func (runtime Runtime) executeHandler(ctx context.Context, request RuntimeRequest) (result Result, err error) {
	if err := runtime.emit(ctx, RuntimeEventActionStarted, runtime.now(), nil, nil); err != nil {
		return Result{}, err
	}

	if runtime.options.RecoverPanics() {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = runtime.recoverPanicError(recovered)
				result = runtime.failureResult(request.startedAt, runtime.now(), err)

				_ = runtime.emit(ctx, RuntimeEventActionCompleted, runtime.now(), &result, err)
			}
		}()
	}

	result, err = runtime.handler.Run(ctx, request)

	if cancelErr := runtimeCancellationError(ctx, err); cancelErr != nil {
		result = runtime.mergeCanceledResult(result, request.startedAt, runtime.now(), cancelErr)
		_ = runtime.emit(ctx, RuntimeEventActionCompleted, runtime.now(), &result, cancelErr)

		return result, cancelErr
	}

	if err != nil {
		_ = runtime.emit(ctx, RuntimeEventActionCompleted, runtime.now(), &result, err)

		return result, fmt.Errorf("%w: %w", ErrRuntimeExecution, err)
	}

	if err := runtime.emit(ctx, RuntimeEventActionCompleted, runtime.now(), &result, nil); err != nil {
		return Result{}, err
	}

	return result, nil
}

// normalizeSuccessResult attaches runtime timing and validates a successful
// handler result.
func (runtime Runtime) normalizeSuccessResult(result Result, startedAt time.Time, finishedAt time.Time) (Result, error) {
	if result.Status().IsZero() {
		result = OKResult("")
	}

	result, err := result.WithTiming(startedAt, finishedAt)
	if err != nil {
		return Result{}, fmt.Errorf("%w: invalid result timing: %w", ErrRuntimeExecution, err)
	}

	if err := result.Validate(); err != nil {
		return Result{}, fmt.Errorf("%w: invalid result: %w", ErrRuntimeExecution, err)
	}

	return result, nil
}

// mergeFailureResult returns a valid failure result, preserving a valid
// handler-produced result when possible.
func (runtime Runtime) mergeFailureResult(result Result, startedAt time.Time, finishedAt time.Time, cause error) Result {
	if cause == nil {
		cause = ErrRuntimeExecution
	}

	if err := result.Validate(); err == nil {
		withTiming, timingErr := result.WithTiming(startedAt, finishedAt)
		if timingErr == nil {
			if withTiming.IsUnsuccessful() {
				return withTiming
			}

			failed, failedErr := withTiming.WithStatus(ResultStatusFailed)
			if failedErr == nil {
				return failed
			}
		}
	}

	return runtime.failureResult(startedAt, finishedAt, cause)
}

// failureResult builds a valid failed Result from an error.
func (runtime Runtime) failureResult(startedAt time.Time, finishedAt time.Time, cause error) Result {
	message := "Runtime execution failed."
	if cause != nil {
		message = cause.Error()
	}

	result, err := NewResult(ResultSpec{
		Status:     ResultStatusFailed,
		Message:    message,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		ExitCode:   intPointer(1),
	})
	if err == nil {
		return result
	}

	return FailedResult("Runtime execution failed.")
}

// canceledResult builds a valid canceled Result from a cancellation cause.
func (runtime Runtime) canceledResult(startedAt time.Time, finishedAt time.Time, cause error) Result {
	message := "Runtime execution canceled."
	if cause != nil {
		message = cause.Error()
	}

	result, err := NewResult(ResultSpec{
		Status:     ResultStatusCanceled,
		Message:    message,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		ExitCode:   intPointer(130),
	})
	if err == nil {
		return result
	}

	return CanceledResult("Runtime execution canceled.")
}

// recoverPanicError converts a panic value into an error.
func (runtime Runtime) recoverPanicError(recovered any) error {
	message := fmt.Sprint(recovered)
	message = normalizeRuntimePanicMessage(message)

	if runtime.options.IncludePanicStack() {
		stack := string(debug.Stack())
		stack = normalizeRuntimePanicMessage(stack)
		if stack != "" {
			message = message + "\n" + stack
		}
	}

	return fmt.Errorf("%w: %s", ErrRuntimePanic, message)
}
