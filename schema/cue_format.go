// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"
	"strings"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/token"

	"github.com/ory/keto/internal/namespace/ast"
)

// FormatCUE renders namespaces as a formatted CUE schema in package keto.
// Parse compiles the result back into the same namespaces.
//
// A relation with types goes to `relations`. A relation with a rewrite and no
// types goes to `permissions`, like an OPL `permits` entry.
func FormatCUE(namespaces []namespace) ([]byte, error) {
	nsFields := make([]any, 0, len(namespaces)*2)
	for _, n := range namespaces {
		nsFields = append(nsFields, cueast.NewStringLabel(n.Name), formatNamespace(n))
	}

	file := &cueast.File{Decls: []cueast.Decl{
		&cueast.Package{Name: cueast.NewIdent("keto")},
		&cueast.Field{Label: cueast.NewIdent("namespaces"), Value: cueast.NewStruct(nsFields...)},
	}}
	out, err := format.Node(file, format.Simplify())
	if err != nil {
		return nil, fmt.Errorf("could not format the CUE schema: %w", err)
	}
	return out, nil
}

func formatNamespace(n namespace) cueast.Expr {
	var relations, permissions []any
	for _, r := range n.Relations {
		switch {
		case len(r.Types) == 0 && r.SubjectSetRewrite != nil:
			permissions = append(permissions, cueast.NewStringLabel(r.Name), formatTopRewrite(r.SubjectSetRewrite))
		case r.SubjectSetRewrite == nil && len(r.Types) > 0:
			// The short form of a relation is its list of types.
			relations = append(relations, cueast.NewStringLabel(r.Name), formatTypes(r.Types))
		default:
			var fields []any
			if len(r.Types) > 0 {
				fields = append(fields, cueast.NewIdent("types"), formatTypes(r.Types))
			}
			if r.SubjectSetRewrite != nil {
				fields = append(fields, cueast.NewIdent("rewrite"), formatTopRewrite(r.SubjectSetRewrite))
			}
			relations = append(relations, cueast.NewStringLabel(r.Name), cueast.NewStruct(fields...))
		}
	}

	var fields []any
	if len(relations) > 0 {
		fields = append(fields, cueast.NewIdent("relations"), cueast.NewStruct(relations...))
	}
	if len(permissions) > 0 {
		fields = append(fields, cueast.NewIdent("permissions"), cueast.NewStruct(permissions...))
	}
	return cueast.NewStruct(fields...)
}

func formatTypes(types []ast.RelationType) cueast.Expr {
	list := make([]cueast.Expr, 0, len(types))
	for _, t := range types {
		s := t.Namespace
		if t.Relation != "" {
			s += "#" + t.Relation
		}
		list = append(list, cueast.NewString(s))
	}
	return cueast.NewList(list...)
}

// formatTopRewrite unwraps the implicit "or" that wraps a single leaf, because
// Parse adds it back.
func formatTopRewrite(rw *ast.SubjectSetRewrite) cueast.Expr {
	if rw.Operation == ast.OperatorOr && len(rw.Children) == 1 {
		return formatRewrite(rw.Children[0])
	}
	return formatRewrite(rw)
}

// formatRewrite renders a rewrite in the short forms. A relation name that a
// short form cannot hold, for example one that contains "->", keeps the long
// form.
func formatRewrite(child ast.Child) cueast.Expr {
	switch c := child.(type) {
	case *ast.SubjectSetRewrite:
		op := shortOr
		if c.Operation == ast.OperatorAnd {
			op = shortAnd
		}
		children := make([]cueast.Expr, 0, len(c.Children))
		allLeaves := true
		for _, ch := range c.Children {
			e := formatRewrite(ch)
			if _, ok := e.(*cueast.BasicLit); !ok {
				allLeaves = false
			}
			children = append(children, e)
		}
		list := cueast.NewList(children...)
		if allLeaves {
			return inlineStruct(cueast.NewIdent(op), list)
		}
		for _, e := range list.Elts {
			cueast.SetRelPos(e, token.Newline)
		}
		list.Rbrack = token.Newline.Pos()
		return cueast.NewStruct(cueast.NewIdent(op), list)
	case *ast.ComputedSubjectSet:
		if isShortName(c.Relation) {
			return cueast.NewString(c.Relation)
		}
		return inlineStruct(
			cueast.NewIdent("op"), cueast.NewString(opDirect),
			cueast.NewIdent("relation"), cueast.NewString(c.Relation),
		)
	case *ast.TupleToSubjectSet:
		if isShortName(c.Relation) && isShortName(c.ComputedSubjectSetRelation) {
			return cueast.NewString(c.Relation + traverseArrow + c.ComputedSubjectSetRelation)
		}
		return inlineStruct(
			cueast.NewIdent("op"), cueast.NewString(opTupleToSubjectSet),
			cueast.NewIdent("tupleSetRelation"), cueast.NewString(c.Relation),
			cueast.NewIdent("computedSubjectSet"), cueast.NewString(c.ComputedSubjectSetRelation),
		)
	case *ast.InvertResult:
		return inlineStruct(cueast.NewIdent(shortNot), formatRewrite(c.Child))
	default:
		panic(fmt.Sprintf("unknown rewrite node %T", child))
	}
}

// isShortName reports whether a relation name fits in a short rewrite.
func isShortName(name string) bool {
	return name != "" && name == strings.TrimSpace(name) && !strings.Contains(name, traverseArrow)
}

// inlineStruct returns a struct that the formatter prints on one line.
func inlineStruct(fields ...any) *cueast.StructLit {
	s := cueast.NewStruct(fields...)
	for _, e := range s.Elts {
		cueast.SetRelPos(e, token.Blank)
	}
	s.Rbrace = token.Blank.Pos()
	return s
}
