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

// WithArtifact returns a validated copy with artifact appended or replaced by
// ID.
func (result Result) WithArtifact(artifact Artifact) (Result, error) {
	spec := result.spec()

	replaced := false
	for index, existing := range spec.Artifacts {
		if existing.ID() == artifact.ID() {
			spec.Artifacts[index] = artifact
			replaced = true

			break
		}
	}

	if !replaced {
		spec.Artifacts = append(spec.Artifacts, artifact)
	}

	return NewResult(spec)
}

// MustWithArtifact returns a validated copy with artifact appended or replaced
// and panics on invalid input.
func (result Result) MustWithArtifact(artifact Artifact) Result {
	next, err := result.WithArtifact(artifact)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutArtifact returns a validated copy without the artifact ID.
func (result Result) WithoutArtifact(id ArtifactID) Result {
	spec := result.spec()
	nextArtifacts := make([]Artifact, 0, len(spec.Artifacts))

	for _, artifact := range spec.Artifacts {
		if artifact.ID() == id {
			continue
		}

		nextArtifacts = append(nextArtifacts, artifact)
	}

	spec.Artifacts = nextArtifacts

	return MustResult(spec)
}
