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
	"fmt"

	"arcoris.dev/devtools/internal/textvalidate"
)

var (
	// ErrEmptyCommandNameSegment reports that a command-name segment is empty.
	ErrEmptyCommandNameSegment = errors.New("command name segment is empty")

	// ErrInvalidCommandNameSegment reports that a command-name segment is malformed.
	ErrInvalidCommandNameSegment = errors.New("command name segment is invalid")
)

// validateCommandNameSegment validates one reusable command-name segment.
//
// Command IDs, command path elements, option names, group keys, and related
// metadata all use the same compact ASCII kebab-case grammar. The shared
// character rules live in internal/textvalidate; this adapter preserves
// command-specific sentinel errors for callers.
func validateCommandNameSegment(segment string) error {
	err := textvalidate.ValidateKebabSegment(segment)
	if err == nil {
		return nil
	}

	if errors.Is(err, textvalidate.ErrEmptyKebabSegment) {
		return fmt.Errorf("%w: %w", ErrEmptyCommandNameSegment, err)
	}

	return fmt.Errorf("%w: %w", ErrInvalidCommandNameSegment, err)
}
