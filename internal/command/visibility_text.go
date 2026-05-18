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

// MarshalText returns the canonical text form of visibility.
func (visibility Visibility) MarshalText() ([]byte, error) {
	if err := visibility.Validate(); err != nil {
		return nil, err
	}

	return []byte(visibility.String()), nil
}

// UnmarshalText parses text into visibility.
func (visibility *Visibility) UnmarshalText(text []byte) error {
	if visibility == nil {
		return fmt.Errorf("%w: cannot unmarshal into nil *Visibility", ErrInvalidVisibility)
	}

	parsed, err := ParseVisibility(string(text))
	if err != nil {
		return err
	}

	*visibility = parsed

	return nil
}
