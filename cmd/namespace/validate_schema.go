// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"

	"github.com/ory/keto/schema"
)

func NewValidateSchemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <schema.cue | schema-dir> [...]",
		Short: "Validate CUE permission schemas",
		Long: `Validate CUE permission schemas.

Each argument is a CUE file or a directory. The *.cue files of a directory
form one CUE package and are validated together, as the server loads them.
Errors name the file, the line and the column.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			failed := false
			for _, arg := range args {
				files, err := readSchemaSources(arg)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not read %q: %+v\n", arg, err)
					return cmdx.FailSilently(cmd)
				}
				namespaces, errs := schema.Parse(files...)
				if len(errs) > 0 {
					failed = true
					for _, e := range errs {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), e.Error())
					}
					continue
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s: valid, %d namespace(s)\n", arg, len(namespaces))
			}
			if failed {
				return cmdx.FailSilently(cmd)
			}
			return nil
		},
	}
}

// readSchemaSources reads a CUE file, or all *.cue files of a directory.
func readSchemaSources(path string) ([]schema.SourceFile, error) {
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	paths := []string{path}
	if info.IsDir() {
		paths, err = filepath.Glob(filepath.Join(path, "*.cue"))
		if err != nil {
			return nil, err
		}
		if len(paths) == 0 {
			return nil, fmt.Errorf("the directory contains no *.cue files")
		}
		sort.Strings(paths)
	}
	files := make([]schema.SourceFile, 0, len(paths))
	for _, p := range paths {
		content, err := os.ReadFile(p) //nolint:gosec // The path comes from the command line of the operator.
		if err != nil {
			return nil, err
		}
		files = append(files, schema.SourceFile{Name: p, Content: content})
	}
	return files, nil
}
