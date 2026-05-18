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
	"testing"
)

// TestOptionPolicyAllowedSourcesCopy verifies detached allowed-source slices.
func TestOptionPolicyAllowedSourcesCopy(t *testing.T) {
	t.Parallel()

	sources := []OptionSource{
		OptionSourceDefault,
		OptionSourceCommandLine,
	}

	policy := MustOptionPolicy(OptionPolicySpec{
		AllowedSources: sources,
	})

	sources[0] = OptionSource("changed")

	if !policy.AllowsSource(OptionSourceDefault) {
		t.Fatalf("policy changed through input source slice")
	}

	out := policy.AllowedSources()
	out[0] = OptionSource("changed")

	if !policy.AllowsSource(OptionSourceDefault) {
		t.Fatalf("policy changed through output source slice")
	}
}

// TestOptionPolicyHighestAndLowestSource verifies source precedence helpers.
func TestOptionPolicyHighestAndLowestSource(t *testing.T) {
	t.Parallel()

	policy := MustOptionPolicy(OptionPolicySpec{
		AllowedSources: []OptionSource{
			OptionSourceConfig,
			OptionSourceCommandLine,
			OptionSourceEnvironment,
		},
	})

	highest, ok := policy.HighestAllowedSource()
	if !ok {
		t.Fatalf("HighestAllowedSource() ok = false, want true")
	}

	if got, want := highest, OptionSourceCommandLine; got != want {
		t.Fatalf("HighestAllowedSource() = %q, want %q", got, want)
	}

	lowest, ok := policy.LowestAllowedSource()
	if !ok {
		t.Fatalf("LowestAllowedSource() ok = false, want true")
	}

	if got, want := lowest, OptionSourceConfig; got != want {
		t.Fatalf("LowestAllowedSource() = %q, want %q", got, want)
	}
}
