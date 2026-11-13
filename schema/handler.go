// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"context"
	"errors"
	"io"
	"net/http"
	"runtime"
	"time"

	"connectrpc.com/connect"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/ory/herodot"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"

	opl "github.com/ory/keto/gen/go/ory/keto/opl/v1alpha1"
	"github.com/ory/keto/gen/go/ory/keto/opl/v1alpha1/oplconnect"
	"github.com/ory/keto/internal/x"
	"github.com/ory/keto/ketoapi"
)

type (
	handlerDependencies interface {
		logrusx.Provider
		httpx.WriterProvider
		x.HandlerOptionsProvider
	}
	Handler struct {
		oplconnect.UnimplementedSyntaxServiceHandler
		d handlerDependencies

		// slots bounds the number of concurrent schema evaluations. CUE
		// evaluation cannot be cancelled, so a slot stays taken until the
		// evaluation ends, also after the request timed out.
		slots chan struct{}
		// timeout bounds how long a request waits for a slot and for the
		// evaluation.
		timeout time.Duration
		parse   func(files ...SourceFile) ([]namespace, []*ParseError)
	}
)

const RouteBase = "/opl/syntax/check"

const (
	// MaxSchemaCheckBytes is the largest schema that the syntax check accepts.
	MaxSchemaCheckBytes = 1 << 20
	// schemaCheckTimeout is the default of Handler.timeout.
	schemaCheckTimeout = 10 * time.Second
)

var (
	errSchemaTooLarge = errors.New("schema too large")
	errSchemaBusy     = errors.New("schema check busy")
	errSchemaTimeout  = errors.New("schema check timed out")
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{
		d:       d,
		slots:   make(chan struct{}, runtime.GOMAXPROCS(0)),
		timeout: schemaCheckTimeout,
		parse:   Parse,
	}
}

func (h *Handler) RegisterSyntaxRoutes(r httprouterx.Router) {
	r.POST(RouteBase, h.postCheckOplSyntax)
	r.Handle(oplconnect.NewSyntaxServiceHandler(h, h.d.HandlerOptions()...))
}

func (h *Handler) ProtoFiles() []protoreflect.FileDescriptor {
	return []protoreflect.FileDescriptor{opl.File_ory_keto_opl_v1alpha1_syntax_service_proto}
}

// checkSchema compiles a single-file CUE schema within the resource limits of
// the syntax check.
func (h *Handler) checkSchema(ctx context.Context, content []byte) ([]*ParseError, error) {
	if len(content) > MaxSchemaCheckBytes {
		return nil, errSchemaTooLarge
	}
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	select {
	case h.slots <- struct{}{}:
	case <-ctx.Done():
		return nil, errSchemaBusy
	}

	done := make(chan []*ParseError, 1)
	go func() {
		defer func() { <-h.slots }()
		_, errs := h.parse(SourceFile{Content: content})
		done <- errs
	}()

	select {
	case errs := <-done:
		return errs, nil
	case <-ctx.Done():
		return nil, errSchemaTimeout
	}
}

func checkError(err error) *herodot.DefaultError {
	switch {
	case errors.Is(err, errSchemaTooLarge):
		return &herodot.DefaultError{
			StatusField:   http.StatusText(http.StatusRequestEntityTooLarge),
			ErrorField:    "The schema is larger than the syntax check accepts",
			CodeField:     http.StatusRequestEntityTooLarge,
			GRPCCodeField: codes.ResourceExhausted,
		}
	case errors.Is(err, errSchemaBusy):
		return &herodot.DefaultError{
			StatusField:   http.StatusText(http.StatusTooManyRequests),
			ErrorField:    "Too many schema checks are running; retry later",
			CodeField:     http.StatusTooManyRequests,
			GRPCCodeField: codes.ResourceExhausted,
		}
	default:
		return &herodot.DefaultError{
			StatusField:   http.StatusText(http.StatusRequestTimeout),
			ErrorField:    "The schema check did not finish in time",
			CodeField:     http.StatusRequestTimeout,
			GRPCCodeField: codes.DeadlineExceeded,
		}
	}
}

func connectCode(err error) connect.Code {
	switch {
	case errors.Is(err, errSchemaTooLarge), errors.Is(err, errSchemaBusy):
		return connect.CodeResourceExhausted
	default:
		return connect.CodeDeadlineExceeded
	}
}

func (h *Handler) Check(ctx context.Context, request *connect.Request[opl.CheckRequest]) (*connect.Response[opl.CheckResponse], error) {
	parseErrors, err := h.checkSchema(ctx, request.Msg.GetContent())
	if err != nil {
		return nil, connect.NewError(connectCode(err), errors.New(checkError(err).ErrorField))
	}
	apiErrors := make([]*opl.ParseError, len(parseErrors))
	for i, e := range parseErrors {
		apiErrors[i] = e.ToProto()
	}
	return connect.NewResponse(&opl.CheckResponse{ParseErrors: apiErrors}), nil
}

// Check Schema Syntax Request Parameters
//
// swagger:parameters checkOplSyntax
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type checkOplSyntax struct {
	// in: body
	Body checkOplSyntaxBody
}

// CUE Permission Schema Document
//
// swagger:model checkOplSyntaxBody
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type checkOplSyntaxBody string

// swagger:route POST /opl/syntax/check relationship checkOplSyntax
//
// # Check the syntax of a CUE permission schema
//
// The CUE schema is expected in the body of the request. The endpoint keeps
// its OPL name for API compatibility.
//
//	Consumes:
//	- text/plain
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: checkOplSyntaxResult
//	  400: errorGeneric
//	  default: errorGeneric
//
//	Extensions:
//	  x-ory-ratelimit-bucket: keto-admin-medium
func (h *Handler) postCheckOplSyntax(w http.ResponseWriter, r *http.Request) {
	data, err := readCapped(w, r)
	if errors.Is(err, errSchemaTooLarge) {
		h.d.Writer().WriteError(w, r, pkgerrors.WithStack(checkError(err)))
		return
	} else if err != nil {
		h.d.Writer().WriteError(w, r, pkgerrors.WithStack(herodot.ErrBadRequest().WithError(err.Error())))
		return
	}
	parseErrors, err := h.checkSchema(r.Context(), data)
	if err != nil {
		h.d.Writer().WriteError(w, r, pkgerrors.WithStack(checkError(err)))
		return
	}
	apiErrors := make([]*ketoapi.ParseError, len(parseErrors))
	for i, e := range parseErrors {
		apiErrors[i] = e.ToAPI()
	}
	h.d.Writer().Write(w, r, &ketoapi.CheckOPLSyntaxResponse{Errors: apiErrors})
}

// readCapped reads the request body up to MaxSchemaCheckBytes.
func readCapped(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, MaxSchemaCheckBytes))
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		return nil, errSchemaTooLarge
	}
	return data, err
}
