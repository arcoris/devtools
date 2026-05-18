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

// WithSyntax returns a validated copy with the primary syntax replaced.
func (usage Usage) WithSyntax(line string) (Usage, error) {
	spec := usage.spec()
	spec.Syntax = line

	return NewUsage(spec)
}

// MustWithSyntax returns a validated copy with the primary syntax replaced and
// panics on invalid input.
func (usage Usage) MustWithSyntax(line string) Usage {
	next, err := usage.WithSyntax(line)
	if err != nil {
		panic(err)
	}

	return next
}

// WithAlternatives returns a validated copy with alternatives replaced.
func (usage Usage) WithAlternatives(lines []string) (Usage, error) {
	spec := usage.spec()
	spec.Alternatives = cloneStringSlice(lines)

	return NewUsage(spec)
}

// MustWithAlternatives returns a validated copy with alternatives replaced and
// panics on invalid input.
func (usage Usage) MustWithAlternatives(lines []string) Usage {
	next, err := usage.WithAlternatives(lines)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAlternatives returns a validated copy without alternative syntax
// lines.
func (usage Usage) WithoutAlternatives() Usage {
	spec := usage.spec()
	spec.Alternatives = nil

	return MustUsage(spec)
}

// WithAlternative returns a validated copy with one alternative usage line
// appended.
func (usage Usage) WithAlternative(line string) (Usage, error) {
	spec := usage.spec()
	spec.Alternatives = append(spec.Alternatives, line)

	return NewUsage(spec)
}

// MustWithAlternative returns a validated copy with one alternative usage line
// appended and panics on invalid input.
func (usage Usage) MustWithAlternative(line string) Usage {
	next, err := usage.WithAlternative(line)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAlternative returns a validated copy without the matching alternative.
//
// If the alternative is not present, the original usage is returned as a
// detached validated copy. The primary syntax is never removed by this method.
func (usage Usage) WithoutAlternative(line string) Usage {
	normalized, err := NewUsageLine(line)
	if err != nil {
		return MustUsage(usage.spec())
	}

	spec := usage.spec()
	spec.Alternatives = spec.Alternatives[:0]
	for _, alternative := range usage.alternatives {
		if alternative == normalized {
			continue
		}

		spec.Alternatives = append(spec.Alternatives, alternative.String())
	}

	return MustUsage(spec)
}
