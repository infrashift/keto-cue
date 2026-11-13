// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// The keto-cue base schema. The compiler unifies every user schema with this
// file before it builds the AST. The definitions are closed, so a misspelled
// field is an error. The top level stays open, so a schema can declare its own
// helper definitions next to `namespaces`.
package keto

// A subject type: a namespace, or a subject set "namespace#relation".
#SubjectType: {
	namespace: string & !=""
	relation?: string & !=""
}

// The short form of a subject type: "Namespace" or "Namespace#relation".
#SubjectTypeRef: =~"^[^#]+(#[^#]+)?$"

#RewriteOp: "or" | "and" | "direct" | "tupleToSubjectSet" | "invert"

// A rewrite is a short form or a node.
#Rewrite: #RewriteRef | #RewriteNode

// The short form of a leaf rewrite:
//   "relation"            direct
//   "tupleSet->computed"  tupleToSubjectSet
#RewriteRef: string & !=""

// A rewrite node. It has exactly one of the fields op, or, and and not. The
// compiler checks which fields each op requires:
//   op: "or", op: "and"          children
//   op: "direct"                 relation
//   op: "tupleToSubjectSet"      tupleSetRelation, computedSubjectSet
//   op: "invert"                 child
//   or, and                      a list of rewrites, with no other field
//   not                          one rewrite, with no other field
#RewriteNode: {
	op?: #RewriteOp

	children?: [#Rewrite, ...#Rewrite]

	relation?: string & !=""

	tupleSetRelation?:   string & !=""
	computedSubjectSet?: string & !=""

	child?: #Rewrite

	or?:  [#Rewrite, ...#Rewrite]
	and?: [#Rewrite, ...#Rewrite]
	not?: #Rewrite
}

// A relation is a list of subject types, or a struct with types and a rewrite.
#Relation: #RelationTypes | #RelationNode

#RelationTypes: [...(#SubjectType | #SubjectTypeRef)]

#RelationNode: {
	types?:   #RelationTypes
	rewrite?: #Rewrite
}

#Namespace: {
	relations?: [string]:   #Relation
	permissions?: [string]: #Rewrite
}

namespaces?: [string]: #Namespace
