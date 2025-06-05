package calc

import (
	"testing"
)

func Test_readNum(t *testing.T) {
	tests := []struct {
		tok      *tokenizer
		expected token
		cursor   int
	}{
		{tok: &tokenizer{data: []rune("16")}, expected: token{numTyp, "16"}, cursor: 2},
		{tok: &tokenizer{data: []rune("16.32")}, expected: token{numTyp, "16.32"}, cursor: 5},
		{tok: &tokenizer{data: []rune("32.16.32")}, expected: token{numTyp, "32.16"}, cursor: 5},
		{tok: &tokenizer{data: []rune(".64")}, expected: token{numTyp, "0.64"}, cursor: 3},
		{tok: &tokenizer{data: []rune(".64.16")}, expected: token{numTyp, "0.64"}, cursor: 3},
		{tok: &tokenizer{data: []rune("0")}, expected: token{numTyp, "0"}, cursor: 1},
		{tok: &tokenizer{data: []rune("0.")}, expected: token{numTyp, "0"}, cursor: 1},
		{tok: &tokenizer{data: []rune("0.64")}, expected: token{numTyp, "0.64"}, cursor: 4},
		{tok: &tokenizer{data: []rune(".")}, expected: token{typ: emptyTyp}, cursor: 0},
		{tok: &tokenizer{data: []rune("6416..16")}, expected: token{numTyp, "6416"}, cursor: 4},
		{tok: &tokenizer{data: []rune("..16")}, expected: token{typ: emptyTyp}, cursor: 0},
		{tok: &tokenizer{data: []rune("16 b6f")}, expected: token{numTyp, "16"}, cursor: 2},
		{tok: &tokenizer{data: []rune("16.32b6f")}, expected: token{numTyp, "16.32"}, cursor: 5},
		{tok: &tokenizer{data: []rune("b6f")}, expected: token{typ: emptyTyp}, cursor: 0},
		{tok: &tokenizer{data: []rune("")}, expected: token{typ: emptyTyp}, cursor: 0},
		{tok: &tokenizer{data: []rune("64.b6f")}, expected: token{numTyp, "64"}, cursor: 2},
	}

	for _, test := range tests {
		got := test.tok.readNum()
		if got != test.expected {
			t.Errorf("tok: expected %v, got %v", test.expected, got)
		}

		if test.tok.cursor != test.cursor {
			t.Errorf("cursor: expected %d, got %d", test.cursor, test.tok.cursor)
		}
	}
}
