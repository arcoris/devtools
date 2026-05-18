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

import "testing"

// TestArgumentWithAllowedValues verifies immutable allowed-value update.
func TestArgumentWithAllowedValues(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text"},
	})

	next := argument.MustWithAllowedValues("text", "json")

	if !next.AllowsValue("json") {
		t.Fatalf("updated argument does not allow json")
	}

	if argument.AllowsValue("json") {
		t.Fatalf("WithAllowedValues mutated original argument")
	}
}

// TestArgumentWithHelpers verifies immutable-style argument updates.
func TestArgumentWithHelpers(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Requirement: ArgumentRequirementOptional,
	}).
		MustWithDocumentation(MustSummaryDocumentation("Package pattern.")).
		MustWithMetadata(MustMetadata(MetadataSpec{Owner: "devtools"})).
		MustWithVisibility(VisibilityHidden).
		MustWithDefaultValues("./...")

	if got, want := argument.Documentation().Summary(), "Package pattern."; got != want {
		t.Fatalf("Documentation().Summary() = %q, want %q", got, want)
	}

	if got, want := argument.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if !argument.Visibility().IsHidden() {
		t.Fatalf("Visibility().IsHidden() = false, want true")
	}

	if !argument.HasDefault() {
		t.Fatalf("HasDefault() = false, want true")
	}

	without := argument.WithoutDefault()
	if without.HasDefault() {
		t.Fatalf("WithoutDefault() still has default")
	}
}
