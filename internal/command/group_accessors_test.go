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

func TestGroupBasicAccessors(t *testing.T) {
	t.Parallel()

	group := MustGroup("diagnostics.perf")

	if got, want := group.String(), "diagnostics.perf"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if got, want := group.Key(), "diagnostics.perf"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}

	if !group.Equal(MustGroup("diagnostics.perf")) {
		t.Fatalf("Equal() = false, want true")
	}

	if group.Equal(MustGroup("diagnostics.trace")) {
		t.Fatalf("Equal(other) = true, want false")
	}
}

func TestGroupIsZero(t *testing.T) {
	t.Parallel()

	var zero Group
	if !zero.IsZero() {
		t.Fatalf("zero group IsZero() = false, want true")
	}

	if MustGroup("benchmark").IsZero() {
		t.Fatalf("non-zero group IsZero() = true, want false")
	}
}

func TestGroupPartsReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	group := MustGroup("diagnostics.perf")

	parts := group.Parts()
	assertStringSlicesEqual(t, parts, []string{"diagnostics", "perf"})

	parts[0] = "changed"

	again := group.Parts()
	assertStringSlicesEqual(t, again, []string{"diagnostics", "perf"})
}

func TestGroupPartsForZeroGroup(t *testing.T) {
	t.Parallel()

	var group Group

	if parts := group.Parts(); parts != nil {
		t.Fatalf("zero group Parts() = %#v, want nil", parts)
	}
}

func TestGroupDepthAndLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		group Group
		want  int
	}{
		{name: "zero", group: "", want: 0},
		{name: "root", group: MustGroup("benchmark"), want: 1},
		{name: "child", group: MustGroup("diagnostics.perf"), want: 2},
		{name: "deep", group: MustGroup("a.b.c"), want: 3},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.group.Depth(); got != test.want {
				t.Fatalf("Depth() = %d, want %d", got, test.want)
			}

			if got := test.group.Len(); got != test.want {
				t.Fatalf("Len() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestGroupAt(t *testing.T) {
	t.Parallel()

	group := MustGroup("diagnostics.perf.stat")

	tests := []struct {
		index int
		want  string
		ok    bool
	}{
		{index: -1, want: "", ok: false},
		{index: 0, want: "diagnostics", ok: true},
		{index: 1, want: "perf", ok: true},
		{index: 2, want: "stat", ok: true},
		{index: 3, want: "", ok: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.want, func(t *testing.T) {
			t.Parallel()

			got, ok := group.At(test.index)
			if ok != test.ok {
				t.Fatalf("At(%d) ok = %v, want %v", test.index, ok, test.ok)
			}

			if got != test.want {
				t.Fatalf("At(%d) = %q, want %q", test.index, got, test.want)
			}
		})
	}
}

func TestGroupLeaf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		group Group
		want  string
	}{
		{name: "zero", group: "", want: ""},
		{name: "root", group: MustGroup("benchmark"), want: "benchmark"},
		{name: "child", group: MustGroup("diagnostics.perf"), want: "perf"},
		{name: "deep", group: MustGroup("a.b.c"), want: "c"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.group.Leaf(); got != test.want {
				t.Fatalf("Leaf() = %q, want %q", got, test.want)
			}
		})
	}
}
