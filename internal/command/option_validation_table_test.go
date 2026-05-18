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

// TestNewOptionRejectsInvalidOption verifies option declaration validation.
func TestNewOptionRejectsInvalidOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec OptionSpec
		err  error
	}{
		{
			name: "empty name",
			spec: OptionSpec{
				Kind: OptionKindString,
			},
			err: ErrEmptyOptionName,
		},
		{
			name: "invalid name",
			spec: OptionSpec{
				Name: "Output",
				Kind: OptionKindString,
			},
			err: ErrInvalidOptionName,
		},
		{
			name: "invalid alias",
			spec: OptionSpec{
				Name:    "output",
				Aliases: []string{"Out"},
				Kind:    OptionKindString,
			},
			err: ErrInvalidOption,
		},
		{
			name: "duplicate alias with name",
			spec: OptionSpec{
				Name:    "output",
				Aliases: []string{"output"},
				Kind:    OptionKindString,
			},
			err: ErrInvalidOption,
		},
		{
			name: "duplicate alias",
			spec: OptionSpec{
				Name:    "output",
				Aliases: []string{"out", "out"},
				Kind:    OptionKindString,
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid shorthand",
			spec: OptionSpec{
				Name:      "output",
				Shorthand: "out",
				Kind:      OptionKindString,
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid kind",
			spec: OptionSpec{
				Name: "output",
				Kind: OptionKind("path"),
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid policy",
			spec: OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Policy: MustOptionPolicy(OptionPolicySpec{
					AllowedSources: []OptionSource{OptionSourceCommandLine},
				}),
				DefaultValues: []string{"out.txt"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid visibility",
			spec: OptionSpec{
				Name:       "output",
				Kind:       OptionKindString,
				Visibility: Visibility("private"),
			},
			err: ErrInvalidOption,
		},
		{
			name: "enum without allowed values",
			spec: OptionSpec{
				Name: "format",
				Kind: OptionKindEnum,
			},
			err: ErrInvalidOption,
		},
		{
			name: "numeric with allowed values",
			spec: OptionSpec{
				Name:          "count",
				Kind:          OptionKindInt,
				AllowedValues: []string{"one"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "default not in allowed values",
			spec: OptionSpec{
				Name:          "format",
				Kind:          OptionKindEnum,
				AllowedValues: []string{"text", "json"},
				DefaultValues: []string{"xml"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "scalar with multiple defaults",
			spec: OptionSpec{
				Name:          "output",
				Kind:          OptionKindString,
				DefaultValues: []string{"a", "b"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid bool default",
			spec: OptionSpec{
				Name:          "verbose",
				Kind:          OptionKindBool,
				DefaultValues: []string{"yes"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid int default",
			spec: OptionSpec{
				Name:          "count",
				Kind:          OptionKindInt,
				DefaultValues: []string{"abc"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid duration default",
			spec: OptionSpec{
				Name:          "timeout",
				Kind:          OptionKindDuration,
				DefaultValues: []string{"soon"},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid documentation",
			spec: OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Documentation: Documentation{
					summary: "bad\x00summary",
				},
			},
			err: ErrInvalidOption,
		},
		{
			name: "invalid metadata",
			spec: OptionSpec{
				Name: "output",
				Kind: OptionKindString,
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
			err: ErrInvalidOption,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewOption(tt.spec)
			if err == nil {
				t.Fatalf("NewOption() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewOption() error = %v, want %v", err, tt.err)
			}
		})
	}
}
