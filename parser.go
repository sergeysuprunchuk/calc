package calc

import (
	"errors"
	"strconv"
)

type parser struct{ tok *tokenizer }

func newParser(data string) *parser {
	return &parser{tok: newTokenizer(data)}
}

const (
	addOp uint8 = iota + 1
	subOp
	mulOp
	divOp
	powOp
)

type node interface{}

type numNode struct{ val float64 }

type unaryNode struct {
	op  uint8
	val node
}

type binaryNode struct {
	op    uint8
	left  node
	right node
}

type errNode struct{ err error }

func isErr(n node) bool {
	_, ok := n.(*errNode)
	return ok
}

/*
чем больше число в конце, тем ниже приоритет (parse0 > parse1),
разбор начинается с большего числа.
parse считывает токен, если токен соответствует шаблону,
то он перемещает курсор на следующий токен
*/
type parse func() node

func (p *parser) parse() node {
	p.tok.nextTok()
	if p.tok.currentTok().typ == eofTyp {
		return nil
	}

	n := p.parse3()
	if isErr(n) {
		return n
	}

	if p.tok.currentTok().typ != eofTyp {
		return &errNode{errors.New("не удалось разобрать выражение")}
	}

	return n
}

func (p *parser) parse0() node {
	tok := p.tok.currentTok()

	if tok.typ == numTyp {
		p.tok.nextTok()
		val, err := strconv.ParseFloat(tok.val, 64)
		if err != nil {
			return &errNode{err}
		}
		return &numNode{val}
	}

	if tok.typ == lParenTyp {
		p.tok.nextTok()
		//parse с самым низким приоритетом
		n := p.parse3()
		if isErr(n) {
			return n
		}

		if p.tok.currentTok().typ != rParenTyp {
			return &errNode{errors.New("ожидалось ')'")}
		}

		p.tok.nextTok()
		return n
	}

	if tok.typ == minusTyp {
		p.tok.nextTok()
		val := p.parse0()
		if isErr(val) {
			return val
		}
		return &unaryNode{subOp, val}
	}

	return &errNode{
		errors.New("ожидалось число | '(' | '-'"),
	}
}

func (p *parser) parse1() node {
	n := p.parse0()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == powerTyp; tok = p.tok.currentTok() {
		p.tok.nextTok()
		right := p.parse0()
		if isErr(right) {
			return right
		}
		n = &binaryNode{powOp, n, right}
	}

	return n
}

func (p *parser) parse2() node {
	n := p.parse1()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == mulTyp || tok.typ == slashTyp; tok = p.tok.currentTok() {
		p.tok.nextTok()

		right := p.parse1()
		if isErr(right) {
			return right
		}

		typ := mulOp
		if tok.typ == slashTyp {
			typ = divOp
		}

		n = &binaryNode{typ, n, right}
	}

	return n
}

func (p *parser) parse3() node {
	n := p.parse2()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == plusTyp || tok.typ == minusTyp; tok = p.tok.currentTok() {
		p.tok.nextTok()

		right := p.parse2()
		if isErr(right) {
			return right
		}

		typ := addOp
		if tok.typ == minusTyp {
			typ = subOp
		}

		n = &binaryNode{typ, n, right}
	}

	return n
}
