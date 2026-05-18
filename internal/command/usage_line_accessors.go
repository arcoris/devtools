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

import "strings"

// String returns the canonical usage line string.
func (line UsageLine) String() string {
	return string(line)
}

// IsZero reports whether the usage line has not been set.
func (line UsageLine) IsZero() bool {
	return line == ""
}

// IsValid reports whether the usage line satisfies the usage-line grammar.
func (line UsageLine) IsValid() bool {
	return line.Validate() == nil
}

// Tokens returns detached usage tokens.
//
// Tokens are split by spaces after canonical usage-line normalization.
func (line UsageLine) Tokens() []string {
	if line == "" {
		return nil
	}

	return cloneStringSlice(strings.Fields(string(line)))
}

// Token returns the usage token at index.
//
// The second return value is false when index is out of range. Token never
// panics.
func (line UsageLine) Token(index int) (string, bool) {
	tokens := line.Tokens()
	if index < 0 || index >= len(tokens) {
		return "", false
	}

	return tokens[index], true
}

// TokenCount returns the number of usage tokens.
func (line UsageLine) TokenCount() int {
	if line == "" {
		return 0
	}

	return len(strings.Fields(string(line)))
}

// FirstToken returns the first usage token and whether it exists.
func (line UsageLine) FirstToken() (string, bool) {
	return line.Token(0)
}

// LastToken returns the last usage token and whether it exists.
func (line UsageLine) LastToken() (string, bool) {
	tokens := line.Tokens()
	if len(tokens) == 0 {
		return "", false
	}

	return tokens[len(tokens)-1], true
}

// HasToken reports whether token appears as a complete usage token.
func (line UsageLine) HasToken(token string) bool {
	for _, candidate := range line.Tokens() {
		if candidate == token {
			return true
		}
	}

	return false
}

// StartsWithToken reports whether token is the first usage token.
func (line UsageLine) StartsWithToken(token string) bool {
	first, ok := line.FirstToken()

	return ok && first == token
}

// EndsWithToken reports whether token is the last usage token.
func (line UsageLine) EndsWithToken(token string) bool {
	last, ok := line.LastToken()

	return ok && last == token
}
