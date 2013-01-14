// Copyright 2013 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bantam

import (
	"fmt"
)

// Lexer defines an interface for lexycal scanners.
//
// We don't have a real lexer implementation in this package, only a dummy one
// in the tests.
type Lexer interface {
	Next() Token
}

// NewStack returns a stack for the given lexer.
func NewStack(lexer Lexer) *Stack {
	return &Stack{lexer: lexer}
}

// Stack is a basic LIFO stack for tokens. It allows forwarding and rewinding.
type Stack struct {
	lexer  Lexer
	tokens []Token
	count  int
}

// Push adds one or more tokens back to the stack.
func (s *Stack) Push(t ...Token) {
	s.tokens = append(s.tokens[:s.count], t...)
	s.count += len(t)
}

// Pop consumes and returns a token from the stack.
func (s *Stack) Pop() Token {
	if s.count == 0 {
		return s.lexer.Next()
	}
	s.count--
	return s.tokens[s.count]
}

// Peek returns without consuming a token at the given index.
func (s *Stack) Peek(index int) Token {
	switch {
	case index == 0:
		t := s.Pop()
		s.Push(t)
		return t
	case index > 0:
		if index < s.count {
			return s.tokens[index]
		}
		t := make([]Token, index+1)
		for index >= 0 {
			t[index] = s.Pop()
			index--
		}
		s.Push(t...)
		return t[0]
	}
	panic("Peek received negative index")
}

// Expect consumes a token and panics if it is not of the expected types.
func (s *Stack) Expect(expected ...TokenType) Token {
	t := s.Pop()
	switch len(expected) {
	case 1:
		if t.Type == expected[0] {
			return t
		}
	default:
		for _, e := range expected {
			if t.Type == e {
				return t
			}
		}
	}
	s.Push(t)
	panic(fmt.Sprintf("expected token %s and found %s", expected, t.Type))
}

// Match consumes a token if it is of the expected type, returning true.
// Otherwise the token is not consumed and it returns false.
func (s *Stack) Match(expected TokenType) bool {
	t := s.Pop()
	if t.Type != expected {
		s.Push(t)
		return false
	}
	return true
}
