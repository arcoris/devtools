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

import "testing"

func TestVisibilityExposurePolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		visibility Visibility
		want       VisibilityExposurePolicy
	}{
		{
			name:       "public",
			visibility: VisibilityPublic,
			want: VisibilityExposurePolicy{
				UserFacing:               true,
				DiscoverableByDefault:    true,
				ShownInDefaultHelp:       true,
				ShownInDefaultDocs:       true,
				AllowsDefaultInvocation:  true,
				RequiresExplicitExposure: false,
			},
		},
		{
			name:       "hidden",
			visibility: VisibilityHidden,
			want: VisibilityExposurePolicy{
				UserFacing:               true,
				DiscoverableByDefault:    false,
				ShownInDefaultHelp:       false,
				ShownInDefaultDocs:       false,
				AllowsDefaultInvocation:  true,
				RequiresExplicitExposure: true,
			},
		},
		{
			name:       "internal",
			visibility: VisibilityInternal,
			want: VisibilityExposurePolicy{
				UserFacing:               false,
				DiscoverableByDefault:    false,
				ShownInDefaultHelp:       false,
				ShownInDefaultDocs:       false,
				AllowsDefaultInvocation:  false,
				RequiresExplicitExposure: true,
			},
		},
		{name: "zero", visibility: "", want: VisibilityExposurePolicy{}},
		{name: "unknown", visibility: Visibility("private"), want: VisibilityExposurePolicy{}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.visibility.ExposurePolicy()
			if got != test.want {
				t.Fatalf("ExposurePolicy() = %+v, want %+v", got, test.want)
			}

			if got.UserFacing != test.visibility.IsUserFacing() {
				t.Fatalf("IsUserFacing() = %v, want %v", test.visibility.IsUserFacing(), got.UserFacing)
			}

			if got.DiscoverableByDefault != test.visibility.IsDiscoverableByDefault() {
				t.Fatalf(
					"IsDiscoverableByDefault() = %v, want %v",
					test.visibility.IsDiscoverableByDefault(),
					got.DiscoverableByDefault,
				)
			}

			if got.ShownInDefaultHelp != test.visibility.IsShownInDefaultHelp() {
				t.Fatalf("IsShownInDefaultHelp() = %v, want %v", test.visibility.IsShownInDefaultHelp(), got.ShownInDefaultHelp)
			}

			if got.ShownInDefaultDocs != test.visibility.IsShownInDefaultDocs() {
				t.Fatalf("IsShownInDefaultDocs() = %v, want %v", test.visibility.IsShownInDefaultDocs(), got.ShownInDefaultDocs)
			}

			if got.AllowsDefaultInvocation != test.visibility.AllowsDefaultInvocation() {
				t.Fatalf(
					"AllowsDefaultInvocation() = %v, want %v",
					test.visibility.AllowsDefaultInvocation(),
					got.AllowsDefaultInvocation,
				)
			}

			if got.RequiresExplicitExposure != test.visibility.RequiresExplicitExposure() {
				t.Fatalf(
					"RequiresExplicitExposure() = %v, want %v",
					test.visibility.RequiresExplicitExposure(),
					got.RequiresExplicitExposure,
				)
			}
		})
	}
}
