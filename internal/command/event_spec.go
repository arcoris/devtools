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

// spec returns a detached construction spec.
func (event Event) spec() EventSpec {
	var result *Result
	if event.result != nil {
		value := *event.result
		result = &value
	}

	id := ""
	if event.hasID {
		id = event.id.String()
	}

	return EventSpec{
		ID:         id,
		Kind:       event.kind,
		Severity:   event.severity,
		OccurredAt: event.occurredAt,
		CommandID:  event.commandID,
		Message:    event.message,
		Fields:     cloneEventStringMap(event.fields),
		Artifacts:  cloneEventArtifacts(event.artifacts),
		Result:     result,
		Labels:     cloneEventStrings(event.labels),
		Metadata:   event.metadata,
		Visibility: event.visibility,
	}
}
