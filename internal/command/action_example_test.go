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

func ExampleExecuteAction() {
	node := MustCommandNode(
		MustID("check"),
		MustPath("check"),
		"check",
	)
	request := MustActionRequest(node, []string{"./..."}, map[string]string{"mode": "smoke"})

	result, _ := ExecuteAction(context.Background(), ActionFunc(func(context.Context, ActionRequest) (ActionResult, error) {
		return ActionResult{
			Message: "checked",
			Fields:  map[string]string{"packages": "all"},
		}, nil
	}), request)

	packages, _ := result.Field("packages")

	fmt.Println(result.Status)
	fmt.Println(result.Message)
	fmt.Println(packages)

	// Output:
	// ok
	// checked
	// all
}
