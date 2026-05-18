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
	"time"
)

func runtimeTestBinding() Binding {
	return MustBinding(BindingSpec{
		Options: []Option{
			MustOption(OptionSpec{
				Name:          "format",
				Kind:          OptionKindEnum,
				AllowedValues: []string{"text", "json"},
				DefaultValues: []string{"text"},
			}),
		},
		Arguments: []Argument{
			MustArgument(ArgumentSpec{
				Name:          "suite",
				Kind:          OptionKindEnum,
				AllowedValues: []string{"smoke", "stable"},
			}),
		},
	})
}

func runtimeTestOKHandler() RuntimeHandler {
	return RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
		return OKResult("ok"), nil
	})
}

func runtimeTestStringOption(name string) Option {
	return MustOption(OptionSpec{
		Name: name,
		Kind: OptionKindString,
	})
}

func runtimeTestTime() time.Time {
	return time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)
}

func boolPointer(value bool) *bool {
	return &value
}
