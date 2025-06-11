package calc

import (
	"errors"
	"math"
	"reflect"
)

type Namespace interface {
	Get(key string) (any, bool)
}

type node interface{ exec(namespace Namespace) any }

type numNode struct{ val float64 }

func (n *numNode) exec(_ Namespace) any { return n.val }

const (
	addOp uint8 = iota + 1
	subOp
	mulOp
	divOp
	powOp
	eqOp
	notEqOp
	lessOp
	lessEqOp
	moreOp
	moreEqOp
	andOp
	orOp
)

type unaryNode struct {
	op  uint8
	val node
}

func (n *unaryNode) exec(namespace Namespace) any {
	val := n.val.exec(namespace)
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

func (n *binaryNode) exec(namespace Namespace) any {
	left := n.left.exec(namespace)
	if _, ok := left.(error); ok {
		return left
	}

	right := n.right.exec(namespace)
	if _, ok := right.(error); ok {
		return right
	}

	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return errors.New("")
	}

	switch n.op {
	case eqOp:
		switch left.(type) {
		case float64:
			return left.(float64) == right.(float64)
		case bool:
			return left.(bool) == right.(bool)
		case string:
			return left.(string) == right.(string)
		default:
			return errors.New("")
		}
	case notEqOp:
		switch left.(type) {
		case float64:
			return left.(float64) != right.(float64)
		case bool:
			return left.(bool) != right.(bool)
		case string:
			return left.(string) != right.(string)
		default:
			return errors.New("")
		}
	case lessOp:
		switch left.(type) {
		case float64:
			return left.(float64) < right.(float64)
		case string:
			return left.(string) < right.(string)
		default:
			return errors.New("")
		}
	case lessEqOp:
		switch left.(type) {
		case float64:
			return left.(float64) <= right.(float64)
		case string:
			return left.(string) <= right.(string)
		default:
			return errors.New("")
		}
	case moreOp:
		switch left.(type) {
		case float64:
			return left.(float64) > right.(float64)
		case string:
			return left.(string) > right.(string)
		default:
			return errors.New("")
		}
	case moreEqOp:
		switch left.(type) {
		case float64:
			return left.(float64) >= right.(float64)
		case string:
			return left.(string) >= right.(string)
		default:
			return errors.New("")
		}

	case andOp:
		if _, ok := left.(bool); !ok {
			return errors.New("")
		}
		return left.(bool) && right.(bool)

	case orOp:
		if _, ok := left.(bool); !ok {
			return errors.New("")
		}
		return left.(bool) || right.(bool)
	}

	if n.op == addOp {
		switch left := left.(type) {
		case float64:
			return left + right.(float64)
		case string:
			return left + right.(string)
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

func (n *errNode) exec(_ Namespace) any { return n.err }

func isErr(n node) bool {
	_, ok := n.(*errNode)
	return ok
}

type ternaryNode struct {
	cond    node
	ifTrue  node
	ifFalse node
}

func (n *ternaryNode) exec(namespace Namespace) any {
	cond := n.cond.exec(namespace)
	if _, ok := cond.(error); ok {
		return cond
	}

	if _, ok := cond.(bool); !ok {
		return errors.New("")
	}

	if cond.(bool) {
		return n.ifTrue.exec(namespace)
	}
	return n.ifFalse.exec(namespace)
}

type strNode struct{ val string }

func (n *strNode) exec(_ Namespace) any { return n.val }

type identNode struct{ val string }

func (n *identNode) exec(namespace Namespace) any {
	val, ok := namespace.Get(n.val)
	if !ok {
		return errors.New("")
	}
	switch v := val.(type) {
	case int:
		val = float64(v)
	case int8:
		val = float64(v)
	case int16:
		val = float64(v)
	case int32:
		val = float64(v)
	case int64:
		val = float64(v)
	case uint:
		val = float64(v)
	case uint8:
		val = float64(v)
	case uint16:
		val = float64(v)
	case uint32:
		val = float64(v)
	case uint64:
		val = float64(v)
	case float32:
		val = float64(v)
	}

	return val
}
