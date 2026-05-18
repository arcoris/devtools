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

// TestNewArgumentRejectsInvalidArgument verifies declaration validation.
func TestNewArgumentRejectsInvalidArgument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec ArgumentSpec
		err  error
	}{
		{
			name: "empty name",
			spec: ArgumentSpec{
				Kind: OptionKindString,
			},
			err: ErrEmptyArgumentName,
		},
		{
			name: "invalid name",
			spec: ArgumentSpec{
				Name: "Package",
				Kind: OptionKindString,
			},
			err: ErrInvalidArgumentName,
		},
		{
			name: "invalid kind",
			spec: ArgumentSpec{
				Name: "package",
				Kind: OptionKind("path"),
			},
			err: ErrInvalidArgument,
		},
		{
			name: "list kind",
			spec: ArgumentSpec{
				Name: "package",
				Kind: OptionKindStringList,
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid requirement",
			spec: ArgumentSpec{
				Name:        "package",
				Kind:        OptionKindString,
				Requirement: ArgumentRequirement("mandatory"),
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid cardinality",
			spec: ArgumentSpec{
				Name:        "package",
				Kind:        OptionKindString,
				Cardinality: ArgumentCardinality("many"),
			},
			err: ErrInvalidArgument,
		},
		{
			name: "bool allows empty",
			spec: ArgumentSpec{
				Name:       "enabled",
				Kind:       OptionKindBool,
				EmptyValue: OptionEmptyValueAllow,
			},
			err: ErrInvalidArgument,
		},
		{
			name: "int allows empty",
			spec: ArgumentSpec{
				Name:       "count",
				Kind:       OptionKindInt,
				EmptyValue: OptionEmptyValueAllow,
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid metavar",
			spec: ArgumentSpec{
				Name:    "package",
				Kind:    OptionKindString,
				Metavar: "PACKAGE PATTERN",
			},
			err: ErrInvalidArgument,
		},
		{
			name: "enum without allowed values",
			spec: ArgumentSpec{
				Name: "format",
				Kind: OptionKindEnum,
			},
			err: ErrInvalidArgument,
		},
		{
			name: "numeric with allowed values",
			spec: ArgumentSpec{
				Name:          "count",
				Kind:          OptionKindInt,
				AllowedValues: []string{"one"},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "required with default",
			spec: ArgumentSpec{
				Name:          "package",
				Kind:          OptionKindString,
				DefaultValues: []string{"./..."},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "single with multiple defaults",
			spec: ArgumentSpec{
				Name:          "package",
				Kind:          OptionKindString,
				Requirement:   ArgumentRequirementOptional,
				DefaultValues: []string{"./...", "./internal/..."},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "default not in allowed values",
			spec: ArgumentSpec{
				Name:          "format",
				Kind:          OptionKindEnum,
				Requirement:   ArgumentRequirementOptional,
				AllowedValues: []string{"text", "json"},
				DefaultValues: []string{"xml"},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid bool default",
			spec: ArgumentSpec{
				Name:          "enabled",
				Kind:          OptionKindBool,
				Requirement:   ArgumentRequirementOptional,
				DefaultValues: []string{"yes"},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid documentation",
			spec: ArgumentSpec{
				Name: "package",
				Kind: OptionKindString,
				Documentation: Documentation{
					summary: "bad\x00summary",
				},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid metadata",
			spec: ArgumentSpec{
				Name: "package",
				Kind: OptionKindString,
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
			err: ErrInvalidArgument,
		},
		{
			name: "invalid visibility",
			spec: ArgumentSpec{
				Name:       "package",
				Kind:       OptionKindString,
				Visibility: Visibility("private"),
			},
			err: ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewArgument(tt.spec)
			if err == nil {
				t.Fatalf("NewArgument() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewArgument() error = %v, want %v", err, tt.err)
			}
		})
	}
}
