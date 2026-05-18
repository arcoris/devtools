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

import "strings"

// Spec returns a detached construction spec.
func (value OptionValue) Spec() OptionValueSpec {
	return OptionValueSpec{
		Name:   value.name.String(),
		Kind:   value.kind,
		Source: value.source,
		Values: cloneStringSlice(value.values),
	}
}

// IsZero reports whether value has not been initialized.
func (value OptionValue) IsZero() bool {
	return value.name == "" &&
		value.kind == "" &&
		value.source == "" &&
		len(value.values) == 0
}

// IsValid reports whether value satisfies structural rules.
func (value OptionValue) IsValid() bool {
	return value.Validate() == nil
}

// Name returns the canonical option name.
func (value OptionValue) Name() OptionName {
	return value.name
}

// Kind returns the resolved value kind.
func (value OptionValue) Kind() OptionKind {
	return value.kind
}

// Source returns the resolved value source.
func (value OptionValue) Source() OptionSource {
	return value.source
}

// Values returns detached raw resolved values.
func (value OptionValue) Values() []string {
	return cloneStringSlice(value.values)
}

// Value returns the first raw value and whether it exists.
//
// Scalar OptionValue instances always return true. List OptionValue instances
// return the first list element.
func (value OptionValue) Value() (string, bool) {
	return value.ValueAt(0)
}

// ValueAt returns one raw value by index.
func (value OptionValue) ValueAt(index int) (string, bool) {
	if index < 0 || index >= len(value.values) {
		return "", false
	}

	return value.values[index], true
}

// MustValue returns the first raw value and panics if it is absent.
func (value OptionValue) MustValue() string {
	raw, ok := value.Value()
	if !ok {
		panic("command option value is absent")
	}

	return raw
}

// Len returns the number of raw values.
func (value OptionValue) Len() int {
	return len(value.values)
}

// IsEmpty reports whether the value contains no raw values.
func (value OptionValue) IsEmpty() bool {
	return len(value.values) == 0
}

// IsScalar reports whether the value kind is scalar.
func (value OptionValue) IsScalar() bool {
	return value.kind.IsScalar()
}

// IsList reports whether the value kind is list-shaped.
func (value OptionValue) IsList() bool {
	return value.kind.IsList()
}

// IsDefault reports whether the value came from declaration defaults.
func (value OptionValue) IsDefault() bool {
	return value.source.IsDefault()
}

// IsExplicit reports whether the value came from an explicit external source.
func (value OptionValue) IsExplicit() bool {
	return value.source.IsExplicit()
}

// String returns a compact string representation of the raw value.
//
// Scalar values are returned as-is. List values are joined with "," for compact
// diagnostics. Use Values when lossless list access is needed.
func (value OptionValue) String() string {
	if len(value.values) == 0 {
		return ""
	}

	if len(value.values) == 1 {
		return value.values[0]
	}

	return strings.Join(value.values, ",")
}
