// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

func LibfuzzerFuzzParser(data []byte) int {
	Parse(SourceFile{Name: "fuzz.cue", Content: data})
	return 0
}
