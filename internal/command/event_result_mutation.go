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

// WithResult returns a validated copy with Result attached.
func (event Event) WithResult(result Result) (Event, error) {
	spec := event.spec()
	spec.Result = &result

	return NewEvent(spec)
}

// MustWithResult returns a validated copy with Result attached and panics on
// invalid input.
func (event Event) MustWithResult(result Result) Event {
	next, err := event.WithResult(result)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutResult returns a validated copy without Result.
func (event Event) WithoutResult() Event {
	spec := event.spec()
	spec.Result = nil

	return MustEvent(spec)
}
