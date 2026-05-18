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

// TestDefaultOptionSource verifies default source helper.
func TestDefaultOptionSource(t *testing.T) {
	t.Parallel()

	if got, want := DefaultOptionSource(), OptionSourceDefault; got != want {
		t.Fatalf("DefaultOptionSource() = %q, want %q", got, want)
	}
}

// TestOptionSourcePredicates verifies direct source predicates.
func TestOptionSourcePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		source      OptionSource
		zero        bool
		known       bool
		def         bool
		inherited   bool
		config      bool
		environment bool
		runtime     bool
		interactive bool
		commandLine bool
		implicit    bool
		explicit    bool
		user        bool
		external    bool
		persistent  bool
		ephemeral   bool
	}{
		{
			name:   "zero",
			source: "",
			zero:   true,
		},
		{
			name:     "default",
			source:   OptionSourceDefault,
			known:    true,
			def:      true,
			implicit: true,
		},
		{
			name:       "inherited",
			source:     OptionSourceInherited,
			known:      true,
			inherited:  true,
			implicit:   true,
			external:   true,
			persistent: true,
		},
		{
			name:       "config",
			source:     OptionSourceConfig,
			known:      true,
			config:     true,
			explicit:   true,
			external:   true,
			persistent: true,
		},
		{
			name:        "environment",
			source:      OptionSourceEnvironment,
			known:       true,
			environment: true,
			explicit:    true,
			external:    true,
			persistent:  true,
		},
		{
			name:      "runtime",
			source:    OptionSourceRuntime,
			known:     true,
			runtime:   true,
			implicit:  true,
			external:  true,
			ephemeral: true,
		},
		{
			name:        "interactive",
			source:      OptionSourceInteractive,
			known:       true,
			interactive: true,
			explicit:    true,
			user:        true,
			external:    true,
			ephemeral:   true,
		},
		{
			name:        "command-line",
			source:      OptionSourceCommandLine,
			known:       true,
			commandLine: true,
			explicit:    true,
			user:        true,
			external:    true,
			ephemeral:   true,
		},
		{
			name:   "unknown",
			source: OptionSource("file"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.source.IsZero(); got != tt.zero {
				t.Fatalf("IsZero() = %v, want %v", got, tt.zero)
			}

			if got := tt.source.IsKnown(); got != tt.known {
				t.Fatalf("IsKnown() = %v, want %v", got, tt.known)
			}

			if got := tt.source.IsDefault(); got != tt.def {
				t.Fatalf("IsDefault() = %v, want %v", got, tt.def)
			}

			if got := tt.source.IsInherited(); got != tt.inherited {
				t.Fatalf("IsInherited() = %v, want %v", got, tt.inherited)
			}

			if got := tt.source.IsConfig(); got != tt.config {
				t.Fatalf("IsConfig() = %v, want %v", got, tt.config)
			}

			if got := tt.source.IsEnvironment(); got != tt.environment {
				t.Fatalf("IsEnvironment() = %v, want %v", got, tt.environment)
			}

			if got := tt.source.IsRuntime(); got != tt.runtime {
				t.Fatalf("IsRuntime() = %v, want %v", got, tt.runtime)
			}

			if got := tt.source.IsInteractive(); got != tt.interactive {
				t.Fatalf("IsInteractive() = %v, want %v", got, tt.interactive)
			}

			if got := tt.source.IsCommandLine(); got != tt.commandLine {
				t.Fatalf("IsCommandLine() = %v, want %v", got, tt.commandLine)
			}

			if got := tt.source.IsImplicit(); got != tt.implicit {
				t.Fatalf("IsImplicit() = %v, want %v", got, tt.implicit)
			}

			if got := tt.source.IsExplicit(); got != tt.explicit {
				t.Fatalf("IsExplicit() = %v, want %v", got, tt.explicit)
			}

			if got := tt.source.IsUserProvided(); got != tt.user {
				t.Fatalf("IsUserProvided() = %v, want %v", got, tt.user)
			}

			if got := tt.source.IsExternal(); got != tt.external {
				t.Fatalf("IsExternal() = %v, want %v", got, tt.external)
			}

			if got := tt.source.IsPersistent(); got != tt.persistent {
				t.Fatalf("IsPersistent() = %v, want %v", got, tt.persistent)
			}

			if got := tt.source.IsEphemeral(); got != tt.ephemeral {
				t.Fatalf("IsEphemeral() = %v, want %v", got, tt.ephemeral)
			}
		})
	}
}
