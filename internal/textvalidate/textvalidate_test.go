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

package textvalidate

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidateKebabSegmentAcceptsValidSegments(t *testing.T) {
	t.Parallel()

	tests := []string{
		"a",
		"abc",
		"a1",
		"a-1",
		"bench",
		"run",
		"release-notes",
		"generate2",
		"config-validate",
	}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			if err := ValidateKebabSegment(segment); err != nil {
				t.Fatalf("ValidateKebabSegment(%q) returned unexpected error: %v", segment, err)
			}
		})
	}
}

func TestValidateKebabSegmentRejectsInvalidSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		segment string
		wantErr error
	}{
		{
			name:    "empty",
			segment: "",
			wantErr: ErrEmptyKebabSegment,
		},
		{
			name:    "uppercase first letter",
			segment: "Bench",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "uppercase continuation",
			segment: "benchRun",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "starts with digit",
			segment: "1bench",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "starts with hyphen",
			segment: "-bench",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "trailing hyphen",
			segment: "bench-",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "repeated hyphen",
			segment: "bench--run",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "underscore",
			segment: "bench_run",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "space",
			segment: "bench run",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "slash",
			segment: "bench/run",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "dot",
			segment: "bench.run",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "unicode",
			segment: "ран",
			wantErr: ErrInvalidKebabSegment,
		},
		{
			name:    "emoji",
			segment: "bench🚀",
			wantErr: ErrInvalidKebabSegment,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateKebabSegment(tt.segment)
			if err == nil {
				t.Fatalf("ValidateKebabSegment(%q) returned nil error", tt.segment)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateKebabSegment(%q) error = %v, want %v", tt.segment, err, tt.wantErr)
			}
		})
	}
}

func TestValidateKebabSegmentRejectsInvalidUTF8(t *testing.T) {
	t.Parallel()

	segment := string([]byte{0xff, 0xfe, 0xfd})

	err := ValidateKebabSegment(segment)
	if err == nil {
		t.Fatalf("ValidateKebabSegment(invalid UTF-8) returned nil error")
	}

	if !errors.Is(err, ErrInvalidKebabSegment) {
		t.Fatalf("ValidateKebabSegment(invalid UTF-8) error = %v, want ErrInvalidKebabSegment", err)
	}
}

func TestValidateKebabSegmentRuneAcceptsValidRunes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		byteOffset int
		r          rune
	}{
		{
			name:       "first lowercase letter",
			byteOffset: 0,
			r:          'a',
		},
		{
			name:       "continuation lowercase letter",
			byteOffset: 1,
			r:          'z',
		},
		{
			name:       "continuation digit",
			byteOffset: 1,
			r:          '9',
		},
		{
			name:       "continuation hyphen",
			byteOffset: 1,
			r:          '-',
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := ValidateKebabSegmentRune(tt.byteOffset, tt.r); err != nil {
				t.Fatalf(
					"ValidateKebabSegmentRune(%d, %q) returned unexpected error: %v",
					tt.byteOffset,
					tt.r,
					err,
				)
			}
		})
	}
}

func TestValidateKebabSegmentRuneRejectsInvalidRunes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		byteOffset int
		r          rune
	}{
		{
			name:       "negative offset",
			byteOffset: -1,
			r:          'a',
		},
		{
			name:       "first uppercase letter",
			byteOffset: 0,
			r:          'A',
		},
		{
			name:       "first digit",
			byteOffset: 0,
			r:          '1',
		},
		{
			name:       "first hyphen",
			byteOffset: 0,
			r:          '-',
		},
		{
			name:       "continuation uppercase letter",
			byteOffset: 1,
			r:          'A',
		},
		{
			name:       "continuation underscore",
			byteOffset: 1,
			r:          '_',
		},
		{
			name:       "continuation dot",
			byteOffset: 1,
			r:          '.',
		},
		{
			name:       "continuation slash",
			byteOffset: 1,
			r:          '/',
		},
		{
			name:       "non ASCII rune",
			byteOffset: 1,
			r:          'я',
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateKebabSegmentRune(tt.byteOffset, tt.r)
			if err == nil {
				t.Fatalf("ValidateKebabSegmentRune(%d, %q) returned nil error", tt.byteOffset, tt.r)
			}

			if !errors.Is(err, ErrInvalidKebabSegment) {
				t.Fatalf(
					"ValidateKebabSegmentRune(%d, %q) error = %v, want ErrInvalidKebabSegment",
					tt.byteOffset,
					tt.r,
					err,
				)
			}
		})
	}
}

func TestValidateKebabSegmentStartAcceptsLowercaseLetters(t *testing.T) {
	t.Parallel()

	for ch := byte('a'); ch <= byte('z'); ch++ {
		ch := ch

		t.Run(string(ch), func(t *testing.T) {
			t.Parallel()

			if err := ValidateKebabSegmentStart(ch); err != nil {
				t.Fatalf("ValidateKebabSegmentStart(%q) returned unexpected error: %v", ch, err)
			}
		})
	}
}

