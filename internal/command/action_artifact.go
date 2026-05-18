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
	"fmt"
	"strings"
)

// ActionArtifact is a lightweight reference to an artifact produced by an
// action.
//
// This type is only a result reference. It does not create files, choose output
// directories, manage retention, or define artifact layout. Those concerns
// belong to the artifact package.
type ActionArtifact struct {
	// Kind is a compact machine-facing artifact kind.
	//
	// Examples: "report", "profile", "coverage", "trace", "benchmark".
	Kind string

	// Path is a repository-relative or absolute artifact path.
	Path string

	// Description is an optional compact human-facing description.
	Description string
}

// NewActionArtifact validates fields and returns an ActionArtifact.
func NewActionArtifact(kind string, path string, description string) (ActionArtifact, error) {
	artifact := ActionArtifact{
		Kind:        kind,
		Path:        path,
		Description: description,
	}

	if err := artifact.Validate(); err != nil {
		return ActionArtifact{}, err
	}

	return artifact, nil
}

// MustActionArtifact validates fields and returns an ActionArtifact.
//
// MustActionArtifact panics on invalid input. It is intended for tests and
// static command wiring.
func MustActionArtifact(kind string, path string, description string) ActionArtifact {
	artifact, err := NewActionArtifact(kind, path, description)
	if err != nil {
		panic(err)
	}

	return artifact
}

// IsZero reports whether artifact has no fields set.
func (artifact ActionArtifact) IsZero() bool {
	return artifact.Kind == "" && artifact.Path == "" && artifact.Description == ""
}

// Validate verifies artifact reference structural rules.
func (artifact ActionArtifact) Validate() error {
	if err := validateActionResultFieldKey("artifact kind", artifact.Kind); err != nil {
		return err
	}

	if strings.TrimSpace(artifact.Path) == "" {
		return fmt.Errorf("%w: artifact path must not be blank", ErrInvalidActionResult)
	}

	if err := validateActionResultText("artifact path", artifact.Path, maxActionArtifactPathLength); err != nil {
		return err
	}

	if artifact.Description != "" {
		if err := validateActionResultText("artifact description", artifact.Description, maxActionDescriptionLength); err != nil {
			return err
		}
	}

	return nil
}
