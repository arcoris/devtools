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

// ActionRequest returns an ActionRequest for this context.
//
// ActionRequest requires the associated node to be NodeCommand. This method is
// a bridge between context-level invocation metadata and the lower-level action
// execution API.
//
// Result fields are built from invocation fields first and context fields
// second. On key conflicts, context fields win because they represent later
// lifecycle metadata.
func (commandContext Context) ActionRequest() (ActionRequest, error) {
	return NewActionRequest(
		commandContext.node,
		commandContext.invocation.Arguments(),
		commandContext.actionRequestFields(),
	)
}

// MustActionRequest returns an ActionRequest for this context and panics on
// invalid input.
func (commandContext Context) MustActionRequest() ActionRequest {
	request, err := commandContext.ActionRequest()
	if err != nil {
		panic(err)
	}

	return request
}

// actionRequestFields returns invocation fields overlaid with context fields.
func (commandContext Context) actionRequestFields() map[string]string {
	fields := commandContext.invocation.Fields()
	for key, value := range commandContext.fields {
		if fields == nil {
			fields = make(map[string]string, len(commandContext.fields))
		}

		fields[key] = value
	}

	return fields
}