func TestValidateKebabSegmentStartRejectsNonLowercaseLetters(t *testing.T) {
	t.Parallel()

	tests := []byte{
		'A',
		'Z',
		'0',
		'9',
		'-',
		'_',
		'.',
		'/',
		' ',
	}

	for _, ch := range tests {
		ch := ch

		t.Run(string(ch), func(t *testing.T) {
			t.Parallel()

			err := ValidateKebabSegmentStart(ch)
			if err == nil {
				t.Fatalf("ValidateKebabSegmentStart(%q) returned nil error", ch)
			}

			if !errors.Is(err, ErrInvalidKebabSegment) {
				t.Fatalf("ValidateKebabSegmentStart(%q) error = %v, want ErrInvalidKebabSegment", ch, err)
			}
		})
	}
}

func TestValidateKebabSegmentContinuationAcceptsAllowedCharacters(t *testing.T) {
	t.Parallel()

	tests := []byte{
		'a',
		'z',
		'0',
		'9',
		'-',
	}

	for _, ch := range tests {
		ch := ch

		t.Run(string(ch), func(t *testing.T) {
			t.Parallel()

			if err := ValidateKebabSegmentContinuation(1, ch); err != nil {
				t.Fatalf("ValidateKebabSegmentContinuation(1, %q) returned unexpected error: %v", ch, err)
			}
		})
	}
}

func TestValidateKebabSegmentContinuationRejectsInvalidCharacters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		byteOffset int
		ch         byte
	}{
		{name: "zero offset", byteOffset: 0, ch: 'a'},
		{name: "uppercase", byteOffset: 1, ch: 'A'},
		{name: "underscore", byteOffset: 1, ch: '_'},
		{name: "dot", byteOffset: 1, ch: '.'},
		{name: "slash", byteOffset: 1, ch: '/'},
		{name: "space", byteOffset: 1, ch: ' '},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateKebabSegmentContinuation(tt.byteOffset, tt.ch)
			if err == nil {
				t.Fatalf("ValidateKebabSegmentContinuation(%d, %q) returned nil error", tt.byteOffset, tt.ch)
			}

			if !errors.Is(err, ErrInvalidKebabSegment) {
				t.Fatalf(
					"ValidateKebabSegmentContinuation(%d, %q) error = %v, want ErrInvalidKebabSegment",
					tt.byteOffset,
					tt.ch,
					err,
				)
			}
		})
	}
}

func TestIsKebabSegmentContinuation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ch   byte
		want bool
	}{
		{ch: 'a', want: true},
		{ch: 'z', want: true},
		{ch: '0', want: true},
		{ch: '9', want: true},
		{ch: '-', want: true},
		{ch: 'A', want: false},
		{ch: '_', want: false},
		{ch: '.', want: false},
		{ch: '/', want: false},
		{ch: ' ', want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(string(tt.ch), func(t *testing.T) {
			t.Parallel()

			if got := IsKebabSegmentContinuation(tt.ch); got != tt.want {
				t.Fatalf("IsKebabSegmentContinuation(%q) = %v, want %v", tt.ch, got, tt.want)
			}
		})
	}
}

func TestIsASCII(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "nul", r: 0, want: true},
		{name: "space", r: ' ', want: true},
		{name: "tilde", r: '~', want: true},
		{name: "delete", r: 127, want: true},
		{name: "negative", r: -1, want: false},
		{name: "unicode", r: 'я', want: false},
		{name: "replacement", r: '\ufffd', want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsASCII(tt.r); got != tt.want {
				t.Fatalf("IsASCII(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestIsASCIILowerLetter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ch   byte
		want bool
	}{
		{ch: 'a', want: true},
		{ch: 'm', want: true},
		{ch: 'z', want: true},
		{ch: 'A', want: false},
		{ch: 'Z', want: false},
		{ch: '0', want: false},
		{ch: '-', want: false},
		{ch: '_', want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(string(tt.ch), func(t *testing.T) {
			t.Parallel()

			if got := IsASCIILowerLetter(tt.ch); got != tt.want {
				t.Fatalf("IsASCIILowerLetter(%q) = %v, want %v", tt.ch, got, tt.want)
			}
		})
	}
}

func TestIsASCIIUpperLetter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ch   byte
		want bool
	}{
		{ch: 'A', want: true},
		{ch: 'M', want: true},
		{ch: 'Z', want: true},
		{ch: 'a', want: false},
		{ch: 'z', want: false},
		{ch: '0', want: false},
		{ch: '_', want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("%q", tt.ch), func(t *testing.T) {
			t.Parallel()

			if got := IsASCIIUpperLetter(tt.ch); got != tt.want {
				t.Fatalf("IsASCIIUpperLetter(%q) = %v, want %v", tt.ch, got, tt.want)
			}
		})
	}
}

func TestIsASCIIDigit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ch   byte
		want bool
	}{
		{ch: '0', want: true},
		{ch: '5', want: true},
		{ch: '9', want: true},
		{ch: 'a', want: false},
		{ch: 'z', want: false},
		{ch: 'A', want: false},
		{ch: '-', want: false},
		{ch: '_', want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(string(tt.ch), func(t *testing.T) {
			t.Parallel()

			if got := IsASCIIDigit(tt.ch); got != tt.want {
				t.Fatalf("IsASCIIDigit(%q) = %v, want %v", tt.ch, got, tt.want)
			}
		})
	}
}
