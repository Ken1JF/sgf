/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/sgf/scanner.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/27/09.
 *	Copyright 2009-2014, all rights reserved.
 *
 *	This package implements reading of SGF game trees.
 *
 */

// Much of this logic is based on the scanner and parser for Go,
// whiich may be found in:
//		${GOROOT}/src/pkg/go/scanner/
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sgf

import (
	"gitHub.com/Ken1JF/ahgo/ah"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// An ErrorHandler may be provided to Scanner.Init. If a syntax error is
// encountered and a handler was installed, the handler is called with a
// ah.Position and an error message. The ah.Position points to the beginning of
// the offending token.
type ErrorHandler func(pos ah.Position, msg string)

// A Scanner holds the scanner's internal state while processing
// a given text.  It can be allocated as part of another data
// structure but must be initialized via Init before use. For
// a sample use, see the implementation of Tokenize.
// TODO: change Scanner to scanner after debugging is complete (Sizeof, etc.)
type Scanner struct {
	// immutable state
	src  []byte       // source
	err  ErrorHandler // error reporting; or nil
	mode uint         // scanning mode

	// scanning state
	pos    ah.Position // previous reading ah.Position (ah.Position before ch)
	offset int         // current reading offset (ah.Position after ch)
	ch     int         // one char look-ahead
	inPV   bool        // between "[" and "]"

	// public state - ok to modify
	ScanErrorCount int // number of errors encountered
}

// Read the next Unicode char into S.ch.
// S.ch < 0 means end-of-file.
func (S *Scanner) next() {
	if S.offset < len(S.src) {
		S.pos.Offset = S.offset
		S.pos.Column++
		r, w := int(S.src[S.offset]), 1
		switch {
		case r == '\n':
			S.pos.Line++
			S.pos.Column = 0
		case r >= 0x80:
			// not ASCII
			re, _ := utf8.DecodeRune(S.src[S.offset:])
			r = int(re)
		}
		S.offset += w
		S.ch = r
	} else {
		S.pos.Offset = len(S.src)
		S.ch = -1 // eof
	}
}

// The mode parameter to the Init function is a set of flags (or 0).
// They control Scanner behavior.
const (
	ScanComments      = 1 << iota // return comments as COMMENT tokens
	AllowIllegalChars             // do not report an error for illegal chars
)

// InitScanner prepares the Scanner S to tokenize the text src. Calls to Scan
// will use the error handler err if they encounter a syntax error and
// err is not nil. Also, for each error encountered, the Scanner field
// ScanErrorCount is incremented by one. The filename parameter is used as
// filename in the ah.Position returned by Scan for each Token. The
// mode parameter determines how comments and illegal characters are
// handled.
func (S *Scanner) InitScanner(filename string, sr []byte, er ErrorHandler, mode uint) {
	// Explicitly initialize all fields since a Scanner may be reused.
	S.src = sr
	S.err = er
	S.mode = mode
	S.pos = ah.Position{Filename: filename, Offset: 0, Line: 1, Column: 0}

	S.ch = ' '
	S.offset = 0
	S.ScanErrorCount = 0
	S.inPV = false
	S.next()
}

func charString(ch int) string {
	var s string
	switch ch {
	case -1:
		return `EOF`
	case '\a':
		s = `\a`
	case '\b':
		s = `\b`
	case '\f':
		s = `\f`
	case '\n':
		s = `\n`
	case '\r':
		s = `\r`
	case '\t':
		s = `\t`
	case '\v':
		s = `\v`
	case '\\':
		s = `\\`
	case '\'':
		s = `\'`
	default:
		s = string(ch)
	}
	return "'" + s + "' (U+" + strconv.FormatInt(int64(ch), 16) + ")"
}

func (S *Scanner) error(pos ah.Position, msg string) {
	if S.err != nil {
		S.err(pos, msg)
	}
	S.ScanErrorCount++
}

func (S *Scanner) expect(ch int) {
	if S.ch != ch {
		S.error(S.pos, "expected "+charString(ch)+", found "+charString(S.ch))
	}
	S.next() // always make progress
}

// Property Identifiers are upper case letters only
func isLetter(ch int) bool { return 'A' <= ch && ch <= 'Z' }

func isDigit(ch int) bool {
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(rune(ch))
}

func (S *Scanner) scanIdentifier() Token {
	//	pos := S.pos.Offset
	for isLetter(S.ch) {
		S.next()
	}
	// no keywords:
	//	return sgfToken.Lookup(S.src[pos:S.pos.Offset])
	return IDENT
}

func digitVal(ch int) int {
	switch {
	case '0' <= ch && ch <= '9':
		return ch - '0'
		// SGF: no hexadecimal values
	}
	return 16 // larger than any legal digit val
}

func (S *Scanner) scanMantissa(base int) {
	for digitVal(S.ch) < base {
		S.next()
	}
}

/* not used?
func (S *Scanner) scanNumber(seen_decimal_point bool) Token {
	tok := INT

	if seen_decimal_point {
		tok = FLOAT
		S.scanMantissa(10)
		goto exit
	}

	// SGF: no hexadecimal or octal

mantissa:
	// decimal int or float
	S.scanMantissa(10)

	if S.ch == '.' {
		// float
		tok = FLOAT
		S.next()
		S.scanMantissa(10)
	}

	// SGF: no exponent

exit:
	return tok
}

end not used? */
func (S *Scanner) scanDigits(base, length int) {
	for length > 0 && digitVal(S.ch) < base {
		S.next()
		length--
	}
	if length > 0 {
		S.error(S.pos, "illegal char escape")
	}
}

func (S *Scanner) scanEscape(quote int) {
	pos := S.pos
	ch := S.ch
	S.next()
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
	// nothing to do
	case '0', '1', '2', '3', '4', '5', '6', '7':
		S.scanDigits(8, 3-1) // 1 char read already
		// SGF: no hexadecimal
	default:
		S.error(pos, "illegal char escape")
	}
}

