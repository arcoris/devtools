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

// OrDefault returns ResultStatusOK when status is zero.
func (status ResultStatus) OrDefault() ResultStatus {
	if status == "" {
		return ResultStatusOK
	}

	return status
}

// String returns the canonical status string.
func (status ResultStatus) String() string {
	return string(status)
}

// IsZero reports whether status has not been set.
func (status ResultStatus) IsZero() bool {
	return status == ""
}

// IsKnown reports whether status is a supported non-zero result status.
func (status ResultStatus) IsKnown() bool {
	switch status {
	case ResultStatusOK,
		ResultStatusSkipped,
		ResultStatusFailed,
		ResultStatusCanceled:
		return true
	default:
		return false
	}
}

// Validate verifies that status is supported and non-zero.
func (status ResultStatus) Validate() error {
	if status == "" {
		return fmt.Errorf("%w: status is empty", ErrInvalidResultStatus)
	}

	if status.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported status %q", ErrInvalidResultStatus, status)
}
