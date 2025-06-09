package calc

import (
	"errors"
	"strconv"
)

type parser struct{ tok *tokenizer }

func newParser(data string) *parser {
	return &parser{tok: newTokenizer(data)}
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

	n := p.parse5()
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
		n := p.parse5()
		if isErr(n) {
			return n
		}

		if p.tok.currentTok().typ != rParenTyp {
			return &errNode{errors.New("ожидалось ')'")}
		}

		p.tok.nextTok()

		return n
	}

	return &errNode{
		errors.New("ожидалось число | '('"),
	}
}

func (p *parser) parse1() node {
	n := p.parse0()
	if isErr(n) {
		return n
	}

	if p.tok.currentTok().typ == powerTyp {
		p.tok.nextTok()

		right := p.parse1()
		if isErr(right) {
			return right
		}

		n = &binaryNode{powOp, n, right}
	}

	return n
}

func (p *parser) parse2() node {
	if p.tok.currentTok().typ == minusTyp {
		p.tok.nextTok()

		val := p.parse1()
		if isErr(val) {
			return val
		}

		return &unaryNode{subOp, val}
	}

	return p.parse1()
}

func (p *parser) parse3() node {
	n := p.parse2()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == mulTyp ||
		tok.typ == slashTyp; tok = p.tok.currentTok() {

		p.tok.nextTok()

		right := p.parse2()
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

func (p *parser) parse4() node {
	n := p.parse3()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == plusTyp ||
		tok.typ == minusTyp; tok = p.tok.currentTok() {

		p.tok.nextTok()

		right := p.parse3()
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

func (p *parser) parse5() node {
	n := p.parse4()
	if isErr(n) {
		return n
	}

	for tok := p.tok.currentTok(); tok.typ == eqTyp ||
		tok.typ == notEqTyp ||
		tok.typ == lessTyp ||
		tok.typ == lessEqTyp ||
		tok.typ == moreTyp ||
		tok.typ == moreEqTyp; tok = p.tok.currentTok() {

		p.tok.nextTok()

		right := p.parse4()
		if isErr(right) {
			return right
		}

		var typ uint8
		switch tok.typ {
		case notEqTyp:
			typ = notEqOp
		case moreTyp:
			typ = moreOp
		case lessTyp:
			typ = lessOp
		case moreEqTyp:
			typ = moreEqOp
		case lessEqTyp:
			typ = lessEqOp
		default:
			typ = eqOp
		}

		n = &binaryNode{typ, n, right}
	}

	return n
}
