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

// WithLabel returns a validated copy with one label appended.
func (event Event) WithLabel(label string) (Event, error) {
	spec := event.spec()
	spec.Labels = append(spec.Labels, label)

	return NewEvent(spec)
}

// MustWithLabel returns a validated copy with one label appended and panics on
// invalid input.
func (event Event) MustWithLabel(label string) Event {
	next, err := event.WithLabel(label)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutLabel returns a validated copy without label.
func (event Event) WithoutLabel(label string) Event {
	spec := event.spec()
	nextLabels := make([]string, 0, len(spec.Labels))

	for _, current := range spec.Labels {
		if current == label {
			continue
		}

		nextLabels = append(nextLabels, current)
	}

	spec.Labels = nextLabels

	return MustEvent(spec)
}
