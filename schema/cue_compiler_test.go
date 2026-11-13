// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ory/x/snapshotx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/keto/internal/namespace/ast"
)

// cueTestCases are CUE ports of the former OPL parser test cases. Their
// snapshots match the ASTs that the OPL parser produced.
var cueTestCases = []struct{ name, input string }{
	{"full example", `
package keto

namespaces: {
	User: relations: manager: types: ["User"]
	Group: relations: members: types: ["User", "Group"]
	Folder: {
		relations: {
			parents: types: ["File"]
			viewers: types: ["Group#members"]
		}
		permissions: view: {op: "direct", relation: "viewers"}
	}
	File: {
		relations: {
			parents: types: ["File", "Folder"]
			viewers: types: ["User", "Group#members"]
			owners: types: ["User", {namespace: "Group", relation: "members"}]
			siblings: types: ["File"]
		}
		// Some comment
		permissions: {
			view: {op: "or", children: [
				{op: "and", children: [
					{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "viewers"},
					{op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"},
				]},
				{op: "or", children: [
					{op: "direct", relation: "viewers"},
					{op: "direct", relation: "viewers"},
					{op: "direct", relation: "viewers"},
				]},
				{op: "direct", relation: "owners"},
			]}
			edit: {op: "direct", relation: "owners"}
			not: {op: "invert", child: {op: "direct", relation: "owners"}}
			rename: {op: "tupleToSubjectSet", tupleSetRelation: "siblings", computedSubjectSet: "edit"}
		}
	}
}
`},
	{"advanced typescript syntax", `
package keto

// A top-level helper definition. The base schema keeps the top level open.
#viaMember: {
	_rel: string
	op:   "tupleToSubjectSet"
	tupleSetRelation:   _rel
	computedSubjectSet: "member"
}

#anyOf: {
	_rels: [...string]
	op: "or"
	children: [for r in _rels {#viaMember & {_rel: r}}]
}

namespaces: {
	Role: relations: member: types: ["Role"]
	Resource: {
		relations: {
			admins: types: ["Role#member"]
			supervisors: types: ["Role#member"]
			annotators: types: ["Role#member"]
			medicalAnnotators: types: ["Role#member"]
		}
		permissions: {
			read: #anyOf & {_rels: ["admins", "annotators", "medicalAnnotators", "supervisors"]}
			comment: {op: "direct", relation: "read"}
			update: #anyOf & {_rels: ["admins", "annotators", "medicalAnnotators", "supervisors"]}
			create: #anyOf & {_rels: ["admins", "annotators", "supervisors"]}
			approve: #anyOf & {_rels: ["admins", "supervisors"]}
			delete: #anyOf & {_rels: ["admins", "supervisors"]}
		}
	}
}
`},
	{"quoted property names", `
package keto

namespaces: Resource: {
	relations: "scope.relation": types: ["Resource"]
	permissions: {
		"scope.action_0": {op: "tupleToSubjectSet", tupleSetRelation: "scope.relation", computedSubjectSet: "scope.action_1"}
		"scope.action_1": {op: "tupleToSubjectSet", tupleSetRelation: "scope.relation", computedSubjectSet: "scope.relation"}
		"scope.action_2": {op: "direct", relation: "scope.action_0"}
	}
}
`},
	{"inverted relationship not case", `
package keto

namespaces: {
	User: {}
	File: {
		relations: {
			viewers: types: ["User"]
			blocklist1: types: ["User"]
			blocklist2: types: ["User"]
		}
		permissions: view: {op: "and", children: [
			{op: "invert", child: {op: "direct", relation: "viewers"}},
			{op: "invert", child: {op: "and", children: [
				{op: "direct", relation: "blocklist1"},
				{op: "direct", relation: "blocklist2"},
			]}},
		]}
	}
}
`},
}

func src(content string) SourceFile {
	return SourceFile{Name: "schema.cue", Content: []byte(content)}
}

func nsMap(ns []namespace) map[string][]ast.Relation {
	m := make(map[string][]ast.Relation)
	for _, n := range ns {
		m[n.Name] = n.Relations
	}
	return m
}

