// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

// typeCheck verifies the references that the compiler recorded. The checks and
// their messages follow the OPL type checks.
func (c *cueCompiler) typeCheck() {
	q := namespaceQuery(c.namespaces)

	for _, ref := range c.typeRefs {
		n, ok := q.find(ref.namespace)
		if !ok {
			c.addErr(ref.at, "namespace %q was not declared", ref.namespace)
			continue
		}
		if ref.relation == "" {
			continue
		}
		if _, ok := relationQuery(n.Relations).find(ref.relation); !ok {
			c.addErr(ref.at, "namespace %q did not declare relation %q",
				ref.namespace, ref.relation)
		}
	}

	for _, ref := range c.relationRefs {
		if _, ok := q.findRelation(ref.namespace, ref.relation); !ok {
			c.addErr(ref.at, "namespace %q did not declare relation %q",
				ref.namespace, ref.relation)
		}
	}

	for _, ref := range c.traverseRefs {
		c.checkAllTypesHaveRelation(ref.at, ref.namespace, ref.relation, ref.computed,
			tupleToSubjectSetTypeCheckMaxDepth)
	}
}

// checkAllTypesHaveRelation checks that every type of namespace#relationType
// declares the relation. A subject-set type is followed recursively.
func (c *cueCompiler) checkAllTypesHaveRelation(at srcRange, namespace, relationType, relation string, depth int) {
	if depth < 0 {
		c.addErr(at, "could not typecheck deeply nested SubjectSet further")
		return
	}
	q := namespaceQuery(c.namespaces)
	r, ok := q.findRelation(namespace, relationType)
	if !ok {
		// The relation reference check reports the missing relation.
		return
	}
	for _, t := range r.Types {
		if t.Relation == "" {
			if _, ok := q.findRelation(t.Namespace, relation); !ok {
				c.addErr(at, "relation %q was not declared in namespace %q",
					relation, t.Namespace)
			}
			continue
		}
		c.checkAllTypesHaveRelation(at, t.Namespace, t.Relation, relation, depth-1)
	}
}
