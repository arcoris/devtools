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

// WithAliases returns a validated copy with aliases replaced.
func (option Option) WithAliases(aliases ...string) (Option, error) {
	spec := option.spec()
	spec.Aliases = cloneStringSlice(aliases)

	return NewOption(spec)
}

// MustWithAliases returns a validated copy with aliases replaced and panics on
// invalid input.
func (option Option) MustWithAliases(aliases ...string) Option {
	next, err := option.WithAliases(aliases...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithAlias returns a validated copy with one alias appended.
func (option Option) WithAlias(alias string) (Option, error) {
	spec := option.spec()
	spec.Aliases = append(spec.Aliases, alias)

	return NewOption(spec)
}

// MustWithAlias returns a validated copy with one alias appended and panics on
// invalid input.
func (option Option) MustWithAlias(alias string) Option {
	next, err := option.WithAlias(alias)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAlias returns a validated copy without the matching alias.
func (option Option) WithoutAlias(alias string) Option {
	name, err := NewOptionName(alias)
	if err != nil {
		return MustOption(option.spec())
	}

	spec := option.spec()
	nextAliases := make([]string, 0, len(spec.Aliases))
	for _, existing := range option.aliases {
		if existing == name {
			continue
		}

		nextAliases = append(nextAliases, existing.String())
	}

	spec.Aliases = nextAliases

	return MustOption(spec)
}

// WithoutAliases returns a validated copy without aliases.
func (option Option) WithoutAliases() Option {
	spec := option.spec()
	spec.Aliases = nil

	return MustOption(spec)
}
