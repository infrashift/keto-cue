// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/keto/internal/namespace/ast"
)

func Test_simplify(t *testing.T) {
	testCases := []struct {
		name            string
		input, expected *ast.SubjectSetRewrite
	}{
		{"empty", nil, nil},
		{
			name: "merge all unions",
			input: &ast.SubjectSetRewrite{
				Operation: ast.OperatorOr,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.SubjectSetRewrite{
						Children: ast.Children{
							&ast.ComputedSubjectSet{Relation: "B"},
							&ast.SubjectSetRewrite{
								Children: ast.Children{
									&ast.ComputedSubjectSet{Relation: "C"},
									&ast.SubjectSetRewrite{
										Children: ast.Children{
											&ast.ComputedSubjectSet{Relation: "D"},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &ast.SubjectSetRewrite{
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.ComputedSubjectSet{Relation: "B"},
					&ast.ComputedSubjectSet{Relation: "C"},
					&ast.ComputedSubjectSet{Relation: "D"},
				},
			},
		},
		{
			// ( ( (A || B) ) ) = A || B
			name: "merges OR children",
			input: &ast.SubjectSetRewrite{
				Operation: ast.OperatorOr,
				Children: ast.Children{
					&ast.SubjectSetRewrite{
						Operation: ast.OperatorOr,
						Children: ast.Children{
							&ast.ComputedSubjectSet{Relation: "A"},
						},
					},
					&ast.SubjectSetRewrite{
						Operation: ast.OperatorOr,
						Children: ast.Children{
							&ast.ComputedSubjectSet{Relation: "B"},
						},
					},
				},
			},
			expected: &ast.SubjectSetRewrite{
				Operation: ast.OperatorOr,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.ComputedSubjectSet{Relation: "B"},
				},
			},
		},
		{
			// ( ( A && ((B)) ) ) = A && B
			name: "merges AND children",
			input: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.SubjectSetRewrite{
						Operation: ast.OperatorOr,
						Children: ast.Children{
							&ast.ComputedSubjectSet{Relation: "A"},
						},
					},
					&ast.SubjectSetRewrite{
						Operation: ast.OperatorOr,
						Children: ast.Children{
							&ast.SubjectSetRewrite{
								Operation: ast.OperatorOr,
								Children: ast.Children{
									&ast.ComputedSubjectSet{Relation: "B"},
								},
							},
						},
					},
				},
			},
			expected: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.ComputedSubjectSet{Relation: "B"},
				},
			},
		},
		{
			// A AND !!B = A AND !!B
			name: "preserve double NOT",
			input: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.InvertResult{
						Child: &ast.InvertResult{
							Child: &ast.SubjectSetRewrite{
								Operation: ast.OperatorOr,
								Children: ast.Children{
									&ast.ComputedSubjectSet{Relation: "B"},
								},
							},
						},
					},
				},
			},
			expected: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.InvertResult{
						Child: &ast.InvertResult{
							Child: &ast.SubjectSetRewrite{
								Operation: ast.OperatorOr,
								Children: ast.Children{
									&ast.ComputedSubjectSet{Relation: "B"},
								},
							},
						},
					},
				},
			},
		},
		{
			// A AND NOT(B OR (C)) -> A AND NOT(B OR C)
			name: "simplifies the expression inside the inverted relationship",
			input: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.InvertResult{
						Child: &ast.SubjectSetRewrite{
							Operation: ast.OperatorOr,
							Children: ast.Children{
								&ast.ComputedSubjectSet{Relation: "B"},
								&ast.SubjectSetRewrite{
									Operation: ast.OperatorOr,
									Children: ast.Children{
										&ast.ComputedSubjectSet{Relation: "C"},
									},
								},
							},
						},
					},
				},
			},
			expected: &ast.SubjectSetRewrite{
				Operation: ast.OperatorAnd,
				Children: ast.Children{
					&ast.ComputedSubjectSet{Relation: "A"},
					&ast.InvertResult{
						Child: &ast.SubjectSetRewrite{
							Operation: ast.OperatorOr,
							Children: ast.Children{
								&ast.ComputedSubjectSet{Relation: "B"},
								&ast.ComputedSubjectSet{Relation: "C"},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, simplifyExpression(tc.input))
		})
	}
}
