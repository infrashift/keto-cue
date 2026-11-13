// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runValidateSchema(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := NewValidateSchemaCmd()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return out.String(), errOut.String(), err
}

func TestValidateSchemaCmd(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o750))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
		return p
	}

	t.Run("case=valid file", func(t *testing.T) {
		f := write("ok.cue", "namespaces: {User: {}, Doc: relations: owners: types: [\"User\"]}\n")
		stdout, _, err := runValidateSchema(t, f)
		require.NoError(t, err)
		assert.Contains(t, stdout, "valid, 2 namespace(s)")
	})

	t.Run("case=invalid file prints line and column", func(t *testing.T) {
		f := write("bad.cue", "namespaces: A: {\n\tpermissions: p: {op: \"direct\", relation: \"nope\"}\n}\n")
		_, stderr, err := runValidateSchema(t, f)
		require.Error(t, err)
		assert.Contains(t, stderr, f+`: error from 2:43 to 2:49: namespace "A" did not declare relation "nope"`)
		assert.Contains(t, stderr, "^~~~~~")
	})

	t.Run("case=directory is one package", func(t *testing.T) {
		write("pkg/a.cue", "package acl\n\n#owned: {op: \"direct\", relation: \"owners\"}\n\nnamespaces: User: {}\n")
		write("pkg/b.cue", "package acl\n\nnamespaces: Doc: {\n\trelations: owners: types: [\"User\"]\n\tpermissions: edit: #owned\n}\n")
		write("pkg/notes.txt", "not part of the schema")
		stdout, stderr, err := runValidateSchema(t, filepath.Join(dir, "pkg"))
		require.NoError(t, err, stderr)
		assert.Contains(t, stdout, "valid, 2 namespace(s)")
	})

	t.Run("case=empty directory", func(t *testing.T) {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "empty"), 0o750))
		_, stderr, err := runValidateSchema(t, filepath.Join(dir, "empty"))
		require.Error(t, err)
		assert.Contains(t, stderr, "no *.cue files")
	})

	t.Run("case=missing file", func(t *testing.T) {
		_, stderr, err := runValidateSchema(t, filepath.Join(dir, "missing.cue"))
		require.Error(t, err)
		assert.Contains(t, stderr, "Could not read")
	})
}
