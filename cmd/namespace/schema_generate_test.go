// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package namespace_test

import (
	"bytes"
	"testing"

	"github.com/ory/x/snapshotx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/keto/cmd/namespace"
	"github.com/ory/keto/schema"
)

func TestGenerateCUEConfig(t *testing.T) {
	cases := []struct {
		name       string
		namespaces []string
	}{{
		name:       "empty",
		namespaces: []string{},
	}, {
		name:       "one",
		namespaces: []string{"one"},
	}, {
		name:       "many",
		namespaces: []string{"one", "two", "three"},
	}, {
		name:       "quoted",
		namespaces: []string{"with space", "with-dash", "_hidden", "#def", "if"},
	}}

	for _, tc := range cases {
		t.Run("case="+tc.name, func(t *testing.T) {
			var out bytes.Buffer
			require.NoError(t, namespace.GenerateCUEConfig(tc.namespaces, &out))
			snapshotx.SnapshotT(t, out.String())

			ns, errs := schema.Parse(schema.SourceFile{Name: "generated.cue", Content: out.Bytes()})
			for _, e := range errs {
				t.Error(e)
			}
			names := make([]string, 0, len(ns))
			for _, n := range ns {
				names = append(names, n.Name)
			}
			assert.ElementsMatch(t, tc.namespaces, names)
		})
	}
}
