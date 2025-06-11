package calc

import (
	"reflect"
	"testing"
)

func Test_node(t *testing.T) {
	tests := []struct {
		n        node
		expected any
	}{
		{n: &numNode{16.}, expected: 16.},
		{n: &numNode{32.}, expected: 32.},
		{n: &numNode{64.64}, expected: 64.64},
		{n: &unaryNode{subOp, &numNode{32.}}, expected: -32.},
		{n: &unaryNode{subOp, &numNode{64.64}}, expected: -64.64},
		{
			n: &binaryNode{
				op: addOp, left: &numNode{32.}, right: &numNode{64.64}},
			expected: 96.64,
		},
		{
			n: &binaryNode{
				op: subOp, left: &numNode{32.}, right: &numNode{16.}},
			expected: 16.,
		},
		{
			n: &binaryNode{
				op: mulOp, left: &numNode{16.}, right: &numNode{64.}},
			expected: 1024.,
		},
		{
			n: &binaryNode{
				op: divOp, left: &numNode{64.}, right: &numNode{16.}},
			expected: 4.,
		},
		{
			n: &binaryNode{
				op: powOp, left: &numNode{16.}, right: &numNode{4.}},
			expected: 65536.,
		},
		{
			n: &binaryNode{
				op: addOp,
				left: &binaryNode{
					op:    mulOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
				right: &numNode{32.},
			},
			expected: 1056.,
		},
		{
			n: &binaryNode{
				op: mulOp,
				left: &binaryNode{
					op:    addOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
				right: &numNode{32.},
			},
			expected: 2560.,
		},
		{
			n: &binaryNode{
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
							right: &numNode{2.},
						},
					},
				},
				right: &binaryNode{
					op:    divOp,
					left:  &numNode{16.},
					right: &numNode{64.},
				},
			},
			expected: 65551.75,
		},
		{
			expected: false,
			n: &binaryNode{
				op:    eqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			expected: true,
			n: &binaryNode{
				op:    notEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			expected: true,
			n: &binaryNode{
				op:    lessEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			expected: false,
			n: &binaryNode{
				op:    moreEqOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			expected: false,
			n: &binaryNode{
				op:    moreOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
		{
			expected: true,
			n: &binaryNode{
				op:    lessOp,
				left:  &numNode{16.},
				right: &numNode{32.},
			},
		},
	}

	for _, test := range tests {
		val := test.n.exec(nil)
		if !reflect.DeepEqual(val, test.expected) {
			t.Errorf("got %v, want %v", val, test.expected)
		}
	}
}
