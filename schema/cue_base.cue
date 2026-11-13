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

// A subject-set rewrite. The compiler checks which fields each op requires:
//   or, and            children
//   direct             relation
//   tupleToSubjectSet  tupleSetRelation, computedSubjectSet
//   invert             child
#Rewrite: {
	op: #RewriteOp

	children?: [#Rewrite, ...#Rewrite]

	relation?: string & !=""

	tupleSetRelation?:   string & !=""
	computedSubjectSet?: string & !=""

	child?: #Rewrite
}

#Relation: {
	types?: [...(#SubjectType | #SubjectTypeRef)]
	rewrite?: #Rewrite
}

#Namespace: {
	relations?: [string]:   #Relation
	permissions?: [string]: #Rewrite
}

namespaces?: [string]: #Namespace
