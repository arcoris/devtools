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

// Package command defines the framework-neutral command kernel for ARCORIS
// developer tooling.
//
// The package owns command identity, declaration, binding, and runtime
// lifecycle value objects. It deliberately does not own Cobra commands, pflag
// values, os.Args parsing, terminal rendering, filesystem artifact writing,
// config loading, environment loading, OpenTelemetry exporters, benchmark
// domain logic, or process termination. Those concerns belong to adapter,
// application, domain, reporting, storage, or observability packages.
//
// Identity and tree structure are modeled by ID, Path, Node, Tree, and
// Registry. ID is the stable machine-facing identifier. Path is the logical
// command-tree location. Node is a declaration for the root, a command family,
// or a leaf command. Tree validates rooted hierarchy invariants. Registry adds
// deterministic lookup by ID, Path, and adapter-facing command segments.
//
// Command input is modeled by Option, Argument, OptionValue, Binding, and
// BoundInput. Option and Argument declare what a command accepts. OptionValue
// is already resolved runtime data with source provenance. Binding validates
// and canonicalizes those values, applies declaration defaults, and returns
// BoundInput. Binding does not parse CLI syntax, inspect parser state, read
// environment variables, load configuration, or prompt users.
//
// Option resolution is a boundary before Binding. ResolveOptionValues combines
// values supplied by command-line, interactive, runtime, environment, config,
// inherited, and declaration-default sources. It respects OptionPolicy allowed
// sources and OptionSource precedence, then returns canonical []OptionValue
// suitable for Binding.Bind. Concrete loaders and parsers stay outside this
// package.
//
// Execution is modeled by Runtime, RuntimeHandler, RuntimeRequest, Result,
// Event, and Artifact. Runtime coordinates validation, binding, lifecycle event
// emission, handler execution, panic recovery, cancellation classification, and
// Result normalization. RuntimeHandler is the canonical executable contract.
// Result is the final lifecycle output; it does not render output or terminate
// a process. Event is an append-only lifecycle observation, not a logger or
// telemetry exporter. Artifact is a validated reference and metadata value; it
// does not create, delete, upload, or check files.
//
// Action, ActionRequest, and ActionResult remain as a lower-level compatibility
// adapter for older command declarations and small tests. New execution code
// should prefer RuntimeHandler, RuntimeRequest, and Result. Use
// RuntimeHandlerFromAction and NewResultFromActionResult when compatibility
// actions must run through the canonical Runtime lifecycle.
//
// Cobra, pflag, shell syntax, terminal output, process exit mapping, config and
// environment loaders, artifact storage, and observability exporters should
// translate to or from these value objects at the package boundary rather than
// leaking adapter-specific state into the command kernel.
package command
