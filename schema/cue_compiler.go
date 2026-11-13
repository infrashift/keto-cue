// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	cueerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"

	"github.com/ory/keto/internal/namespace/ast"
	"github.com/ory/keto/ketoapi"
)

// baseSchemaFilename is the file name that errors in the base schema carry.
const baseSchemaFilename = "keto_base.cue"

//go:embed cue_base.cue
var baseSchema []byte

// SourceFile is one CUE source file of a schema.
type SourceFile struct {
	// Name is the file name. It appears in errors and must be unique within
	// one call to Parse.
	Name    string
	Content []byte
}

// Rewrite operations of the CUE schema.
const (
	opOr                = "or"
	opAnd               = "and"
	opDirect            = "direct"
	opTupleToSubjectSet = "tupleToSubjectSet"
	opInvert            = "invert"

	// The short operator fields of a rewrite node.
	shortOr  = "or"
	shortAnd = "and"
	shortNot = "not"

	// traverseArrow separates the two relations of a short tupleToSubjectSet
	// rewrite: "tupleSetRelation->computedSubjectSet".
	traverseArrow = "->"
)

// shortOps are the short operator fields, in the order that errors list them.
var shortOps = []string{shortOr, shortAnd, shortNot}

// rewriteFields lists the optional fields of a #RewriteNode and the ops that
// use each of them. The short operator fields belong to no op.
var rewriteFields = []struct {
	name string
	ops  []string
}{
	{"children", []string{opOr, opAnd}},
	{"relation", []string{opDirect}},
	{"tupleSetRelation", []string{opTupleToSubjectSet}},
	{"computedSubjectSet", []string{opTupleToSubjectSet}},
	{"child", []string{opInvert}},
	{shortOr, nil},
	{shortAnd, nil},
	{shortNot, nil},
}

type (
	// srcRange is a source range of a CUE value in one of the user files.
	srcRange struct {
		file       string
		start, end token.Pos
	}

	// relationRef refers to a relation of the namespace that contains it.
	relationRef struct {
		namespace, relation string
		at                  srcRange
	}

	// typeRef refers to a namespace, or to a relation of a namespace, from a
	// relation type.
	typeRef struct {
		namespace, relation string
		at                  srcRange
	}

	// traverseRef is the reference of a tupleToSubjectSet rewrite: all types
	// of namespace#relation must declare computed.
	traverseRef struct {
		namespace, relation, computed string
		at                            srcRange
	}

	cueCompiler struct {
		files      map[string]SourceFile
		user       cue.Value
		namespaces []namespace
		errors     []*ParseError

		relationRefs []relationRef
		typeRefs     []typeRef
		traverseRefs []traverseRef
	}
)

// Parse compiles the CUE files of one schema into namespaces. All files
// belong to one CUE package. Only standard-library imports are allowed. The
// files are unified with the embedded base schema before compilation.
func Parse(files ...SourceFile) (namespaces []namespace, errs []*ParseError) {
	c := &cueCompiler{files: make(map[string]SourceFile, len(files))}
	defer func() {
		// The CUE evaluator must not take down the server on a bad schema.
		if r := recover(); r != nil {
			namespaces = nil
			errs = []*ParseError{newParseError(
				fmt.Sprintf("internal error while evaluating the schema: %v", r),
				"", "", ketoapi.SourcePosition{}, ketoapi.SourcePosition{})}
		}
	}()
	return c.compile(files)
}

