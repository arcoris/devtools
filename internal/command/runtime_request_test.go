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

// TestRuntimeRequestAccessors verifies request helper methods.
func TestRuntimeRequestAccessors(t *testing.T) {
	t.Parallel()

	startedAt := runtimeTestTime()
	binding := runtimeTestBinding()
	bound := binding.MustBind(BindingValueSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})

	request := RuntimeRequest{
		runtimeName: "test-runtime",
		commandID:   MustID("bench.run"),
		binding:     binding,
		input:       bound,
		startedAt:   startedAt,
		metadata:    MustMetadata(MetadataSpec{Owner: "devtools"}),
	}

	if got, want := request.RuntimeName(), "test-runtime"; got != want {
		t.Fatalf("RuntimeName() = %q, want %q", got, want)
	}

	if commandID, ok := request.CommandID(); !ok || commandID != MustID("bench.run") {
		t.Fatalf("CommandID() = %q, %v; want bench.run, true", commandID, ok)
	}

	if got := request.StartedAt(); !got.Equal(startedAt) {
		t.Fatalf("StartedAt() = %v, want %v", got, startedAt)
	}

	if got, want := request.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if _, ok := request.Option(MustOptionName("format")); !ok {
		t.Fatalf("Option(format) ok = false, want true")
	}

	if _, ok := request.Argument(MustArgumentName("suite")); !ok {
		t.Fatalf("Argument(suite) ok = false, want true")
	}
}
