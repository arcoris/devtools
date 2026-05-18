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

func TestWalkOrderStringAndValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		order WalkOrder
		valid bool
		text  string
	}{
		{name: "pre order", order: WalkPreOrder, valid: true, text: "pre-order"},
		{name: "post order", order: WalkPostOrder, valid: true, text: "post-order"},
		{name: "unknown", order: WalkOrder("unknown"), valid: false, text: "unknown"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.order.IsValid(); got != test.valid {
				t.Fatalf("IsValid() = %v, want %v", got, test.valid)
			}

			if got := test.order.String(); got != test.text {
				t.Fatalf("String() = %q, want %q", got, test.text)
			}
		})
	}
}

func TestTreeWalkPreOrder(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	var got []string
	err := tree.Walk(func(path Path, _ Node) error {
		got = append(got, path.String())

		return nil
	})
	if err != nil {
		t.Fatalf("Walk() returned unexpected error: %v", err)
	}

	assertStringSlicesEqual(t, got, []string{"", "bench", "bench run", "bench compare"})
}

func TestTreeWalkPostOrder(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	var got []string
	err := tree.WalkOrder(WalkPostOrder, func(path Path, _ Node) error {
		got = append(got, path.String())

		return nil
	})
	if err != nil {
		t.Fatalf("WalkOrder() returned unexpected error: %v", err)
	}

	assertStringSlicesEqual(t, got, []string{"bench run", "bench compare", "bench", ""})
}

func TestTreeWalkRejectsUnknownOrder(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	err := tree.WalkOrder(WalkOrder("unknown"), func(Path, Node) error {
		return nil
	})
	if err == nil {
		t.Fatalf("WalkOrder() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("WalkOrder() error = %v, want ErrInvalidTree", err)
	}
}

func TestTreeWalkPropagatesErrors(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)
	expected := errors.New("stop")

	err := tree.Walk(func(Path, Node) error {
		return expected
	})

	if !errors.Is(err, expected) {
		t.Fatalf("Walk() error = %v, want %v", err, expected)
	}
}

func TestTreeWalkNilFuncIsNoop(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	if err := tree.Walk(nil); err != nil {
		t.Fatalf("Walk(nil) error = %v, want nil", err)
	}

	if err := tree.WalkOrder(WalkOrder("unknown"), nil); err != nil {
		t.Fatalf("WalkOrder(unknown, nil) error = %v, want nil", err)
	}
}

func TestZeroTreeWalkReturnsInvalidTree(t *testing.T) {
	t.Parallel()

	var tree Tree

	err := tree.Walk(func(Path, Node) error {
		return nil
	})
	if err == nil {
		t.Fatalf("zero tree Walk() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("zero tree Walk() error = %v, want ErrInvalidTree", err)
	}
}
