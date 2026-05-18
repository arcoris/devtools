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

// TestEventSeverityPredicates verifies severity predicates.
func TestEventSeverityPredicates(t *testing.T) {
	t.Parallel()

	if !MustSimpleEvent(EventKindDiagnostic, "").MustWithSeverity(EventSeverityTrace).IsTrace() {
		t.Fatalf("IsTrace() = false, want true")
	}

	if !MustSimpleEvent(EventKindDiagnostic, "").MustWithSeverity(EventSeverityDebug).IsDebug() {
		t.Fatalf("IsDebug() = false, want true")
	}

	if !MustSimpleEvent(EventKindDiagnostic, "").IsInfo() {
		t.Fatalf("IsInfo() = false, want true")
	}

	if !MustSimpleEvent(EventKindWarning, "").MustWithSeverity(EventSeverityWarning).IsWarning() {
		t.Fatalf("IsWarning() = false, want true")
	}

	if !MustSimpleEvent(EventKindDiagnostic, "").MustWithSeverity(EventSeverityError).IsError() {
		t.Fatalf("IsError() = false, want true")
	}
}

// TestEventSeverityValidation verifies EventSeverity behavior.
func TestEventSeverityValidation(t *testing.T) {
	t.Parallel()

	valid := []EventSeverity{
		EventSeverityTrace,
		EventSeverityDebug,
		EventSeverityInfo,
		EventSeverityWarning,
		EventSeverityError,
	}

	for _, severity := range valid {
		severity := severity

		t.Run(severity.String(), func(t *testing.T) {
			t.Parallel()

			if !severity.IsKnown() {
				t.Fatalf("IsKnown() = false, want true")
			}

			if err := severity.Validate(); err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}
		})
	}

	if got, want := EventSeverity("").OrDefault(), EventSeverityInfo; got != want {
		t.Fatalf("OrDefault() = %q, want %q", got, want)
	}

	if err := EventSeverity("").Validate(); !errors.Is(err, ErrInvalidEventSeverity) {
		t.Fatalf("zero Validate() error = %v, want ErrInvalidEventSeverity", err)
	}

	if err := EventSeverity("fatal").Validate(); !errors.Is(err, ErrInvalidEventSeverity) {
		t.Fatalf("fatal Validate() error = %v, want ErrInvalidEventSeverity", err)
	}
}
