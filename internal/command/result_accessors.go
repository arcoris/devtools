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

// Status returns the final lifecycle status.
func (result Result) Status() ResultStatus {
	return result.status
}

// Message returns the compact human-facing result message.
func (result Result) Message() string {
	return result.message
}

// HasMessage reports whether Message is set.
func (result Result) HasMessage() bool {
	return result.message != ""
}

// StartedAt returns the execution start timestamp and whether it is set.
func (result Result) StartedAt() (time.Time, bool) {
	return result.startedAt, !result.startedAt.IsZero()
}

// FinishedAt returns the execution finish timestamp and whether it is set.
func (result Result) FinishedAt() (time.Time, bool) {
	return result.finishedAt, !result.finishedAt.IsZero()
}

// HasTiming reports whether both StartedAt and FinishedAt are set.
func (result Result) HasTiming() bool {
	return !result.startedAt.IsZero() && !result.finishedAt.IsZero()
}

// Duration returns FinishedAt - StartedAt and whether both timestamps are set.
func (result Result) Duration() (time.Duration, bool) {
	if !result.HasTiming() {
		return 0, false
	}

	return result.finishedAt.Sub(result.startedAt), true
}

// ExitCode returns the explicitly declared adapter-facing exit code and whether
// it is set.
func (result Result) ExitCode() (int, bool) {
	return result.exitCode, result.hasExit
}

// RecommendedExitCode returns a deterministic process exit code.
//
// If an explicit exit code is set, it is returned. Otherwise the code is derived
// from Status:
//
//   - ok: 0;
//   - skipped: 0;
//   - failed: 1;
//   - canceled: 130.
//
// The result layer does not call os.Exit. CLI adapters may use this method when
// converting Result into process behavior.
func (result Result) RecommendedExitCode() int {
	if result.hasExit {
		return result.exitCode
	}

	switch result.status {
	case ResultStatusOK, ResultStatusSkipped:
		return 0
	case ResultStatusFailed:
		return 1
	case ResultStatusCanceled:
		return 130
	default:
		return 1
	}
}

// Artifacts returns detached artifact references.
func (result Result) Artifacts() []Artifact {
	return cloneResultArtifacts(result.artifacts)
}

// HasArtifacts reports whether artifacts are present.
func (result Result) HasArtifacts() bool {
	return len(result.artifacts) > 0
}

// Artifact returns an artifact by ID.
func (result Result) Artifact(id ArtifactID) (Artifact, bool) {
	for _, artifact := range result.artifacts {
		if artifact.ID() == id {
			return artifact, true
		}
	}

	return Artifact{}, false
}

// HasArtifact reports whether an artifact with id exists.
func (result Result) HasArtifact(id ArtifactID) bool {
	_, ok := result.Artifact(id)

	return ok
}

// ArtifactIDs returns artifact IDs in declaration order.
func (result Result) ArtifactIDs() []ArtifactID {
	ids := make([]ArtifactID, len(result.artifacts))
	for index, artifact := range result.artifacts {
		ids[index] = artifact.ID()
	}

	return ids
}

// SortedArtifactIDs returns artifact IDs in deterministic lexical order.
func (result Result) SortedArtifactIDs() []ArtifactID {
	ids := result.ArtifactIDs()
	sort.Slice(ids, func(i int, j int) bool {
		return ids[i].String() < ids[j].String()
	})

	return ids
}

// Warnings returns detached result warnings.
func (result Result) Warnings() []ResultWarning {
	return cloneResultWarnings(result.warnings)
}

// HasWarnings reports whether warnings are present.
func (result Result) HasWarnings() bool {
	return len(result.warnings) > 0
}

// WarningCount returns the number of warnings.
func (result Result) WarningCount() int {
	return len(result.warnings)
}

// Fields returns a detached copy of result fields.
func (result Result) Fields() map[string]string {
	return cloneResultStringMap(result.fields)
}

// Field returns one result metadata field and whether it exists.
func (result Result) Field(key string) (string, bool) {
	value, ok := result.fields[key]

	return value, ok
}

// HasField reports whether one result metadata field exists.
func (result Result) HasField(key string) bool {
	_, ok := result.Field(key)

	return ok
}

// FieldKeys returns result metadata field keys in deterministic lexical order.
func (result Result) FieldKeys() []string {
	keys := make([]string, 0, len(result.fields))
	for key := range result.fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// Metadata returns result metadata.
func (result Result) Metadata() Metadata {
	return result.metadata
}

// Visibility returns result visibility.
func (result Result) Visibility() Visibility {
	return result.visibility
}

// IsVisibleByDefault reports whether default reports/docs/discovery should
// expose the result.
func (result Result) IsVisibleByDefault() bool {
	return result.visibility.IsDiscoverableByDefault()
}

// IsOK reports whether Status is ok.
func (result Result) IsOK() bool {
	return result.status == ResultStatusOK
}

// IsSkipped reports whether Status is skipped.
func (result Result) IsSkipped() bool {
	return result.status == ResultStatusSkipped
}

// IsFailed reports whether Status is failed.
func (result Result) IsFailed() bool {
	return result.status == ResultStatusFailed
}

// IsCanceled reports whether Status is canceled.
func (result Result) IsCanceled() bool {
	return result.status == ResultStatusCanceled
}

// IsSuccessful reports whether Status represents a non-failure outcome.
func (result Result) IsSuccessful() bool {
	switch result.status {
	case ResultStatusOK, ResultStatusSkipped:
		return true
	default:
		return false
	}
}

// IsUnsuccessful reports whether Status represents a failure-like outcome.
func (result Result) IsUnsuccessful() bool {
	switch result.status {
	case ResultStatusFailed, ResultStatusCanceled:
		return true
	default:
		return false
	}
}
