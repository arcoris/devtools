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
	"errors"
	"testing"
)

// TestValidateBindingArgumentsDirect verifies direct argument validation helper.
func TestValidateBindingArgumentsDirect(t *testing.T) {
	t.Parallel()

	if err := validateBindingArguments([]Argument{bindingTestStringArgument("package")}); err != nil {
		t.Fatalf("validateBindingArguments(valid) returned unexpected error: %v", err)
	}

	err := validateBindingArguments([]Argument{
		bindingTestStringArgument("package"),
		bindingTestStringArgument("package"),
	})
	if err == nil {
		t.Fatalf("validateBindingArguments(duplicate) returned nil error")
	}

	if !errors.Is(err, ErrInvalidBinding) {
		t.Fatalf("validateBindingArguments(duplicate) error = %v, want ErrInvalidBinding", err)
	}
}

// TestValidateBindingOptionsDirect verifies direct option validation helper.
func TestValidateBindingOptionsDirect(t *testing.T) {
	t.Parallel()

	if err := validateBindingOptions([]Option{bindingTestStringOption("output")}); err != nil {
		t.Fatalf("validateBindingOptions(valid) returned unexpected error: %v", err)
	}

	err := validateBindingOptions([]Option{
		bindingTestStringOption("output"),
		bindingTestStringOption("output"),
	})
	if err == nil {
		t.Fatalf("validateBindingOptions(duplicate) returned nil error")
	}

	if !errors.Is(err, ErrInvalidBinding) {
		t.Fatalf("validateBindingOptions(duplicate) error = %v, want ErrInvalidBinding", err)
	}
}
