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
	"time"
)

// Execute runs the runtime lifecycle for already parsed input values.
//
// Execute returns a Result even when the handler fails, panics, or the context
// is canceled. The returned error explains lifecycle-level failure to callers
// that need fail-fast behavior. CLI adapters can usually render the Result and
// use Result.RecommendedExitCode for process behavior.
func (runtime Runtime) Execute(ctx context.Context, spec RuntimeExecutionSpec) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := runtime.Validate(); err != nil {
		return Result{}, err
	}

	startedAt := runtime.now()

	if err := runtime.emit(ctx, RuntimeEventCommandStarted, startedAt, nil, nil); err != nil {
		return Result{}, err
	}

	bound, err := runtime.bind(ctx, spec, startedAt)
	if err != nil {
		result := runtime.failureResult(startedAt, runtime.now(), err)
		_ = runtime.emit(ctx, RuntimeEventCommandCompleted, runtime.now(), &result, err)

		return result, err
	}

	if err := ctx.Err(); err != nil {
		result := runtime.canceledResult(startedAt, runtime.now(), err)
		_ = runtime.emit(ctx, RuntimeEventCommandCompleted, runtime.now(), &result, err)

		return result, fmt.Errorf("%w: %w", ErrRuntimeCanceled, err)
	}

	request := RuntimeRequest{
		runtimeName: runtime.name,
		commandID:   runtime.commandID,
		binding:     runtime.binding,
		input:       bound,
		startedAt:   startedAt,
		metadata:    spec.Metadata,
	}

	result, err := runtime.executeHandler(ctx, request)
	finishedAt := runtime.now()

	if err != nil {
		result = runtime.mergeFailureResult(result, startedAt, finishedAt, err)
		_ = runtime.emit(ctx, RuntimeEventCommandCompleted, finishedAt, &result, err)

		return result, err
	}

	result, err = runtime.normalizeSuccessResult(result, startedAt, finishedAt)
	if err != nil {
		failed := runtime.failureResult(startedAt, finishedAt, err)
		_ = runtime.emit(ctx, RuntimeEventCommandCompleted, finishedAt, &failed, err)

		return failed, err
	}

	if err := runtime.emit(ctx, RuntimeEventCommandCompleted, finishedAt, &result, nil); err != nil {
		failed := runtime.failureResult(startedAt, finishedAt, err)

		return failed, err
	}

	return result, nil
}

// MustExecute runs Execute and panics on error.
func (runtime Runtime) MustExecute(ctx context.Context, spec RuntimeExecutionSpec) Result {
	result, err := runtime.Execute(ctx, spec)
	if err != nil {
		panic(err)
	}

	return result
}

// bind validates input values and emits binding lifecycle events.
func (runtime Runtime) bind(ctx context.Context, spec RuntimeExecutionSpec, startedAt time.Time) (BoundInput, error) {
	_ = startedAt

	now := runtime.now()
	if err := runtime.emit(ctx, RuntimeEventBindingStarted, now, nil, nil); err != nil {
		return BoundInput{}, err
	}

	bound, err := runtime.binding.Bind(BindingValueSpec{
		OptionValues:     spec.OptionValues,
		PositionalValues: spec.PositionalValues,
	})

	now = runtime.now()
	if err != nil {
		_ = runtime.emit(ctx, RuntimeEventBindingCompleted, now, nil, err)

		return BoundInput{}, fmt.Errorf("%w: %w", ErrRuntimeExecution, err)
	}

	if err := runtime.emit(ctx, RuntimeEventBindingCompleted, now, nil, nil); err != nil {
		return BoundInput{}, err
	}

	return bound, nil
}