func (c *cueCompiler) compile(files []SourceFile) ([]namespace, []*ParseError) {
	if len(files) == 0 {
		return nil, nil
	}

	inst := build.NewContext().NewInstance("", nil)
	for _, f := range files {
		if f.Name == baseSchemaFilename {
			c.addErr(srcRange{}, "the source file name %q is reserved", f.Name)
			return nil, c.errors
		}
		if _, dup := c.files[f.Name]; dup {
			c.addErr(srcRange{}, "duplicate source file name %q", f.Name)
			return nil, c.errors
		}
		c.files[f.Name] = f
		file, err := parser.ParseFile(f.Name, f.Content, parser.ParseComments)
		if err != nil {
			c.addCUEErr(err)
			continue
		}
		if err := inst.AddSyntax(file); err != nil {
			c.addCUEErr(err)
		}
	}
	if len(c.errors) > 0 {
		return nil, c.errors
	}
	if inst.Err != nil {
		c.addCUEErr(inst.Err)
		return nil, c.errors
	}

	ctx := cuecontext.New()
	base := ctx.CompileBytes(baseSchema, cue.Filename(baseSchemaFilename))
	if err := base.Err(); err != nil {
		panic(fmt.Sprintf("the embedded base schema does not compile: %v", err))
	}
	c.user = ctx.BuildInstance(inst)
	if err := c.user.Err(); err != nil {
		c.addCUEErr(err)
		return nil, c.errors
	}

	val := base.Unify(c.user)
	if err := val.Validate(); err != nil {
		c.addCUEErr(err)
		return nil, c.errors
	}
	nsVal := val.LookupPath(cue.MakePath(cue.Str("namespaces")))
	if !nsVal.Exists() {
		return nil, nil
	}
	if err := nsVal.Validate(cue.Concrete(true)); err != nil {
		c.addCUEErr(err)
		return nil, c.errors
	}

	c.compileNamespaces(nsVal)
	if len(c.errors) > 0 {
		return nil, c.errors
	}
	c.typeCheck()
	if len(c.errors) > 0 {
		return nil, c.errors
	}
	return c.namespaces, nil
}

func (c *cueCompiler) compileNamespaces(nsVal cue.Value) {
	iter, err := nsVal.Fields()
	if err != nil {
		c.addCUEErr(err)
		return
	}
	for iter.Next() {
		name := iter.Selector().Unquoted()
		c.namespaces = append(c.namespaces, c.compileNamespace(name, iter.Value()))
	}
}

func (c *cueCompiler) compileNamespace(name string, v cue.Value) namespace {
	ns := namespace{Name: name}
	declared := make(map[string]bool)

	if rels := v.LookupPath(cue.MakePath(cue.Str("relations"))); rels.Exists() {
		iter, err := rels.Fields()
		if err != nil {
			c.addCUEErr(err)
			return ns
		}
		for iter.Next() {
			relName := iter.Selector().Unquoted()
			declared[relName] = true
			ns.Relations = append(ns.Relations, c.compileRelation(name, relName, iter.Value()))
		}
	}

	if perms := v.LookupPath(cue.MakePath(cue.Str("permissions"))); perms.Exists() {
		iter, err := perms.Fields()
		if err != nil {
			c.addCUEErr(err)
			return ns
		}
		for iter.Next() {
			permName := iter.Selector().Unquoted()
			if declared[permName] {
				c.addErr(c.rangeOf(iter.Value()),
					"namespace %q declares %q as a relation and as a permission",
					name, permName)
				continue
			}
			declared[permName] = true
			ns.Relations = append(ns.Relations, ast.Relation{
				Name:              permName,
				SubjectSetRewrite: c.compileTopRewrite(name, iter.Value()),
			})
		}
	}

	return ns
}

func (c *cueCompiler) compileRelation(nsName, relName string, v cue.Value) ast.Relation {
	rel := ast.Relation{Name: relName}

	// The short form of a relation is its list of types.
	types := v
	if v.IncompleteKind() != cue.ListKind {
		types = v.LookupPath(cue.MakePath(cue.Str("types")))
	}
	if types.Exists() {
		iter, err := types.List()
		if err != nil {
			c.addCUEErr(err)
			return rel
		}
		for iter.Next() {
			if t, ok := c.compileType(iter.Value()); ok {
				rel.Types = append(rel.Types, t)
			}
		}
	}

	if rw := v.LookupPath(cue.MakePath(cue.Str("rewrite"))); rw.Exists() {
		rel.SubjectSetRewrite = c.compileTopRewrite(nsName, rw)
	}
	return rel
}

