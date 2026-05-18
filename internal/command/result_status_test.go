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
	"errors"
	"testing"
)

// TestResultExplicitExitCodeOverridesRecommended verifies explicit exit code behavior.
func TestResultExplicitExitCodeOverridesRecommended(t *testing.T) {
	t.Parallel()

	result := FailedResult("failed").MustWithExitCode(42)

	if got, want := result.RecommendedExitCode(), 42; got != want {
		t.Fatalf("RecommendedExitCode() = %d, want %d", got, want)
	}

	without := result.WithoutExitCode()

	if _, ok := without.ExitCode(); ok {
		t.Fatalf("WithoutExitCode() still has explicit exit code")
	}

	if got, want := without.RecommendedExitCode(), 1; got != want {
		t.Fatalf("RecommendedExitCode() after WithoutExitCode() = %d, want %d", got, want)
	}
}

// TestResultPredicates verifies result status predicates.
func TestResultPredicates(t *testing.T) {
	t.Parallel()

	ok := OKResult("ok")
	if !ok.IsOK() || !ok.IsSuccessful() || ok.IsUnsuccessful() {
		t.Fatalf("ok predicates are invalid")
	}

	skipped := SkippedResult("skipped")
	if !skipped.IsSkipped() || !skipped.IsSuccessful() || skipped.IsUnsuccessful() {
		t.Fatalf("skipped predicates are invalid")
	}

	failed := FailedResult("failed")
	if !failed.IsFailed() || failed.IsSuccessful() || !failed.IsUnsuccessful() {
		t.Fatalf("failed predicates are invalid")
	}

	canceled := CanceledResult("canceled")
	if !canceled.IsCanceled() || canceled.IsSuccessful() || !canceled.IsUnsuccessful() {
		t.Fatalf("canceled predicates are invalid")
	}
}

// TestResultStatusValidation verifies status helpers.
func TestResultStatusValidation(t *testing.T) {
	t.Parallel()

	valid := []ResultStatus{
		ResultStatusOK,
		ResultStatusSkipped,
		ResultStatusFailed,
		ResultStatusCanceled,
	}

	for _, status := range valid {
		status := status

		t.Run(status.String(), func(t *testing.T) {
			t.Parallel()

			if !status.IsKnown() {
				t.Fatalf("IsKnown() = false, want true")
			}

			if err := status.Validate(); err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}
		})
	}

	if got, want := ResultStatus("").OrDefault(), ResultStatusOK; got != want {
		t.Fatalf("OrDefault() = %q, want %q", got, want)
	}

	if err := ResultStatus("").Validate(); !errors.Is(err, ErrInvalidResultStatus) {
		t.Fatalf("zero Validate() error = %v, want ErrInvalidResultStatus", err)
	}

	if err := ResultStatus("bad").Validate(); !errors.Is(err, ErrInvalidResultStatus) {
		t.Fatalf("bad Validate() error = %v, want ErrInvalidResultStatus", err)
	}
}
