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
)

// TestMustEventPanicsForInvalidEvent verifies fail-fast construction.
func TestMustEventPanicsForInvalidEvent(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustEvent did not panic")
		}
	}()

	_ = MustEvent(EventSpec{})
}

// TestNewEventRejectsInvalidEvent verifies event validation.
func TestNewEventRejectsInvalidEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec EventSpec
		err  error
	}{
		{
			name: "invalid id",
			spec: EventSpec{
				ID:   "Command.Started",
				Kind: EventKindCommandStarted,
			},
			err: ErrInvalidEventID,
		},
		{
			name: "empty kind",
			spec: EventSpec{},
			err:  ErrEmptyEventKind,
		},
		{
			name: "invalid kind",
			spec: EventSpec{
				Kind: EventKind("Command.Started"),
			},
			err: ErrInvalidEventKind,
		},
		{
			name: "invalid severity",
			spec: EventSpec{
				Kind:     EventKindCommandStarted,
				Severity: EventSeverity("fatal"),
			},
			err: ErrInvalidEventSeverity,
		},
		{
			name: "invalid command id",
			spec: EventSpec{
				Kind:      EventKindCommandStarted,
				CommandID: ID("Bad.Command"),
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid message",
			spec: EventSpec{
				Kind:    EventKindCommandStarted,
				Message: "bad\x00message",
			},
			err: ErrInvalidEvent,
		},
		{
			name: "too long message",
			spec: EventSpec{
				Kind:    EventKindCommandStarted,
				Message: strings.Repeat("x", maxEventMessageLength+1),
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid field key",
			spec: EventSpec{
				Kind: EventKindCommandStarted,
				Fields: map[string]string{
					"Bad": "value",
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid field value",
			spec: EventSpec{
				Kind: EventKindCommandStarted,
				Fields: map[string]string{
					"duration": "bad\x00value",
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid artifact",
			spec: EventSpec{
				Kind: EventKindArtifactProduced,
				Artifacts: []Artifact{
					{},
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "duplicate artifact",
			spec: EventSpec{
				Kind: EventKindArtifactProduced,
				Artifacts: []Artifact{
					eventTestArtifact("bench.report"),
					eventTestArtifact("bench.report"),
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid result",
			spec: EventSpec{
				Kind: EventKindResultProduced,
				Result: &Result{
					status: ResultStatus("bad"),
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid label",
			spec: EventSpec{
				Kind:   EventKindCommandStarted,
				Labels: []string{"Bad"},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "duplicate label",
			spec: EventSpec{
				Kind:   EventKindCommandStarted,
				Labels: []string{"bench", "bench"},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid metadata",
			spec: EventSpec{
				Kind: EventKindCommandStarted,
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
			err: ErrInvalidEvent,
		},
		{
			name: "invalid visibility",
			spec: EventSpec{
				Kind:       EventKindCommandStarted,
				Visibility: Visibility("private"),
			},
			err: ErrInvalidEvent,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewEvent(tt.spec)
			if err == nil {
				t.Fatalf("NewEvent() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewEvent() error = %v, want %v", err, tt.err)
			}
		})
	}
}
