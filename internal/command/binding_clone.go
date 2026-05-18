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

// cloneBindingOptions returns a detached copy of option declarations.
func cloneBindingOptions(values []Option) []Option {
	if values == nil {
		return nil
	}

	out := make([]Option, len(values))
	copy(out, values)

	return out
}

// cloneBindingArguments returns a detached copy of argument declarations.
func cloneBindingArguments(values []Argument) []Argument {
	if values == nil {
		return nil
	}

	out := make([]Argument, len(values))
	copy(out, values)

	return out
}

// cloneBindingOptionValues returns a detached copy of option values.
func cloneBindingOptionValues(values []OptionValue) []OptionValue {
	if values == nil {
		return nil
	}

	out := make([]OptionValue, len(values))
	copy(out, values)

	return out
}

// cloneBoundArguments returns a detached copy of bound arguments.
func cloneBoundArguments(values []BoundArgument) []BoundArgument {
	if values == nil {
		return nil
	}

	out := make([]BoundArgument, len(values))
	for index, value := range values {
		out[index] = BoundArgument{
			argument: value.argument,
			values:   cloneBindingStrings(value.values),
		}
	}

	return out
}

// cloneBindingStrings returns a detached copy of strings.
func cloneBindingStrings(values []string) []string {
	if values == nil {
		return nil
	}

	out := make([]string, len(values))
	copy(out, values)

	return out
}
