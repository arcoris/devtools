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

// WithStatus returns a validated copy with Status replaced.
func (result ActionResult) WithStatus(status ActionStatus) (ActionResult, error) {
	next := result.Normalize()
	next.Status = status

	return NewActionResult(next)
}

// MustWithStatus returns a validated copy with Status replaced and panics on
// invalid input.
func (result ActionResult) MustWithStatus(status ActionStatus) ActionResult {
	next, err := result.WithStatus(status)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMessage returns a validated copy with Message replaced.
func (result ActionResult) WithMessage(message string) (ActionResult, error) {
	next := result.Normalize()
	next.Message = message

	return NewActionResult(next)
}

// MustWithMessage returns a validated copy with Message replaced and panics on
// invalid input.
func (result ActionResult) MustWithMessage(message string) ActionResult {
	next, err := result.WithMessage(message)
	if err != nil {
		panic(err)
	}

	return next
}

// WithData returns a validated copy with Data replaced.
//
// Data is not deep-copied or inspected by the command kernel.
func (result ActionResult) WithData(data any) (ActionResult, error) {
	next := result.Normalize()
	next.Data = data

	return NewActionResult(next)
}

// MustWithData returns a validated copy with Data replaced and panics on
// invalid input.
func (result ActionResult) MustWithData(data any) ActionResult {
	next, err := result.WithData(data)
	if err != nil {
		panic(err)
	}

	return next
}

// WithArtifacts returns a validated copy with artifact references replaced.
func (result ActionResult) WithArtifacts(artifacts []ActionArtifact) (ActionResult, error) {
	next := result.Normalize()
	next.Artifacts = cloneActionArtifacts(artifacts)

	return NewActionResult(next)
}

// MustWithArtifacts returns a validated copy with artifact references replaced
// and panics on invalid input.
func (result ActionResult) MustWithArtifacts(artifacts []ActionArtifact) ActionResult {
	next, err := result.WithArtifacts(artifacts)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutArtifacts returns a validated copy without artifact references.
func (result ActionResult) WithoutArtifacts() ActionResult {
	next := result.Normalize()
	next.Artifacts = nil

	return MustActionResult(next)
}

// WithArtifact returns a validated copy with one artifact appended.
func (result ActionResult) WithArtifact(artifact ActionArtifact) (ActionResult, error) {
	next := result.Normalize()
	next.Artifacts = append(next.Artifacts, artifact)

	return NewActionResult(next)
}

// MustWithArtifact returns a validated copy with one artifact appended and
// panics on invalid input.
func (result ActionResult) MustWithArtifact(artifact ActionArtifact) ActionResult {
	next, err := result.WithArtifact(artifact)
	if err != nil {
		panic(err)
	}

	return next
}

// WithWarnings returns a validated copy with warnings replaced.
func (result ActionResult) WithWarnings(warnings []ActionWarning) (ActionResult, error) {
	next := result.Normalize()
	next.Warnings = cloneActionWarnings(warnings)

	return NewActionResult(next)
}

// MustWithWarnings returns a validated copy with warnings replaced and panics
// on invalid input.
func (result ActionResult) MustWithWarnings(warnings []ActionWarning) ActionResult {
	next, err := result.WithWarnings(warnings)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutWarnings returns a validated copy without warnings.
func (result ActionResult) WithoutWarnings() ActionResult {
	next := result.Normalize()
	next.Warnings = nil

	return MustActionResult(next)
}

// WithWarning returns a validated copy with one warning appended.
func (result ActionResult) WithWarning(warning ActionWarning) (ActionResult, error) {
	next := result.Normalize()
	next.Warnings = append(next.Warnings, warning)

	return NewActionResult(next)
}

// MustWithWarning returns a validated copy with one warning appended and panics
// on invalid input.
func (result ActionResult) MustWithWarning(warning ActionWarning) ActionResult {
	next, err := result.WithWarning(warning)
	if err != nil {
		panic(err)
	}

	return next
}

// WithFields returns a validated copy with result metadata fields replaced.
func (result ActionResult) WithFields(fields map[string]string) (ActionResult, error) {
	next := result.Normalize()
	next.Fields = cloneActionStringMap(fields)

	return NewActionResult(next)
}

// MustWithFields returns a validated copy with result metadata fields replaced
// and panics on invalid input.
func (result ActionResult) MustWithFields(fields map[string]string) ActionResult {
	next, err := result.WithFields(fields)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutFields returns a validated copy without result metadata fields.
func (result ActionResult) WithoutFields() ActionResult {
	next := result.Normalize()
	next.Fields = nil

	return MustActionResult(next)
}

// WithField returns a validated copy with one result metadata field added or
// replaced.
func (result ActionResult) WithField(key string, value string) (ActionResult, error) {
	next := result.Normalize()
	if next.Fields == nil {
		next.Fields = make(map[string]string)
	}

	next.Fields[key] = value

	return NewActionResult(next)
}

// MustWithField returns a validated copy with one field added or replaced and
// panics on invalid input.
func (result ActionResult) MustWithField(key string, value string) ActionResult {
	next, err := result.WithField(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutField returns a validated copy without one result metadata field.
func (result ActionResult) WithoutField(key string) ActionResult {
	next := result.Normalize()
	delete(next.Fields, key)

	return MustActionResult(next)
}