func TestParseCUE(t *testing.T) {
	t.Run("suite=snapshots", func(t *testing.T) {
		for _, tc := range cueTestCases {
			t.Run(tc.name, func(t *testing.T) {
				ns, errs := Parse(src(tc.input))
				for _, err := range errs {
					t.Error(err)
				}
				snapshotx.SnapshotT(t, nsMap(ns))
			})
		}
	})

	t.Run("case=keeps declaration order of namespaces and relations", func(t *testing.T) {
		ns, errs := Parse(src(`
namespaces: {
	Z: relations: {b: types: ["A"], a: types: ["A"]}
	A: permissions: {y: {op: "direct", relation: "x"}, x: {op: "direct", relation: "y"}}
}`))
		require.Empty(t, errs)
		require.Len(t, ns, 2)
		assert.Equal(t, "Z", ns[0].Name)
		assert.Equal(t, "b", ns[0].Relations[0].Name)
		assert.Equal(t, "y", ns[1].Relations[0].Name)
	})

	t.Run("case=empty schema", func(t *testing.T) {
		ns, errs := Parse(src("package keto\n"))
		assert.Empty(t, errs)
		assert.Empty(t, ns)
		ns, errs = Parse()
		assert.Empty(t, errs)
		assert.Empty(t, ns)
	})

	t.Run("case=relation with types and rewrite", func(t *testing.T) {
		ns, errs := Parse(src(`
namespaces: {
	User: {}
	Doc: relations: {
		owners: types: ["User"]
		editors: {types: ["User"], rewrite: {op: "direct", relation: "owners"}}
	}
}`))
		require.Empty(t, errs)
		editors := ns[1].Relations[1]
		assert.Equal(t, []ast.RelationType{{Namespace: "User"}}, editors.Types)
		assert.Equal(t, &ast.SubjectSetRewrite{Children: ast.Children{&ast.ComputedSubjectSet{Relation: "owners"}}}, editors.SubjectSetRewrite)
	})

	t.Run("case=standard library import", func(t *testing.T) {
		ns, errs := Parse(src(`
package keto

import "strings"

namespaces: (strings.ToUpper("user")): {}
`))
		require.Empty(t, errs)
		require.Len(t, ns, 1)
		assert.Equal(t, "USER", ns[0].Name)
	})

	t.Run("case=multi-file package", func(t *testing.T) {
		ns, errs := Parse(
			SourceFile{Name: "a.cue", Content: []byte(`
package acl

#owned: {op: "direct", relation: "owners"}

namespaces: User: {}
`)},
			SourceFile{Name: "b.cue", Content: []byte(`
package acl

namespaces: Doc: {
	relations: owners: types: ["User"]
	permissions: edit: #owned
}
`)},
		)
		require.Empty(t, errs)
		m := nsMap(ns)
		require.Contains(t, m, "User")
		require.Contains(t, m, "Doc")
		assert.Equal(t, "edit", m["Doc"][1].Name)
	})
}

func TestParseCUEErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		files   []SourceFile
		message string
		// line is the expected line of the first error; 0 skips the check.
		line int
	}{
		{
			name:    "syntax error",
			files:   []SourceFile{src("namespaces: {\n  User: {\n}\n")},
			message: "expected",
		},
		{
			name: "bad op",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: {op: "xor"}
}`)},
			message: "op",
			line:    3,
		},
		{
			name: "missing op field",
			files: []SourceFile{src(`namespaces: A: {
	permissions: p: {op: "direct"}
}`)},
			message: `a rewrite with op "direct" requires the field "relation"`,
			line:    2,
		},
		{
			name: "field of another op",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: {op: "direct", relation: "r",
		child: {op: "direct", relation: "r"}}
}`)},
			message: `field "child" is not allowed in a rewrite with op "direct"`,
			line:    4,
		},
		{
			name: "misspelled field",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: {op: "direct", relaton: "r"}
}`)},
			message: "relaton",
			line:    3,
		},
		{
			name: "non-concrete value",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: [string]
}`)},
			message: "incomplete",
		},
		{
			name: "dangling relation",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: {op: "direct", relation: "nope"}
}`)},
			message: `namespace "A" did not declare relation "nope"`,
			line:    3,
		},
		{
			name: "unknown namespace",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["Nope"]
}`)},
			message: `namespace "Nope" was not declared`,
			line:    2,
		},
		{
			name: "bad subject set relation",
			files: []SourceFile{src(`namespaces: {
	Group: relations: members: types: ["Group"]
	A: relations: r: types: [{namespace: "Group", relation: "nope"}]
}`)},
			message: `namespace "Group" did not declare relation "nope"`,
			line:    3,
		},
		{
			name: "bad type short form",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A#b#c"]
}`)},
			message: "types",
		},
		{
			name: "traversal to missing relation",
			files: []SourceFile{src(`namespaces: {
	Folder: relations: owners: types: ["Folder"]
	File: {
		relations: parents: types: ["Folder"]
		permissions: view: {op: "tupleToSubjectSet", tupleSetRelation: "parents", computedSubjectSet: "view"}
	}
}`)},
			message: `relation "view" was not declared in namespace "Folder"`,
			line:    5,
		},
		{
			name: "duplicate name",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: r: {op: "direct", relation: "r"}
}`)},
			message: `namespace "A" declares "r" as a relation and as a permission`,
			line:    3,
		},
		{
			name:    "empty children",
			files:   []SourceFile{src(`namespaces: A: permissions: p: {op: "or", children: []}`)},
			message: "children",
		},
		{
			name: "nested too deeply",
			files: []SourceFile{src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: ` + strings.Repeat(`{op: "invert", child: `, 12) + `{op: "direct", relation: "r"}` + strings.Repeat("}", 12) + `
}`)},
			message: "expression nested too deeply; maximal nesting depth is 10",
			line:    3,
		},
		{
			name: "mixed package clauses",
			files: []SourceFile{
				{Name: "a.cue", Content: []byte("package a\n")},
				{Name: "b.cue", Content: []byte("package b\n")},
			},
			message: "conflicts",
		},
		{
			name:    "foreign import",
			files:   []SourceFile{src("package keto\n\nimport \"example.com/foo\"\n\nnamespaces: A: foo.x\n")},
			message: "example.com/foo",
		},
		{
			name: "duplicate file names",
			files: []SourceFile{
				{Name: "a.cue", Content: []byte("package a\n")},
				{Name: "a.cue", Content: []byte("package a\n")},
			},
			message: `duplicate source file name "a.cue"`,
		},
	} {
		t.Run("case="+tc.name, func(t *testing.T) {
			ns, errs := Parse(tc.files...)
			assert.Nil(t, ns)
			require.NotEmpty(t, errs)
			var msgs []string
			for _, e := range errs {
				msgs = append(msgs, e.Error())
			}
			all := strings.Join(msgs, "\n")
			assert.Contains(t, all, tc.message)
			if tc.line > 0 {
				assert.Equal(t, tc.line, errs[0].Start().Line, all)
				assert.Equal(t, "schema.cue", errs[0].Filename())
			}
			t.Log(all)
		})
	}
}

func TestParseDisjunctionErrorIsSingle(t *testing.T) {
	_, errs := Parse(src("namespaces: A: {\n\tpermissions: p: {op: \"xor\"}\n}\n"))
	require.Len(t, errs, 1)
	assert.Equal(t, `namespaces.A.permissions.p.op: "xor" is not allowed; use "or", "and", "direct", "tupleToSubjectSet" or "invert"`, errs[0].Message())
	assert.Equal(t, 2, errs[0].Start().Line)
	assert.Equal(t, 23, errs[0].Start().Col)

	_, errs = Parse(src(`namespaces: A: relations: r: types: ["A#b#c"]`))
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Message(), `"A#b#c" is not allowed`)
}

func TestParseErrorFormat(t *testing.T) {
	_, errs := Parse(src(`namespaces: A: {
	relations: r: types: ["A"]
	permissions: p: {op: "direct", relation: "nope"}
}`))
	require.Len(t, errs, 1)
	e := errs[0]
	assert.Equal(t, 3, e.Start().Line)
	assert.Equal(t, 43, e.Start().Col)
	assert.Equal(t, 49, e.End().Col)
	assert.Equal(t, `schema.cue: error from 3:43 to 3:49: namespace "A" did not declare relation "nope"

   2 | 	relations: r: types: ["A"]
   3 | 	permissions: p: {op: "direct", relation: "nope"}
     | 	                                         ^~~~~~
   4 | }

`, e.Error())
	assert.Equal(t, `schema.cue: namespace "A" did not declare relation "nope"`, e.ToAPI().Message)
	assert.EqualValues(t, 43, e.ToProto().Start.Column)
}

func FuzzParseCUE(f *testing.F) {
	for _, tc := range cueTestCases {
		f.Add(tc.input)
	}
	f.Fuzz(func(t *testing.T, input string) {
		ns, errs := Parse(src(input))
		if len(errs) > 0 && ns != nil {
			t.Fatal("Parse returned namespaces and errors")
		}
		for _, e := range errs {
			_ = e.Error()
			_ = e.ToAPI()
			_ = e.ToProto()
		}
	})
}

// TestRepositorySchemasCompile compiles every CUE schema in the repository,
// for example the contrib examples and the benchmark fixture.
func TestRepositorySchemasCompile(t *testing.T) {
	var n int
	for _, dir := range []string{"../contrib", "../internal/check/testfixtures"} {
		require.NoError(t, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(path) != ".cue" {
				return err
			}
			n++
			content, err := os.ReadFile(path)
			require.NoError(t, err)
			ns, errs := Parse(SourceFile{Name: path, Content: content})
			for _, e := range errs {
				t.Error(e)
			}
			assert.NotEmpty(t, ns, path)
			return nil
		}))
	}
	assert.GreaterOrEqual(t, n, 5)
}
