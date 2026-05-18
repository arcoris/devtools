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
	"errors"
	"fmt"
	"time"
)

// runtimeCancellationError returns ErrRuntimeCanceled when ctx or err indicates
// cancellation.
func runtimeCancellationError(ctx context.Context, err error) error {
	if errors.Is(err, ErrRuntimeCanceled) {
		return err
	}

	if ctx != nil && ctx.Err() != nil {
		cause := context.Cause(ctx)
		if cause == nil {
			cause = ctx.Err()
		}

		return fmt.Errorf("%w: %w", ErrRuntimeCanceled, cause)
	}

	if isContextCancellation(err) {
		return fmt.Errorf("%w: %w", ErrRuntimeCanceled, err)
	}

	return nil
}

// isRuntimeCancellation reports whether err describes runtime cancellation.
func isRuntimeCancellation(err error) bool {
	return errors.Is(err, ErrRuntimeCanceled) || isContextCancellation(err)
}

// isContextCancellation reports whether err is a standard context cancellation
// or deadline error.
func isContextCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// mergeCanceledResult returns a valid canceled Result, preserving handler
// artifacts and warnings when the handler returned a valid partial result.
func (runtime Runtime) mergeCanceledResult(result Result, startedAt time.Time, finishedAt time.Time, cause error) Result {
	if cause == nil {
		cause = ErrRuntimeCanceled
	}

	if err := result.Validate(); err == nil {
		withTiming, timingErr := result.WithTiming(startedAt, finishedAt)
		if timingErr == nil {
			canceled, statusErr := withTiming.WithStatus(ResultStatusCanceled)
			if statusErr == nil {
				withMessage, messageErr := canceled.WithMessage(cause.Error())
				if messageErr == nil {
					return withMessage
				}

				return canceled
			}
		}
	}

	return runtime.canceledResult(startedAt, finishedAt, cause)
}
