package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimplePrintln(t *testing.T) {
	toks := runTokenizerTest("println(42)")
	require.Equal(t, Token{KindIdent, 0, 7, "println"}, toks[0])
	require.Equal(t, Token{KindLParen, 7, 8, nil}, toks[1])
	require.Equal(t, Token{KindInteger, 8, 10, int64(42)}, toks[2])
	require.Equal(t, Token{KindRParen, 10, 11, nil}, toks[3])
}

func runTokenizerTest(input string) []Token {
	strm := NewTokenStream(input)
	var toks []Token
	for tok := strm.Next(); tok.Kind != KindNil; tok = strm.Next() {
		toks = append(toks, tok)
	}
	return toks
}
