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
	"fmt"
	"strings"

	"arcoris.dev/devtools/internal/textvalidate"
)

// Validate verifies runtime wiring.
func (runtime Runtime) Validate() error {
	if runtime.name == "" {
		return fmt.Errorf("%w: name must not be empty", ErrInvalidRuntime)
	}

	if err := validateRuntimeName(runtime.name); err != nil {
		return err
	}

	if runtime.commandID != "" {
		if err := runtime.commandID.Validate(); err != nil {
			return fmt.Errorf("%w: invalid command id: %w", ErrInvalidRuntime, err)
		}
	}

	if err := runtime.binding.Validate(); err != nil {
		return fmt.Errorf("%w: invalid binding: %w", ErrInvalidRuntime, err)
	}

	if runtime.handler == nil {
		return fmt.Errorf("%w: handler must not be nil", ErrInvalidRuntime)
	}

	if runtime.clock == nil {
		return fmt.Errorf("%w: clock must not be nil", ErrInvalidRuntime)
	}

	if err := runtime.options.Validate(); err != nil {
		return err
	}

	if err := runtime.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidRuntime, err)
	}

	if err := runtime.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidRuntime, err)
	}

	return nil
}

// validateRuntimeName validates runtime name text.
func validateRuntimeName(raw string) error {
	if raw == "" {
		return fmt.Errorf("%w: name must not be empty", ErrInvalidRuntime)
	}

	if raw != strings.TrimSpace(raw) {
		return fmt.Errorf("%w: name is not canonical", ErrInvalidRuntime)
	}

	if err := textvalidate.ValidateSingleLineText(raw, maxRuntimeNameLength); err != nil {
		return fmt.Errorf("%w: invalid name: %w", ErrInvalidRuntime, err)
	}

	return nil
}