func (c *cueCompiler) compileType(v cue.Value) (ast.RelationType, bool) {
	var t ast.RelationType
	at := c.rangeOf(v)

	switch v.IncompleteKind() {
	case cue.StringKind:
		s, err := v.String()
		if err != nil {
			c.addCUEErr(err)
			return t, false
		}
		t.Namespace, t.Relation, _ = strings.Cut(s, "#")
	case cue.StructKind:
		var err error
		if t.Namespace, err = v.LookupPath(cue.ParsePath("namespace")).String(); err != nil {
			c.addCUEErr(err)
			return t, false
		}
		if r := v.LookupPath(cue.ParsePath("relation")); r.Exists() {
			if t.Relation, err = r.String(); err != nil {
				c.addCUEErr(err)
				return t, false
			}
		}
	default:
		c.addErr(at, "a relation type must be a string or a struct, got %s", v.IncompleteKind())
		return t, false
	}

	c.typeRefs = append(c.typeRefs, typeRef{namespace: t.Namespace, relation: t.Relation, at: at})
	return t, true
}

// compileTopRewrite compiles the rewrite of a relation or permission. Like the
// OPL parser, it wraps a leaf in a rewrite and flattens nested operations.
func (c *cueCompiler) compileTopRewrite(nsName string, v cue.Value) *ast.SubjectSetRewrite {
	child := c.compileRewrite(nsName, v, 0)
	if child == nil {
		return nil
	}
	return simplifyExpression(child.AsRewrite())
}

// compileRewrite compiles one rewrite node. The nesting is the number of
// operator nodes (or, and, invert) above v.
func (c *cueCompiler) compileRewrite(nsName string, v cue.Value, nesting int) ast.Child {
	at := c.rangeOf(v)
	if nesting > expressionNestingMaxDepth {
		c.addErr(at, "expression nested too deeply; maximal nesting depth is %d",
			expressionNestingMaxDepth)
		return nil
	}

	if v.IncompleteKind() == cue.StringKind {
		return c.compileShortRef(nsName, v)
	}

	var forms []string
	for _, f := range append([]string{"op"}, shortOps...) {
		if v.LookupPath(cue.ParsePath(f)).Exists() {
			forms = append(forms, f)
		}
	}
	switch {
	case len(forms) == 0:
		c.addErr(at, `a rewrite requires one of the fields "op", "or", "and" or "not"`)
		return nil
	case len(forms) > 1:
		c.addErr(at, "a rewrite can have only one of the fields %q and %q", forms[0], forms[1])
		return nil
	case forms[0] != "op":
		return c.compileShortOp(nsName, v, forms[0], nesting)
	}

	op, err := v.LookupPath(cue.ParsePath("op")).String()
	if err != nil {
		c.addCUEErr(err)
		return nil
	}

	ok := true
	for _, f := range rewriteFields {
		if !v.LookupPath(cue.ParsePath(f.name)).Exists() || contains(f.ops, op) {
			continue
		}
		c.addErr(c.rangeOf(v.LookupPath(cue.ParsePath(f.name))),
			"field %q is not allowed in a rewrite with op %q", f.name, op)
		ok = false
	}
	if !ok {
		return nil
	}

	switch op {
	case opOr, opAnd:
		children := v.LookupPath(cue.ParsePath("children"))
		if !children.Exists() {
			c.addErr(at, "a rewrite with op %q requires the field %q", op, "children")
			return nil
		}
		iter, err := children.List()
		if err != nil {
			c.addCUEErr(err)
			return nil
		}
		rw := &ast.SubjectSetRewrite{Operation: ast.OperatorOr}
		if op == opAnd {
			rw.Operation = ast.OperatorAnd
		}
		for iter.Next() {
			child := c.compileRewrite(nsName, iter.Value(), nesting+1)
			if child == nil {
				return nil
			}
			rw.Children = append(rw.Children, child)
		}
		return rw

	case opDirect:
		rel, relAt, ok := c.requiredString(v, op, "relation")
		if !ok {
			return nil
		}
		c.relationRefs = append(c.relationRefs, relationRef{namespace: nsName, relation: rel, at: relAt})
		return &ast.ComputedSubjectSet{Relation: rel}

	case opTupleToSubjectSet:
		rel, relAt, ok := c.requiredString(v, op, "tupleSetRelation")
		if !ok {
			return nil
		}
		computed, _, ok := c.requiredString(v, op, "computedSubjectSet")
		if !ok {
			return nil
		}
		c.relationRefs = append(c.relationRefs, relationRef{namespace: nsName, relation: rel, at: relAt})
		c.traverseRefs = append(c.traverseRefs, traverseRef{namespace: nsName, relation: rel, computed: computed, at: relAt})
		return &ast.TupleToSubjectSet{Relation: rel, ComputedSubjectSetRelation: computed}

	case opInvert:
		childVal := v.LookupPath(cue.ParsePath("child"))
		if !childVal.Exists() {
			c.addErr(at, "a rewrite with op %q requires the field %q", op, "child")
			return nil
		}
		child := c.compileRewrite(nsName, childVal, nesting+1)
		if child == nil {
			return nil
		}
		return &ast.InvertResult{Child: child}

	default:
		// The base schema restricts op, so this is not reached.
		c.addErr(at, "unknown rewrite op %q", op)
		return nil
	}
}

