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

// NewUsage creates a Usage from spec.
func NewUsage(spec UsageSpec) (Usage, error) {
	syntax, err := NewUsageLine(spec.Syntax)
	if err != nil {
		return Usage{}, err
	}

	alternatives, err := newUsageLines(spec.Alternatives)
	if err != nil {
		return Usage{}, err
	}

	usage := Usage{
		syntax:       syntax,
		alternatives: cloneUsageLines(alternatives),
	}

	if err := usage.Validate(); err != nil {
		return Usage{}, err
	}

	return usage, nil
}

// MustUsage creates a Usage from spec and panics on invalid input.
//
// MustUsage is intended for static command definitions and tests where invalid
// usage declarations are programmer errors.
func MustUsage(spec UsageSpec) Usage {
	usage, err := NewUsage(spec)
	if err != nil {
		panic(err)
	}

	return usage
}

// NewSimpleUsage creates a Usage with one primary syntax line and no
// alternatives.
func NewSimpleUsage(syntax string) (Usage, error) {
	return NewUsage(UsageSpec{
		Syntax: syntax,
	})
}

// MustSimpleUsage creates a Usage with one primary syntax line and panics on
// invalid input.
func MustSimpleUsage(syntax string) Usage {
	return MustUsage(UsageSpec{
		Syntax: syntax,
	})
}

// NewUsageLine validates raw and returns a normalized UsageLine.
func NewUsageLine(raw string) (UsageLine, error) {
	if err := validateUsageLineInput(raw); err != nil {
		return "", err
	}

	line := UsageLine(normalizeUsageLine(raw))
	if err := line.Validate(); err != nil {
		return "", err
	}

	return line, nil
}

// MustUsageLine validates raw and returns a normalized UsageLine.
//
// MustUsageLine panics on invalid input. It is intended for static command
// definitions and tests.
func MustUsageLine(raw string) UsageLine {
	line, err := NewUsageLine(raw)
	if err != nil {
		panic(err)
	}

	return line
}

// newUsageLines validates raw lines and returns normalized usage lines.
func newUsageLines(rawLines []string) ([]UsageLine, error) {
	lines := make([]UsageLine, 0, len(rawLines))
	for index, raw := range rawLines {
		line, err := NewUsageLine(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: alternative %d: %w", ErrInvalidUsage, index, err)
		}

		lines = append(lines, line)
	}

	return lines, nil
}
