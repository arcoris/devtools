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

// IsBool reports whether kind describes a boolean option.
func (kind OptionKind) IsBool() bool {
	return kind == OptionKindBool
}

// IsString reports whether kind describes an unconstrained string option.
func (kind OptionKind) IsString() bool {
	return kind == OptionKindString
}

// IsEnum reports whether kind describes an enum option.
func (kind OptionKind) IsEnum() bool {
	return kind == OptionKindEnum
}

// IsEnumLike reports whether kind describes enum values.
func (kind OptionKind) IsEnumLike() bool {
	switch kind {
	case OptionKindEnum, OptionKindEnumList:
		return true
	default:
		return false
	}
}

// IsInteger reports whether kind describes a signed or unsigned integer scalar.
func (kind OptionKind) IsInteger() bool {
	switch kind {
	case OptionKindInt, OptionKindInt64, OptionKindUint, OptionKindUint64:
		return true
	default:
		return false
	}
}

// IsSignedInteger reports whether kind describes a signed integer scalar.
func (kind OptionKind) IsSignedInteger() bool {
	switch kind {
	case OptionKindInt, OptionKindInt64:
		return true
	default:
		return false
	}
}

// IsUnsignedInteger reports whether kind describes an unsigned integer scalar.
func (kind OptionKind) IsUnsignedInteger() bool {
	switch kind {
	case OptionKindUint, OptionKindUint64:
		return true
	default:
		return false
	}
}

// IsFloat reports whether kind describes a floating-point scalar.
func (kind OptionKind) IsFloat() bool {
	return kind == OptionKindFloat64
}

// IsDuration reports whether kind describes a duration scalar.
func (kind OptionKind) IsDuration() bool {
	return kind == OptionKindDuration
}

// IsNumeric reports whether kind describes a numeric scalar.
func (kind OptionKind) IsNumeric() bool {
	return kind.IsInteger() || kind.IsFloat()
}

// IsScalar reports whether kind describes exactly one value.
func (kind OptionKind) IsScalar() bool {
	switch kind {
	case OptionKindBool,
		OptionKindString,
		OptionKindEnum,
		OptionKindInt,
		OptionKindInt64,
		OptionKindUint,
		OptionKindUint64,
		OptionKindFloat64,
		OptionKindDuration:
		return true
	default:
		return false
	}
}

// IsList reports whether kind describes multiple values.
func (kind OptionKind) IsList() bool {
	switch kind {
	case OptionKindStringList,
		OptionKindEnumList,
		OptionKindIntList,
		OptionKindInt64List,
		OptionKindUintList,
		OptionKindUint64List,
		OptionKindFloat64List,
		OptionKindDurationList:
		return true
	default:
		return false
	}
}

// IsRepeatable reports whether adapters may accept the option more than once.
func (kind OptionKind) IsRepeatable() bool {
	return kind.IsList()
}

// RequiresValue reports whether normal CLI syntax should require a value.
func (kind OptionKind) RequiresValue() bool {
	return kind.IsKnown() && kind != OptionKindBool
}

// AllowsImplicitBoolean reports whether adapters may support implicit true when
// the option appears without a value.
func (kind OptionKind) AllowsImplicitBoolean() bool {
	return kind == OptionKindBool
}

// ElementKind returns the scalar kind for a list kind.
func (kind OptionKind) ElementKind() OptionKind {
	switch kind {
	case OptionKindStringList:
		return OptionKindString
	case OptionKindEnumList:
		return OptionKindEnum
	case OptionKindIntList:
		return OptionKindInt
	case OptionKindInt64List:
		return OptionKindInt64
	case OptionKindUintList:
		return OptionKindUint
	case OptionKindUint64List:
		return OptionKindUint64
	case OptionKindFloat64List:
		return OptionKindFloat64
	case OptionKindDurationList:
		return OptionKindDuration
	default:
		return kind
	}
}

// ListKind returns the list kind for a scalar kind.
func (kind OptionKind) ListKind() OptionKind {
	switch kind {
	case OptionKindString:
		return OptionKindStringList
	case OptionKindEnum:
		return OptionKindEnumList
	case OptionKindInt:
		return OptionKindIntList
	case OptionKindInt64:
		return OptionKindInt64List
	case OptionKindUint:
		return OptionKindUintList
	case OptionKindUint64:
		return OptionKindUint64List
	case OptionKindFloat64:
		return OptionKindFloat64List
	case OptionKindDuration:
		return OptionKindDurationList
	default:
		return kind
	}
}

// CanHaveAllowedValues reports whether an option may attach allowed values.
func (kind OptionKind) CanHaveAllowedValues() bool {
	switch kind {
	case OptionKindString, OptionKindStringList, OptionKindEnum, OptionKindEnumList:
		return true
	default:
		return false
	}
}

// RequiresAllowedValues reports whether a declaration must provide allowed
// values for this kind.
func (kind OptionKind) RequiresAllowedValues() bool {
	switch kind {
	case OptionKindEnum, OptionKindEnumList:
		return true
	default:
		return false
	}
}

// CanHaveRange reports whether numeric or duration range constraints are
// meaningful for this kind.
func (kind OptionKind) CanHaveRange() bool {
	switch kind {
	case OptionKindInt,
		OptionKindInt64,
		OptionKindUint,
		OptionKindUint64,
		OptionKindFloat64,
		OptionKindDuration,
		OptionKindIntList,
		OptionKindInt64List,
		OptionKindUintList,
		OptionKindUint64List,
		OptionKindFloat64List,
		OptionKindDurationList:
		return true
	default:
		return false
	}
}

// ValueMetavar returns a conventional metavar for help and documentation.
func (kind OptionKind) ValueMetavar() string {
	switch kind {
	case OptionKindBool:
		return "BOOL"
	case OptionKindString:
		return "STRING"
	case OptionKindEnum:
		return "VALUE"
	case OptionKindInt, OptionKindInt64:
		return "INT"
	case OptionKindUint, OptionKindUint64:
		return "UINT"
	case OptionKindFloat64:
		return "FLOAT"
	case OptionKindDuration:
		return "DURATION"
	case OptionKindStringList:
		return "STRING"
	case OptionKindEnumList:
		return "VALUE"
	case OptionKindIntList, OptionKindInt64List:
		return "INT"
	case OptionKindUintList, OptionKindUint64List:
		return "UINT"
	case OptionKindFloat64List:
		return "FLOAT"
	case OptionKindDurationList:
		return "DURATION"
	default:
		return "VALUE"
	}
}
