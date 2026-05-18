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

import "fmt"

// OptionOccurrence describes how many times an option may appear in one
// invocation.
type OptionOccurrence string

const (
	// OptionOccurrenceSingle means the option may appear at most once.
	OptionOccurrenceSingle OptionOccurrence = "single"

	// OptionOccurrenceMultiple means the option may appear multiple times.
	OptionOccurrenceMultiple OptionOccurrence = "multiple"
)

// NewOptionOccurrence validates raw and returns it as an OptionOccurrence.
func NewOptionOccurrence(raw string) (OptionOccurrence, error) {
	occurrence := OptionOccurrence(raw)
	if err := occurrence.Validate(); err != nil {
		return "", err
	}

	return occurrence, nil
}

// ParseOptionOccurrence is an alias for NewOptionOccurrence.
func ParseOptionOccurrence(raw string) (OptionOccurrence, error) {
	return NewOptionOccurrence(raw)
}

// MustOptionOccurrence validates raw and returns it as an OptionOccurrence.
func MustOptionOccurrence(raw string) OptionOccurrence {
	occurrence, err := NewOptionOccurrence(raw)
	if err != nil {
		panic(err)
	}

	return occurrence
}

// String returns the canonical string representation of the occurrence policy.
func (occurrence OptionOccurrence) String() string {
	return string(occurrence)
}

// IsZero reports whether the occurrence policy has not been set.
func (occurrence OptionOccurrence) IsZero() bool {
	return occurrence == ""
}

// OrDefaultForKind returns a default occurrence policy for kind.
func (occurrence OptionOccurrence) OrDefaultForKind(kind OptionKind) OptionOccurrence {
	if !occurrence.IsZero() {
		return occurrence
	}

	if kind.IsList() {
		return OptionOccurrenceMultiple
	}

	return OptionOccurrenceSingle
}

// IsKnown reports whether occurrence is one of the supported non-zero states.
func (occurrence OptionOccurrence) IsKnown() bool {
	switch occurrence {
	case OptionOccurrenceSingle, OptionOccurrenceMultiple:
		return true
	default:
		return false
	}
}

// IsValid reports whether occurrence satisfies policy rules.
func (occurrence OptionOccurrence) IsValid() bool {
	return occurrence.Validate() == nil
}

// Validate verifies that occurrence is a supported non-zero state.
func (occurrence OptionOccurrence) Validate() error {
	if occurrence == "" {
		return fmt.Errorf("%w: occurrence is empty", ErrInvalidOptionPolicy)
	}

	if occurrence.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported occurrence %q", ErrInvalidOptionPolicy, occurrence)
}

// IsSingle reports whether the option may appear at most once.
func (occurrence OptionOccurrence) IsSingle() bool {
	return occurrence == OptionOccurrenceSingle
}

// IsMultiple reports whether the option may appear multiple times.
func (occurrence OptionOccurrence) IsMultiple() bool {
	return occurrence == OptionOccurrenceMultiple
}
