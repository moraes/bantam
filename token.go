// Copyright 2013 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bantam

import (
	"fmt"
)

const (
	TokenEOF TokenType = iota
	// Variable
	TokenName
	// Operators
	TokenAsterisk    // *
	TokenSlash       // /
	TokenPlus        // +
	TokenMinus       // -
	TokenCaret       // ^
	TokenTilde       // ~
	TokenAssignment  // =
	TokenQuestion    // ?
	TokenExclamation // !
	TokenParenL      // (
	TokenParenR      // )
	TokenColon       // :
	TokenComma       // ,
)

var tokenNames = map[TokenType]string{
	TokenEOF:         "EOF",
	TokenAsterisk:    "*",
	TokenSlash:       "/",
	TokenPlus:        "+",
	TokenMinus:       "-",
	TokenCaret:       "^",
	TokenTilde:       "~",
	TokenAssignment:  "=",
	TokenQuestion:    "?",
	TokenExclamation: "!",
	TokenParenL:      "(",
	TokenParenR:      ")",
	TokenColon:       ":",
	TokenComma:       ",",
}

// TokenType identifies the type of Tokens.
type TokenType int

func (t TokenType) String() string {
	if s, ok := tokenNames[t]; ok {
		return s
	}
	return fmt.Sprintf("<%s>", t)
}

type Token struct {
	Type TokenType
	Text string
}

func (t Token) String() string {
	if s, ok := tokenNames[t.Type]; ok {
		return s
	}
	return t.Text
}
