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
	"strings"
	"time"
)

// now returns current runtime time in UTC.
func (runtime Runtime) now() time.Time {
	return runtime.clock.Now().UTC()
}

// normalizeRuntimePanicMessage returns compact panic text.
func normalizeRuntimePanicMessage(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return "panic"
	}

	if len(raw) <= maxRuntimePanicMessageLength {
		return raw
	}

	return raw[:maxRuntimePanicMessageLength]
}

// intPointer returns a pointer to value.
func intPointer(value int) *int {
	return &value
}
