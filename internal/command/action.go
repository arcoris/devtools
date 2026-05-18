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
)

var (
	// ErrMissingAction reports that a runnable command node has no action.
	ErrMissingAction = errors.New("command action is missing")

	// ErrInvalidActionRequest reports that an action request is malformed.
	ErrInvalidActionRequest = errors.New("command action request is invalid")

	// ErrInvalidActionResult reports that an action result is malformed.
	ErrInvalidActionResult = errors.New("command action result is invalid")
)

// Action executes a command node as a compatibility declaration adapter.
//
// Action is framework-neutral. It MUST NOT know about Cobra, pflag, os.Args,
// process exit codes, terminal rendering, or any specific CLI adapter.
//
// RuntimeHandler, RuntimeRequest, and Result are the canonical command
// lifecycle execution contract. Action remains useful for older declarations
// and very small tests, but it intentionally has no independent lifecycle,
// event, cancellation-status, or artifact model. Use RuntimeHandlerFromAction
// when an Action must participate in the canonical Runtime pipeline.
//
// The compatibility action layer is intentionally split:
//
//   - Node describes command-tree structure;
//   - Action describes lower-level executable behavior;
//   - ActionRequest describes one compatibility invocation;
//   - ActionResult adapts to Result for lifecycle output;
//   - errors describe failures.
//
// Domain packages such as check, bench, profile, trace, or perf should be
// called from an Action implementation, but the command kernel itself should
// not depend on those packages.
type Action interface {
	// Run executes the action for one invocation.
	//
	// Implementations MUST respect ctx cancellation where practical.
	// Implementations SHOULD return structured ActionResult values on success
	// and ordinary errors on failure. The adapter or application layer is
	// responsible for mapping errors to diagnostics, output, and exit codes.
	Run(ctx context.Context, request ActionRequest) (ActionResult, error)
}

// ActionFunc adapts a function to Action.
type ActionFunc func(ctx context.Context, request ActionRequest) (ActionResult, error)

// Run executes f as an Action.
func (f ActionFunc) Run(ctx context.Context, request ActionRequest) (ActionResult, error) {
	if f == nil {
		return ActionResult{}, ErrMissingAction
	}

	return f(ctx, request)
}

// NoopAction returns an action that succeeds without side effects.
//
// It is useful for early command-tree wiring, tests, generated placeholder
// commands, and command families that temporarily need a structural action in
// internal validation scenarios. Public leaf commands should normally use a
// domain-specific action instead.
func NoopAction() Action {
	return ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		return ActionResult{
			Status:  ActionStatusOK,
			Message: "no operation",
		}, nil
	})
}

// ExecuteAction validates request, executes action, normalizes the returned
// compatibility result, and validates the normalized result.
//
// ExecuteAction does not emit lifecycle events, recover panics, bind options,
// or normalize a final Result. Runtime.Execute owns those canonical lifecycle
// semantics. Use RuntimeHandlerFromAction to run an Action through Runtime.
func ExecuteAction(ctx context.Context, action Action, request ActionRequest) (ActionResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if action == nil {
		return ActionResult{}, ErrMissingAction
	}

	if err := request.Validate(); err != nil {
		return ActionResult{}, err
	}

	if err := ctx.Err(); err != nil {
		return ActionResult{}, err
	}

	result, err := action.Run(ctx, request)
	if err != nil {
		return ActionResult{}, err
	}

	if err := ctx.Err(); err != nil {
		return ActionResult{}, err
	}

	result = result.Normalize()
	if err := result.Validate(); err != nil {
		return ActionResult{}, err
	}

	return result, nil
}
