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
	"fmt"
)

// RuntimeHandlerFromAction adapts a compatibility Action to the canonical
// RuntimeHandler execution contract.
//
// The supplied node is the command node exposed to the ActionRequest. Runtime
// remains responsible for binding, lifecycle events, panic recovery, timing,
// cancellation classification, and final Result normalization.
func RuntimeHandlerFromAction(action Action, node Node) RuntimeHandler {
	return RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
		actionRequest, err := NewActionRequestFromRuntimeRequest(node, request)
		if err != nil {
			return Result{}, err
		}

		actionResult, err := ExecuteAction(ctx, action, actionRequest)
		if err != nil {
			return Result{}, err
		}

		result, err := NewResultFromActionResult(actionResult)
		if err != nil {
			return Result{}, err
		}

		return result, nil
	})
}

// NewActionRequestFromRuntimeRequest adapts canonical bound runtime input into
// the older ActionRequest shape.
//
// Positional arguments are flattened in declaration order. Bound options are
// exposed as compact fields using the "option.<name>" key so legacy actions can
// inspect resolved values without depending on parser-specific flags.
func NewActionRequestFromRuntimeRequest(node Node, request RuntimeRequest) (ActionRequest, error) {
	input := request.Input()

	arguments := make([]string, 0)
	for _, argument := range input.Arguments() {
		arguments = append(arguments, argument.Values()...)
	}

	fields := map[string]string{
		"runtime.name": request.RuntimeName(),
	}

	if commandID, ok := request.CommandID(); ok {
		fields["command.id"] = commandID.String()
	}

	for _, value := range input.Options() {
		fields["option."+value.Name().String()] = value.String()
	}

	return NewActionRequest(node, arguments, fields)
}

// NewResultFromActionResult converts a compatibility ActionResult into the
// canonical lifecycle Result value.
func NewResultFromActionResult(actionResult ActionResult) (Result, error) {
	normalized, err := NewActionResult(actionResult)
	if err != nil {
		return Result{}, err
	}

	artifacts := make([]Artifact, len(normalized.Artifacts))
	for index, artifact := range normalized.Artifacts {
		converted, err := NewArtifactFromActionArtifact(index, artifact)
		if err != nil {
			return Result{}, err
		}

		artifacts[index] = converted
	}

	warnings := make([]ResultWarning, len(normalized.Warnings))
	for index, warning := range normalized.Warnings {
		converted, err := NewResultWarningFromActionWarning(warning)
		if err != nil {
			return Result{}, err
		}

		warnings[index] = converted
	}

	fields := normalized.FieldMap()
	if normalized.HasData() {
		if fields == nil {
			fields = make(map[string]string)
		}

		fields["action.data.present"] = "true"
	}

	return NewResult(ResultSpec{
		Status:    ResultStatusFromActionStatus(normalized.Status),
		Message:   normalized.Message,
		Artifacts: artifacts,
		Warnings:  warnings,
		Fields:    fields,
	})
}

// ResultStatusFromActionStatus maps compatibility action status to canonical
// result status.
func ResultStatusFromActionStatus(status ActionStatus) ResultStatus {
	switch status.OrDefault() {
	case ActionStatusSkipped:
		return ResultStatusSkipped
	case ActionStatusFailed:
		return ResultStatusFailed
	default:
		return ResultStatusOK
	}
}

// NewArtifactFromActionArtifact converts a compatibility artifact reference
// into the canonical Artifact value.
func NewArtifactFromActionArtifact(index int, actionArtifact ActionArtifact) (Artifact, error) {
	if err := actionArtifact.Validate(); err != nil {
		return Artifact{}, err
	}

	kind, err := NewArtifactKind(actionArtifact.Kind)
	if err != nil {
		return Artifact{}, fmt.Errorf("%w: action artifact %d kind: %w", ErrInvalidActionResult, index, err)
	}

	artifact, err := NewArtifact(ArtifactSpec{
		ID:          fmt.Sprintf("action.artifact-%d", index+1),
		Kind:        kind,
		Location:    actionArtifact.Path,
		Description: actionArtifact.Description,
	})
	if err != nil {
		return Artifact{}, fmt.Errorf("%w: action artifact %d: %w", ErrInvalidActionResult, index, err)
	}

	return artifact, nil
}

// NewResultWarningFromActionWarning converts a compatibility warning into the
// canonical ResultWarning value.
func NewResultWarningFromActionWarning(actionWarning ActionWarning) (ResultWarning, error) {
	if err := actionWarning.Validate(); err != nil {
		return ResultWarning{}, err
	}

	warning, err := NewResultWarning(ResultWarningSpec{
		Kind:    actionWarning.Kind,
		Message: actionWarning.Message,
	})
	if err != nil {
		return ResultWarning{}, fmt.Errorf("%w: action warning: %w", ErrInvalidActionResult, err)
	}

	return warning, nil
}
