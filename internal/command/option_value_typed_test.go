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
	"testing"
	"time"
)

// TestOptionValueListTypedAccessors verifies list typed accessors.
func TestOptionValueListTypedAccessors(t *testing.T) {
	t.Parallel()

	textValue := MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./...", "./internal/...")
	textValues, err := textValue.TextValues()
	if err != nil {
		t.Fatalf("TextValues() returned unexpected error: %v", err)
	}

	assertOptionValueStrings(t, textValues, []string{"./...", "./internal/..."})

	enumValue := MustListOptionValue("format", OptionKindEnumList, OptionSourceCommandLine, "text", "json")
	enumValues, err := enumValue.TextValues()
	if err != nil {
		t.Fatalf("TextValues(enum-list) returned unexpected error: %v", err)
	}

	assertOptionValueStrings(t, enumValues, []string{"text", "json"})

	intValue := MustListOptionValue("count", OptionKindIntList, OptionSourceCommandLine, "1", "2")
	intValues, err := intValue.IntValues()
	if err != nil {
		t.Fatalf("IntValues() returned unexpected error: %v", err)
	}

	if len(intValues) != 2 || intValues[0] != 1 || intValues[1] != 2 {
		t.Fatalf("IntValues() = %v, want [1 2]", intValues)
	}

	int64Value := MustListOptionValue("size", OptionKindInt64List, OptionSourceCommandLine, "-1", "2")
	int64Values, err := int64Value.Int64Values()
	if err != nil {
		t.Fatalf("Int64Values() returned unexpected error: %v", err)
	}

	if len(int64Values) != 2 || int64Values[0] != -1 || int64Values[1] != 2 {
		t.Fatalf("Int64Values() = %v, want [-1 2]", int64Values)
	}

	uintValue := MustListOptionValue("count", OptionKindUintList, OptionSourceCommandLine, "1", "2")
	uintValues, err := uintValue.UintValues()
	if err != nil {
		t.Fatalf("UintValues() returned unexpected error: %v", err)
	}

	if len(uintValues) != 2 || uintValues[0] != 1 || uintValues[1] != 2 {
		t.Fatalf("UintValues() = %v, want [1 2]", uintValues)
	}

	uint64Value := MustListOptionValue("size", OptionKindUint64List, OptionSourceCommandLine, "1", "2")
	uint64Values, err := uint64Value.Uint64Values()
	if err != nil {
		t.Fatalf("Uint64Values() returned unexpected error: %v", err)
	}

	if len(uint64Values) != 2 || uint64Values[0] != 1 || uint64Values[1] != 2 {
		t.Fatalf("Uint64Values() = %v, want [1 2]", uint64Values)
	}

	floatValue := MustListOptionValue("ratio", OptionKindFloat64List, OptionSourceCommandLine, "1.5", "2.25")
	floatValues, err := floatValue.Float64Values()
	if err != nil {
		t.Fatalf("Float64Values() returned unexpected error: %v", err)
	}

	if len(floatValues) != 2 || floatValues[0] != 1.5 || floatValues[1] != 2.25 {
		t.Fatalf("Float64Values() = %v, want [1.5 2.25]", floatValues)
	}

	durationValue := MustListOptionValue("timeout", OptionKindDurationList, OptionSourceCommandLine, "1s", "2s")
	durationValues, err := durationValue.DurationValues()
	if err != nil {
		t.Fatalf("DurationValues() returned unexpected error: %v", err)
	}

	if len(durationValues) != 2 || durationValues[0] != time.Second || durationValues[1] != 2*time.Second {
		t.Fatalf("DurationValues() = %v, want [1s 2s]", durationValues)
	}
}

// TestOptionValueTypedAccessorsRejectWrongKind verifies accessor kind checks.
func TestOptionValueTypedAccessorsRejectWrongKind(t *testing.T) {
	t.Parallel()

	value := MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "out.txt")

	if _, err := value.Bool(); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("Bool() error = %v, want ErrInvalidOptionValue", err)
	}

	if _, err := MustListOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./...").Text(); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("Text(list) error = %v, want ErrInvalidOptionValue", err)
	}

	if _, err := value.IntValues(); !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("IntValues() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueTypedAccessors verifies scalar typed accessors.
func TestOptionValueTypedAccessors(t *testing.T) {
	t.Parallel()

	boolValue := MustScalarOptionValue("verbose", OptionKindBool, OptionSourceCommandLine, "true")
	if got, err := boolValue.Bool(); err != nil || !got {
		t.Fatalf("Bool() = %v, %v; want true, nil", got, err)
	}

	textValue := MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "out.txt")
	if got, err := textValue.Text(); err != nil || got != "out.txt" {
		t.Fatalf("Text() = %q, %v; want out.txt, nil", got, err)
	}

	enumValue := MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json")
	if got, err := enumValue.Text(); err != nil || got != "json" {
		t.Fatalf("Text(enum) = %q, %v; want json, nil", got, err)
	}

	intValue := MustScalarOptionValue("count", OptionKindInt, OptionSourceCommandLine, "-10")
	if got, err := intValue.Int(); err != nil || got != -10 {
		t.Fatalf("Int() = %d, %v; want -10, nil", got, err)
	}

	int64Value := MustScalarOptionValue("count", OptionKindInt64, OptionSourceCommandLine, "-10")
	if got, err := int64Value.Int64(); err != nil || got != -10 {
		t.Fatalf("Int64() = %d, %v; want -10, nil", got, err)
	}

	uintValue := MustScalarOptionValue("count", OptionKindUint, OptionSourceCommandLine, "10")
	if got, err := uintValue.Uint(); err != nil || got != 10 {
		t.Fatalf("Uint() = %d, %v; want 10, nil", got, err)
	}

	uint64Value := MustScalarOptionValue("count", OptionKindUint64, OptionSourceCommandLine, "10")
	if got, err := uint64Value.Uint64(); err != nil || got != 10 {
		t.Fatalf("Uint64() = %d, %v; want 10, nil", got, err)
	}

	floatValue := MustScalarOptionValue("ratio", OptionKindFloat64, OptionSourceCommandLine, "3.5")
	if got, err := floatValue.Float64(); err != nil || got != 3.5 {
		t.Fatalf("Float64() = %f, %v; want 3.5, nil", got, err)
	}

	durationValue := MustScalarOptionValue("timeout", OptionKindDuration, OptionSourceCommandLine, "10s")
	if got, err := durationValue.Duration(); err != nil || got != 10*time.Second {
		t.Fatalf("Duration() = %v, %v; want 10s, nil", got, err)
	}
}
