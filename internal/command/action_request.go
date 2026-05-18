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
	"fmt"
	"sort"
)

// ActionRequest describes one command action invocation.
//
// ActionRequest is not a CLI parser result. It is the framework-neutral input
// passed to executable command behavior after adapter parsing and command-tree
// resolution have already happened.
//
// Option decoding is intentionally not modeled here yet. The future option
// layer can add typed option values and binders without changing the Action
// interface shape.
type ActionRequest struct {
	// Node is the command node being executed.
	//
	// It MUST be a valid NodeCommand. Family and root nodes are structural and
	// should not be executed as leaf actions.
	Node Node

	// Arguments contains adapter-neutral positional arguments after parsing.
	//
	// The slice is copied by constructors, accessors, and With* helpers.
	Arguments []string

	// Fields contains optional invocation metadata.
	//
	// Fields are machine-facing key-value pairs. Keys use dotted kebab-case.
	// Values are compact UTF-8 text without disallowed control characters.
	Fields map[string]string
}

// NewActionRequest validates values and returns an ActionRequest.
//
// The returned value owns detached copies of mutable input values.
func NewActionRequest(node Node, arguments []string, fields map[string]string) (ActionRequest, error) {
	request := ActionRequest{
		Node:      node,
		Arguments: cloneStringSlice(arguments),
		Fields:    cloneActionStringMap(fields),
	}

	if err := request.Validate(); err != nil {
		return ActionRequest{}, err
	}

	return request, nil
}

// MustActionRequest validates input and returns an ActionRequest.
//
// MustActionRequest panics on invalid input. It is intended for static tests or
// controlled command wiring.
func MustActionRequest(node Node, arguments []string, fields map[string]string) ActionRequest {
	request, err := NewActionRequest(node, arguments, fields)
	if err != nil {
		panic(err)
	}

	return request
}

// Validate verifies action request structural rules.
func (request ActionRequest) Validate() error {
	if request.Node.Kind() == "" {
		return fmt.Errorf("%w: node is not set", ErrInvalidActionRequest)
	}

	if err := request.Node.Validate(); err != nil {
		return fmt.Errorf("%w: invalid node: %w", ErrInvalidActionRequest, err)
	}

	if !request.Node.IsCommand() {
		return fmt.Errorf(
			"%w: node %q must be a command node, got kind %q",
			ErrInvalidActionRequest,
			request.Node.Path(),
			request.Node.Kind(),
		)
	}

	for index, argument := range request.Arguments {
		if err := validateActionRequestText(fmt.Sprintf("argument %d", index), argument, maxActionArgumentLength); err != nil {
			return err
		}
	}

	return validateActionRequestFields(request.Fields)
}

// Args returns a detached copy of request positional arguments.
func (request ActionRequest) Args() []string {
	return cloneStringSlice(request.Arguments)
}

// HasArguments reports whether request contains positional arguments.
func (request ActionRequest) HasArguments() bool {
	return len(request.Arguments) > 0
}

// ArgCount returns the number of positional arguments.
func (request ActionRequest) ArgCount() int {
	return len(request.Arguments)
}

// Argument returns the positional argument at index.
//
// The second return value is false when index is out of range. Argument never
// panics.
func (request ActionRequest) Argument(index int) (string, bool) {
	if index < 0 || index >= len(request.Arguments) {
		return "", false
	}

	return request.Arguments[index], true
}

// Field returns an invocation metadata field and whether it exists.
func (request ActionRequest) Field(key string) (string, bool) {
	value, ok := request.Fields[key]

	return value, ok
}

// HasField reports whether an invocation metadata field exists.
func (request ActionRequest) HasField(key string) bool {
	_, ok := request.Field(key)

	return ok
}

// HasFields reports whether request contains invocation metadata fields.
func (request ActionRequest) HasFields() bool {
	return len(request.Fields) > 0
}

// FieldCount returns the number of invocation metadata fields.
func (request ActionRequest) FieldCount() int {
	return len(request.Fields)
}

// FieldMap returns a detached copy of invocation metadata fields.
func (request ActionRequest) FieldMap() map[string]string {
	return cloneActionStringMap(request.Fields)
}

// FieldKeys returns invocation metadata keys in deterministic lexical order.
func (request ActionRequest) FieldKeys() []string {
	keys := make([]string, 0, len(request.Fields))
	for key := range request.Fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// WithNode returns a validated copy with Node replaced.
func (request ActionRequest) WithNode(node Node) (ActionRequest, error) {
	return NewActionRequest(node, request.Arguments, request.Fields)
}

// MustWithNode returns a validated copy with Node replaced and panics on
// invalid input.
func (request ActionRequest) MustWithNode(node Node) ActionRequest {
	next, err := request.WithNode(node)
	if err != nil {
		panic(err)
	}

	return next
}

// WithArguments returns a validated copy with arguments replaced.
func (request ActionRequest) WithArguments(arguments ...string) (ActionRequest, error) {
	return NewActionRequest(request.Node, arguments, request.Fields)
}

// MustWithArguments returns a validated copy with arguments replaced and panics
// on invalid input.
func (request ActionRequest) MustWithArguments(arguments ...string) ActionRequest {
	next, err := request.WithArguments(arguments...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutArguments returns a validated copy without positional arguments.
func (request ActionRequest) WithoutArguments() ActionRequest {
	return MustActionRequest(request.Node, nil, request.Fields)
}

// WithFields returns a validated copy with invocation metadata fields replaced.
func (request ActionRequest) WithFields(fields map[string]string) (ActionRequest, error) {
	return NewActionRequest(request.Node, request.Arguments, fields)
}

// MustWithFields returns a validated copy with invocation metadata fields
// replaced and panics on invalid input.
func (request ActionRequest) MustWithFields(fields map[string]string) ActionRequest {
	next, err := request.WithFields(fields)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutFields returns a validated copy without invocation metadata fields.
func (request ActionRequest) WithoutFields() ActionRequest {
	return MustActionRequest(request.Node, request.Arguments, nil)
}

// WithField returns a validated copy with one invocation metadata field added
// or replaced.
func (request ActionRequest) WithField(key string, value string) (ActionRequest, error) {
	fields := cloneActionStringMap(request.Fields)
	if fields == nil {
		fields = make(map[string]string)
	}

	fields[key] = value

	return NewActionRequest(request.Node, request.Arguments, fields)
}

// MustWithField returns a validated copy with one field added or replaced and
// panics on invalid input.
func (request ActionRequest) MustWithField(key string, value string) ActionRequest {
	next, err := request.WithField(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutField returns a validated copy without one invocation metadata field.
func (request ActionRequest) WithoutField(key string) ActionRequest {
	fields := cloneActionStringMap(request.Fields)
	delete(fields, key)

	return MustActionRequest(request.Node, request.Arguments, fields)
}
