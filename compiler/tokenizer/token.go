package tokenizer

import (
	"fmt"
	"math"
)

type TokenKind uint
type Position uint64

const (
	KindNil     TokenKind = 0
	KindBlank   TokenKind = 1
	KindInteger TokenKind = 2
	KindFloat   TokenKind = 3

	KindError TokenKind = math.MaxInt64
)

var kindNames []string = []string{
	"KindNil",
	"KindBlank",
	"KindInteger",
	"KindFloat",
}

func (k TokenKind) String() string {
	if k == KindError {
		return "KindError"
	} else if int(k) >= len(kindNames) {
		return fmt.Sprintf("%d", k)
	} else {
		return kindNames[k]
	}
}

type Token struct {
	Kind  TokenKind
	Start Position
	End   Position
	Value interface{}
}

func NewToken(kind TokenKind, start Position, end Position, value interface{}) Token {
	return Token{kind, start, end, value}
}
