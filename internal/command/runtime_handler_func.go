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

// RuntimeHandler executes a command after input binding.
//
// RuntimeHandler is deliberately small and framework-neutral. It receives
// context.Context for cancellation and RuntimeRequest for command-kernel data.
// It returns a Result and an error because a handler may produce partial
// artifacts or warnings even when execution fails.
type RuntimeHandler interface {
	Run(ctx context.Context, request RuntimeRequest) (Result, error)
}

// RuntimeHandlerFunc adapts a function to RuntimeHandler.
type RuntimeHandlerFunc func(ctx context.Context, request RuntimeRequest) (Result, error)

// Run executes fn.
func (fn RuntimeHandlerFunc) Run(ctx context.Context, request RuntimeRequest) (Result, error) {
	if fn == nil {
		return Result{}, fmt.Errorf("%w: nil runtime handler function", ErrInvalidRuntime)
	}

	return fn(ctx, request)
}