// compileShortRef compiles the short form of a leaf rewrite: "relation" or
// "tupleSetRelation->computedSubjectSet".
func (c *cueCompiler) compileShortRef(nsName string, v cue.Value) ast.Child {
	at := c.rangeOf(v)
	s, err := v.String()
	if err != nil {
		c.addCUEErr(err)
		return nil
	}

	rel, computed, isTraverse := strings.Cut(s, traverseArrow)
	rel, computed = strings.TrimSpace(rel), strings.TrimSpace(computed)
	if rel == "" || (isTraverse && (computed == "" || strings.Contains(computed, traverseArrow))) {
		c.addErr(at, `the rewrite %q is not valid; use "relation" or "tupleSetRelation%scomputedSubjectSet"`, s, traverseArrow)
		return nil
	}

	c.relationRefs = append(c.relationRefs, relationRef{namespace: nsName, relation: rel, at: at})
	if !isTraverse {
		return &ast.ComputedSubjectSet{Relation: rel}
	}
	c.traverseRefs = append(c.traverseRefs, traverseRef{namespace: nsName, relation: rel, computed: computed, at: at})
	return &ast.TupleToSubjectSet{Relation: rel, ComputedSubjectSetRelation: computed}
}

// compileShortOp compiles a rewrite node with a short operator field: {or:
// [...]}, {and: [...]} or {not: rewrite}. The node has no other fields.
func (c *cueCompiler) compileShortOp(nsName string, v cue.Value, op string, nesting int) ast.Child {
	ok := true
	for _, f := range rewriteFields {
		if f.name == op || !v.LookupPath(cue.ParsePath(f.name)).Exists() {
			continue
		}
		c.addErr(c.rangeOf(v.LookupPath(cue.ParsePath(f.name))),
			"field %q is not allowed in a rewrite with the field %q", f.name, op)
		ok = false
	}
	if !ok {
		return nil
	}

	operand := v.LookupPath(cue.ParsePath(op))
	if op == shortNot {
		child := c.compileRewrite(nsName, operand, nesting+1)
		if child == nil {
			return nil
		}
		return &ast.InvertResult{Child: child}
	}

	iter, err := operand.List()
	if err != nil {
		c.addCUEErr(err)
		return nil
	}
	rw := &ast.SubjectSetRewrite{Operation: ast.OperatorOr}
	if op == shortAnd {
		rw.Operation = ast.OperatorAnd
	}
	for iter.Next() {
		child := c.compileRewrite(nsName, iter.Value(), nesting+1)
		if child == nil {
			return nil
		}
		rw.Children = append(rw.Children, child)
	}
	return rw
}

