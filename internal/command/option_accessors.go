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

// Spec returns a detached construction spec.
func (option Option) Spec() OptionSpec {
	return option.spec()
}

// IsZero reports whether option has not been initialized.
func (option Option) IsZero() bool {
	return option.name == "" &&
		len(option.aliases) == 0 &&
		option.shorthand == "" &&
		option.kind == "" &&
		option.metavar == "" &&
		len(option.defaultValues) == 0 &&
		len(option.allowedValues) == 0 &&
		isZeroOptionPolicy(option.policy) &&
		option.documentation.IsZero() &&
		option.metadata.IsZero() &&
		option.visibility == ""
}

// IsValid reports whether option satisfies structural rules.
func (option Option) IsValid() bool {
	return option.Validate() == nil
}

// Name returns the canonical option name.
func (option Option) Name() OptionName {
	return option.name
}

// LongFlag returns the canonical long CLI flag spelling.
func (option Option) LongFlag() string {
	return option.name.LongFlag()
}

// Aliases returns detached alias option names.
func (option Option) Aliases() []OptionName {
	return cloneOptionNames(option.aliases)
}

// Alias returns the alias at index.
//
// The second return value is false when index is out of range. Alias never
// panics.
func (option Option) Alias(index int) (OptionName, bool) {
	if index < 0 || index >= len(option.aliases) {
		return "", false
	}

	return option.aliases[index], true
}

// AliasCount returns the number of aliases.
func (option Option) AliasCount() int {
	return len(option.aliases)
}

// HasAliases reports whether aliases are declared.
func (option Option) HasAliases() bool {
	return len(option.aliases) > 0
}

// AliasFlags returns detached long CLI flag spellings for aliases.
func (option Option) AliasFlags() []string {
	flags := make([]string, len(option.aliases))
	for index, alias := range option.aliases {
		flags[index] = alias.LongFlag()
	}

	return flags
}

// LongFlags returns the canonical long flag followed by alias long flags.
func (option Option) LongFlags() []string {
	if option.name == "" {
		return nil
	}

	flags := make([]string, 0, 1+len(option.aliases))
	flags = append(flags, option.LongFlag())
	flags = append(flags, option.AliasFlags()...)

	return flags
}

// HasAlias reports whether name is registered as an alias.
func (option Option) HasAlias(name OptionName) bool {
	for _, alias := range option.aliases {
		if alias == name {
			return true
		}
	}

	return false
}

// MatchesName reports whether name matches the canonical name or one alias.
func (option Option) MatchesName(name OptionName) bool {
	if option.name == name {
		return true
	}

	return option.HasAlias(name)
}

// MatchesNameString reports whether raw matches the canonical name or one
// alias. Invalid raw names never match.
func (option Option) MatchesNameString(raw string) bool {
	name, err := NewOptionName(raw)
	if err != nil {
		return false
	}

	return option.MatchesName(name)
}

// Shorthand returns the optional short option name without "-".
func (option Option) Shorthand() string {
	return option.shorthand
}

// ShortFlag returns the optional short CLI flag spelling.
//
// If no shorthand is set, ShortFlag returns an empty string.
func (option Option) ShortFlag() string {
	if option.shorthand == "" {
		return ""
	}

	return "-" + option.shorthand
}

// HasShorthand reports whether a shorthand option name is set.
func (option Option) HasShorthand() bool {
	return option.shorthand != ""
}

// Kind returns the option value kind.
func (option Option) Kind() OptionKind {
	return option.kind
}

// Metavar returns the value placeholder used by help and documentation.
func (option Option) Metavar() string {
	return option.metavar
}

// DefaultValues returns detached declaration default values.
func (option Option) DefaultValues() []string {
	return cloneStringSlice(option.defaultValues)
}

// DefaultValue returns the scalar default value and whether it is present.
//
// For list-shaped options, use DefaultValues.
func (option Option) DefaultValue() (string, bool) {
	return option.DefaultValueAt(0)
}

// DefaultValueAt returns a default value by declaration-order index.
func (option Option) DefaultValueAt(index int) (string, bool) {
	if index < 0 || index >= len(option.defaultValues) {
		return "", false
	}

	return option.defaultValues[index], true
}

// DefaultValueCount returns the number of declaration default values.
func (option Option) DefaultValueCount() int {
	return len(option.defaultValues)
}

// HasDefault reports whether declaration default values are present.
func (option Option) HasDefault() bool {
	return len(option.defaultValues) > 0
}

// AllowedValues returns detached allowed values.
func (option Option) AllowedValues() []string {
	return cloneStringSlice(option.allowedValues)
}

// AllowedValue returns an allowed value by declaration-order index.
func (option Option) AllowedValue(index int) (string, bool) {
	if index < 0 || index >= len(option.allowedValues) {
		return "", false
	}

	return option.allowedValues[index], true
}

// AllowedValueCount returns the number of allowed values.
func (option Option) AllowedValueCount() int {
	return len(option.allowedValues)
}

// HasAllowedValues reports whether allowed values are declared.
func (option Option) HasAllowedValues() bool {
	return len(option.allowedValues) > 0
}

// AllowsValue reports whether value is allowed by the declaration-level
// allowed-value set.
//
// If no allowed values are declared, AllowsValue returns true. Type parsing and
// range validation are separate concerns.
func (option Option) AllowsValue(value string) bool {
	if len(option.allowedValues) == 0 {
		return true
	}

	return containsString(option.allowedValues, value)
}

// Policy returns the option declaration policy.
func (option Option) Policy() OptionPolicy {
	return option.policy
}

// Documentation returns option documentation.
func (option Option) Documentation() Documentation {
	return option.documentation
}

// Metadata returns option metadata.
func (option Option) Metadata() Metadata {
	return option.metadata
}

// Visibility returns option visibility.
func (option Option) Visibility() Visibility {
	return option.visibility
}

// IsVisibleByDefault reports whether default help/docs/discovery should expose
// the option.
func (option Option) IsVisibleByDefault() bool {
	return option.visibility.IsDiscoverableByDefault()
}

// IsRequired reports whether the option must resolve to a value.
func (option Option) IsRequired() bool {
	return option.policy.IsRequired()
}

// IsRepeatable reports whether the option may appear multiple times.
func (option Option) IsRepeatable() bool {
	return option.policy.IsRepeatable()
}

// spec returns a detached construction spec.
func (option Option) spec() OptionSpec {
	aliases := make([]string, len(option.aliases))
	for index, alias := range option.aliases {
		aliases[index] = alias.String()
	}

	return OptionSpec{
		Name:          option.name.String(),
		Aliases:       aliases,
		Shorthand:     option.shorthand,
		Kind:          option.kind,
		Metavar:       option.metavar,
		DefaultValues: cloneStringSlice(option.defaultValues),
		AllowedValues: cloneStringSlice(option.allowedValues),
		Policy:        option.policy,
		Documentation: option.documentation,
		Metadata:      option.metadata,
		Visibility:    option.visibility,
	}
}
