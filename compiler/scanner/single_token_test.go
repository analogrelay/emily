package scanner

import (
	"testing"

	"github.com/anurse/emily/compiler/token"
	"github.com/stretchr/testify/require"
)

func TestEof(t *testing.T) {
	toks := NewScanner("")
	require.False(t, toks.Scan())
	require.Equal(t, token.KindNil, toks.Kind())
	require.Equal(t, token.Position(0), toks.Token().Start)
	require.Equal(t, token.Position(0), toks.Token().End)
	require.Nil(t, toks.Token().Value)
}

func TestBlankToken(t *testing.T) {
	runSingleTokenTest(t, " \t  \r\n  \t", token.KindBlank, nil)
}

func TestDecimalIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "1234_56789", token.KindInteger, int64(123456789))
	runSingleTokenTest(t, "-1234_56789", token.KindInteger, int64(-123456789))
}

func TestHexIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0xB33F_cafe", token.KindInteger, int64(0xB33FCAFE))
}

func TestOctalIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0o1234_567", token.KindInteger, int64(0o1234567))
}

func TestBinaryIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0b1010_101", token.KindInteger, int64(0b1010101))
}

func TestFloatToken(t *testing.T) {
	runSingleTokenTest(t, "3.14159", token.KindFloat, float64(3.14159))
	runSingleTokenTest(t, "03.14159", token.KindFloat, float64(3.14159))
	runSingleTokenTest(t, "-03.14159", token.KindFloat, float64(-3.14159))
	runSingleTokenTest(t, "-.14159", token.KindFloat, float64(-0.14159))
	runSingleTokenTest(t, ".14159", token.KindFloat, float64(0.14159))
	runSingleTokenTest(t, "1e10", token.KindFloat, float64(1e10))
	runSingleTokenTest(t, "-3.14E10", token.KindFloat, float64(-3.14e10))
}

func TestIdentifier(t *testing.T) {
	runSingleTokenTest(t, "anIdentifier", token.KindIdent, "anIdentifier")
	runSingleTokenTest(t, "with_underscores", token.KindIdent, "with_underscores")
	runSingleTokenTest(t, "_leadingUnderscore", token.KindIdent, "_leadingUnderscore")
	runSingleTokenTest(t, "including_1_digit", token.KindIdent, "including_1_digit")
	runSingleTokenTest(t, "unicode_ʔ_is_allowed", token.KindIdent, "unicode_ʔ_is_allowed")
	runLeadingTokenTest(t, "1cant_start_with_a_digit_though", "1", token.KindInteger, int64(1))
}

func TestSymbols(t *testing.T) {
	runSingleTokenTest(t, "(", token.KindLParen, nil)
	runSingleTokenTest(t, ")", token.KindRParen, nil)
}

func runSingleTokenTest(t *testing.T, input string, kind token.TokenKind, value interface{}) {
	toks := runLeadingTokenTest(t, input, input, kind, value)

	require.False(t, toks.Scan())
}

func runLeadingTokenTest(t *testing.T, input, tokenContent string, kind token.TokenKind, value interface{}) *Scanner {
	toks := NewScanner(input)
	require.True(t, toks.Scan())

	require.NoError(t, toks.Err())
	require.Equal(t, kind, toks.Kind())
	require.Equal(t, value, toks.Token().Value)
	require.Equal(t, tokenContent, input[toks.Token().Start:toks.Token().End])
	require.Equal(t, token.Position(0), toks.Token().Start)
	require.Equal(t, token.Position(len(tokenContent)), toks.Token().End)

	return toks
}
