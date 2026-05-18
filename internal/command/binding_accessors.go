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

import "sort"

// Options returns detached option declarations.
func (binding Binding) Options() []Option {
	return cloneBindingOptions(binding.options)
}

// Arguments returns detached positional argument declarations.
func (binding Binding) Arguments() []Argument {
	return cloneBindingArguments(binding.arguments)
}

// OptionCount returns the number of option declarations.
func (binding Binding) OptionCount() int {
	return len(binding.options)
}

// ArgumentCount returns the number of positional argument declarations.
func (binding Binding) ArgumentCount() int {
	return len(binding.arguments)
}

// IsZero reports whether binding declares no inputs.
func (binding Binding) IsZero() bool {
	return len(binding.options) == 0 && len(binding.arguments) == 0
}

// HasOptions reports whether binding declares options.
func (binding Binding) HasOptions() bool {
	return len(binding.options) > 0
}

// HasArguments reports whether binding declares positional arguments.
func (binding Binding) HasArguments() bool {
	return len(binding.arguments) > 0
}

// Option returns an option declaration by canonical name or alias.
func (binding Binding) Option(name OptionName) (Option, bool) {
	for _, option := range binding.options {
		if option.MatchesName(name) {
			return option, true
		}
	}

	return Option{}, false
}

// OptionByName parses raw and returns an option declaration by canonical name
// or alias.
func (binding Binding) OptionByName(raw string) (Option, bool) {
	name, err := NewOptionName(raw)
	if err != nil {
		return Option{}, false
	}

	return binding.Option(name)
}

// HasOption reports whether name matches a declared option canonical name or
// alias.
func (binding Binding) HasOption(name OptionName) bool {
	_, ok := binding.Option(name)

	return ok
}

// OptionByShorthand returns an option declaration by shorthand name.
func (binding Binding) OptionByShorthand(shorthand string) (Option, bool) {
	for _, option := range binding.options {
		if option.Shorthand() == shorthand {
			return option, true
		}
	}

	return Option{}, false
}

// HasOptionShorthand reports whether shorthand is declared.
func (binding Binding) HasOptionShorthand(shorthand string) bool {
	_, ok := binding.OptionByShorthand(shorthand)

	return ok
}

// Argument returns an argument declaration by name.
func (binding Binding) Argument(name ArgumentName) (Argument, bool) {
	for _, argument := range binding.arguments {
		if argument.Name() == name {
			return argument, true
		}
	}

	return Argument{}, false
}

// ArgumentByName parses raw and returns an argument declaration by name.
func (binding Binding) ArgumentByName(raw string) (Argument, bool) {
	name, err := NewArgumentName(raw)
	if err != nil {
		return Argument{}, false
	}

	return binding.Argument(name)
}

// HasArgument reports whether name matches a declared positional argument.
func (binding Binding) HasArgument(name ArgumentName) bool {
	_, ok := binding.Argument(name)

	return ok
}

// OptionNames returns canonical option names in declaration order.
func (binding Binding) OptionNames() []OptionName {
	names := make([]OptionName, len(binding.options))
	for index, option := range binding.options {
		names[index] = option.Name()
	}

	return names
}

// ArgumentNames returns argument names in declaration order.
func (binding Binding) ArgumentNames() []ArgumentName {
	names := make([]ArgumentName, len(binding.arguments))
	for index, argument := range binding.arguments {
		names[index] = argument.Name()
	}

	return names
}

// SortedOptionNames returns canonical option names in lexical order.
func (binding Binding) SortedOptionNames() []OptionName {
	names := binding.OptionNames()
	sort.Slice(names, func(i int, j int) bool {
		return names[i].String() < names[j].String()
	})

	return names
}

// SortedArgumentNames returns argument names in lexical order.
func (binding Binding) SortedArgumentNames() []ArgumentName {
	names := binding.ArgumentNames()
	sort.Slice(names, func(i int, j int) bool {
		return names[i].String() < names[j].String()
	})

	return names
}
