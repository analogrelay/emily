package ast

import "github.com/anurse/emily/compiler/token"

type Node interface {
	Start() token.Position
	End() token.Position
}

type Expr interface {
	Node

	// Requiring this impl means only Exprs can be assigned to Expr.
	exprSentinel()
}

type Stmt interface {
	Node

	// Requiring this impl means only Stmts can be assigned to Stmt.
	stmtSentinel()
}

type File struct {
	List []Stmt
}

type ExprStmt struct {
	X Expr
}

var _ Stmt = &ExprStmt{}

func (s *ExprStmt) Start() token.Position { return s.X.Start() }
func (s *ExprStmt) End() token.Position   { return s.X.End() }
func (s *ExprStmt) stmtSentinel()         {}

type CallExpr struct {
	Function Expr
	Lparen   token.Token
	Args     []Expr
	Rparen   token.Token
}

var _ Expr = &CallExpr{}

func (e *CallExpr) Start() token.Position { return e.Function.Start() }
func (e *CallExpr) End() token.Position   { return e.Rparen.End }
func (e *CallExpr) exprSentinel()         {}

type LiteralExpr struct {
	Value token.Token
}

var _ Expr = &LiteralExpr{}

func (e *LiteralExpr) Start() token.Position { return e.Value.Start }
func (e *LiteralExpr) End() token.Position   { return e.Value.End }
func (e *LiteralExpr) exprSentinel()         {}

type IdentExpr struct {
	Identifier token.Token
}

var _ Expr = &IdentExpr{}

func (e *IdentExpr) Start() token.Position { return e.Identifier.Start }
func (e *IdentExpr) End() token.Position   { return e.Identifier.End }
func (e *IdentExpr) exprSentinel()         {}
