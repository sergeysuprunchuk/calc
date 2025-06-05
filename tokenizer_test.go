package calc

import (
	"testing"
)

func getTok(data string) *tokenizer { return &tokenizer{data: []rune(data)} }

type testRow struct {
	tok      *tokenizer
	expected token
	cursor   int
}

func tr(data string, expected token, cursor int) testRow {
	return testRow{
		tok:      getTok(data),
		expected: expected,
		cursor:   cursor,
	}
}

func (test *testRow) check(t *testing.T, got token) {
	if got != test.expected {
		t.Errorf("tok: expected %v, got %v", test.expected, got)
	}

	if test.tok.cursor != test.cursor {
		t.Errorf("cursor: expected %d, got %d", test.cursor, test.tok.cursor)
	}
}

func Test_readNum(t *testing.T) {
	tests := []testRow{
		tr("16", token{numTyp, "16"}, 2),
		tr("16.32", token{numTyp, "16.32"}, 5),
		tr("32.16.32", token{numTyp, "32.16"}, 5),
		tr(".64", token{numTyp, "0.64"}, 3),
		tr(".64.16", token{numTyp, "0.64"}, 3),
		tr("0", token{numTyp, "0"}, 1),
		tr("0.", token{numTyp, "0"}, 1),
		tr("0.64", token{numTyp, "0.64"}, 4),
		tr(".", token{typ: emptyTyp}, 0),
		tr("6416..16", token{numTyp, "6416"}, 4),
		tr("..16", token{typ: emptyTyp}, 0),
		tr("16 b6f", token{numTyp, "16"}, 2),
		tr("16.32b6f", token{numTyp, "16.32"}, 5),
		tr("b6f", token{typ: emptyTyp}, 0),
		tr("", token{typ: emptyTyp}, 0),
		tr("64.b6f", token{numTyp, "64"}, 2),
		tr("+", token{typ: emptyTyp}, 0),
		tr("-", token{typ: emptyTyp}, 0),
		tr("*", token{typ: emptyTyp}, 0),
		tr("**", token{typ: emptyTyp}, 0),
		tr("/", token{typ: emptyTyp}, 0),
	}

	for _, test := range tests {
		test.check(t, test.tok.readNum())
	}
}

func Test_readOperator(t *testing.T) {
	tests := []testRow{
		tr("+", token{typ: plusTyp}, 1),
		tr("-", token{typ: minusTyp}, 1),
		tr("*", token{typ: mulTyp}, 1),
		tr("**", token{typ: powerTyp}, 2),
		tr("/", token{typ: slashTyp}, 1),
		tr("++", token{typ: plusTyp}, 1),
		tr("--", token{typ: minusTyp}, 1),
		tr("**", token{typ: powerTyp}, 2),
		tr("****", token{typ: powerTyp}, 2),
		tr("//", token{typ: slashTyp}, 1),
		tr("+-", token{typ: plusTyp}, 1),
		tr("-+", token{typ: minusTyp}, 1),
		tr("*+", token{typ: mulTyp}, 1),
		tr("**+", token{typ: powerTyp}, 2),
		tr("/+", token{typ: slashTyp}, 1),
		tr("", token{typ: emptyTyp}, 0),
		tr("16.32", token{typ: emptyTyp}, 0),
		tr(" +", token{typ: emptyTyp}, 0),
		tr("	+", token{typ: emptyTyp}, 0),
		tr("32+", token{typ: emptyTyp}, 0),
		tr("*+*", token{typ: mulTyp}, 1),
	}

	for _, test := range tests {
		test.check(t, test.tok.readOperator())
	}
}

func Test_read(t *testing.T) {
	type item struct {
		tok    token
		cursor int
	}

	tests := []struct {
		tok      *tokenizer
		expected []item
	}{
		{
			tok: getTok("16.32 + -.32"),
			expected: []item{
				{token{numTyp, "16.32"}, 5},
				{token{typ: plusTyp}, 7},
				{token{typ: minusTyp}, 9},
				{token{numTyp, "0.32"}, 12},
			},
		},
		{
			tok: getTok("16.32 + -.32**2"),
			expected: []item{
				{token{numTyp, "16.32"}, 5},
				{token{typ: plusTyp}, 7},
				{token{typ: minusTyp}, 9},
				{token{numTyp, "0.32"}, 12},
				{token{typ: powerTyp}, 14},
				{token{numTyp, "2"}, 15},
			},
		},
		{
			tok: getTok("16.32 + -.32**-2.001"),
			expected: []item{
				{token{numTyp, "16.32"}, 5},
				{token{typ: plusTyp}, 7},
				{token{typ: minusTyp}, 9},
				{token{numTyp, "0.32"}, 12},
				{token{typ: powerTyp}, 14},
				{token{typ: minusTyp}, 15},
				{token{numTyp, "2.001"}, 20},
			},
		},
		{
			tok: getTok("-0.032*-2.001"),
			expected: []item{
				{token{typ: minusTyp}, 1},
				{token{numTyp, "0.032"}, 6},
				{token{typ: mulTyp}, 7},
				{token{typ: minusTyp}, 8},
				{token{numTyp, "2.001"}, 13},
			},
		},
		{
			tok: getTok("	16.32 	+ -	.32	"),
			expected: []item{
				{token{numTyp, "16.32"}, 6},
				{token{typ: plusTyp}, 9},
				{token{typ: minusTyp}, 11},
				{token{numTyp, "0.32"}, 15},
				{token{typ: eofTyp}, 16},
			},
		},
		{
			tok: getTok(""),
			expected: []item{
				{token{typ: eofTyp}, 0},
				{token{typ: eofTyp}, 0},
				{token{typ: eofTyp}, 0},
			},
		},
		{
			tok: getTok("32.	0"),
			expected: []item{
				{token{numTyp, "32"}, 2},
				{token{errTyp, "неизвестный символ ."}, 2},
				{token{errTyp, "неизвестный символ ."}, 2},
				{token{errTyp, "неизвестный символ ."}, 2},
			},
		},
	}

	for _, test := range tests {
		for index, i := range test.expected {
			got := test.tok.read()

			if got != i.tok {
				t.Errorf("%d: tok: expected %v, got %v", index, i.tok, got)
			}

			if test.tok.cursor != i.cursor {
				t.Errorf("%d: cursor: expected %d, got %d", index, i.cursor, test.tok.cursor)
			}
		}
	}
}
