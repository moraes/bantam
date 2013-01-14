// Copyright 2013 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bantam

import (
	"fmt"
)

// PrefixParser is of the two interfaces used by the Pratt parser.
// A PrefixParser is associated with a token that appears at the beginning of
// an expression. Its Parse() method will be called with the consumed leading
// token, and is responsible for parsing anything that comes after that token.
// This interface is also used for single-token expressions like variables, in
// which case Parse() simply doesn't consume any more tokens.
type PrefixParser interface {
	Parse(*Parser, Token) Node
}

// InfixParser is one of the two parselet interfaces used by the Pratt parser.
// An InfixParser is associated with a token that appears in the middle of the
// expression it parses. Its Parse() method will be called after the left-hand
// side has been parsed, and it in turn is responsible for parsing everything
// that comes after the token. This is also used for postfix expressions, in
// which case it simply doesn't consume any more tokens in its Parse() call.
type InfixParser interface {
	Parse(*Parser, Node, Token) Node
	Precedence() int
}

// ----------------------------------------------------------------------------

// Default prefix parsers.
var PrefixParsers = map[TokenType]PrefixParser{
	TokenName:        NameParser{},
	TokenParenL:      GroupParser{},
	TokenPlus:        UnaryParser(6),
	TokenMinus:       UnaryParser(6),
	TokenTilde:       UnaryParser(6),
	TokenExclamation: UnaryParser(6),
}

// Default infix parsers.
var InfixParsers = map[TokenType]InfixParser{
	TokenAssignment:  AssignParser(1),
	TokenQuestion:    TernaryParser(2),
	TokenPlus:        BinaryParser(3),
	TokenMinus:       BinaryParser(3),
	TokenAsterisk:    BinaryParser(4),
	TokenSlash:       BinaryParser(4),
	TokenCaret:       BinaryRightParser(5),
	TokenExclamation: UnaryPostfixParser(7),
	TokenParenL:      FunctionParser(8),
}

// ----------------------------------------------------------------------------

type Parser struct {
	*Stack
	PrefixParsers map[TokenType]PrefixParser
	InfixParsers  map[TokenType]InfixParser
}

func NewParser(stack *Stack) *Parser {
	return &Parser{
		Stack:         stack,
		PrefixParsers: make(map[TokenType]PrefixParser),
		InfixParsers:  make(map[TokenType]InfixParser),
	}
}

func (p *Parser) Parse(precedence int) Node {
	token := p.Pop()
	prefix, ok := PrefixParsers[token.Type]
	if !ok {
		// TODO: use Parser.errorf()
		panic(fmt.Sprintf("could not parse %s", token))
	}
	left := prefix.Parse(p, token)
	for precedence < p.precedence() {
		token = p.Pop()
		infix, ok := p.InfixParsers[token.Type]
		if !ok {
			// TODO: use Parser.errorf()
			panic(fmt.Sprintf("could not parse %s", token))
		}
		left = infix.Parse(p, left, token)
	}
	return left
}

func (p *Parser) precedence() int {
	if parser, ok := p.InfixParsers[p.Peek(0).Type]; ok {
		return parser.Precedence()
	}
	return 0
}

// ----------------------------------------------------------------------------

// NameParser is a simple parser for a named variable like "abc".
type NameParser struct{}

func (NameParser) Parse(parser *Parser, token Token) Node {
	return NewNameNode(token.Text)
}

// ----------------------------------------------------------------------------

// GroupParser parses parentheses used to group expressions,
// like "a * (b + c)".
type GroupParser struct{}

func (GroupParser) Parse(parser *Parser, token Token) Node {
	n := parser.Parse(0)
	parser.Expect(TokenParenR)
	return n
}

// ----------------------------------------------------------------------------

// UnaryParser parses an unary prefix operator.
type UnaryParser int

func (p UnaryParser) Parse(parser *Parser, token Token) Node {
	right := parser.Parse(int(p))
	return NewUnaryNode(token.Type, right)
}

// ----------------------------------------------------------------------------

// UnaryPostfixParser parses an unary postfix operator.
type UnaryPostfixParser int

func (p UnaryPostfixParser) Parse(parser *Parser, left Node, token Token) Node {
	return NewUnaryPostfixNode(left, token.Type)
}

func (p UnaryPostfixParser) Precedence() int {
	return int(p)
}

// ----------------------------------------------------------------------------

// AssignParser parses assignment expressions like "a = b". The left side of
// an assignment expression must be a simple name like "a", and expressions are
// right-associative. (In other words, "a = b = c" is parsed as "a = (b = c)").
type AssignParser int

func (p AssignParser) Parse(parser *Parser, left Node, token Token) Node {
	l, ok := left.(*NameNode)
	if !ok {
		// TODO: use Parser.errorf()
		panic("the left-hand side of an assignment must be a name")
	}
	right := parser.Parse(int(p) - 1);
	return NewAssignNode(l.Name, right)
}

func (p AssignParser) Precedence() int {
	return int(p)
}

// ----------------------------------------------------------------------------

// FunctionParser parses a function call like "a(b, c, d)".
type FunctionParser int

func (p FunctionParser) Parse(parser *Parser, left Node, token Token) Node {
	// Parse the comma-separated arguments until we hit, ")".
	// There may be no arguments at all.
	args := NewListNode()
	if !parser.Match(TokenParenR) {
		for {
			args.Append(parser.Parse(0))
			if !parser.Match(TokenComma) {
				break
			}
		}
		parser.Expect(TokenParenR)
	}
	return NewFunctionNode(left, args)
}

func (p FunctionParser) Precedence() int {
	return int(p)
}

// ----------------------------------------------------------------------------

// BinaryParser parses a left-associative binary operator.
type BinaryParser int

func (p BinaryParser) Parse(parser *Parser, left Node, token Token) Node {
	right := parser.Parse(int(p))
	return NewBinaryNode(left, token.Type, right)
}

func (p BinaryParser) Precedence() int {
	return int(p)
}

// ----------------------------------------------------------------------------

// BinaryRightParser parses a right-associative binary operator.
type BinaryRightParser int

func (p BinaryRightParser) Parse(parser *Parser, left Node, token Token) Node {
	// To handle right-associative operators like "^", we allow a slightly
	// lower precedence when parsing the right-hand side. This will let a
	// parser with the same precedence appear on the right, which will then
	// take *this* parser's result as its left-hand argument.
	right := parser.Parse(int(p) - 1)
	return NewBinaryNode(left, token.Type, right)
}

func (p BinaryRightParser) Precedence() int {
	return int(p)
}

// ----------------------------------------------------------------------------

// TernaryParser parses a ternary operator.
type TernaryParser int

func (p TernaryParser) Parse(parser *Parser, left Node, token Token) Node {
	node := parser.Parse(0)
	parser.Expect(TokenColon)
	elseNode := parser.Parse(int(p) - 1)
	return NewTernaryNode(left, listNode(node), listNode(elseNode))
}

func (p TernaryParser) Precedence() int {
	return int(p)
}