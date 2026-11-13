// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"
)

func namespaceNames(t *testing.T, w *oplConfigWatcher) []string {
	t.Helper()
	nn, err := w.Namespaces(context.Background())
	require.NoError(t, err)
	names := make([]string, 0, len(nn))
	for _, n := range nn {
		names = append(names, n.Name)
	}
	sort.Strings(names)
	return names
}

func TestSchemaWatcherDirectory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	cfg, err := NewDefault(ctx, nil, logrusx.New("", ""), configx.SkipValidation())
	require.NoError(t, err)

	dir := t.TempDir()
	write := func(name, content string) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0600))
	}
	// b.cue uses a definition from a.cue: both files form one package.
	write("a.cue", "package acl\n\n#owned: {op: \"direct\", relation: \"owners\"}\n\nnamespaces: User: {}\n")
	write("b.cue", "package acl\n\nnamespaces: Doc: {\n\trelations: owners: types: [\"User\"]\n\tpermissions: edit: #owned\n}\n")
	write("README.md", "This file is not part of the schema.\n")

	w, err := newOPLConfigWatcher(ctx, cfg.l, cfg.Fetcher, "file://"+dir)
	require.NoError(t, err)
	assert.True(t, w.isDir)
	assert.Equal(t, []string{"Doc", "User"}, namespaceNames(t, w))

	doc, err := w.GetNamespaceByName(ctx, "Doc")
	require.NoError(t, err)
	require.NotNil(t, doc.FindRelation("edit"))
	assert.NotNil(t, doc.FindRelation("edit").SubjectSetRewrite)

	// A change to one file compiles all files again. Before keto-cue, the
	// watcher kept drained readers, so the other files lost their content.
	write("a.cue", "package acl\n\n#owned: {op: \"direct\", relation: \"owners\"}\n\nnamespaces: {\n\tUser: {}\n\tGroup: {}\n}\n")
	require.Eventually(t, func() bool {
		names := namespaceNames(t, w)
		return len(names) == 3
	}, 5*time.Second, 20*time.Millisecond)
	assert.Equal(t, []string{"Doc", "Group", "User"}, namespaceNames(t, w))

	// A broken change keeps the last good namespaces.
	write("b.cue", "package acl\n\nnamespaces: Doc: permissions: edit: {op: \"direct\", relation: \"nope\"}\n")
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, []string{"Doc", "Group", "User"}, namespaceNames(t, w))
}

func TestSchemaWatcherSingleFileAnyExtension(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	cfg, err := NewDefault(ctx, nil, logrusx.New("", ""), configx.SkipValidation())
	require.NoError(t, err)

	f := filepath.Join(t.TempDir(), "keto_namespaces")
	require.NoError(t, os.WriteFile(f, []byte("namespaces: {User: {}, Doc: {}}\n"), 0600))

	w, err := newOPLConfigWatcher(ctx, cfg.l, cfg.Fetcher, "file://"+f)
	require.NoError(t, err)
	assert.False(t, w.isDir)
	assert.Equal(t, []string{"Doc", "User"}, namespaceNames(t, w))
}

func TestSchemaWatcherBase64(t *testing.T) {
	ctx := context.Background()
	cfg, err := NewDefault(ctx, nil, logrusx.New("", ""), configx.SkipValidation())
	require.NoError(t, err)

	target := "base64://" + base64.StdEncoding.EncodeToString([]byte("namespaces: {User: {}}\n"))
	w, err := newOPLConfigWatcher(ctx, cfg.l, cfg.Fetcher, target)
	require.NoError(t, err)
	assert.Equal(t, []string{"User"}, namespaceNames(t, w))
}
