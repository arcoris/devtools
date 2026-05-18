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

// IsDefault reports whether source is a declaration default.
func (source OptionSource) IsDefault() bool {
	return source == OptionSourceDefault
}

// IsInherited reports whether source is inherited from a wider option scope.
func (source OptionSource) IsInherited() bool {
	return source == OptionSourceInherited
}

// IsConfig reports whether source is configuration-backed.
func (source OptionSource) IsConfig() bool {
	return source == OptionSourceConfig
}

// IsEnvironment reports whether source is environment-backed.
func (source OptionSource) IsEnvironment() bool {
	return source == OptionSourceEnvironment
}

// IsRuntime reports whether source is runtime-supplied.
func (source OptionSource) IsRuntime() bool {
	return source == OptionSourceRuntime
}

// IsInteractive reports whether source was supplied through an interactive
// prompt.
func (source OptionSource) IsInteractive() bool {
	return source == OptionSourceInteractive
}

// IsCommandLine reports whether source was explicitly supplied through CLI
// syntax.
func (source OptionSource) IsCommandLine() bool {
	return source == OptionSourceCommandLine
}

// IsImplicit reports whether the value was not directly supplied for this
// concrete option by an explicit user input channel.
func (source OptionSource) IsImplicit() bool {
	switch source {
	case OptionSourceDefault, OptionSourceInherited, OptionSourceRuntime:
		return true
	default:
		return false
	}
}

// IsExplicit reports whether the value came from an explicit external input.
func (source OptionSource) IsExplicit() bool {
	switch source {
	case OptionSourceConfig,
		OptionSourceEnvironment,
		OptionSourceInteractive,
		OptionSourceCommandLine:
		return true
	default:
		return false
	}
}

// IsUserProvided reports whether source represents direct user-provided input.
func (source OptionSource) IsUserProvided() bool {
	switch source {
	case OptionSourceInteractive, OptionSourceCommandLine:
		return true
	default:
		return false
	}
}

// IsExternal reports whether source comes from outside the compiled command
// declaration.
func (source OptionSource) IsExternal() bool {
	return source.IsKnown() && source != OptionSourceDefault
}

// IsPersistent reports whether source can usually survive across invocations.
func (source OptionSource) IsPersistent() bool {
	switch source {
	case OptionSourceInherited, OptionSourceConfig, OptionSourceEnvironment:
		return true
	default:
		return false
	}
}

// IsEphemeral reports whether source is normally tied to one invocation.
func (source OptionSource) IsEphemeral() bool {
	switch source {
	case OptionSourceRuntime, OptionSourceInteractive, OptionSourceCommandLine:
		return true
	default:
		return false
	}
}

// Precedence returns the default precedence rank for source.
func (source OptionSource) Precedence() int {
	switch source {
	case OptionSourceDefault:
		return 10
	case OptionSourceInherited:
		return 20
	case OptionSourceConfig:
		return 30
	case OptionSourceEnvironment:
		return 40
	case OptionSourceRuntime:
		return 50
	case OptionSourceInteractive:
		return 60
	case OptionSourceCommandLine:
		return 70
	default:
		return 0
	}
}

// Overrides reports whether source has strictly higher default precedence than
// other.
func (source OptionSource) Overrides(other OptionSource) bool {
	return source.Precedence() > other.Precedence()
}

// CanBeOverriddenBy reports whether other has strictly higher default
// precedence than source.
func (source OptionSource) CanBeOverriddenBy(other OptionSource) bool {
	return other.Overrides(source)
}

// SamePrecedence reports whether both sources have the same non-zero precedence.
func (source OptionSource) SamePrecedence(other OptionSource) bool {
	left := source.Precedence()
	right := other.Precedence()

	return left != 0 && left == right
}
