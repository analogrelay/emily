package scanner

import (
	"testing"

	"github.com/anurse/emily/compiler/token"
	"github.com/stretchr/testify/require"
)

func TestSimplePrintln(t *testing.T) {
	toks := runScannerTest("println(42)")
	require.Equal(t, token.Token{token.KindIdent, 0, 7, "println"}, toks[0])
	require.Equal(t, token.Token{token.KindLParen, 7, 8, nil}, toks[1])
	require.Equal(t, token.Token{token.KindInteger, 8, 10, int64(42)}, toks[2])
	require.Equal(t, token.Token{token.KindRParen, 10, 11, nil}, toks[3])
}

func runScannerTest(input string) []token.Token {
	strm := NewScanner(input)
	var toks []token.Token
	for strm.Scan() {
		toks = append(toks, strm.Token())
	}
	return toks
}
