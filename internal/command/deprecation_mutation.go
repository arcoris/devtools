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

// WithSince returns a validated copy with Since replaced.
func (deprecation Deprecation) WithSince(since string) (Deprecation, error) {
	spec := deprecation.spec()
	spec.Since = since

	return NewDeprecation(spec)
}

// MustWithSince returns a validated copy with Since replaced and panics on
// invalid input.
func (deprecation Deprecation) MustWithSince(since string) Deprecation {
	next, err := deprecation.WithSince(since)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMessage returns a validated copy with Message replaced.
func (deprecation Deprecation) WithMessage(message string) (Deprecation, error) {
	spec := deprecation.spec()
	spec.Message = message

	return NewDeprecation(spec)
}

// MustWithMessage returns a validated copy with Message replaced and panics on
// invalid input.
func (deprecation Deprecation) MustWithMessage(message string) Deprecation {
	next, err := deprecation.WithMessage(message)
	if err != nil {
		panic(err)
	}

	return next
}

// WithReplacement returns a validated copy with Replacement replaced.
func (deprecation Deprecation) WithReplacement(replacement Path) (Deprecation, error) {
	spec := deprecation.spec()
	spec.Replacement = replacement

	return NewDeprecation(spec)
}

// MustWithReplacement returns a validated copy with Replacement replaced and
// panics on invalid input.
func (deprecation Deprecation) MustWithReplacement(replacement Path) Deprecation {
	next, err := deprecation.WithReplacement(replacement)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutReplacement returns a validated copy without a replacement path.
func (deprecation Deprecation) WithoutReplacement() Deprecation {
	spec := deprecation.spec()
	spec.Replacement = RootPath()

	return MustDeprecation(spec)
}
