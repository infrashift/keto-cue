// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/ory/x/fetcher"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/urlx"
	"github.com/ory/x/watcherx"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"

	"github.com/ory/keto/internal/namespace"
	"github.com/ory/keto/schema"
)

type (
	// configFiles holds the content of each schema source by its path or URL.
	// The watcher keeps bytes, not readers, because every change compiles
	// all sources again.
	configFiles struct {
		byPath map[string][]byte
		sync.Mutex
	}

	oplConfigWatcher struct {
		logger *logrusx.Logger
		target string
		files  configFiles
		// isDir is true when the target is a local directory. Then only the
		// *.cue files of the directory belong to the schema.
		isDir bool
		// ignored records the non-CUE files of a directory target that the
		// watcher already warned about.
		ignored map[string]bool

		memoryNamespaceManager
	}
)

var (
	_ namespace.Manager = (*oplConfigWatcher)(nil)

	// cache holds the parsed namespaces per remote OPL location. Content at a location is
	// assumed to be immutable, therefore, a change to the OPL should create a new file location,
	// therefore a new cache entry.
	cache, _ = ristretto.NewCache(&ristretto.Config[string, []*namespace.Namespace]{
		MaxCost:     20_000_000, // 20 MB of estimated heap, each item ca. 10 KB => max 2000 items
		NumCounters: 20_000,     // max 2000 items => 20000 counters
		BufferItems: 64,
		Cost:        namespace.EstimatedSize,
	})

	fetchGroup singleflight.Group
)

func newOPLConfigWatcher(ctx context.Context, l *logrusx.Logger, newFetcher func() *fetcher.Fetcher, target string) (*oplConfigWatcher, error) {
	nw := &oplConfigWatcher{
		logger:                 l,
		target:                 target,
		files:                  configFiles{byPath: make(map[string][]byte)},
		ignored:                make(map[string]bool),
		memoryNamespaceManager: *NewMemoryNamespaceManager(),
	}

	targetUrl, err := urlx.Parse(target)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch targetUrl.Scheme {
	case "file", "":
		if info, err := os.Stat(targetUrl.Path); err == nil && info.IsDir() { //nolint:gosec // Path is from a validated configuration URL, not user input.
			nw.isDir = true
		}
		return nw, watchTarget(ctx, target, nw, l)
	case "base64":
		file, err := newFetcher().FetchContext(ctx, target)
		if err != nil {
			return nil, err
		}
		nw.files.Lock()
		defer nw.files.Unlock()
		// The base64 URL is not a useful file name in errors.
		nw.files.byPath[""] = file.Bytes()
		nw.parseFiles()
		return nw, err
	case "http", "https":
		if namespaces, ok := cache.Get(target); ok {
			nw.set(namespaces)
			return nw, nil
		}
		v, err, _ := fetchGroup.Do(target, func() (any, error) {
			// Detach from the leader's context so that its cancellation does not
			// fail the concurrent callers sharing this fetch.
			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
			defer cancel()

			buf, err := newFetcher().FetchContext(ctx, target)
			if err != nil {
				return nil, err
			}
			b := buf.Bytes()
			nw.files.Lock()
			defer nw.files.Unlock()
			nw.files.byPath[targetUrl.String()] = b
			namespaces, ok := nw.parseFiles()
			if !ok {
				// Cache a parse failure as an empty result, so that a broken
				// OPL file does not trigger a fetch per construction.
				namespaces = []*namespace.Namespace{}
			}
			// Cost 0 lets the cache's Cost function compute the entry size.
			cache.SetWithTTL(target, namespaces, 0, 30*time.Minute)
			return namespaces, nil
		})
		if err != nil {
			return nil, err
		}
		if namespaces, ok := v.([]*namespace.Namespace); ok {
			nw.set(namespaces)
		}
		return nw, nil
	default:
		return nil, fmt.Errorf("unexpected url scheme: %q", targetUrl.Scheme)
	}
}

func (nw *oplConfigWatcher) handleChange(e *watcherx.ChangeEvent) {
	// the lock is acquired before parsing to ensure that the getters are
	// waiting for the updated values
	content, err := io.ReadAll(e.Reader())
	if err != nil {
		nw.logger.WithError(err).Errorf("Could not read schema file %s.", e.Source())
		return
	}
	nw.files.Lock()
	defer nw.files.Unlock()
	nw.files.byPath[e.Source()] = content
	nw.parseFiles()
}

func (nw *oplConfigWatcher) handleRemove(e *watcherx.RemoveEvent) {
	nw.files.Lock()
	defer nw.files.Unlock()
	delete(nw.files.byPath, e.Source())
	nw.parseFiles()
}

func (nw *oplConfigWatcher) handleError(e *watcherx.ErrorEvent) {
	nw.logger.
		WithError(e).
		Errorf("Received error while watching OPL config files at target %s.",
			nw.target)
}

// parseFiles compiles all schema sources as one CUE package. It sets and
// returns the namespaces only if there were no errors.
//
// The caller must hold the lock to nw.files.
func (nw *oplConfigWatcher) parseFiles() ([]*namespace.Namespace, bool) {
	paths := make([]string, 0, len(nw.files.byPath))
	for path := range nw.files.byPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	sources := make([]schema.SourceFile, 0, len(paths))
	for _, path := range paths {
		if nw.isDir && filepath.Ext(path) != ".cue" {
			if !nw.ignored[path] {
				nw.ignored[path] = true
				nw.logger.Warnf("Ignoring the file %s in the schema directory %s: only *.cue files belong to the schema.", path, nw.target)
			}
			continue
		}
		sources = append(sources, schema.SourceFile{Name: path, Content: nw.files.byPath[path]})
	}

	nn, errs := schema.Parse(sources...)
	if len(errs) > 0 {
		for _, err := range errs {
			nw.logger.
				WithError(err).
				Errorf("Failed to parse the schema at target %s.", nw.target)
		}
		return nil, false
	}
	namespaces := make([]*namespace.Namespace, 0, len(nn))
	for _, n := range nn {
		namespaces = append(namespaces, &n)
	}
	nw.set(namespaces)
	return namespaces, true
}
