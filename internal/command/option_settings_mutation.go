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

// WithPolicy returns a validated copy with Policy replaced.
func (option Option) WithPolicy(policy OptionPolicy) (Option, error) {
	spec := option.spec()
	spec.Policy = policy

	return NewOption(spec)
}

// MustWithPolicy returns a validated copy with Policy replaced and panics on
// invalid input.
func (option Option) MustWithPolicy(policy OptionPolicy) Option {
	next, err := option.WithPolicy(policy)
	if err != nil {
		panic(err)
	}

	return next
}
