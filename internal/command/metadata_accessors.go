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

// Spec returns a detached construction spec for metadata.
func (metadata Metadata) Spec() MetadataSpec {
	return metadata.spec()
}

// IsZero reports whether metadata has no fields set.
func (metadata Metadata) IsZero() bool {
	return metadata.owner == "" &&
		metadata.area == "" &&
		metadata.since == "" &&
		metadata.deprecation == nil &&
		len(metadata.annotations) == 0
}

// Owner returns the stable owner key.
func (metadata Metadata) Owner() string {
	return metadata.owner
}

// HasOwner reports whether Owner is set.
func (metadata Metadata) HasOwner() bool {
	return metadata.owner != ""
}

// Area returns the stable area key.
func (metadata Metadata) Area() string {
	return metadata.area
}

// HasArea reports whether Area is set.
func (metadata Metadata) HasArea() bool {
	return metadata.area != ""
}

// Since returns the compact introduction marker.
func (metadata Metadata) Since() string {
	return metadata.since
}

// HasSince reports whether Since is set.
func (metadata Metadata) HasSince() bool {
	return metadata.since != ""
}

// Deprecation returns deprecation metadata and whether it is set.
//
// The returned Deprecation is a value copy.
func (metadata Metadata) Deprecation() (Deprecation, bool) {
	if metadata.deprecation == nil {
		return Deprecation{}, false
	}

	return *metadata.deprecation, true
}

// HasDeprecation reports whether deprecation metadata is set.
func (metadata Metadata) HasDeprecation() bool {
	return metadata.deprecation != nil
}

// IsDeprecated reports whether deprecation metadata is set.
func (metadata Metadata) IsDeprecated() bool {
	return metadata.HasDeprecation()
}

// spec returns a detached construction spec.
func (metadata Metadata) spec() MetadataSpec {
	var deprecation *DeprecationSpec
	if metadata.deprecation != nil {
		spec := metadata.deprecation.spec()
		deprecation = &spec
	}

	return MetadataSpec{
		Owner:       metadata.owner,
		Area:        metadata.area,
		Since:       metadata.since,
		Deprecation: deprecation,
		Annotations: cloneStringMap(metadata.annotations),
	}
}
