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
	"errors"
	"testing"
)

func TestGroupParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		group     Group
		want      Group
		wantFound bool
	}{
		{name: "zero", group: "", want: "", wantFound: false},
		{name: "root", group: MustGroup("benchmark"), want: "", wantFound: false},
		{name: "child", group: MustGroup("diagnostics.perf"), want: MustGroup("diagnostics"), wantFound: true},
		{name: "deep", group: MustGroup("a.b.c"), want: MustGroup("a.b"), wantFound: true},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, found := test.group.Parent()
			if found != test.wantFound {
				t.Fatalf("Parent() found = %v, want %v", found, test.wantFound)
			}

			if got != test.want {
				t.Fatalf("Parent() group = %q, want %q", got, test.want)
			}

			if test.group.HasParent() != test.wantFound {
				t.Fatalf("HasParent() = %v, want %v", test.group.HasParent(), test.wantFound)
			}
		})
	}
}

func TestGroupHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		group  Group
		prefix Group
		want   bool
	}{
		{name: "same group", group: MustGroup("diagnostics.perf"), prefix: MustGroup("diagnostics.perf"), want: true},
		{name: "parent prefix", group: MustGroup("diagnostics.perf"), prefix: MustGroup("diagnostics"), want: true},
		{name: "grandparent prefix", group: MustGroup("a.b.c"), prefix: MustGroup("a"), want: true},
		{name: "similar text is not hierarchy", group: MustGroup("diagnostic-tools.perf"), prefix: MustGroup("diagnostics"), want: false},
		{name: "sibling is false", group: MustGroup("diagnostics.perf"), prefix: MustGroup("diagnostics.trace"), want: false},
		{name: "empty group", group: "", prefix: MustGroup("diagnostics"), want: false},
		{name: "empty prefix", group: MustGroup("diagnostics.perf"), prefix: "", want: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.group.HasPrefix(test.prefix); got != test.want {
				t.Fatalf("%q.HasPrefix(%q) = %v, want %v", test.group, test.prefix, got, test.want)
			}
		})
	}
}

func TestGroupTrimPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		group  Group
		prefix Group
		want   Group
		ok     bool
	}{
		{name: "parent", group: MustGroup("diagnostics.perf.stat"), prefix: MustGroup("diagnostics"), want: MustGroup("perf.stat"), ok: true},
		{name: "self", group: MustGroup("diagnostics.perf"), prefix: MustGroup("diagnostics.perf"), want: "", ok: true},
		{name: "missing", group: MustGroup("diagnostics.perf"), prefix: MustGroup("benchmark"), want: "", ok: false},
		{name: "empty prefix", group: MustGroup("diagnostics.perf"), prefix: "", want: "", ok: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, ok := test.group.TrimPrefix(test.prefix)
			if ok != test.ok {
				t.Fatalf("TrimPrefix() ok = %v, want %v", ok, test.ok)
			}

			if got != test.want {
				t.Fatalf("TrimPrefix() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGroupAppend(t *testing.T) {
	t.Parallel()

	parent := MustGroup("diagnostics")

	child, err := parent.Append("perf")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustGroup("diagnostics.perf"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

func TestGroupAppendToZeroGroup(t *testing.T) {
	t.Parallel()

	var group Group

	child, err := group.Append("benchmark")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustGroup("benchmark"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

func TestGroupAppendRejectsInvalidSegment(t *testing.T) {
	t.Parallel()

	tests := []string{"", "Perf", "perf.stat", "perf_stat", "1perf", "-perf", "перф"}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			_, err := MustGroup("diagnostics").Append(segment)
			if err == nil {
				t.Fatalf("Append(%q) returned nil error", segment)
			}

			if !errors.Is(err, ErrInvalidGroup) {
				t.Fatalf("Append(%q) error = %v, want ErrInvalidGroup", segment, err)
			}
		})
	}
}

func TestGroupAppendRejectsInvalidParent(t *testing.T) {
	t.Parallel()

	_, err := Group("Diagnostics").Append("perf")
	if err == nil {
		t.Fatalf("Append returned nil error")
	}

	if !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("Append error = %v, want ErrInvalidGroup", err)
	}
}

func TestGroupMustAppendPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustGroup("diagnostics").MustAppend("Perf")
	})
}

func TestGroupJoin(t *testing.T) {
	t.Parallel()

	joined, err := MustGroup("diagnostics").Join(MustGroup("perf.stat"))
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustGroup("diagnostics.perf.stat"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
}

func TestGroupJoinWithZeroSide(t *testing.T) {
	t.Parallel()

	joined, err := Group("").Join(MustGroup("benchmark"))
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustGroup("benchmark"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}

	joined, err = MustGroup("benchmark").Join("")
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustGroup("benchmark"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
}

func TestGroupJoinRejectsInvalidSide(t *testing.T) {
	t.Parallel()

	_, err := MustGroup("diagnostics").Join(Group("Perf"))
	if err == nil {
		t.Fatalf("Join() returned nil error")
	}

	if !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("Join() error = %v, want ErrInvalidGroup", err)
	}
}

func TestGroupMustJoinPanicsForInvalidSide(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustGroup("diagnostics").MustJoin(Group("Perf"))
	})
}
