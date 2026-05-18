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
	"context"
	"errors"
	"testing"
)

func TestExecuteActionRunsValidAction(t *testing.T) {
	t.Parallel()

	request := actionTestRequest(t)

	result, err := ExecuteAction(context.Background(), ActionFunc(func(ctx context.Context, got ActionRequest) (ActionResult, error) {
		if ctx == nil {
			t.Fatalf("ctx is nil")
		}

		if got.Node.ID() != request.Node.ID() {
			t.Fatalf("request node ID = %q, want %q", got.Node.ID(), request.Node.ID())
		}

		return ActionResult{
			Message: "done",
		}, nil
	}), request)
	if err != nil {
		t.Fatalf("ExecuteAction() returned unexpected error: %v", err)
	}

	if !result.IsOK() {
		t.Fatalf("result IsOK() = false, want true")
	}

	if got, want := result.Message, "done"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}

func TestExecuteActionUsesBackgroundForNilContext(t *testing.T) {
	t.Parallel()

	called := false

	_, err := ExecuteAction(nil, ActionFunc(func(ctx context.Context, _ ActionRequest) (ActionResult, error) {
		called = true
		if ctx == nil {
			t.Fatalf("ctx is nil")
		}

		return ActionResult{}, nil
	}), actionTestRequest(t))
	if err != nil {
		t.Fatalf("ExecuteAction(nil ctx) returned unexpected error: %v", err)
	}

	if !called {
		t.Fatalf("action was not called")
	}
}

func TestExecuteActionRejectsMissingAction(t *testing.T) {
	t.Parallel()

	_, err := ExecuteAction(context.Background(), nil, actionTestRequest(t))
	if err == nil {
		t.Fatalf("ExecuteAction() returned nil error")
	}

	if !errors.Is(err, ErrMissingAction) {
		t.Fatalf("ExecuteAction() error = %v, want ErrMissingAction", err)
	}
}

func TestActionFuncRejectsNilFunction(t *testing.T) {
	t.Parallel()

	var action ActionFunc

	_, err := action.Run(context.Background(), actionTestRequest(t))
	if err == nil {
		t.Fatalf("ActionFunc.Run() returned nil error")
	}

	if !errors.Is(err, ErrMissingAction) {
		t.Fatalf("ActionFunc.Run() error = %v, want ErrMissingAction", err)
	}
}

func TestExecuteActionRejectsInvalidRequestBeforeRun(t *testing.T) {
	t.Parallel()

	called := false

	_, err := ExecuteAction(context.Background(), ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		called = true
		return ActionResult{}, nil
	}), ActionRequest{})
	if err == nil {
		t.Fatalf("ExecuteAction() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("ExecuteAction() error = %v, want ErrInvalidActionRequest", err)
	}

	if called {
		t.Fatalf("action was called for invalid request")
	}
}

func TestExecuteActionRejectsInvalidResultAfterRun(t *testing.T) {
	t.Parallel()

	_, err := ExecuteAction(context.Background(), ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		return ActionResult{
			Status: ActionStatus("unknown"),
		}, nil
	}), actionTestRequest(t))
	if err == nil {
		t.Fatalf("ExecuteAction() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("ExecuteAction() error = %v, want ErrInvalidActionResult", err)
	}
}

func TestExecuteActionPropagatesActionError(t *testing.T) {
	t.Parallel()

	expected := errors.New("domain failure")

	_, err := ExecuteAction(context.Background(), ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		return ActionResult{}, expected
	}), actionTestRequest(t))

	if !errors.Is(err, expected) {
		t.Fatalf("ExecuteAction() error = %v, want %v", err, expected)
	}
}

func TestExecuteActionRespectsCanceledContextBeforeRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false

	_, err := ExecuteAction(ctx, ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		called = true
		return ActionResult{}, nil
	}), actionTestRequest(t))

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ExecuteAction() error = %v, want context.Canceled", err)
	}

	if called {
		t.Fatalf("action was called for canceled context")
	}
}

func TestExecuteActionRespectsCanceledContextAfterRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	_, err := ExecuteAction(ctx, ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		cancel()
		return ActionResult{Message: "done"}, nil
	}), actionTestRequest(t))

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ExecuteAction() error = %v, want context.Canceled", err)
	}
}

func TestNoopAction(t *testing.T) {
	t.Parallel()

	result, err := ExecuteAction(context.Background(), NoopAction(), actionTestRequest(t))
	if err != nil {
		t.Fatalf("ExecuteAction(NoopAction()) returned unexpected error: %v", err)
	}

	if !result.IsOK() {
		t.Fatalf("NoopAction result IsOK() = false, want true")
	}

	if got, want := result.Message, "no operation"; got != want {
		t.Fatalf("NoopAction message = %q, want %q", got, want)
	}
}
