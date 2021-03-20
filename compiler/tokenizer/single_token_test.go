package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEof(t *testing.T) {
	toks := NewTokenStream("")
	tok := toks.Next()
	require.Equal(t, KindNil, tok.Kind)
	require.Equal(t, Position(0), tok.Start)
	require.Equal(t, Position(0), tok.End)
	require.Nil(t, tok.Value)
}

func TestBlankToken(t *testing.T) {
	runSingleTokenTest(t, " \t  \r\n  \t", KindBlank, nil)
}

func TestDecimalIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "1234_56789", KindInteger, int64(123456789))
	runSingleTokenTest(t, "-1234_56789", KindInteger, int64(-123456789))
}

func TestHexIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0xB33F_cafe", KindInteger, int64(0xB33FCAFE))
}

func TestOctalIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0o1234_567", KindInteger, int64(0o1234567))
}

func TestBinaryIntegerToken(t *testing.T) {
	runSingleTokenTest(t, "0b1010_101", KindInteger, int64(0b1010101))
}

func TestFloatToken(t *testing.T) {
	runSingleTokenTest(t, "3.14159", KindFloat, float64(3.14159))
	runSingleTokenTest(t, "03.14159", KindFloat, float64(3.14159))
	runSingleTokenTest(t, "-03.14159", KindFloat, float64(-3.14159))
	runSingleTokenTest(t, "-.14159", KindFloat, float64(-0.14159))
	runSingleTokenTest(t, ".14159", KindFloat, float64(0.14159))
	runSingleTokenTest(t, "1e10", KindFloat, float64(1e10))
	runSingleTokenTest(t, "-3.14E10", KindFloat, float64(-3.14e10))
}

func TestIdentifier(t *testing.T) {
	runSingleTokenTest(t, "anIdentifier", KindIdent, "anIdentifier")
	runSingleTokenTest(t, "with_underscores", KindIdent, "with_underscores")
	runSingleTokenTest(t, "_leadingUnderscore", KindIdent, "_leadingUnderscore")
	runSingleTokenTest(t, "including_1_digit", KindIdent, "including_1_digit")
	runSingleTokenTest(t, "unicode_ʔ_is_allowed", KindIdent, "unicode_ʔ_is_allowed")
	runLeadingTokenTest(t, "1cant_start_with_a_digit_though", "1", KindInteger, int64(1))
}

func TestSymbols(t *testing.T) {
	runSingleTokenTest(t, "(", KindLParen, nil)
	runSingleTokenTest(t, ")", KindRParen, nil)
}

func runSingleTokenTest(t *testing.T, input string, kind TokenKind, value interface{}) {
	toks := runLeadingTokenTest(t, input, input, kind, value)

	tok := toks.Next()
	require.Equal(t, KindNil, tok.Kind)
}

func runLeadingTokenTest(t *testing.T, input, token string, kind TokenKind, value interface{}) *TokenStream {
	toks := NewTokenStream(input)
	tok := toks.Next()

	if tok.Kind == KindError {
		require.NoError(t, tok.Value.(error))
	}

	require.Equal(t, kind, tok.Kind)
	require.Equal(t, value, tok.Value)
	require.Equal(t, token, input[tok.Start:tok.End])
	require.Equal(t, Position(0), tok.Start)
	require.Equal(t, Position(len(token)), tok.End)

	return toks
}
