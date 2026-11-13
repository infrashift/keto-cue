// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"

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
		default:
			var fields []any
			if len(r.Types) > 0 {
				types := make([]cueast.Expr, 0, len(r.Types))
				for _, t := range r.Types {
					s := t.Namespace
					if t.Relation != "" {
						s += "#" + t.Relation
					}
					types = append(types, cueast.NewString(s))
				}
				fields = append(fields, cueast.NewIdent("types"), cueast.NewList(types...))
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

// formatTopRewrite unwraps the implicit "or" that wraps a single leaf, because
// Parse adds it back.
func formatTopRewrite(rw *ast.SubjectSetRewrite) cueast.Expr {
	if rw.Operation == ast.OperatorOr && len(rw.Children) == 1 {
		return formatRewrite(rw.Children[0])
	}
	return formatRewrite(rw)
}

func formatRewrite(child ast.Child) cueast.Expr {
	switch c := child.(type) {
	case *ast.SubjectSetRewrite:
		op := opOr
		if c.Operation == ast.OperatorAnd {
			op = opAnd
		}
		children := make([]cueast.Expr, 0, len(c.Children))
		for _, ch := range c.Children {
			children = append(children, formatRewrite(ch))
		}
		list := cueast.NewList(children...)
		for _, e := range list.Elts {
			cueast.SetRelPos(e, token.Newline)
		}
		list.Rbrack = token.Newline.Pos()
		return cueast.NewStruct(
			cueast.NewIdent("op"), cueast.NewString(op),
			cueast.NewIdent("children"), list,
		)
	case *ast.ComputedSubjectSet:
		return inlineStruct(
			cueast.NewIdent("op"), cueast.NewString(opDirect),
			cueast.NewIdent("relation"), cueast.NewString(c.Relation),
		)
	case *ast.TupleToSubjectSet:
		return inlineStruct(
			cueast.NewIdent("op"), cueast.NewString(opTupleToSubjectSet),
			cueast.NewIdent("tupleSetRelation"), cueast.NewString(c.Relation),
			cueast.NewIdent("computedSubjectSet"), cueast.NewString(c.ComputedSubjectSetRelation),
		)
	case *ast.InvertResult:
		return cueast.NewStruct(
			cueast.NewIdent("op"), cueast.NewString(opInvert),
			cueast.NewIdent("child"), formatRewrite(c.Child),
		)
	default:
		panic(fmt.Sprintf("unknown rewrite node %T", child))
	}
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
