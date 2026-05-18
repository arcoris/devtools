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

// TestMustOptionPanicsForInvalidOption verifies fail-fast construction.
func TestMustOptionPanicsForInvalidOption(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustOption did not panic")
		}
	}()

	_ = MustOption(OptionSpec{
		Name: "Output",
		Kind: OptionKindString,
	})
}

// TestNewOptionAcceptsEnumOption verifies enum allowed-value validation.
func TestNewOptionAcceptsEnumOption(t *testing.T) {
	t.Parallel()

	option, err := NewOption(OptionSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json", "markdown"},
		DefaultValues: []string{"text"},
	})
	if err != nil {
		t.Fatalf("NewOption() returned unexpected error: %v", err)
	}

	if !option.HasAllowedValues() {
		t.Fatalf("HasAllowedValues() = false, want true")
	}

	if !option.AllowsValue("json") {
		t.Fatalf("AllowsValue(json) = false, want true")
	}

	if option.AllowsValue("xml") {
		t.Fatalf("AllowsValue(xml) = true, want false")
	}
}

// TestNewOptionAcceptsValidOption verifies full option declaration construction.
func TestNewOptionAcceptsValidOption(t *testing.T) {
	t.Parallel()

	option, err := NewOption(OptionSpec{
		Name:      "output",
		Aliases:   []string{"out"},
		Shorthand: "o",
		Kind:      OptionKindString,
		Metavar:   "PATH",
		DefaultValues: []string{
			"report.txt",
		},
		Documentation: MustSummaryDocumentation("Write report to path."),
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewOption() returned unexpected error: %v", err)
	}

	if got, want := option.Name(), MustOptionName("output"); got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if got, want := option.LongFlag(), "--output"; got != want {
		t.Fatalf("LongFlag() = %q, want %q", got, want)
	}

	if got, want := option.ShortFlag(), "-o"; got != want {
		t.Fatalf("ShortFlag() = %q, want %q", got, want)
	}

	if got, want := option.Metavar(), "PATH"; got != want {
		t.Fatalf("Metavar() = %q, want %q", got, want)
	}

	if !option.HasDefault() {
		t.Fatalf("HasDefault() = false, want true")
	}

	if got, ok := option.DefaultValue(); !ok || got != "report.txt" {
		t.Fatalf("DefaultValue() = %q, %v; want report.txt, true", got, ok)
	}

	if !option.MatchesName(MustOptionName("out")) {
		t.Fatalf("MatchesName(alias) = false, want true")
	}
}

// TestNewOptionDefaultsListOccurrence verifies kind-aware default policy for list options.
func TestNewOptionDefaultsListOccurrence(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "package",
		Kind: OptionKindStringList,
	})

	if !option.IsRepeatable() {
		t.Fatalf("string-list option IsRepeatable() = false, want true")
	}
}

// TestNewOptionDefaults verifies default policy, visibility, and metavar.
func TestNewOptionDefaults(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "verbose",
		Kind: OptionKindBool,
	})

	if got, want := option.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if got, want := option.Metavar(), "BOOL"; got != want {
		t.Fatalf("Metavar() = %q, want %q", got, want)
	}

	if !option.Policy().AllowsDefaultSource() {
		t.Fatalf("default policy should allow default source")
	}

	if option.IsRepeatable() {
		t.Fatalf("bool option IsRepeatable() = true, want false")
	}
}
