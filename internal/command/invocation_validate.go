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

// Validate verifies invocation metadata structural rules.
func (invocation Invocation) Validate() error {
	for index, argument := range invocation.arguments {
		if err := validateInvocationText(fmt.Sprintf("argument %d", index), argument, maxInvocationArgumentLength); err != nil {
			return err
		}
	}

	for index, argument := range invocation.rawArguments {
		if err := validateInvocationText(fmt.Sprintf("raw argument %d", index), argument, maxInvocationRawArgumentLength); err != nil {
			return err
		}
	}

	if invocation.workingDir != "" {
		if strings.TrimSpace(invocation.workingDir) == "" {
			return fmt.Errorf("%w: working directory must not be blank", ErrInvalidInvocation)
		}

		if err := validateInvocationText("working directory", invocation.workingDir, maxInvocationWorkingDirLength); err != nil {
			return err
		}
	}

	if err := validateInvocationEnv(invocation.env); err != nil {
		return err
	}

	return validateInvocationFields(invocation.fields)
}

// validateInvocationEnv validates selected invocation environment metadata.
func validateInvocationEnv(env map[string]string) error {
	for name, value := range env {
		if err := validateInvocationEnvName(name); err != nil {
			return err
		}

		if err := validateInvocationText("environment "+name, value, maxInvocationEnvValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateInvocationEnvName validates one captured environment variable name.
func validateInvocationEnvName(name string) error {
	if err := textvalidate.ValidateEnvName(name, maxInvocationEnvNameLength); err != nil {
		return fmt.Errorf("%w: invalid environment name %q: %w", ErrInvalidInvocation, name, err)
	}

	return nil
}

// validateInvocationFields validates invocation metadata fields.
func validateInvocationFields(fields map[string]string) error {
	for key, value := range fields {
		if err := validateInvocationFieldKey(key); err != nil {
			return err
		}

		if err := validateInvocationText("field "+key, value, maxInvocationFieldValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateInvocationFieldKey validates a dot-separated invocation metadata key.
func validateInvocationFieldKey(raw string) error {
	if err := textvalidate.ValidateDottedKebabKey(raw, maxInvocationFieldKeyLength, maxInvocationFieldKeyDepth); err != nil {
		return fmt.Errorf("%w: invalid field key %q: %w", ErrInvalidInvocation, raw, err)
	}

	return nil
}

// validateInvocationText validates compact UTF-8 invocation text.
func validateInvocationText(field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidInvocation, field, err)
	}

	return nil
}
