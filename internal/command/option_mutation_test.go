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

// TestOptionWithAllowedValues verifies immutable allowed-value update.
func TestOptionWithAllowedValues(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "format",
		Kind: OptionKindEnum,
		AllowedValues: []string{
			"text",
		},
	})

	next := option.MustWithAllowedValues("text", "json")

	if !next.AllowsValue("json") {
		t.Fatalf("updated option does not allow json")
	}

	if option.AllowsValue("json") {
		t.Fatalf("WithAllowedValues mutated original option")
	}
}

// TestOptionWithHelpers verifies immutable-style option updates.
func TestOptionWithHelpers(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "output",
		Kind: OptionKindString,
	}).
		MustWithDocumentation(MustSummaryDocumentation("Write output.")).
		MustWithMetadata(MustMetadata(MetadataSpec{Owner: "devtools"})).
		MustWithVisibility(VisibilityHidden).
		MustWithDefaultValues("report.txt")

	if got, want := option.Documentation().Summary(), "Write output."; got != want {
		t.Fatalf("Documentation().Summary() = %q, want %q", got, want)
	}

	if got, want := option.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if !option.Visibility().IsHidden() {
		t.Fatalf("Visibility().IsHidden() = false, want true")
	}

	if !option.HasDefault() {
		t.Fatalf("HasDefault() = false, want true")
	}

	without := option.WithoutDefault()
	if without.HasDefault() {
		t.Fatalf("WithoutDefault() still has default")
	}
}
