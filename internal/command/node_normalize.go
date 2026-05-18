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

const nodeExampleNotePrefix = "Example: "

// newNodeDocumentation merges legacy documentation fields into the structured
// Documentation value used by Node.
func newNodeDocumentation(spec NodeSpec) (Documentation, error) {
	docSpec := spec.Documentation.Spec()

	if spec.Short != "" && docSpec.Summary != "" {
		short := normalizeDocumentationSingleLine(spec.Short)
		if short != docSpec.Summary {
			return Documentation{}, fmt.Errorf(
				"%w: short field conflicts with documentation summary",
				ErrInvalidNode,
			)
		}
	}

	if docSpec.Summary == "" {
		docSpec.Summary = spec.Short
	}

	if spec.Long != "" && docSpec.Description != "" {
		long := normalizeDocumentationBlock(spec.Long)
		if long != docSpec.Description {
			return Documentation{}, fmt.Errorf(
				"%w: long field conflicts with documentation description",
				ErrInvalidNode,
			)
		}
	}

	if docSpec.Description == "" {
		docSpec.Description = spec.Long
	}

	if !spec.Usage.IsZero() && !docSpec.Usage.IsZero() && !usageEqual(spec.Usage, docSpec.Usage) {
		return Documentation{}, fmt.Errorf(
			"%w: usage field conflicts with documentation usage",
			ErrInvalidNode,
		)
	}

	if docSpec.Usage.IsZero() {
		docSpec.Usage = spec.Usage
	}

	if spec.Example != "" {
		note := nodeExampleNotePrefix + normalizeDocumentationBlock(spec.Example)
		for _, existing := range docSpec.Notes {
			if strings.HasPrefix(existing, nodeExampleNotePrefix) && existing != note {
				return Documentation{}, fmt.Errorf(
					"%w: example field conflicts with documentation example note",
					ErrInvalidNode,
				)
			}
		}

		if !containsNodeString(docSpec.Notes, note) {
			docSpec.Notes = append(docSpec.Notes, note)
		}
	}

	documentation, err := NewDocumentation(docSpec)
	if err != nil {
		return Documentation{}, fmt.Errorf("%w: invalid documentation: %w", ErrInvalidNode, err)
	}

	return documentation, nil
}

// newNodeMetadata merges legacy deprecation text into structured Metadata.
func newNodeMetadata(spec NodeSpec) (Metadata, error) {
	metadataSpec := spec.Metadata.Spec()

	if spec.Deprecated != "" {
		if metadataSpec.Deprecation != nil && metadataSpec.Deprecation.Message != spec.Deprecated {
			return Metadata{}, fmt.Errorf(
				"%w: deprecated field conflicts with metadata deprecation message",
				ErrInvalidNode,
			)
		}

		if metadataSpec.Deprecation == nil {
			metadataSpec.Deprecation = &DeprecationSpec{Message: spec.Deprecated}
		}
	}

	metadata, err := NewMetadata(metadataSpec)
	if err != nil {
		return Metadata{}, fmt.Errorf("%w: invalid metadata: %w", ErrInvalidNode, err)
	}

	return metadata, nil
}

// newNodeVisibility converts the legacy Hidden field into explicit Visibility.
func newNodeVisibility(spec NodeSpec) Visibility {
	if !spec.Visibility.IsZero() {
		return spec.Visibility
	}

	return VisibilityFromHidden(spec.Hidden)
}

// nodeExampleFromDocumentation extracts the compatibility Example field from
// structured documentation notes.
func nodeExampleFromDocumentation(documentation Documentation) string {
	for _, note := range documentation.Notes() {
		if strings.HasPrefix(note, nodeExampleNotePrefix) {
			return strings.TrimPrefix(note, nodeExampleNotePrefix)
		}
	}

	return ""
}

// containsNodeString reports whether values contains target.
func containsNodeString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

// usageEqual reports whether two Usage declarations contain the same lines in
// the same declaration order.
func usageEqual(left Usage, right Usage) bool {
	leftLines := left.LineStrings()
	rightLines := right.LineStrings()
	if len(leftLines) != len(rightLines) {
		return false
	}

	for index := range leftLines {
		if leftLines[index] != rightLines[index] {
			return false
		}
	}

	return true
}

// cloneTopics returns a detached copy of topics.
func cloneTopics(topics []Topic) []Topic {
	if topics == nil {
		return nil
	}

	out := make([]Topic, len(topics))
	copy(out, topics)

	return out
}
