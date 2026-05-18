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

// TestUsageLineValidateRejectsInvalidStoredLines verifies validation of raw
// stored values, including cases constructors normalize away.
func TestUsageLineValidateRejectsInvalidStoredLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		line       UsageLine
		want       error
		wantShared error
	}{
		{
			name: "empty",
			line: "",
			want: ErrEmptyUsage,
		},
		{
			name: "blank",
			line: "   ",
			want: ErrEmptyUsage,
		},
		{
			name: "not canonical",
			line: "bench   run",
			want: ErrInvalidUsage,
		},
		{
			name:       "tab",
			line:       "bench\trun",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "newline",
			line:       "bench\nrun",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "control rune",
			line:       "bench\x00run",
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "invalid utf8",
			line:       UsageLine(string([]byte{0xff, 0xfe})),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "too long line",
			line:       UsageLine(strings.Repeat("x", maxUsageLineLength+1)),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name:       "too long token",
			line:       UsageLine(strings.Repeat("x", maxUsageTokenLength+1)),
			want:       ErrInvalidUsage,
			wantShared: textvalidate.ErrInvalidTokenText,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.line.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.want)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("Validate() error = %v, want shared sentinel %v", err, tt.wantShared)
			}
		})
	}
}

// TestUsageValidateRejectsInvalidStoredUsage verifies aggregate validation on
// manually constructed values.
func TestUsageValidateRejectsInvalidStoredUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		usage Usage
		want  error
	}{
		{
			name:  "empty",
			usage: Usage{},
			want:  ErrEmptyUsage,
		},
		{
			name: "invalid syntax",
			usage: Usage{
				syntax: "bench\nrun",
			},
			want: ErrInvalidUsage,
		},
		{
			name: "invalid alternative",
			usage: Usage{
				syntax:       "bench run",
				alternatives: []UsageLine{"bench\nsmoke"},
			},
			want: ErrInvalidUsage,
		},
		{
			name: "duplicate alternative",
			usage: Usage{
				syntax:       "bench run",
				alternatives: []UsageLine{"bench smoke", "bench smoke"},
			},
			want: ErrInvalidUsage,
		},
		{
			name: "alternative duplicates primary",
			usage: Usage{
				syntax:       "bench run",
				alternatives: []UsageLine{"bench run"},
			},
			want: ErrInvalidUsage,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.usage.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// TestUsageNormalizationHelper verifies the local canonical form.
func TestUsageNormalizationHelper(t *testing.T) {
	t.Parallel()

	if got, want := normalizeUsageLine("  bench   run   [flags]  "), "bench run [flags]"; got != want {
		t.Fatalf("normalizeUsageLine() = %q, want %q", got, want)
	}
}
