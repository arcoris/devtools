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

// VisibilityExposurePolicy captures adapter-neutral exposure decisions for a
// visibility state.
type VisibilityExposurePolicy struct {
	// UserFacing reports whether the value belongs to the user-facing command
	// contract.
	UserFacing bool

	// DiscoverableByDefault reports whether default command-discovery surfaces
	// should include the command.
	DiscoverableByDefault bool

	// ShownInDefaultHelp reports whether default help output should include the
	// command.
	ShownInDefaultHelp bool

	// ShownInDefaultDocs reports whether default generated documentation should
	// include the command.
	ShownInDefaultDocs bool

	// AllowsDefaultInvocation reports whether normal adapters should allow direct
	// invocation without an explicit internal/debug mode.
	AllowsDefaultInvocation bool

	// RequiresExplicitExposure reports whether an adapter or generator mode must
	// explicitly opt in before exposing the command.
	RequiresExplicitExposure bool
}

// ExposurePolicy returns the adapter-neutral exposure policy for visibility.
func (visibility Visibility) ExposurePolicy() VisibilityExposurePolicy {
	switch visibility {
	case VisibilityPublic:
		return VisibilityExposurePolicy{
			UserFacing:               true,
			DiscoverableByDefault:    true,
			ShownInDefaultHelp:       true,
			ShownInDefaultDocs:       true,
			AllowsDefaultInvocation:  true,
			RequiresExplicitExposure: false,
		}
	case VisibilityHidden:
		return VisibilityExposurePolicy{
			UserFacing:               true,
			DiscoverableByDefault:    false,
			ShownInDefaultHelp:       false,
			ShownInDefaultDocs:       false,
			AllowsDefaultInvocation:  true,
			RequiresExplicitExposure: true,
		}
	case VisibilityInternal:
		return VisibilityExposurePolicy{
			UserFacing:               false,
			DiscoverableByDefault:    false,
			ShownInDefaultHelp:       false,
			ShownInDefaultDocs:       false,
			AllowsDefaultInvocation:  false,
			RequiresExplicitExposure: true,
		}
	default:
		return VisibilityExposurePolicy{}
	}
}

// IsUserFacing reports whether visibility belongs to the user-facing command
// contract.
func (visibility Visibility) IsUserFacing() bool {
	return visibility.ExposurePolicy().UserFacing
}

// IsDiscoverableByDefault reports whether default command-discovery surfaces
// should include this visibility.
func (visibility Visibility) IsDiscoverableByDefault() bool {
	return visibility.ExposurePolicy().DiscoverableByDefault
}

// IsShownInDefaultHelp reports whether default help output should include this
// visibility.
func (visibility Visibility) IsShownInDefaultHelp() bool {
	return visibility.ExposurePolicy().ShownInDefaultHelp
}

// IsShownInDefaultDocs reports whether default generated documentation should
// include this visibility.
func (visibility Visibility) IsShownInDefaultDocs() bool {
	return visibility.ExposurePolicy().ShownInDefaultDocs
}

// AllowsDefaultInvocation reports whether normal adapters should allow direct
// invocation without an explicit internal/debug mode.
func (visibility Visibility) AllowsDefaultInvocation() bool {
	return visibility.ExposurePolicy().AllowsDefaultInvocation
}

// RequiresExplicitExposure reports whether the visibility should require an
// explicit adapter or generator mode before being exposed.
func (visibility Visibility) RequiresExplicitExposure() bool {
	return visibility.ExposurePolicy().RequiresExplicitExposure
}
