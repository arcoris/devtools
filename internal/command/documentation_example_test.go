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

func ExampleDocumentation() {
	documentation := MustDocumentation(DocumentationSpec{
		Summary: "Run configured benchmarks.",
		Usage:   MustSimpleUsage("bench run [flags]"),
		References: []DocumentationReference{
			MustDocumentationReference(DocumentationReferenceSpec{
				Key:    "benchmark-methodology",
				Kind:   DocumentationReferenceDocument,
				Label:  "Benchmark methodology",
				Target: "docs/benchmark-methodology.md",
			}),
		},
	})

	fmt.Println(documentation.Summary())
	fmt.Println(documentation.HasUsage())
	fmt.Println(documentation.ReferenceKeys()[0])

	// Output:
	// Run configured benchmarks.
	// true
	// benchmark-methodology
}
