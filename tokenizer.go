package calc

import (
	"strings"
	"unicode"
)

type tokenizer struct {
	data   []rune
	cursor int
	tok    token //последний прочитанный токен
}

func newTokenizer(data string) *tokenizer {
	return &tokenizer{data: []rune(data)}
}

func (t *tokenizer) char() rune {
	if t.cursor >= len(t.data) {
		return 0
	}
	return t.data[t.cursor]
}

func (t *tokenizer) nextChar() rune {
	if t.cursor+1 >= len(t.data) {
		return 0
	}
	return t.data[t.cursor+1]
}

func (t *tokenizer) next() { t.cursor++ }

func (t *tokenizer) skipSpace() {
	for {
		if t.char() == 0 || unicode.IsSpace(t.char()) == false {
			break
		}
		t.next()
	}
}

const (
	emptyTyp uint8 = iota + 1
	errTyp
	eofTyp
	numTyp
	plusTyp
	minusTyp
	mulTyp
	slashTyp
	powerTyp
	lParenTyp
	rParenTyp
	eqTyp
	notEqTyp
	moreTyp
	lessTyp
	moreEqTyp
	lessEqTyp
	andTyp
	orTyp
	questionTyp
	colonTyp
	strTyp
)

type token struct {
	typ uint8
	val string
}

/*
читает один токен из последовательности символов, начиная с текущей позиции.
если текущий символ соответствует входному символу целевого токена (кавычка для строки и т.д.), функция продолжает
считывать символы, пока они соответствуют шаблону токена, и возвращает сформированный токен.
если текущий символ не является входным для целевого токена, возвращает emptyTyp токен.
*/
type reader func() token

// считывает один токен, каждый вызов возвращает новый токен!
func (t *tokenizer) nextTok() (tok token) {
	defer func() { t.tok = tok }()

	t.skipSpace()

	if t.char() == 0 {
		return token{typ: eofTyp}
	}

	for _, r := range []reader{
		t.readNum,
		t.readOperator,
		t.readStr,
	} {
		tok := r()
		if tok.typ != emptyTyp {
			return tok
		}
	}

	return token{errTyp, "неизвестный символ " + string(t.char())}
}

func (t *tokenizer) readStr() token {
	quote := t.char()

	if quote != '"' && quote != '\'' {
		return token{typ: emptyTyp}
	}
	var builder strings.Builder

	for {
		t.next()
		if t.char() == quote {
			t.next()
			break
		}

		if t.char() == 0 {
			return token{errTyp, "ожидалось " + string(quote)}
		}

		builder.WriteRune(t.char())
	}

	return token{strTyp, builder.String()}
}

func (t *tokenizer) readNum() token {
	if t.char() != '.' && unicode.IsDigit(t.char()) == false {
		return token{typ: emptyTyp}
	}

	var val strings.Builder
	var float bool

	if t.char() == '0' {
		val.WriteRune('0')
		t.next()
		if t.char() != '.' || unicode.IsDigit(t.nextChar()) == false {
			return token{numTyp, val.String()}
		}
	} else if t.char() == '.' {
		if unicode.IsDigit(t.nextChar()) == false {
			return token{typ: emptyTyp}
		}
		val.WriteRune('0')
	}

	for ; ; t.next() {
		if unicode.IsDigit(t.char()) {
			val.WriteRune(t.char())
			continue
		}

		if t.char() == '.' && unicode.IsDigit(t.nextChar()) && float == false {
			val.WriteRune('.')
			float = true
			continue
		}

		break
	}

	return token{numTyp, val.String()}
}

func (t *tokenizer) readOperator() token {
	var tok token

	switch t.char() {
	case '?':
		t.next()
		return token{typ: questionTyp}
	case ':':
		t.next()
		return token{typ: colonTyp}

	case '=':
		if t.nextChar() != '=' {
			return token{typ: emptyTyp}
		}
		t.next()
		t.next()
		return token{typ: eqTyp}

	case '!':
		if t.nextChar() != '=' {
			return token{typ: emptyTyp}
		}
		t.next()
		t.next()
		return token{typ: notEqTyp}

	case '<':
		t.next()
		if t.char() != '=' {
			return token{typ: lessTyp}
		}
		t.next()
		return token{typ: lessEqTyp}

	case '>':
		t.next()
		if t.char() != '=' {
			return token{typ: moreTyp}
		}
		t.next()
		return token{typ: moreEqTyp}

	case '&':
		if t.nextChar() != '&' {
			return token{typ: emptyTyp}
		}
		t.next()
		t.next()
		return token{typ: andTyp}

	case '|':
		if t.nextChar() != '|' {
			return token{typ: emptyTyp}
		}
		t.next()
		t.next()
		return token{typ: orTyp}

	case '+':
		tok.typ = plusTyp
	case '-':
		tok.typ = minusTyp
	case '*':
		if t.nextChar() == '*' {
			t.next()
			tok.typ = powerTyp
			break
		}
		tok.typ = mulTyp
	case '/':
		tok.typ = slashTyp
	case '(':
		tok.typ = lParenTyp
	case ')':
		tok.typ = rParenTyp
	default:
		return token{typ: emptyTyp}
	}

	t.next()

	return tok
}

func (t *tokenizer) currentTok() token { return t.tok }