func (S *Scanner) scanString(pos ah.Position) {
	// '[' already consumed

	for S.ch != ']' {
		ch := S.ch
		S.next()
		if ch < 0 { // SGF: "strings" or property values can cross lines
			S.error(pos, "string not terminated")
			break
		}
		if ch == '\\' {
			S.scanEscape(']')
		}
	}
	S.inPV = false // seen terminating ']'
}

// Scan scans the next Token and returns the Token ah.Position pos,
// the Token tok, and the literal text lit corresponding to the
// Token. The source end is indicated by EOF.
//
// For more tolerant parsing, Scan will return a valid Token if
// possible even if a syntax error was encountered. Thus, even
// if the resulting Token sequence contains no illegal tokens,
// a client may not assume that no error occurred. Instead it
// must check the Scanner's ScanErrorCount or the number of calls
// of the error handler, if there was one installed.
func (S *Scanner) Scan() (pos ah.Position, tok Token, lit []byte) {
	if S.inPV {
		// current Token start
		pos, tok = S.pos, ILLEGAL

		switch ch := S.ch; {

		case ch == ']':
			tok = STRING
			S.inPV = false

		case ch == -1:
			tok = STRING
			S.inPV = false

		default:
			tok = STRING
			S.scanString(pos)
		}
	} else {
		// skip white space
		for S.ch == ' ' || S.ch == '\t' || S.ch == '\n' || S.ch == '\r' {
			S.next()
		}
		// current Token start
		pos, tok = S.pos, ILLEGAL

		// determine Token value
		switch ch := S.ch; {

		case isLetter(ch):
			tok = S.scanIdentifier()

		default:
			S.next() // always make progress
			switch ch {

			case -1:
				tok = EOF

			case ';':
				tok = SEMICOLON

			case '(':
				tok = LPAREN

			case ')':
				tok = RPAREN

			case '[':
				tok = LBRACK // start of property value
				S.inPV = true

			case ']':
				tok = RBRACK

			default:
				if S.mode&AllowIllegalChars == 0 {
					S.error(pos, "illegal character "+charString(ch))
				}
			}
		}
	}
	return pos, tok, S.src[pos.Offset:S.pos.Offset]
}

// Tokenize calls a function f with the Token ah.Position, Token value, and Token
// text for each Token in the source src. The other parameters have the same
// meaning as for the Init function. Tokenize keeps scanning until f returns
// false (usually when the Token value is EOF).
// The result is the number of tokens scanned and the number of errors encountered.
func Tokenize(filename string, src []byte, err ErrorHandler, mode uint, f func(pos ah.Position, tok Token, lit []byte) bool) (nTok int, nErr int) {
	var s Scanner
	s.InitScanner(filename, src, err, mode)
	for f(s.Scan()) {
		nTok++
		// action happens in f
	}
	return nTok, s.ScanErrorCount
}