func (c *cueCompiler) requiredString(v cue.Value, op, field string) (string, srcRange, bool) {
	f := v.LookupPath(cue.ParsePath(field))
	if !f.Exists() {
		c.addErr(c.rangeOf(v), "a rewrite with op %q requires the field %q", op, field)
		return "", srcRange{}, false
	}
	s, err := f.String()
	if err != nil {
		c.addCUEErr(err)
		return "", srcRange{}, false
	}
	return s, c.rangeOf(f), true
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// rangeOf returns the source range of v in the user files. It looks up the
// same path in the user value, so that the range points at what the user
// wrote and not at the base schema.
func (c *cueCompiler) rangeOf(v cue.Value) srcRange {
	for _, cand := range []cue.Value{c.user.LookupPath(v.Path()), v} {
		if !cand.Exists() {
			continue
		}
		src := cand.Source()
		if f, ok := src.(*cueast.Field); ok {
			// Point at the value, not at the label.
			src = f.Value
		}
		if src != nil && c.isUserPos(src.Pos()) {
			return srcRange{file: src.Pos().Filename(), start: src.Pos(), end: src.End()}
		}
		if p := cand.Pos(); c.isUserPos(p) {
			return srcRange{file: p.Filename(), start: p, end: p}
		}
	}
	return srcRange{}
}

func (c *cueCompiler) isUserPos(p token.Pos) bool {
	if !p.IsValid() {
		return false
	}
	_, ok := c.files[p.Filename()]
	return ok
}

// addErr adds an error. It drops an error with the same message at the same
// position, for example one that a comprehension makes for each entry.
func (c *cueCompiler) addErr(at srcRange, format string, a ...any) {
	e := c.newError(fmt.Sprintf(format, a...), at)
	for _, prev := range c.errors {
		if prev.Message() == e.Message() && prev.Filename() == e.Filename() &&
			prev.Start() == e.Start() && prev.End() == e.End() {
			return
		}
	}
	c.errors = append(c.errors, e)
}

func (c *cueCompiler) newError(msg string, at srcRange) *ParseError {
	if !at.start.IsValid() {
		return newParseError(msg, "", "", ketoapi.SourcePosition{}, ketoapi.SourcePosition{})
	}
	return newParseError(msg, at.file, string(c.files[at.file].Content), toSourcePos(at.start), toSourcePos(at.end))
}

func toSourcePos(p token.Pos) ketoapi.SourcePosition {
	if !p.IsValid() {
		return ketoapi.SourcePosition{}
	}
	return ketoapi.SourcePosition{Line: p.Line(), Col: p.Column()}
}

// addCUEErr converts a CUE error list into parse errors. Each error takes the
// first position that is in a user file.
func (c *cueCompiler) addCUEErr(err error) {
	errs := cueerrors.Errors(err)

	// CUE reports an empty disjunction once as a summary and once per
	// alternative. The summary becomes one clear message, so drop the
	// per-alternative conflicts at the same path.
	// When an error inside the value explains the conflict, for example a
	// misspelled field of a struct, drop the summary instead.
	disjunctions := make(map[string]bool)
	var paths []string
	for _, e := range errs {
		path := strings.Join(e.Path(), ".")
		paths = append(paths, path)
		if format, _ := e.Msg(); strings.Contains(format, "empty disjunction") {
			disjunctions[path] = true
		}
	}
	hasInner := func(path string) bool {
		for _, p := range paths {
			if strings.HasPrefix(p, path+".") {
				return true
			}
		}
		return false
	}

	seen := make(map[string]bool)
	for _, e := range errs {
		e = cueerrors.Sanitize(e)
		format, args := e.Msg()
		path := strings.Join(e.Path(), ".")
		isSummary := strings.Contains(format, "empty disjunction")
		if disjunctions[path] && (!isSummary || hasInner(path)) {
			continue
		}
		msg := fmt.Sprintf(format, args...)
		if path := e.Path(); len(path) > 0 {
			msg = strings.Join(path, ".") + ": " + msg
		}

		var at srcRange
		for _, p := range cueerrors.Positions(e) {
			if c.isUserPos(p) {
				at = srcRange{file: p.Filename(), start: p, end: p}
				break
			}
		}
		// Errors against the base schema, for example an empty disjunction,
		// can carry base positions only. Find the user value by its path.
		if userVal, ok := c.userValueAt(e.Path()); ok {
			if !at.start.IsValid() {
				at = c.rangeOf(userVal)
			}
			if strings.Contains(format, "empty disjunction") {
				msg = disjunctionMessage(e.Path(), userVal)
			}
		}

		key := fmt.Sprintf("%s|%s|%d|%d", msg, at.file, at.start.Line(), at.start.Column())
		if !at.start.IsValid() {
			key = msg
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		c.errors = append(c.errors, c.newError(msg, at))
	}
}

// userValueAt returns the user value at a path of CUE selector strings.
func (c *cueCompiler) userValueAt(path []string) (cue.Value, bool) {
	if len(path) == 0 || !c.user.Exists() {
		return cue.Value{}, false
	}
	sels := make([]cue.Selector, 0, len(path))
	for _, elem := range path {
		if n, err := strconv.Atoi(elem); err == nil {
			sels = append(sels, cue.Index(n))
			continue
		}
		p := cue.ParsePath(elem)
		if p.Err() != nil {
			return cue.Value{}, false
		}
		sels = append(sels, p.Selectors()...)
	}
	v := c.user.LookupPath(cue.MakePath(sels...))
	return v, v.Exists()
}

// disjunctionMessage explains which values the base schema allows when a
// value matches none of the alternatives of a disjunction.
func disjunctionMessage(path []string, v cue.Value) string {
	var allowed string
	switch {
	case path[len(path)-1] == "op":
		allowed = `"or", "and", "direct", "tupleToSubjectSet" or "invert"`
	case len(path) >= 2 && path[len(path)-2] == "types":
		allowed = `"Namespace", "Namespace#relation" or {namespace: "Namespace", relation?: "relation"}`
	case len(path) >= 2 && path[len(path)-2] == "relations":
		allowed = `a list of types or {types?: [...], rewrite?: ...}`
	case isRewritePath(path):
		allowed = `"relation", "tupleSetRelation->computedSubjectSet", {or: [...]}, {and: [...]}, {not: ...} or {op: ...}`
	default:
		allowed = "a value that the base schema allows"
	}
	return fmt.Sprintf("%s: %v is not allowed; use %s", strings.Join(path, "."), v, allowed)
}

// isRewritePath reports whether a CUE path ends at a rewrite.
func isRewritePath(path []string) bool {
	n := len(path)
	switch {
	case n >= 1 && (path[n-1] == "rewrite" || path[n-1] == "child" || path[n-1] == shortNot):
		return true
	case n >= 2 && path[n-2] == "permissions":
		return true
	case n >= 2 && (path[n-2] == "children" || path[n-2] == shortOr || path[n-2] == shortAnd):
		_, err := strconv.Atoi(path[n-1])
		return err == nil
	}
	return false
}

// simplifyExpression rewrites the expression to use n-ary set operations
// instead of binary ones.
func simplifyExpression(root *ast.SubjectSetRewrite) *ast.SubjectSetRewrite {
	if root == nil {
		return nil
	}
	var newChildren []ast.Child
	for _, child := range root.Children {
		switch ch := child.(type) {
		case *ast.SubjectSetRewrite:
			simplifyExpression(ch)
			if len(ch.Children) == 1 || ch.Operation == root.Operation {
				// Flatten the child to parent, when:
				// - it's a single-element wrapper (AND(OR(A)) ≡ AND(A))
				// - shares the root's operation (OR(OR(A,B),C) ≡ OR(A,B,C)).
				newChildren = append(newChildren, ch.Children...)
				continue
			}
			newChildren = append(newChildren, child)
		case *ast.InvertResult:
			if innerCh, ok := ch.Child.(*ast.SubjectSetRewrite); ok {
				simplifyExpression(innerCh)
			}
			newChildren = append(newChildren, child)
		default:
			// Leaf node or different operation with multiple children; keep as-is.
			newChildren = append(newChildren, child)
		}
	}
	root.Children = newChildren

	return root
}
