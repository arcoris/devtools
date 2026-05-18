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

import "errors"

var (
	// ErrEmptyOptionKind reports that an option kind was not provided.
	ErrEmptyOptionKind = errors.New("command option kind is empty")

	// ErrInvalidOptionKind reports that an option kind is not one of the
	// supported option value kinds.
	ErrInvalidOptionKind = errors.New("command option kind is invalid")
)

// OptionKind describes the value shape accepted by a command option.
//
// OptionKind is intentionally a closed kernel-level enum. Unlike Group or
// Topic, option value kinds are infrastructure semantics that parsers,
// validators, help renderers, documentation generators, and shell-completion
// adapters must interpret consistently.
//
// OptionKind does not describe a concrete option. It only describes what kind
// of value an option accepts. Concrete option identity, names, defaults,
// validation rules, deprecation policy, and documentation belong to Option and
// higher layers.
type OptionKind string

const (
	// OptionKindBool describes a boolean option.
	OptionKindBool OptionKind = "bool"

	// OptionKindString describes an unconstrained string option.
	OptionKindString OptionKind = "string"

	// OptionKindEnum describes a string option constrained by allowed values.
	OptionKindEnum OptionKind = "enum"

	// OptionKindInt describes an int option.
	OptionKindInt OptionKind = "int"

	// OptionKindInt64 describes an int64 option.
	OptionKindInt64 OptionKind = "int64"

	// OptionKindUint describes a uint option.
	OptionKindUint OptionKind = "uint"

	// OptionKindUint64 describes a uint64 option.
	OptionKindUint64 OptionKind = "uint64"

	// OptionKindFloat64 describes a float64 option.
	OptionKindFloat64 OptionKind = "float64"

	// OptionKindDuration describes a duration option.
	OptionKindDuration OptionKind = "duration"

	// OptionKindStringList describes a list of string values.
	OptionKindStringList OptionKind = "string-list"

	// OptionKindEnumList describes a list of enum values.
	OptionKindEnumList OptionKind = "enum-list"

	// OptionKindIntList describes a list of int values.
	OptionKindIntList OptionKind = "int-list"

	// OptionKindInt64List describes a list of int64 values.
	OptionKindInt64List OptionKind = "int64-list"

	// OptionKindUintList describes a list of uint values.
	OptionKindUintList OptionKind = "uint-list"

	// OptionKindUint64List describes a list of uint64 values.
	OptionKindUint64List OptionKind = "uint64-list"

	// OptionKindFloat64List describes a list of float64 values.
	OptionKindFloat64List OptionKind = "float64-list"

	// OptionKindDurationList describes a list of duration values.
	OptionKindDurationList OptionKind = "duration-list"
)

var knownOptionKinds = []OptionKind{
	OptionKindBool,
	OptionKindString,
	OptionKindEnum,
	OptionKindInt,
	OptionKindInt64,
	OptionKindUint,
	OptionKindUint64,
	OptionKindFloat64,
	OptionKindDuration,
	OptionKindStringList,
	OptionKindEnumList,
	OptionKindIntList,
	OptionKindInt64List,
	OptionKindUintList,
	OptionKindUint64List,
	OptionKindFloat64List,
	OptionKindDurationList,
}
