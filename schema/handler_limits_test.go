// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSchemaLimits(t *testing.T) {
	newHandler := func(slots int, timeout time.Duration, parse func(...SourceFile) ([]namespace, []*ParseError)) *Handler {
		return &Handler{slots: make(chan struct{}, slots), timeout: timeout, parse: parse}
	}

	t.Run("case=too large", func(t *testing.T) {
		h := newHandler(1, time.Second, Parse)
		_, err := h.checkSchema(t.Context(), bytes.Repeat([]byte(" "), MaxSchemaCheckBytes+1))
		assert.ErrorIs(t, err, errSchemaTooLarge)
		assert.Equal(t, http.StatusRequestEntityTooLarge, checkError(err).CodeField)
	})

	t.Run("case=busy", func(t *testing.T) {
		h := newHandler(1, 50*time.Millisecond, Parse)
		h.slots <- struct{}{}
		_, err := h.checkSchema(t.Context(), []byte("namespaces: A: {}"))
		assert.ErrorIs(t, err, errSchemaBusy)
		assert.Equal(t, http.StatusTooManyRequests, checkError(err).CodeField)
	})

	t.Run("case=timeout keeps the slot until the evaluation ends", func(t *testing.T) {
		release := make(chan struct{})
		h := newHandler(1, 50*time.Millisecond, func(...SourceFile) ([]namespace, []*ParseError) {
			<-release
			return nil, nil
		})
		_, err := h.checkSchema(t.Context(), []byte("namespaces: A: {}"))
		assert.ErrorIs(t, err, errSchemaTimeout)
		assert.Equal(t, http.StatusRequestTimeout, checkError(err).CodeField)
		assert.Len(t, h.slots, 1, "the running evaluation keeps its slot")

		close(release)
		require.Eventually(t, func() bool { return len(h.slots) == 0 }, time.Second, 5*time.Millisecond)
	})

	t.Run("case=ok", func(t *testing.T) {
		h := newHandler(1, time.Second, Parse)
		errs, err := h.checkSchema(context.Background(), []byte(`namespaces: A: permissions: p: {op: "direct", relation: "x"}`))
		require.NoError(t, err)
		require.Len(t, errs, 1)
		assert.Equal(t, `namespace "A" did not declare relation "x"`, errs[0].ToAPI().Message)
		assert.Empty(t, h.slots)
	})

	t.Run("case=REST body cap", func(t *testing.T) {
		h := newHandler(1, time.Second, Parse)
		h.d = nil
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, RouteBase, bytes.NewReader(bytes.Repeat([]byte(" "), MaxSchemaCheckBytes+1)))
		data, err := readCapped(w, r)
		assert.Nil(t, data)
		assert.ErrorIs(t, err, errSchemaTooLarge)
	})
}
