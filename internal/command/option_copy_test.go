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

// TestOptionCopySemantics verifies detached alias/default/allowed slices.
func TestOptionCopySemantics(t *testing.T) {
	t.Parallel()

	aliases := []string{"out"}
	defaults := []string{"text"}
	allowed := []string{"text", "json"}

	option := MustOption(OptionSpec{
		Name:          "format",
		Aliases:       aliases,
		Kind:          OptionKindEnum,
		DefaultValues: defaults,
		AllowedValues: allowed,
	})

	aliases[0] = "changed"
	defaults[0] = "changed"
	allowed[0] = "changed"

	if option.HasAlias(MustOptionName("changed")) {
		t.Fatalf("alias changed through input slice")
	}

	if got, want := option.DefaultValues()[0], "text"; got != want {
		t.Fatalf("default changed through input slice: got %q, want %q", got, want)
	}

	if got, want := option.AllowedValues()[0], "text"; got != want {
		t.Fatalf("allowed changed through input slice: got %q, want %q", got, want)
	}

	outAliases := option.Aliases()
	outAliases[0] = MustOptionName("changed")

	if option.HasAlias(MustOptionName("changed")) {
		t.Fatalf("alias changed through output slice")
	}

	outDefaults := option.DefaultValues()
	outDefaults[0] = "changed"

	if got, want := option.DefaultValues()[0], "text"; got != want {
		t.Fatalf("default changed through output slice: got %q, want %q", got, want)
	}
}
