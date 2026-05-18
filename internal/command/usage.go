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

const (
	// maxUsageLineLength is the maximum allowed byte length of one usage line.
	//
	// Usage lines are intended to be compact command syntax declarations, not
	// long-form documentation. Long explanations belong to help templates,
	// generated docs, examples, or other documentation layers.
	maxUsageLineLength = 512

	// maxUsageTokenLength is the maximum allowed byte length of one usage token.
	maxUsageTokenLength = 128
)

var (
	// ErrEmptyUsage reports that a command usage declaration was not provided.
	ErrEmptyUsage = errors.New("command usage is empty")

	// ErrInvalidUsage reports that a command usage declaration is malformed.
	ErrInvalidUsage = errors.New("command usage is invalid")
)

// UsageSpec describes command usage syntax before validation.
//
// UsageSpec is a construction DTO. Usage stores detached and normalized values,
// so callers cannot mutate a constructed value through shared slices.
//
// Usage is intentionally not an executable parser. It describes human-facing
// syntax for help, generated documentation, command discovery, tests, and
// adapter rendering.
//
// Examples of valid syntax:
//
//   - "check [flags]"
//   - "bench run [flags]"
//   - "bench compare <old> <new> [flags]"
//   - "config print [field]"
//   - "profile cpu -- <go test args>"
type UsageSpec struct {
	// Syntax is the primary canonical usage syntax.
	//
	// Syntax MUST be a single logical line. It is normalized by trimming leading
	// and trailing spaces and collapsing internal ASCII spaces into single
	// spaces.
	Syntax string

	// Alternatives contains optional alternative syntax forms.
	//
	// Alternatives are useful when one command has several meaningful invocation
	// shapes. They are not aliases. Aliases belong to Node. Alternatives
	// describe usage syntax only.
	Alternatives []string
}

// Usage is a validated command usage declaration.
//
// Usage is adapter-neutral. A Cobra adapter may use it when rendering command
// help, but this type does not depend on Cobra or any other CLI framework.
//
// Usage is immutable-style:
//
//   - constructors normalize and copy input values;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - duplicate usage lines are rejected.
type Usage struct {
	// syntax stores the primary canonical usage syntax line.
	syntax UsageLine

	// alternatives stores optional canonical alternative syntax lines.
	alternatives []UsageLine
}

// UsageLine is one validated command usage syntax line.
//
// UsageLine is human-facing syntax, not a stable machine identifier. It may
// contain placeholders and adapter syntax such as:
//
//   - "[flags]"
//   - "<path>"
//   - "-- <args>"
//   - "<old> <new>"
//
// UsageLine is deliberately less strict than ID, Path, Group, or Topic because
// it describes syntax rather than identity.
type UsageLine string
