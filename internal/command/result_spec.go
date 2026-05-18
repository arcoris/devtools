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

// spec returns a detached construction spec.
func (result Result) spec() ResultSpec {
	var exitCode *int
	if result.hasExit {
		value := result.exitCode
		exitCode = &value
	}

	return ResultSpec{
		Status:     result.status,
		Message:    result.message,
		StartedAt:  result.startedAt,
		FinishedAt: result.finishedAt,
		ExitCode:   exitCode,
		Artifacts:  cloneResultArtifacts(result.artifacts),
		Warnings:   cloneResultWarnings(result.warnings),
		Fields:     cloneResultStringMap(result.fields),
		Metadata:   result.metadata,
		Visibility: result.visibility,
	}
}
