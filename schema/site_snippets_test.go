// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	cueSnippet  = regexp.MustCompile("(?s)```cue\n(.*?)```")
	fileComment = regexp.MustCompile(`^// (\S+\.cue)\n`)
)

// TestSiteSnippets compiles each CUE snippet of the documentation site, so a
// page cannot show a schema that the compiler rejects. A snippet that starts
// with a "// name.cue" comment is one file of a package: the consecutive
// snippets of such files on a page compile together.
func TestSiteSnippets(t *testing.T) {
	var pages []string
	require.NoError(t, filepath.WalkDir("../docs/src/content/docs", func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".mdx")) {
			pages = append(pages, path)
		}
		return err
	}))
	require.NotEmpty(t, pages)

	for _, page := range pages {
		content, err := os.ReadFile(page) //nolint:gosec // Test code reads the repository docs.
		require.NoError(t, err)

		var group []SourceFile
		check := func(files ...SourceFile) {
			if len(files) == 0 {
				return
			}
			_, errs := Parse(files...)
			for _, e := range errs {
				t.Errorf("%s: %s", page, e)
			}
		}
		for i, m := range cueSnippet.FindAllStringSubmatch(string(content), -1) {
			if f := fileComment.FindStringSubmatch(m[1]); f != nil {
				group = append(group, SourceFile{Name: f[1], Content: []byte(m[1])})
				continue
			}
			check(group...)
			group = nil
			check(SourceFile{Name: filepath.Base(page) + "#" + strconv.Itoa(i), Content: []byte(m[1])})
		}
		check(group...)
	}
}
