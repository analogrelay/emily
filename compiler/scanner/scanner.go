package scanner

import (
	"strconv"
	"unicode"

	"github.com/anurse/emily/compiler/token"
	"github.com/pkg/errors"
)

var symbolMappings = map[rune]token.TokenKind{
	'(': token.KindLParen,
	')': token.KindRParen,
}

type Scanner struct {
	win     window
	current token.Token
}

func NewScanner(input string) *Scanner {
	return &Scanner{newWindow(input), token.Token{}}
}

func (t *Scanner) Token() token.Token {
	return t.current
}

func (t *Scanner) Kind() token.TokenKind {
	return t.current.Kind
}

func (t *Scanner) Err() error {
	if t.current.Kind == token.KindError {
		return t.current.Value.(error)
	}
	return nil
}

// Scan processes the next token in the stream, returning true if there is one and false if EOF has been reached
func (t *Scanner) Scan() bool {
	r := t.win.Next()
	if r == 0 {
		t.current = token.Token{}
		return false
	}

	p := t.win.Peek()

	if kind, ok := symbolMappings[r]; ok {
		return t.emit(kind, nil)
	}

	switch {
	case isIdentStart(r):
		return t.readIdent()
	case unicode.IsSpace(r):
		return t.readBlank()
	case r == '0' && (p == 'o' || p == 'x' || p == 'b'):
		return t.readPrefixedInteger()
	case r == '+', r == '-', r == '.', isDigit(r):
		return t.readDecimal(r == '.')
	default:
		return t.emit(token.KindError, errors.Errorf("unexpected %q", r))
	}
}

func (t *Scanner) readBlank() bool {
	t.win.TakeWhile(unicode.IsSpace)
	return t.emit(token.KindBlank, nil)
}

func (t *Scanner) readPrefixedInteger() bool {
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

func (t *Scanner) readDecimal(seenDot bool) bool {
	if !seenDot {
		t.win.TakeWhile(isDigit)
		seenDot = t.win.TakeIf('.')
	}

	if seenDot {
		t.win.TakeWhile(isDigit)
	}

	seenE := t.win.TakeIf('e', 'E')
	if seenE {
		t.win.TakeWhile(isDigit)
	}

	if seenE || seenDot {
		// Parse and return a float
		val, err := strconv.ParseFloat(t.win.String(), 64)
		if err != nil {
			return t.emit(token.KindError, errors.Wrapf(err, "error parsing floating-point number"))
		}
		return t.emit(token.KindFloat, val)
	} else {
		// Parse and return an integer
		return t.emitInteger()
	}
}

func (t *Scanner) readIdent() bool {
	t.win.TakeWhile(isIdentPart)
	return t.emit(token.KindIdent, t.win.String())
}

func (t *Scanner) emitInteger() bool {
	// Parse and return
	val, err := strconv.ParseInt(t.win.String(), 0, 64)
	if err != nil {
		return t.emit(token.KindError, errors.Wrapf(err, "error parsing integer"))
	}
	return t.emit(token.KindInteger, val)
}

// emit creates a new token using the current start/end value and resets them
func (t *Scanner) emit(kind token.TokenKind, value interface{}) bool {
	start, end := t.win.Advance()
	tok := token.NewToken(kind, token.Position(start), token.Position(end), value)
	t.current = tok
	return true
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

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}
