// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"embed"
	"io"
	"text/template"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/literal"
	"github.com/pkg/errors"
)

//go:embed config_template/*
var configTemplate embed.FS

// GenerateCUEConfig derives a CUE permission schema from the legacy namespace
// names and writes it to out. The schema declares each namespace without
// relations, so it is functionally equivalent to the list of namespaces.
func GenerateCUEConfig(namespaces []string, out io.Writer) error {
	t, err := template.New("config_template").ParseFS(configTemplate, "config_template/*")
	if err != nil {
		return errors.WithStack(err)
	}
	labels := make([]string, len(namespaces))
	for i, n := range namespaces {
		labels[i] = cueLabel(n)
	}
	return errors.WithStack(t.ExecuteTemplate(out,
		"namespaces.cue.tmpl",
		struct{ Namespaces []string }{Namespaces: labels}))
}

// cueLabel returns the CUE field label for a namespace name.
func cueLabel(name string) string {
	if cueast.StringLabelNeedsQuoting(name) {
		return literal.String.Quote(name)
	}
	return name
}
