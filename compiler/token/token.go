package token

import (
	"fmt"
	"math"
)

type TokenKind struct{ uint32 }
type Position uint64

var (
	KindNil     = TokenKind{0}
	KindBlank   = TokenKind{1}
	KindInteger = TokenKind{2}
	KindFloat   = TokenKind{3}
	KindIdent   = TokenKind{4}
	KindLParen  = TokenKind{5}
	KindRParen  = TokenKind{6}

	KindError = TokenKind{math.MaxInt32}
)

var kindNames []string = []string{
	"KindNil",
	"KindBlank",
	"KindInteger",
	"KindFloat",
	"KindIdent",
	"KindLParen",
	"KindRParen",
}

func (k TokenKind) String() string {
	if k == KindError {
		return "KindError"
	} else if int(k.uint32) >= len(kindNames) {
		return fmt.Sprintf("%d", k.uint32)
	} else {
		return kindNames[k.uint32]
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
