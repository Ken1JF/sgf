/*
 *  File:		src/github.com/Ken1JF/sgf/token.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/16/09.
 *	Copyright 2009-2014, all rights reserved.
 *
 *	This package defines constants representing lexical tokens
 *	of SGF file and basic operations on tokens (printing, predicates).
 *
 */

// Much of this logic is based on the scanner and parser for Go,
// whiich may be found in:
//		${GOROOT}/src/pkg/go/scanner/
//		${GOROOT}/src/pkg/go/parser/
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sgf

import (
	"strconv"
)

// Token is the set of lexical tokens of the SGF (FF4) language
type Token uint8

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // SZ
	INT    // 19 special case of STRING
	FLOAT  // 6.5 special case of STRING
	STRING // [abc]
	literal_end

	operator_beg
	// Operators and delimiters
	LPAREN    // (
	LBRACK    // [
	PERIOD    // . found in STRING
	RPAREN    // )
	RBRACK    // ]
	SEMICOLON // ;
	COLON     // : found in STRING
	operator_end

	keyword_beg
	// SGF is a strange language! no keywords...
	keyword_end
)

// map of Token values to strings
var tokens = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	LPAREN:    "(",
	LBRACK:    "[",
	PERIOD:    ".",
	RPAREN:    ")",
	RBRACK:    "]",
	SEMICOLON: ";",
	COLON:     ":",
}

// String returns the string corresponding to the Token tok.
// For operators and delimiters, the string is the actual
// Token character sequence (e.g., for the Token ADD, the string is
// "+"). For all other tokens the string corresponds to the Token
// constant name (e.g. for the Token IDENT, the string is "IDENT").
func (tok Token) String() string {
	if str, exists := tokens[tok]; exists {
		return str
	}
	return "Token(" + strconv.Itoa(int(tok)) + ")"
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; returns false otherwise.
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; returns false otherwise.
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }
