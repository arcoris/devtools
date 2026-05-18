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
	"fmt"
)

var (
	// ErrInvalidOptionResolution reports that option values could not be
	// resolved into a canonical Binding input.
	ErrInvalidOptionResolution = errors.New("command option resolution is invalid")
)

// OptionResolutionSpec describes framework-neutral option resolution input.
//
// The resolver boundary accepts values that have already been parsed or loaded
// by adapters. It does not inspect os.Args, Cobra flags, pflag state,
// environment variables, config files, interactive prompts, or runtime
// injection mechanisms. Those layers translate their data into OptionValue
// values and pass them here.
type OptionResolutionSpec struct {
	// Binding declares the options that may be resolved.
	Binding Binding

	// InheritedValues contains values inherited from parent commands or shared
	// command context.
	InheritedValues []OptionValue

	// ConfigValues contains values loaded by configuration adapters.
	ConfigValues []OptionValue

	// EnvironmentValues contains values loaded by environment adapters.
	EnvironmentValues []OptionValue

	// RuntimeValues contains values injected by application/runtime wiring.
	RuntimeValues []OptionValue

	// InteractiveValues contains values supplied by an explicit prompt layer.
	InteractiveValues []OptionValue

	// CommandLineValues contains values parsed from command-line syntax by an
	// adapter.
	CommandLineValues []OptionValue
}

// ResolveOptionValues canonicalizes and resolves option values according to
// Binding declarations, OptionPolicy allowed sources, and default source
// precedence.
//
// Same-source duplicate values are rejected for scalar options. List-shaped
// options whose policy allows multiple occurrences are merged in source order
// and revalidated against the declaration. Higher-precedence sources still
// override lower-precedence sources for the same option.
//
// The returned values are suitable for Binding.Bind. They are ordered by the
// Binding option declaration order. Declaration defaults are included when no
// higher-precedence value exists and the option policy allows the default
// source.
func ResolveOptionValues(spec OptionResolutionSpec) ([]OptionValue, error) {
	resolver := optionResolver{binding: spec.Binding}
	if err := resolver.binding.Validate(); err != nil {
		return nil, fmt.Errorf("%w: invalid binding: %w", ErrInvalidOptionResolution, err)
	}

	if err := resolver.addValues(OptionSourceInherited, spec.InheritedValues); err != nil {
		return nil, err
	}

	if err := resolver.addValues(OptionSourceConfig, spec.ConfigValues); err != nil {
		return nil, err
	}

	if err := resolver.addValues(OptionSourceEnvironment, spec.EnvironmentValues); err != nil {
		return nil, err
	}

	if err := resolver.addValues(OptionSourceRuntime, spec.RuntimeValues); err != nil {
		return nil, err
	}

	if err := resolver.addValues(OptionSourceInteractive, spec.InteractiveValues); err != nil {
		return nil, err
	}

	if err := resolver.addValues(OptionSourceCommandLine, spec.CommandLineValues); err != nil {
		return nil, err
	}

	return resolver.values()
}

type optionResolver struct {
	binding  Binding
	resolved map[OptionName]OptionValue
}

// addValues validates and records source-specific values.
func (resolver *optionResolver) addValues(expectedSource OptionSource, values []OptionValue) error {
	for index, value := range values {
		if err := value.Validate(); err != nil {
			return fmt.Errorf(
				"%w: %s value %d: %w",
				ErrInvalidOptionResolution,
				expectedSource,
				index,
				err,
			)
		}

		if value.Source() != expectedSource {
			return fmt.Errorf(
				"%w: %s value %d for option %q has source %q",
				ErrInvalidOptionResolution,
				expectedSource,
				index,
				value.Name(),
				value.Source(),
			)
		}

		option, ok := resolver.binding.Option(value.Name())
		if !ok {
			return fmt.Errorf(
				"%w: %s value %d references unknown option %q",
				ErrInvalidOptionResolution,
				expectedSource,
				index,
				value.Name(),
			)
		}

		canonical, err := NewOptionValueFromOption(option, expectedSource, value.Values()...)
		if err != nil {
			return fmt.Errorf(
				"%w: %s value %d for option %q: %w",
				ErrInvalidOptionResolution,
				expectedSource,
				index,
				option.Name(),
				err,
			)
		}

		if resolver.resolved == nil {
			resolver.resolved = make(map[OptionName]OptionValue)
		}

		existing, exists := resolver.resolved[option.Name()]
		if exists && existing.Source() == canonical.Source() {
			merged, err := mergeSameSourceOptionValues(option, existing, canonical)
			if err != nil {
				return fmt.Errorf(
					"%w: duplicate %s value for option %q: %w",
					ErrInvalidOptionResolution,
					expectedSource,
					option.Name(),
					err,
				)
			}

			resolver.resolved[option.Name()] = merged

			continue
		}

		if !exists || canonical.Source().Overrides(existing.Source()) {
			resolver.resolved[option.Name()] = canonical
		}
	}

	return nil
}

// values returns resolved values in binding declaration order.
func (resolver optionResolver) values() ([]OptionValue, error) {
	resolved := make([]OptionValue, 0, resolver.binding.OptionCount())

	for _, option := range resolver.binding.Options() {
		if value, ok := resolver.resolved[option.Name()]; ok {
			resolved = append(resolved, value)
			continue
		}

		if option.HasDefault() && option.Policy().AllowsDefaultSource() {
			value, err := NewDefaultOptionValue(option)
			if err != nil {
				return nil, fmt.Errorf(
					"%w: default value for option %q: %w",
					ErrInvalidOptionResolution,
					option.Name(),
					err,
				)
			}

			resolved = append(resolved, value)
		}
	}

	return resolved, nil
}

// mergeSameSourceOptionValues merges repeatable list-shaped values from one
// source. Scalar values remain single-occurrence and are rejected.
func mergeSameSourceOptionValues(option Option, existing OptionValue, next OptionValue) (OptionValue, error) {
	if !optionAllowsSameSourceMerge(option) {
		return OptionValue{}, fmt.Errorf(
			"option kind %q with occurrence %q does not allow same-source duplicates",
			option.Kind(),
			option.Policy().Occurrence(),
		)
	}

	values := existing.Values()
	values = append(values, next.Values()...)

	merged, err := NewOptionValueFromOption(option, existing.Source(), values...)
	if err != nil {
		return OptionValue{}, err
	}

	return merged, nil
}

// optionAllowsSameSourceMerge reports whether the option can represent repeated
// same-source occurrences without losing data.
func optionAllowsSameSourceMerge(option Option) bool {
	return option.Kind().IsList() && option.IsRepeatable()
}
