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
	"sort"
	"time"
)

// ID returns the event ID and whether it is present.
func (event Event) ID() (EventID, bool) {
	if !event.hasID {
		return "", false
	}

	return event.id, true
}

// MustID returns the event ID and panics if it is absent.
func (event Event) MustID() EventID {
	if !event.hasID {
		panic("command event id is absent")
	}

	return event.id
}

// HasID reports whether the event has a stable ID.
func (event Event) HasID() bool {
	return event.hasID
}

// Kind returns the event kind.
func (event Event) Kind() EventKind {
	return event.kind
}

// Severity returns the event severity.
func (event Event) Severity() EventSeverity {
	return event.severity
}

// OccurredAt returns the event timestamp.
func (event Event) OccurredAt() time.Time {
	return event.occurredAt
}

// CommandID returns the optional command ID and whether it is set.
func (event Event) CommandID() (ID, bool) {
	if event.commandID == "" {
		return "", false
	}

	return event.commandID, true
}

// HasCommandID reports whether the event has a command ID.
func (event Event) HasCommandID() bool {
	return event.commandID != ""
}

// Message returns the event message.
func (event Event) Message() string {
	return event.message
}

// HasMessage reports whether Message is set.
func (event Event) HasMessage() bool {
	return event.message != ""
}

// Fields returns a detached copy of event fields.
func (event Event) Fields() map[string]string {
	return cloneEventStringMap(event.fields)
}

// Field returns one event field and whether it exists.
func (event Event) Field(key string) (string, bool) {
	value, ok := event.fields[key]

	return value, ok
}

// HasField reports whether one event field exists.
func (event Event) HasField(key string) bool {
	_, ok := event.Field(key)

	return ok
}

// FieldKeys returns event field keys in deterministic lexical order.
func (event Event) FieldKeys() []string {
	keys := make([]string, 0, len(event.fields))
	for key := range event.fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// Artifacts returns detached artifact references.
func (event Event) Artifacts() []Artifact {
	return cloneEventArtifacts(event.artifacts)
}

// HasArtifacts reports whether artifact references are present.
func (event Event) HasArtifacts() bool {
	return len(event.artifacts) > 0
}

// Artifact returns an artifact by ID.
func (event Event) Artifact(id ArtifactID) (Artifact, bool) {
	for _, artifact := range event.artifacts {
		if artifact.ID() == id {
			return artifact, true
		}
	}

	return Artifact{}, false
}

// HasArtifact reports whether an artifact with id exists.
func (event Event) HasArtifact(id ArtifactID) bool {
	_, ok := event.Artifact(id)

	return ok
}

// ArtifactIDs returns artifact IDs in event declaration order.
func (event Event) ArtifactIDs() []ArtifactID {
	ids := make([]ArtifactID, len(event.artifacts))
	for index, artifact := range event.artifacts {
		ids[index] = artifact.ID()
	}

	return ids
}

// SortedArtifactIDs returns artifact IDs in deterministic lexical order.
func (event Event) SortedArtifactIDs() []ArtifactID {
	ids := event.ArtifactIDs()
	sort.Slice(ids, func(i int, j int) bool {
		return ids[i].String() < ids[j].String()
	})

	return ids
}

// Result returns the attached result and whether it is present.
func (event Event) Result() (Result, bool) {
	if event.result == nil {
		return Result{}, false
	}

	return *event.result, true
}

// HasResult reports whether the event carries a result.
func (event Event) HasResult() bool {
	return event.result != nil
}

// Labels returns detached event labels.
func (event Event) Labels() []string {
	return cloneEventStrings(event.labels)
}

// HasLabels reports whether event labels are present.
func (event Event) HasLabels() bool {
	return len(event.labels) > 0
}

// HasLabel reports whether label is present.
func (event Event) HasLabel(label string) bool {
	for _, current := range event.labels {
		if current == label {
			return true
		}
	}

	return false
}

// SortedLabels returns event labels in deterministic lexical order.
func (event Event) SortedLabels() []string {
	labels := event.Labels()
	sort.Strings(labels)

	return labels
}

// Metadata returns event metadata.
func (event Event) Metadata() Metadata {
	return event.metadata
}

// Visibility returns event visibility.
func (event Event) Visibility() Visibility {
	return event.visibility
}

// IsVisibleByDefault reports whether default reports/logs/discovery should
// expose the event.
func (event Event) IsVisibleByDefault() bool {
	return event.visibility.IsDiscoverableByDefault()
}

// IsTrace reports whether Severity is trace.
func (event Event) IsTrace() bool {
	return event.severity == EventSeverityTrace
}

// IsDebug reports whether Severity is debug.
func (event Event) IsDebug() bool {
	return event.severity == EventSeverityDebug
}

// IsInfo reports whether Severity is info.
func (event Event) IsInfo() bool {
	return event.severity == EventSeverityInfo
}

// IsWarning reports whether Severity is warning.
func (event Event) IsWarning() bool {
	return event.severity == EventSeverityWarning
}

// IsError reports whether Severity is error.
func (event Event) IsError() bool {
	return event.severity == EventSeverityError
}
