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

	"arcoris.dev/devtools/internal/textvalidate"
)

// TestNewUsageLineAcceptsValidLines verifies constructor normalization and
// representative command syntax forms.
func TestNewUsageLineAcceptsValidLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "simple",
			raw:  "check [flags]",
			want: "check [flags]",
		},
		{
			name: "subcommand",
			raw:  "bench run [flags]",
			want: "bench run [flags]",
		},
		{
			name: "positional arguments",
			raw:  "bench compare <old> <new> [flags]",
			want: "bench compare <old> <new> [flags]",
		},
		{
			name: "double dash passthrough",
			raw:  "profile cpu -- <go test args>",
			want: "profile cpu -- <go test args>",
		},
		{
			name: "spaces normalized",
			raw:  "  bench   run   [flags]  ",
			want: "bench run [flags]",
		},
		{
			name: "unicode display token",
			raw:  "notes print <имя>",
			want: "notes print <имя>",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			line, err := NewUsageLine(tt.raw)
			if err != nil {
				t.Fatalf("NewUsageLine(%q) returned unexpected error: %v", tt.raw, err)
			}

			if got := line.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}

			if !line.IsValid() {
				t.Fatalf("IsValid() = false, want true")
			}
		})
	}
}

// TestNewUsageLineRejectsInvalidRawInput verifies that constructor input is
// checked before normalization can hide unsupported characters.
func TestNewUsageLineRejectsInvalidRawInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		raw        string
		want       error
		wantShared error
	}{
		{
			name: "empty",
			raw:  "",
			want: ErrEmptyUsage,
		},
		{
			name: "blank",
			raw:  "   ",
			want: ErrEmptyUsage,
		},
		{
			name:       "tab",
			raw:        "bench\trun",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "newline",
			raw:        "bench\nrun",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "carriage return",
			raw:        "bench\rrun",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "control rune",
			raw:        "bench\x00run",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "invalid utf8",
			raw:        string([]byte{0xff, 0xfe}),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "too long line",
			raw:        strings.Repeat("x", maxUsageLineLength+1),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "too long token",
			raw:        strings.Repeat("x", maxUsageTokenLength+1),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidTokenText,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewUsageLine(tt.raw)
			if err == nil {
				t.Fatalf("NewUsageLine(%q) returned nil error", tt.raw)
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("NewUsageLine(%q) error = %v, want %v", tt.raw, err, tt.want)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("NewUsageLine(%q) error = %v, want shared sentinel %v", tt.raw, err, tt.wantShared)
			}
		})
	}
}

// TestMustUsageLinePanicsForInvalidLine verifies fail-fast static construction.
func TestMustUsageLinePanicsForInvalidLine(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustUsageLine("bench\nrun")
	})
}
