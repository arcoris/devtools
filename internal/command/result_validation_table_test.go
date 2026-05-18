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
	"strings"
	"testing"
	"time"
)

// TestMustResultPanicsForInvalidResult verifies fail-fast construction.
func TestMustResultPanicsForInvalidResult(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustResult did not panic")
		}
	}()

	_ = MustResult(ResultSpec{
		Status: ResultStatus("unknown"),
	})
}

// TestNewResultRejectsInvalidResult verifies result validation.
func TestNewResultRejectsInvalidResult(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(-time.Second)
	negativeExitCode := -1
	tooLargeExitCode := maxResultExitCode + 1

	tests := []struct {
		name string
		spec ResultSpec
		err  error
	}{
		{
			name: "invalid status",
			spec: ResultSpec{
				Status: ResultStatus("unknown"),
			},
			err: ErrInvalidResultStatus,
		},
		{
			name: "invalid message",
			spec: ResultSpec{
				Message: "bad\x00message",
			},
			err: ErrInvalidResult,
		},
		{
			name: "too long message",
			spec: ResultSpec{
				Message: strings.Repeat("x", maxResultMessageLength+1),
			},
			err: ErrInvalidResult,
		},
		{
			name: "finished before started",
			spec: ResultSpec{
				StartedAt:  startedAt,
				FinishedAt: finishedAt,
			},
			err: ErrInvalidResult,
		},
		{
			name: "negative exit code",
			spec: ResultSpec{
				ExitCode: &negativeExitCode,
			},
			err: ErrInvalidResult,
		},
		{
			name: "too large exit code",
			spec: ResultSpec{
				ExitCode: &tooLargeExitCode,
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid artifact",
			spec: ResultSpec{
				Artifacts: []Artifact{
					{},
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "duplicate artifact",
			spec: ResultSpec{
				Artifacts: []Artifact{
					resultTestArtifact("bench.report"),
					resultTestArtifact("bench.report"),
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid warning",
			spec: ResultSpec{
				Warnings: []ResultWarning{
					{},
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid field key",
			spec: ResultSpec{
				Fields: map[string]string{
					"Bad": "value",
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid field value",
			spec: ResultSpec{
				Fields: map[string]string{
					"mode": "bad\x00value",
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid metadata",
			spec: ResultSpec{
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
			err: ErrInvalidResult,
		},
		{
			name: "invalid visibility",
			spec: ResultSpec{
				Visibility: Visibility("private"),
			},
			err: ErrInvalidResult,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewResult(tt.spec)
			if err == nil {
				t.Fatalf("NewResult() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewResult() error = %v, want %v", err, tt.err)
			}
		})
	}
}
