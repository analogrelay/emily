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

func runSingleTokenTest(t *testing.T, input string, kind TokenKind, value interface{}) {
	toks := NewTokenStream(input)
	tok := toks.Next()

	if tok.Kind == KindError {
		require.NoError(t, tok.Value.(error))
	}

	require.Equal(t, kind, tok.Kind)
	require.Equal(t, Position(0), tok.Start)
	require.Equal(t, Position(len(input)), tok.End)
	require.Equal(t, value, tok.Value)

	tok = toks.Next()
	require.Equal(t, KindNil, tok.Kind)
}
