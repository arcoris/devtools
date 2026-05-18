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

import "fmt"

// ActionStatus describes the high-level non-error outcome of an action.
//
// Failures should normally be represented as errors, not as ActionStatusFailed.
// The failed status is still available for adapters or lifecycle layers that
// need to serialize an already-classified failure result.
type ActionStatus string

const (
	// ActionStatusOK means the action completed successfully.
	ActionStatusOK ActionStatus = "ok"

	// ActionStatusSkipped means the action intentionally did not run or had no
	// work to perform.
	ActionStatusSkipped ActionStatus = "skipped"

	// ActionStatusFailed means a failure has already been classified into a
	// structured result by a higher layer.
	ActionStatusFailed ActionStatus = "failed"
)

// String returns the stable text form of status.
func (status ActionStatus) String() string {
	return string(status)
}

// IsZero reports whether status has not been set.
func (status ActionStatus) IsZero() bool {
	return status == ""
}

// OrDefault returns ActionStatusOK when status is zero.
func (status ActionStatus) OrDefault() ActionStatus {
	if status.IsZero() {
		return ActionStatusOK
	}

	return status
}

// IsKnown reports whether status is a supported non-zero action status.
func (status ActionStatus) IsKnown() bool {
	switch status {
	case ActionStatusOK, ActionStatusSkipped, ActionStatusFailed:
		return true
	default:
		return false
	}
}

// Validate verifies that status is a supported non-zero action status.
func (status ActionStatus) Validate() error {
	if status == "" {
		return fmt.Errorf("%w: status is empty", ErrInvalidActionResult)
	}

	if status.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported status %q", ErrInvalidActionResult, status)
}
