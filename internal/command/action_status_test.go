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

func TestActionStatusValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status ActionStatus
		valid  bool
	}{
		{name: "ok", status: ActionStatusOK, valid: true},
		{name: "skipped", status: ActionStatusSkipped, valid: true},
		{name: "failed", status: ActionStatusFailed, valid: true},
		{name: "empty", status: "", valid: false},
		{name: "unknown", status: ActionStatus("unknown"), valid: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.status.Validate()
			if test.valid && err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}

			if !test.valid {
				if err == nil {
					t.Fatalf("Validate() returned nil error")
				}

				if !errors.Is(err, ErrInvalidActionResult) {
					t.Fatalf("Validate() error = %v, want ErrInvalidActionResult", err)
				}
			}
		})
	}
}

func TestActionStatusDefaultsAndText(t *testing.T) {
	t.Parallel()

	var status ActionStatus
	if !status.IsZero() {
		t.Fatalf("zero IsZero() = false, want true")
	}

	if got, want := status.OrDefault(), ActionStatusOK; got != want {
		t.Fatalf("zero OrDefault() = %q, want %q", got, want)
	}

	if got, want := ActionStatusSkipped.String(), "skipped"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if !ActionStatusFailed.IsKnown() {
		t.Fatalf("ActionStatusFailed.IsKnown() = false, want true")
	}
}
