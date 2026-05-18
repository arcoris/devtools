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

// TestOptionSourceOverrides verifies precedence comparison helpers.
func TestOptionSourceOverrides(t *testing.T) {
	t.Parallel()

	if !OptionSourceCommandLine.Overrides(OptionSourceEnvironment) {
		t.Fatalf("command-line should override environment")
	}

	if !OptionSourceEnvironment.Overrides(OptionSourceConfig) {
		t.Fatalf("environment should override config")
	}

	if OptionSourceConfig.Overrides(OptionSourceEnvironment) {
		t.Fatalf("config should not override environment")
	}

	if OptionSourceDefault.Overrides(OptionSourceDefault) {
		t.Fatalf("same source must not strictly override itself")
	}

	if OptionSource("unknown").Overrides(OptionSourceDefault) {
		t.Fatalf("unknown source must not override known source")
	}

	if !OptionSourceConfig.CanBeOverriddenBy(OptionSourceEnvironment) {
		t.Fatalf("config should be overridable by environment")
	}

	if !OptionSourceConfig.SamePrecedence(OptionSourceConfig) {
		t.Fatalf("same known source should have same precedence")
	}

	if OptionSource("").SamePrecedence(OptionSource("")) {
		t.Fatalf("zero sources should not report same precedence")
	}
}

// TestOptionSourcePrecedence verifies default precedence ordering.
func TestOptionSourcePrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		source OptionSource
		want   int
	}{
		{source: "", want: 0},
		{source: OptionSource("unknown"), want: 0},
		{source: OptionSourceDefault, want: 10},
		{source: OptionSourceInherited, want: 20},
		{source: OptionSourceConfig, want: 30},
		{source: OptionSourceEnvironment, want: 40},
		{source: OptionSourceRuntime, want: 50},
		{source: OptionSourceInteractive, want: 60},
		{source: OptionSourceCommandLine, want: 70},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.source.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.source.Precedence(); got != tt.want {
				t.Fatalf("%q.Precedence() = %d, want %d", tt.source, got, tt.want)
			}
		})
	}
}
