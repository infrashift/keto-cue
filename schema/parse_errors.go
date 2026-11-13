// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	opl "github.com/ory/keto/gen/go/ory/keto/opl/v1alpha1"
	"github.com/ory/keto/ketoapi"
)

type (
	// ParseError is an error in a schema source. Line and Col of the start and
	// end positions are 1-based. A zero Line means the error has no position.
	ParseError struct {
		msg      string
		filename string
		src      string
		start    ketoapi.SourcePosition
		end      ketoapi.SourcePosition
	}
)

func newParseError(msg, filename, src string, start, end ketoapi.SourcePosition) *ParseError {
	if end.Line < start.Line || (end.Line == start.Line && end.Col < start.Col) {
		end = start
	}
	return &ParseError{msg: msg, filename: filename, src: src, start: start, end: end}
}

// Message returns the error message without the source position.
func (e *ParseError) Message() string { return e.msg }

// Filename returns the name of the source file the error belongs to, if any.
func (e *ParseError) Filename() string { return e.filename }

// Start returns the start position of the error.
func (e *ParseError) Start() ketoapi.SourcePosition { return e.start }

// End returns the end position of the error.
func (e *ParseError) End() ketoapi.SourcePosition { return e.end }

func (e *ParseError) Error() string {
	var s strings.Builder
	if e.filename != "" {
		s.WriteString(e.filename)
		s.WriteString(": ")
	}
	if e.start.Line <= 0 {
		s.WriteString(e.msg)
		s.WriteRune('\n')
		return s.String()
	}

	_, _ = fmt.Fprintf(&s, "error from %d:%d to %d:%d: %s\n\n",
		e.start.Line, e.start.Col,
		e.end.Line, e.end.Col,
		e.msg)

	rows := strings.Split(e.src, "\n")
	if len(rows) < e.start.Line {
		s.WriteString("meta error: could not find source position in input\n")
		return s.String()
	}
	errorLineIdx := e.start.Line - 1
	startLineIdx := max(errorLineIdx-1, 0)

	for line := startLineIdx; line <= errorLineIdx; line++ {
		_, _ = fmt.Fprintf(&s, "%4d | %s\n", line+1, rows[line])
	}

	// Columns are 1-based; the caret marks the first column of the error and
	// the tildes mark the rest of the range on the same line.
	startCol := e.start.Col - 1
	endCol := startCol
	if e.end.Line == e.start.Line {
		endCol = e.end.Col - 1
	} else {
		endCol = len(rows[errorLineIdx])
	}
	s.WriteString("     | ")
	for i, r := range rows[errorLineIdx] {
		if i > startCol && i >= endCol {
			break
		}
		switch {
		case i == startCol:
			s.WriteRune('^')
		case startCol < i && i < endCol:
			s.WriteRune('~')
		case unicode.IsSpace(r):
			s.WriteRune(r)
		default:
			s.WriteRune(' ')
		}
	}
	if startCol >= len(rows[errorLineIdx]) {
		s.WriteRune('^')
	}
	s.WriteRune('\n')

	if errorLineIdx+1 < len(rows) {
		_, _ = fmt.Fprintf(&s, "%4d | %s\n", errorLineIdx+2, rows[errorLineIdx+1])
		s.WriteRune('\n')
	}

	return s.String()
}

func (e *ParseError) ToAPI() *ketoapi.ParseError {
	return &ketoapi.ParseError{
		Message: e.apiMessage(),
		Start:   e.start,
		End:     e.end,
	}
}

func (e *ParseError) ToProto() *opl.ParseError {
	return &opl.ParseError{
		Message: e.apiMessage(),
		Start: &opl.SourcePosition{
			Line:   clampUint32(e.start.Line),
			Column: clampUint32(e.start.Col),
		},
		End: &opl.SourcePosition{
			Line:   clampUint32(e.end.Line),
			Column: clampUint32(e.end.Col),
		},
	}
}

// apiMessage prefixes the message with the file name, because the API types
// have no field for it.
func (e *ParseError) apiMessage() string {
	if e.filename == "" {
		return e.msg
	}
	return e.filename + ": " + e.msg
}

func clampUint32(v int) uint32 {
	switch {
	case v < 0:
		return 0
	case uint64(v) > math.MaxUint32:
		return math.MaxUint32
	default:
		return uint32(v)
	}
}
