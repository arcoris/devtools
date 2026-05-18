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

// WithOption returns a validated copy with one option declaration appended or
// replaced by canonical name.
func (binding Binding) WithOption(option Option) (Binding, error) {
	next := Binding{
		options:   cloneBindingOptions(binding.options),
		arguments: cloneBindingArguments(binding.arguments),
	}

	replaced := false
	for index, existing := range next.options {
		if existing.Name() == option.Name() {
			next.options[index] = option
			replaced = true

			break
		}
	}

	if !replaced {
		next.options = append(next.options, option)
	}

	if err := next.Validate(); err != nil {
		return Binding{}, err
	}

	return next, nil
}

// MustWithOption returns a validated copy with one option declaration appended
// or replaced and panics on invalid input.
func (binding Binding) MustWithOption(option Option) Binding {
	next, err := binding.WithOption(option)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutOption returns a validated copy without the option matching canonical
// name or alias.
func (binding Binding) WithoutOption(name OptionName) Binding {
	next := Binding{
		options:   make([]Option, 0, len(binding.options)),
		arguments: cloneBindingArguments(binding.arguments),
	}

	for _, option := range binding.options {
		if option.MatchesName(name) {
			continue
		}

		next.options = append(next.options, option)
	}

	return MustBinding(BindingSpec{
		Options:   next.options,
		Arguments: next.arguments,
	})
}

// WithArgument returns a validated copy with one positional argument appended or
// replaced by name.
func (binding Binding) WithArgument(argument Argument) (Binding, error) {
	next := Binding{
		options:   cloneBindingOptions(binding.options),
		arguments: cloneBindingArguments(binding.arguments),
	}

	replaced := false
	for index, existing := range next.arguments {
		if existing.Name() == argument.Name() {
			next.arguments[index] = argument
			replaced = true

			break
		}
	}

	if !replaced {
		next.arguments = append(next.arguments, argument)
	}

	if err := next.Validate(); err != nil {
		return Binding{}, err
	}

	return next, nil
}

// MustWithArgument returns a validated copy with one positional argument
// appended or replaced and panics on invalid input.
func (binding Binding) MustWithArgument(argument Argument) Binding {
	next, err := binding.WithArgument(argument)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutArgument returns a validated copy without the argument matching name.
func (binding Binding) WithoutArgument(name ArgumentName) Binding {
	next := Binding{
		options:   cloneBindingOptions(binding.options),
		arguments: make([]Argument, 0, len(binding.arguments)),
	}

	for _, argument := range binding.arguments {
		if argument.Name() == name {
			continue
		}

		next.arguments = append(next.arguments, argument)
	}

	return MustBinding(BindingSpec{
		Options:   next.options,
		Arguments: next.arguments,
	})
}
