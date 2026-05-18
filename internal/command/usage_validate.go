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

// Validate verifies usage structural rules.
func (usage Usage) Validate() error {
	if usage.syntax == "" {
		return ErrEmptyUsage
	}

	if err := usage.syntax.Validate(); err != nil {
		return err
	}

	seen := map[UsageLine]struct{}{
		usage.syntax: {},
	}

	for index, alternative := range usage.alternatives {
		if err := alternative.Validate(); err != nil {
			return fmt.Errorf("%w: alternative %d: %w", ErrInvalidUsage, index, err)
		}

		if _, exists := seen[alternative]; exists {
			return fmt.Errorf("%w: duplicate usage line %q", ErrInvalidUsage, alternative)
		}

		seen[alternative] = struct{}{}
	}

	return nil
}

// cloneUsageLines returns a detached copy of usage lines.
func cloneUsageLines(values []UsageLine) []UsageLine {
	if values == nil {
		return nil
	}

	out := make([]UsageLine, len(values))
	copy(out, values)

	return out
}
