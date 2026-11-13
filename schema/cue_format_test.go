// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"encoding/json"
	goast "go/ast"
	goparser "go/parser"
	gotoken "go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keto "github.com/ory/keto/internal/namespace/ast"
)

// oplFixtures collects every OPL schema in the repository: the parser test
// cases, the .ts files, the CLI snapshots and the string literals in Go test
// files.
func oplFixtures(t *testing.T) map[string]string {
	t.Helper()
	fixtures := make(map[string]string)
	for _, tc := range parserTestCases {
		fixtures["parser_test.go/"+tc.name] = tc.input
	}

	root := ".."
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", "namespace-type-lib", "gen", "docs":
				return filepath.SkipDir
			}
			return nil
		}
		switch {
		case strings.HasSuffix(path, ".ts"):
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(b), "implements Namespace") {
				fixtures[path] = string(b)
			}
		case strings.HasSuffix(path, ".json") && strings.Contains(path, filepath.Join("cmd", "namespace", ".snapshots")):
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var s string
			if json.Unmarshal(b, &s) == nil && strings.Contains(s, "implements Namespace") {
				fixtures[path] = s
			}
		case strings.HasSuffix(path, "_test.go") && !strings.HasPrefix(path, filepath.Join("..", "schema")):
			f, err := goparser.ParseFile(gotoken.NewFileSet(), path, nil, goparser.SkipObjectResolution)
			if err != nil {
				return err
			}
			i := 0
			goast.Inspect(f, func(n goast.Node) bool {
				lit, ok := n.(*goast.BasicLit)
				if !ok || lit.Kind != gotoken.STRING || !strings.Contains(lit.Value, "implements Namespace") {
					return true
				}
				s, err := strconv.Unquote(lit.Value)
				require.NoError(t, err)
				i++
				fixtures[path+"#"+strconv.Itoa(i)] = s
				return true
			})
		}
		return nil
	})
	require.NoError(t, err)
	return fixtures
}

func normalize(ns []namespace) map[string]map[string]keto.Relation {
	m := make(map[string]map[string]keto.Relation)
	for _, n := range ns {
		if m[n.Name] == nil {
			m[n.Name] = make(map[string]keto.Relation)
		}
		for _, r := range n.Relations {
			m[n.Name][r.Name] = r
		}
	}
	return m
}

func TestFormatCUEEquivalence(t *testing.T) {
	fixtures := oplFixtures(t)
	var roundTrips, invalid int

	for name, input := range fixtures {
		want, errs := Parse(input)
		if len(errs) > 0 {
			// An OPL error fixture has no AST to compare.
			invalid++
			continue
		}
		roundTrips++
		t.Run(name, func(t *testing.T) {
			out, err := FormatCUE(want)
			require.NoError(t, err)
			got, cerrs := ParseCUE(SourceFile{Name: "formatted.cue", Content: out})
			for _, e := range cerrs {
				t.Error(e)
			}
			if !assert.Equal(t, normalize(want), normalize(got)) {
				t.Logf("formatted CUE:\n%s", out)
			}
		})
	}

	t.Logf("OPL fixtures: %d round trips, %d invalid OPL inputs", roundTrips, invalid)
	assert.GreaterOrEqual(t, roundTrips, 60, "the fixture scan found too few OPL schemas")
}

func TestFormatCUE(t *testing.T) {
	ns, errs := Parse(parserTestCases[0].input)
	require.Empty(t, errs)
	out, err := FormatCUE(ns)
	require.NoError(t, err)
	t.Logf("\n%s", out)
	s := string(out)
	assert.True(t, strings.HasPrefix(s, "package keto\n"))
	assert.Contains(t, s, `edit: {op: "direct", relation: "owners"}`)
	assert.Contains(t, s, `viewers: types: ["User", "Group#members"]`)
}
