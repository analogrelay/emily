package tokenizer

import (
	"strconv"
	"unicode"

	"github.com/pkg/errors"
)

type TokenStream struct {
	win window
}

func NewTokenStream(input string) *TokenStream {
	return &TokenStream{newWindow(input)}
}

// Next returns the next token found in the reader, or an error if an I/O or parsing error occurs.
//
// When the end of the input stream is reached, a token of kind `KindNil`, and an `io.EOF` error are returned.
func (t *TokenStream) Next() Token {
	r := t.win.Next()
	if r == 0 {
		return t.emit(KindNil, nil)
	}

	p := t.win.Peek()

	switch {
	case unicode.IsSpace(r):
		return t.readBlank()
	case r == '0' && (p == 'o' || p == 'x' || p == 'b'):
		return t.readPrefixedInteger()
	case r == '+', r == '-', isDigit(r):
		return t.readDecimal()
	default:
		return t.emit(KindError, errors.Errorf("unexpected '%q'", r))
	}
}

func (t *TokenStream) readBlank() Token {
	t.win.TakeWhile(unicode.IsSpace)
	return t.emit(KindBlank, nil)
}

func (t *TokenStream) readPrefixedInteger() Token {
	switch t.win.Next() {
	case 'o':
		t.win.TakeWhile(isOctalDigit)
	case 'x':
		t.win.TakeWhile(isHexDigit)
	case 'b':
		t.win.TakeWhile(isBinaryDigit)
	}
	return t.emitInteger()
}

func (t *TokenStream) readDecimal() Token {
	t.win.TakeWhile(isDigit)
	return t.emitInteger()
}

func (t *TokenStream) emitInteger() Token {
	// Parse and return
	val, err := strconv.ParseInt(t.win.String(), 0, 64)
	if err != nil {
		return t.emit(KindError, errors.Wrapf(err, "error parsing integer"))
	}
	return t.emit(KindInteger, val)
}

// emit creates a new token using the current start/end value and resets them
func (t *TokenStream) emit(kind TokenKind, value interface{}) Token {
	start, end := t.win.Advance()
	tok := NewToken(kind, Position(start), Position(end), value)
	return tok
}

func isDigit(r rune) bool {
	return ('0' <= r && r <= '9') || r == '_'
}

func isHexDigit(r rune) bool {
	return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F') || r == '_'
}

func isOctalDigit(r rune) bool {
	return ('0' <= r && r <= '7') || r == '_'
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1' || r == '_'
}
