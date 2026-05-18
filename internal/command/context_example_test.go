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

import "fmt"

func ExampleContext_ActionRequest() {
	commandContext := MustContext(ContextSpec{
		Node: contextTestCommandNode(),
		Invocation: MustInvocation(InvocationSpec{
			Arguments: []string{"stable"},
			Fields:    map[string]string{"mode": "smoke"},
		}),
		Fields: map[string]string{"source": "example"},
	})

	request := commandContext.MustActionRequest()
	mode, _ := request.Field("mode")

	fmt.Println(request.Node.ID())
	fmt.Println(request.Args()[0])
	fmt.Println(mode)

	// Output:
	// bench.run
	// stable
	// smoke
}
