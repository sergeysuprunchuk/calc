package calc

import (
	"errors"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		data     string
		expected node
	}{
		{
			data:     "16",
			expected: &numNode{16.},
		},
		{
			data: "16	 +32",
			expected: &binaryNode{
				op:    addOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "16*64	 + 		32",
			expected: &binaryNode{
				op: addOp,
				left: &binaryNode{
					op:    mulOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
				right: &numNode{32.},
			},
		},
		{
			data: "16 ** 32 	** 64	",
			expected: &binaryNode{
				op: powOp,
				right: &binaryNode{
					op:    powOp,
					left:  &numNode{32.},
					right: &numNode{64.},
				},
				left: &numNode{16.},
			},
		},
		{
			data: "(		16 + 64	) *	 32",
			expected: &binaryNode{
				op: mulOp,
				left: &binaryNode{
					op:    addOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
				right: &numNode{32.},
			},
		},
		{
			data: "	16 + 64	 * 32  **	 64 - 16 / 64	 ",
			expected: &binaryNode{
				op: subOp,
				left: &binaryNode{
					op:   addOp,
					left: &numNode{16.},
					right: &binaryNode{
						op:   mulOp,
						left: &numNode{64.},
						right: &binaryNode{
							op:    powOp,
							left:  &numNode{32.},
							right: &numNode{64.},
						},
					},
				},
				right: &binaryNode{
					op:    divOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
			},
		},
		{
			data:     "16 ++ 32",
			expected: &errNode{errors.New("ожидалось число | '('")},
		},
		{
			data:     "32 * (16 + 64",
			expected: &errNode{errors.New("ожидалось ')'")},
		},
		{data: "", expected: nil},
		{
			data: "16	 ==	32",
			expected: &binaryNode{
				op:    eqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "	16	 !=	32",
			expected: &binaryNode{
				op:    notEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "16	 <=	32		",
			expected: &binaryNode{
				op:    lessEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "16	 >= 	32",
			expected: &binaryNode{
				op:    moreEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: " 16		 >	32 ",
			expected: &binaryNode{
				op:    moreOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "16	 <	32	",
			expected: &binaryNode{
				op:    lessOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			data: "16+	16 ==	32",
			expected: &binaryNode{
				op: eqOp,
				left: &binaryNode{
					op:    addOp,
					left:  &numNode{16.},
					right: &numNode{16.},
				},
				right: &numNode{32.},
			},
		},
		{
			data: "16+	16 !=	32+16",
			expected: &binaryNode{
				op: notEqOp,
				left: &binaryNode{
					op:    addOp,
					left:  &numNode{16.},
					right: &numNode{16.},
				},
				right: &binaryNode{
					op:    addOp,
					left:  &numNode{32.},
					right: &numNode{16.},
				},
			},
		},
		{
			data: "16	 <	32	&&	 16	!=		32",
			expected: &binaryNode{
				op: andOp,
				left: &binaryNode{
					op:    lessOp,
					left:  &numNode{16.},
					right: &numNode{32.},
				},
				right: &binaryNode{
					op:    notEqOp,
					left:  &numNode{16.},
					right: &numNode{32.},
				},
			},
		},
		{
			data: "16	 <	32	||	 16	!=		32",
			expected: &binaryNode{
				op: orOp,
				left: &binaryNode{
					op:    lessOp,
					left:  &numNode{16.},
					right: &numNode{32.},
				},
				right: &binaryNode{
					op:    notEqOp,
					left:  &numNode{16.},
					right: &numNode{32.},
				},
			},
		},
		{
			data:     `'привет мир'`,
			expected: &strNode{`привет мир`},
		},
		{
			data:     `"привет мир"`,
			expected: &strNode{`привет мир`},
		},
		{
			data:     `"привет ' мир"`,
			expected: &strNode{`привет ' мир`},
		},
		{
			data:     `'привет " мир'`,
			expected: &strNode{`привет " мир`},
		},
	}

	for _, test := range tests {
		n := newParser(test.data).parse()
		if !reflect.DeepEqual(n, test.expected) {
			t.Errorf("parse(%q): got %#v, want %#v", test.data, n, test.expected)
		}
	}
}
