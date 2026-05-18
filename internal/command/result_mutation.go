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
func (result Result) WithStatus(status ResultStatus) (Result, error) {
	spec := result.spec()
	spec.Status = status

	return NewResult(spec)
}

// MustWithStatus returns a validated copy with Status replaced and panics on
// invalid input.
func (result Result) MustWithStatus(status ResultStatus) Result {
	next, err := result.WithStatus(status)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMessage returns a validated copy with Message replaced.
func (result Result) WithMessage(message string) (Result, error) {
	spec := result.spec()
	spec.Message = message

	return NewResult(spec)
}

// MustWithMessage returns a validated copy with Message replaced and panics on
// invalid input.
func (result Result) MustWithMessage(message string) Result {
	next, err := result.WithMessage(message)
	if err != nil {
		panic(err)
	}

	return next
}
