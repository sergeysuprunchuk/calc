package calc

import (
	"context"
	"errors"
	"math"
	"reflect"
)

type node interface{ exec(ctx context.Context) any }

type numNode struct{ val float64 }

func (n *numNode) exec(_ context.Context) any { return n.val }

const (
	addOp uint8 = iota + 1
	subOp
	mulOp
	divOp
	powOp
)

type unaryNode struct {
	op  uint8
	val node
}

func (n *unaryNode) exec(ctx context.Context) any {
	val := n.val.exec(ctx)
	if _, ok := val.(error); ok {
		return val
	}

	switch n.op {
	case subOp:
		if _, ok := val.(float64); !ok {
			return errors.New("")
		}
		return -val.(float64)

	default:
		return errors.New("")
	}
}

type binaryNode struct {
	op    uint8
	left  node
	right node
}

func (n *binaryNode) exec(ctx context.Context) any {
	left := n.left.exec(ctx)
	if _, ok := left.(error); ok {
		return left
	}

	right := n.right.exec(ctx)
	if _, ok := right.(error); ok {
		return right
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return errors.New("")
	}

	if n.op == addOp {
		switch left := left.(type) {
		case float64:
			return left + right.(float64)
		default:
			return errors.New("")
		}
	}

	if _, ok := left.(float64); !ok {
		return errors.New("")
	}

	switch n.op {
	case subOp:
		return left.(float64) - right.(float64)
	case mulOp:
		return left.(float64) * right.(float64)
	case divOp:
		return left.(float64) / right.(float64)
	case powOp:
		return math.Pow(left.(float64), right.(float64))
	default:
		return errors.New("")
	}
}

type errNode struct{ err error }

func (n *errNode) exec(_ context.Context) any { return n.err }

func isErr(n node) bool {
	_, ok := n.(*errNode)
	return ok
}
