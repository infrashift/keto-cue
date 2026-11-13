// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/keto/internal/namespace/ast"
)

func normalize(ns []namespace) map[string]map[string]ast.Relation {
	m := make(map[string]map[string]ast.Relation)
	for _, n := range ns {
		if m[n.Name] == nil {
			m[n.Name] = make(map[string]ast.Relation)
		}
		for _, r := range n.Relations {
			m[n.Name][r.Name] = r
		}
	}
	return m
}

// TestFormatCUERoundTrip checks Parse -> FormatCUE -> Parse for the test cases
// and every schema in the repository. Phase 2 of the migration ran the same
// check from OPL for all 62 OPL fixtures.
func TestFormatCUERoundTrip(t *testing.T) {
	inputs := make(map[string]string)
	for _, tc := range cueTestCases {
		inputs[tc.name] = tc.input
	}
	for _, dir := range []string{"../contrib", "../internal/check/testfixtures"} {
		require.NoError(t, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".cue" {
				return err
			}
			b, err := os.ReadFile(path)
			inputs[path] = string(b)
			return err
		}))
	}

	for name, input := range inputs {
		t.Run(name, func(t *testing.T) {
			want, errs := Parse(SourceFile{Name: "input.cue", Content: []byte(input)})
			require.Empty(t, errs)
			out, err := FormatCUE(want)
			require.NoError(t, err)
			got, errs := Parse(SourceFile{Name: "formatted.cue", Content: out})
			for _, e := range errs {
				t.Error(e)
			}
			if !assert.Equal(t, normalize(want), normalize(got)) {
				t.Logf("formatted CUE:\n%s", out)
			}
		})
	}
}

func TestFormatCUE(t *testing.T) {
	ns, errs := Parse(SourceFile{Name: "input.cue", Content: []byte(cueTestCases[0].input)})
	require.Empty(t, errs)
	out, err := FormatCUE(ns)
	require.NoError(t, err)
	t.Logf("\n%s", out)
	s := string(out)
	assert.True(t, strings.HasPrefix(s, "package keto\n"))
	assert.Contains(t, s, `edit: {op: "direct", relation: "owners"}`)
	assert.Contains(t, s, `viewers: types: ["User", "Group#members"]`)
}
