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

import "fmt"

// MarshalText returns the canonical text form of group.
func (group Group) MarshalText() ([]byte, error) {
	if err := group.Validate(); err != nil {
		return nil, err
	}

	return []byte(group.String()), nil
}

// UnmarshalText parses text into group.
func (group *Group) UnmarshalText(text []byte) error {
	if group == nil {
		return fmt.Errorf("%w: cannot unmarshal into nil *Group", ErrInvalidGroup)
	}

	parsed, err := ParseGroup(string(text))
	if err != nil {
		return err
	}

	*group = parsed

	return nil
}
