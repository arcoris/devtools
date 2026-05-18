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

// TestBindingWithHelpers verifies immutable-style declaration updates.
func TestBindingWithHelpers(t *testing.T) {
	t.Parallel()

	binding := EmptyBinding().
		MustWithOption(bindingTestStringOption("output")).
		MustWithArgument(bindingTestStringArgument("package"))

	if !binding.HasOption(MustOptionName("output")) {
		t.Fatalf("HasOption(output) = false, want true")
	}

	if !binding.HasArgument(MustArgumentName("package")) {
		t.Fatalf("HasArgument(package) = false, want true")
	}

	withoutOption := binding.WithoutOption(MustOptionName("output"))
	if withoutOption.HasOption(MustOptionName("output")) {
		t.Fatalf("WithoutOption() still has output")
	}

	withoutArgument := binding.WithoutArgument(MustArgumentName("package"))
	if withoutArgument.HasArgument(MustArgumentName("package")) {
		t.Fatalf("WithoutArgument() still has package")
	}
}
