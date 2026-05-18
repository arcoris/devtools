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

import "time"

// RuntimeClock provides timestamps for runtime lifecycle records.
type RuntimeClock interface {
	Now() time.Time
}

// SystemRuntimeClock uses time.Now.
type SystemRuntimeClock struct{}

// Now returns current wall-clock time.
func (SystemRuntimeClock) Now() time.Time {
	return time.Now().UTC()
}

// FixedRuntimeClock returns a stable timestamp.
//
// FixedRuntimeClock is primarily useful in tests.
type FixedRuntimeClock struct {
	Time time.Time
}

// Now returns the configured timestamp, or Unix epoch when Time is zero.
func (clock FixedRuntimeClock) Now() time.Time {
	if clock.Time.IsZero() {
		return time.Unix(0, 0).UTC()
	}

	return clock.Time.UTC()
}
