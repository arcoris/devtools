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

// Validate verifies usage-line structural rules.
func (line UsageLine) Validate() error {
	raw := string(line)

	if strings.TrimSpace(raw) == "" {
		return ErrEmptyUsage
	}

	if err := validateUsageLineText(raw); err != nil {
		return err
	}

	if raw != normalizeUsageLine(raw) {
		return fmt.Errorf("%w: line %q is not canonical", ErrInvalidUsage, raw)
	}

	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return ErrEmptyUsage
	}

	for index, token := range tokens {
		if err := validateUsageToken(index, token); err != nil {
			return err
		}
	}

	return nil
}

// validateUsageLineInput validates raw constructor input before normalization.
func validateUsageLineInput(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return ErrEmptyUsage
	}

	return validateUsageLineText(raw)
}

// validateUsageLineText validates a bounded single-line usage string.
func validateUsageLineText(raw string) error {
	if err := textvalidate.ValidateSingleLineText(raw, maxUsageLineLength); err != nil {
		return fmt.Errorf("%w: invalid usage line: %w", ErrInvalidUsage, err)
	}

	return nil
}

// validateUsageToken validates one usage token.
func validateUsageToken(index int, token string) error {
	if err := textvalidate.ValidateTokenText(token, maxUsageTokenLength); err != nil {
		return fmt.Errorf("%w: invalid token %d %q: %w", ErrInvalidUsage, index, token, err)
	}

	return nil
}

// normalizeUsageLine returns the canonical single-line usage syntax.
func normalizeUsageLine(raw string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
}
