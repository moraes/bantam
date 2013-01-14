// Copyright 2013 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bantam

import (
	"testing"
)

type lexer struct {
	src string
	pos int
}

var stringToToken = map[string]TokenType{
	"*": TokenAsterisk,
	"/": TokenSlash,
	"+": TokenPlus,
	"-": TokenMinus,
	"^": TokenCaret,
	"~": TokenTilde,
	"=": TokenAssignment,
	"?": TokenQuestion,
	"!": TokenExclamation,
	"(": TokenParenL,
	")": TokenParenR,
	":": TokenColon,
	",": TokenComma,
}

// stupendously weak lexer, just for testing.
func (l *lexer) Next() Token {
	for l.pos < len(l.src) {
		s := string(l.src[l.pos])
		l.pos++
		if s == " " {
			continue
		}
		if t, ok := stringToToken[s]; ok {
			return Token{Type: t}
		}
		return Token{Type: TokenName, Text: s}
	}
	return Token{Type: TokenEOF}
}

func TestParser(t *testing.T) {
	type parserTest struct {
		source string
		result string
	}

	tests := []parserTest{
		// Function call.
		{"a()", "a()"},
		{"a(b)", "a(b)"},
		{"a(b, c)", "a(b, c)"},
		{"a(b)(c)", "a(b)(c)"},
		{"a(b) + c(d)", "(a(b) + c(d))"},
		{"a(b ? c : d, e + f)", "a((b ? c : d), (e + f))"},
		// Unary precedence.
		{"~!-+a", "(~(!(-(+a))))"},
		{"a!!!", "(((a!)!)!)"},
		// Unary and binary predecence.
		{"-a * b", "((-a) * b)"},
		{"!a + b", "((!a) + b)"},
		{"~a ^ b", "((~a) ^ b)"},
		{"-a!", "(-(a!))"},
		{"!a!", "(!(a!))"},
		// Binary precedence.
		{"a = b + c * d ^ e - f / g", "(a = ((b + (c * (d ^ e))) - (f / g)))"},
		// Binary associativity.
		{"a = b = c", "(a = (b = c))"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a ^ b ^ c", "(a ^ (b ^ c))"},
		// Conditional operator.
		{"a ? b : c ? d : e", "(a ? b : (c ? d : e))"},
		{"a ? b ? c : d : e", "(a ? (b ? c : d) : e)"},
		{"a + b ? c * d : e / f", "((a + b) ? (c * d) : (e / f))"},
		// Grouping.
		{"a + (b + c) + d", "((a + (b + c)) + d)"},
		{"a ^ (b + c)", "(a ^ (b + c))"},
		{"(!a)!", "((!a)!)"},
	}

	for _, test := range tests {
		l := &lexer{src: test.source}
		s := &Stack{lexer: l}
		p := &Parser{s, PrefixParsers, InfixParsers}
		n, e := p.Parse()
		if e != nil {
			t.Errorf("%q: error parsing: %v", test.source, e)
			continue
		}
		r := n.String()
		if r != test.result {
			t.Errorf("%q: expected %q, got %q", test.source, test.result, r)
			continue
		}
	}

	/*
	l := &lexer{src: "-a * b"}
	s := &Stack{lexer: l}
	for {
		token := s.Pop()
		t.Errorf("%s", token)
		if token.Type == TokenEOF {
			break
		}
	}
	*/
}
