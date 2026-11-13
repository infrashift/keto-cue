// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package check_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/keto/internal/check"
	"github.com/ory/keto/internal/check/trace"
	"github.com/ory/keto/internal/driver"
	"github.com/ory/keto/internal/namespace"
	"github.com/ory/keto/internal/testhelpers"
	"github.com/ory/keto/schema"
)

func newRegistryForSchema(t *testing.T, schemaSrc string, opts ...driver.TestRegistryOption) *driver.RegistryDefault {
	t.Helper()
	parsed, errs := schema.Parse(schema.SourceFile{Name: "schema.cue", Content: []byte(schemaSrc)})
	require.Empty(t, errs)
	nss := make([]*namespace.Namespace, len(parsed))
	for i := range parsed {
		nss[i] = &parsed[i]
	}
	return driver.NewSqliteTestRegistry(t, append([]driver.TestRegistryOption{driver.WithNamespaces(nss)}, opts...)...)
}

// requireCyclicError asserts err is the call-limit backstop, not just any error
// (e.g. a context cancellation), by checking its reason.
func requireCyclicError(t *testing.T, err error) {
	t.Helper()
	var re interface{ Reason() string }
	require.ErrorAs(t, err, &re)
	require.Contains(t, re.Reason(), "cyclic permit reference")
}

func TestCallLimit(t *testing.T) {
	t.Parallel()

	type outcome int
	const (
		cyclicError outcome = iota
		member
		depthLimited
	)

	for _, tc := range []struct {
		name      string
		schemaSrc       string
		opts      []driver.TestRegistryOption
		tuples    []string
		check     string
		restDepth int
		want      outcome
	}{
		{
			name: "invert cycle aborts regardless of configured depth",
			schemaSrc: `
				namespaces: {
					User: {}
					Resource: permissions: {
						edit: {
							op: "invert"
							child: {op: "direct", relation: "view"}
						}
						view: {
							op: "invert"
							child: {op: "direct", relation: "edit"}
						}
					}
				}
			`,
			opts:      []driver.TestRegistryOption{driver.WithConfig("limit.max_read_depth", 65535)},
			check:     "Resource:doc#edit@User:alice",
			restDepth: 65535,
			want:      cyclicError,
		},
		{
			name: "branching cycle aborts decisively",
			schemaSrc: `
				namespaces: {
					User: {}
					R: permissions: {
						a: {
							op: "and"
							children: [
								{op: "direct", relation: "b"},
								{op: "direct", relation: "c"},
							]
						}
						b: {
							op: "and"
							children: [
								{op: "direct", relation: "a"},
								{op: "direct", relation: "a"},
							]
						}
						c: {
							op: "and"
							children: [
								{op: "direct", relation: "a"},
								{op: "direct", relation: "a"},
							]
						}
					}
				}
			`,
			check:     "R:x#a@User:u",
			restDepth: 5,
			want:      cyclicError,
		},
		{
			name: "self-loop is bounded by RestDepth",
			schemaSrc: `
				namespaces: {
					User: {}
					Resource: permissions: edit: {op: "direct", relation: "edit"}
				}
			`,
			check:     "Resource:doc#edit@User:alice",
			restDepth: 5,
			want:      depthLimited,
		},
		{
			name: "legit role chain resolves",
			schemaSrc: `
				namespaces: {
					User: {}
					Doc: {
						relations: {
							l0: types: ["User"]
							l1: types: ["User"]
							l2: types: ["User"]
							l3: types: ["User"]
						}
						permissions: {
							p0: {
								op: "or"
								children: [
									{op: "direct", relation: "p1"},
									{op: "direct", relation: "l0"},
								]
							}
							p1: {
								op: "or"
								children: [
									{op: "direct", relation: "p2"},
									{op: "direct", relation: "l1"},
								]
							}
							p2: {
								op: "or"
								children: [
									{op: "direct", relation: "p3"},
									{op: "direct", relation: "l2"},
								]
							}
							p3: {op: "direct", relation: "l3"}
						}
					}
				}
			`,
			tuples:    []string{"Doc:d#l3@User:alice"},
			check:     "Doc:d#p0@User:alice",
			restDepth: 5,
			want:      member,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reg := newRegistryForSchema(t, tc.schemaSrc, tc.opts...)
			if len(tc.tuples) > 0 {
				testhelpers.MapAndInsertTuplesFromString(t, reg, tc.tuples)
			}
			rt := testhelpers.TupleFromString(t, tc.check)

			assertResult := func(t *testing.T, res check.Result) {
				switch tc.want {
				case cyclicError:
					requireCyclicError(t, res.Err)
					require.Equal(t, check.MembershipUnknown, res.Membership)
				case member:
					require.NoError(t, res.Err)
					require.Equal(t, check.IsMember, res.Membership)
				case depthLimited:
					require.NoError(t, res.Err)
					require.Equal(t, check.MembershipUnknown, res.Membership)
					require.Equal(t, check.LimitationMaxDepthExceeded, res.Limitation)
				}
			}

			// Both executors wire the middleware separately.
			t.Run("check engine", func(t *testing.T) {
				assertResult(t, check.NewEngine(reg).CheckRelationTuple(t.Context(), rt, tc.restDepth))
			})
			t.Run("trace engine", func(t *testing.T) {
				res, tree := trace.NewEngine(reg).CheckRelationTupleWithTrace(t.Context(), rt, tc.restDepth)
				require.NotNil(t, tree)
				assertResult(t, res)
			})
		})
	}
}
