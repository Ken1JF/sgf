/*
 *  File:		src/github.com/Ken1JF/sgf/parser.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/08/09.
 *	Copyright 2009-2014, all rights reserved.
 *
 *	This package implements reading of SGF game trees.
 *
 */

// Much of this logic is based on the scanner and parser for Go,
// whiich may be found in:
//		${GOROOT}/src/pkg/go/parser/
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sgf

import (
	"bytes"
	"fmt"
	"github.com/Ken1JF/ah"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// The mode parameter to the Parse* functions is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// Parser functionality.

const (
	ParseComments uint = 1 << iota // parse comments and add them to tree
	Trace                          // print a trace of parsed productions
	Play                           // Do and Undo Board Moves while reading file
	GoGoD                          // apply GoGoD error checks
)

// The Parser structure holds the Parser's internal state,
// as well as the parse tree slice.

type Parser struct {
	// Errors and Warnings
	errors   ah.ErrorList
	warnings ah.ErrorList

	UnknownProperty Property // most recent unknown property

	// scanner state
	scanner Scanner

	// Tracing/debugging
	mode   uint // parsing mode
	trace  bool // == (mode & Trace != 0)
	play   bool // == (mode & Play != 0)
	indent uint // indentation used for tracing output

	// moveLimit, 0 => no moveLimit
	moveLimit    int
	limitReached bool

	// Next token
	pos ah.Position // token ah.Position
	tok Token       // one token look-ahead
	lit []byte      // token literal

	// Parse tree
	GameTree
}

// Add an error
func (p *Parser) Error(pos ah.Position, msg string) {
	p.errors.Add(pos, msg)
}

// ReportException prints values of Properties that cannot be understood
func (p *Parser) ReportException(idx PropertyDefIdx, str []byte, err string) {
	if idx != TM_idx {
		fmt.Printf("BAD Property Value: %s:%d:%d: %s[%s] %s\n", p.pos.Filename, p.pos.Line, p.pos.Column, string(GetProperty(idx).ID), str, err)
	}
}

// addProp maintains a variable sized array of properties.
// properties are maintained as a circular linked list.
func (p *Parser) addProp(n TreeNodeIdx, pv PropertyValue) {
	err := p.AddAProp(n, pv)
	if len(err) != 0 {
		p.errors.Add(p.pos, err[0].Msg)
	}
}

func (p *Parser) addNode(par TreeNodeIdx, ty TreeNodeType) TreeNodeIdx {
	newIdx, err := p.AddChild(par, ty, p.Board.GetMovDepth())
	if len(err) != 0 {
		p.errors.Add(p.pos, "adding node "+err[0].Msg)
		// TODO: need to exit, cannot continue without updating p.treeNodes etc.
	}
	return newIdx
}

// scannerMode returns the scanner mode bits given the Parser's mode bits.
func scannerMode(mode uint) uint {
	var m uint
	if mode&ParseComments != 0 {
		m |= ScanComments
	}
	return m
}

// Advance to the next token.
func (p *Parser) next0() {
	// Because of one-token look-ahead, print the previous token
	// when tracing as it provides a more readable output. The
	// very first token (p.pos.Line == 0) is not initialized (it
	// is token.ILLEGAL), so don't print it .
	if p.trace && p.pos.Line > 0 {
		s := p.tok.String()
		switch {
		case p.tok.IsLiteral():
			p.printTrace(s, string(p.lit))
		case p.tok.IsOperator():
			p.printTrace("\"" + s + "\"")
		default:
			p.printTrace(s)
		}
	}

	p.pos, p.tok, p.lit = p.scanner.Scan()
}

// Advance to the next token.
func (p *Parser) next() {
	p.next0()
}

// initParser must be called before a Parser can be used
func (p *Parser) initParser(filename string, src []byte, mode uint, fileLimit int) {

	eh := func(pos ah.Position, msg string) { p.errors.Add(pos, msg) }

	p.scanner.InitScanner(filename, src, eh, scannerMode(mode))
	p.mode = mode
	p.moveLimit = fileLimit
	// for convenience (used frequently)
	p.trace = (mode&Trace != 0) || ah.GetAHTrace()
	p.play = (mode&Play != 0)

	p.next()
	p.GameTree.initGameTree()
}

func (p *Parser) errorExpected(pos ah.Position, msg string) {
	if p.limitReached != true {
		msg = "expected " + msg
		if pos.Offset == p.pos.Offset {
			// the error happened at the current ah.Position;
			// make the error message more specific
			msg += ", found '" + p.tok.String() + "'"
			if p.tok.IsLiteral() {
				msg += " " + string(p.lit)
			}
		}
		p.errors.Add(pos, msg)
	}
}

func (p *Parser) warningPropertyType(pos ah.Position, m error) {
	msg := "PropertyType warning " + m.Error()
	if pos.Offset == p.pos.Offset {
		// the warning happened at the current ah.Position;
		// make the warning message more specific
		msg += ", at '" + p.tok.String() + "'"
		if p.tok.IsLiteral() {
			msg += " " + string(p.lit)
		}
	}
	p.warnings.Add(pos, msg)
}

func (p *Parser) expect(tok Token) ah.Position {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress in any case
	return pos
}

func (p *Parser) expect2(tok1, tok2 Token) ah.Position {
	pos := p.pos
	if p.tok != tok1 && p.tok != tok2 {
		p.errorExpected(pos, "'"+tok1.String()+"' OR '"+tok2.String()+"'")
	}
	p.next() // make progress in any case
	return pos
}

// ----------------------------------------------------------------------------
// Parsing support

func (p *Parser) printTrace(a ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . " +
		". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = uint(len(dots))
	fmt.Printf("%5d:%3d: ", p.pos.Line, p.pos.Column)
	i := 2 * p.indent
	for ; i > n; i -= n {
		fmt.Print(dots)
	}
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

// Usage pattern: defer un(trace(p, "..."));
func un(p *Parser) {
	p.indent--
	p.printTrace(")")
}

// ----------------------------------------------------------------------------
// Source files

func (p *Parser) parsePropValue(val PropValueType) (pv PropertyValue) {
	if p.trace {
		defer un(trace(p, "parsePropValue"))
	}
	pv.ValType = val

	if len(p.lit) > 0 {
		pv.StrValue = p.lit
	}

	switch val {
	case Unknown, SimpText, Text, CompSimpText_simpText:
		pv.ValType = val
		if p.tok == STRING {
			p.next()
			p.expect(RBRACK)
		} else if p.tok == RBRACK { // empty string
			pv.ValType = None
			p.next()
		} else {
			p.errorExpected(p.pos, ValueNames[val])
			p.next()
			p.expect(RBRACK)
		}

	case None_OR_compNum_simpText:
		if p.tok == RBRACK {
			pv.ValType = None
			p.next()
		} else if p.tok == STRING {
			pv.ValType = CompNum_simpText
			p.next()
			p.expect(RBRACK)
		}

	case None:
		if p.tok == RBRACK {
			pv.ValType = None
			p.next()
		} else {
			p.errorExpected(p.pos, ValueNames[val])
			p.next()
			p.expect(RBRACK)
		}

	case CompressedListOfPoint:
		if (len(pv.StrValue) % 2) != 0 {
			p.errors.Add(p.pos, "CompressedListOfPoint not even:"+string(pv.StrValue))
		}
		p.next()
		p.expect(RBRACK)

	case ListOfCompPoint_simpTest:
		pv.ValType = ListOfCompPoint_simpTest
		p.next()
		p.expect(RBRACK)
		for p.tok == LBRACK { // more than one, concatenate together, with separator brackets
			var sep = make([]byte, 2)
			sep[0] = ']'
			sep[1] = '['
			p.next()
			pv.StrValue = append(pv.StrValue, sep...)
			pv.StrValue = append(pv.StrValue, p.lit...)
			p.next()
			p.expect(RBRACK)
		}

	case ListOfCompPoint_Point:
		pv.ValType = ListOfCompPoint_Point
		p.next()
		p.expect(RBRACK)
		for p.tok == LBRACK { // more than one, concatenate together. User can process
			p.next()
			pv.StrValue = append(pv.StrValue, p.lit...)
			p.next()
			p.expect(RBRACK)
		}

	case EListOfPoint:
		if p.tok == RBRACK {
			pv.ValType = None
			p.next()
		} else {
			pv.ValType = Point
			_, err := SGFPoint(pv.StrValue)
			if len(err) != 0 {
				p.errors.Add(p.pos, "Bad "+ValueNames[val]+": "+err[0].Msg+": from "+string(pv.StrValue))
			}
			p.next()
			p.expect(RBRACK)
			for p.tok == LBRACK { // more than one, make it a CompressedListOfPoint
				pv.ValType = ListOfPoint
				p.next()
				pv.StrValue = append(pv.StrValue, p.lit...)
				if (len(pv.StrValue) % 2) != 0 {
					p.errors.Add(p.pos, "EListOfPoint not even:"+string(pv.StrValue))
				}
				p.next()
				p.expect(RBRACK)
			}
		}

	case Point, Move, Stone:
		_, err := SGFPoint(pv.StrValue)
		if len(err) != 0 {
			p.errors.Add(p.pos, "Bad "+ValueNames[val]+": "+err[0].Msg+": from "+string(pv.StrValue))
		}
		p.next()
		p.expect(RBRACK)

	case ListOfPoint, ListOfStone:
		pv.ValType = Point
		_, err := SGFPoint(pv.StrValue)
		if len(err) != 0 {
			p.errors.Add(p.pos, "Bad "+ValueNames[val]+": "+err[0].Msg+": from "+string(pv.StrValue))
		}
		p.next()
		p.expect(RBRACK)
		for p.tok == LBRACK { // more than one, make it a ListOfPoint
			pv.ValType = ListOfPoint
			p.next()
			_, err := SGFPoint(p.lit)
			if len(err) != 0 {
				p.errors.Add(p.pos, "Bad "+ValueNames[val]+": "+err[0].Msg+": from "+string(pv.StrValue))
			}
			pv.StrValue = append(pv.StrValue, p.lit...)
			if (len(pv.StrValue) % 2) != 0 {
				p.errors.Add(p.pos, "CompressedListOfPoint not even:"+string(pv.StrValue))
			}
			p.next()
			p.expect(RBRACK)
		}

	case Num_OR_compNum_num: // only used for SZ, don't have any rectangular board tests, for now
		if p.tok == STRING {
			var err error
			idx_colon := strings.Index(string(p.lit), ":")
			if idx_colon > 0 {
				_, err = strconv.Atoi(string(p.lit[0:idx_colon]))
				if err != nil {
					p.errors.Add(p.pos, "Error in composite (column): "+err.Error()+string(p.lit[0:idx_colon]))
				}
				_, err = strconv.Atoi(string(p.lit[idx_colon+1:]))
				if err != nil {
					p.errors.Add(p.pos, "Error in composite (row): "+err.Error()+string(p.lit[idx_colon+1:]))
				}
			} else {
				_, err = strconv.Atoi(string(p.lit))
				if err != nil {
					p.errors.Add(p.pos, "Error in number: "+err.Error()+string(p.lit))
				}
			}
		} else {
			p.errorExpected(p.pos, ValueNames[val])
		}
		p.next()
		p.expect(RBRACK)

	case Num_0_3, Num_1_4, Num_1_5_or_7_16, Num:
		if p.tok == STRING {
			i, _ := strconv.Atoi(string(p.lit))
			if val != Num {
				switch val {
				case Num_0_3:
					if i < 0 || i > 3 {
						p.errors.Add(p.pos, "not in range 0-3: "+string(p.lit))
					}
				case Num_1_4:
					if i < 1 || i > 4 {
						p.errors.Add(p.pos, "not in range 1-4: "+string(p.lit))
					}
				case Num_1_5_or_7_16:
					if i < 1 || i > 16 || i == 6 {
						p.errors.Add(p.pos, "not in range 1-5 or 7-15: "+string(p.lit))
					}
				}
			}
		} else {
			p.errorExpected(p.pos, ValueNames[val])
		}
		p.next()
		p.expect(RBRACK)

	case Real:
		if p.tok == STRING { // pass as a string
			p.next()
			p.expect(RBRACK)
		} else if p.tok == RBRACK { // empty string
			pv.ValType = None
			p.next()
		} else {
			p.errorExpected(p.pos, ValueNames[val])
			p.next()
			p.expect(RBRACK)
		}

	case Double:
		p.errors.Add(p.pos, "Not Implemented: "+ValueNames[val])
		p.next()
		p.expect(RBRACK)

	case Color:
		p.errors.Add(p.pos, "Not Implemented: "+ValueNames[val])
		p.next()
		p.expect(RBRACK)

		// not possible?		default:
	}
	if p.trace {
		fmt.Printf("ValType: %s String: %s\n", ValueNames[pv.ValType], string(pv.StrValue))
	}
	return pv
}

var ID_Counts ID_CountArray // TODO: what to do with these globals?
var Unkn_Count int          // TODO: what to do with these globals?

var HA_map map[string]int = make(map[string]int, 100) // TODO: what to do with these globals?
var OH_map map[string]int = make(map[string]int, 100) // TODO: what to do with these globals?

// Break RE into value and comment:
var RE_map map[string]int = make(map[string]int, 100) // TODO: what to do with these globals?
var RC_map map[string]int = make(map[string]int, 100) // TODO: what to do with these globals?

var RU_map map[string]int = make(map[string]int, 100)     // TODO: what to do with these globals?
var BWRank_map map[string]int = make(map[string]int, 100) // TODO: what to do with these globals?

type PlayerInfo struct {
	NGames    int
	FirstGame string
	FirstRank string
	LastGame  string
	LastRank  string
}

var BWPlayer_map map[string]PlayerInfo = make(map[string]PlayerInfo, 100) // TODO: what to do with these globals?

func GameName(fileName string) string {
	var name []byte
	name = []byte(fileName)
	i := bytes.LastIndex(name, []byte("/"))
	name = name[i+1:]
	i = bytes.Index(name, []byte(".sgf"))
	if i > 0 {
		name = name[0:i]
	}
	return string(name)
}

func check_OH(strVal []byte) (err string) {
	s := string(strVal)
	if s == "1" {
		err = "check value: OH[1]"
	}
	/*
		if s == "BWB" {
			err = "check value: BWB (must be B?)"
		}
		if s == "BBW" {
			err = "check value: BBW (must be B?)"
		}
		if s == "B2B" {
			err = "check value: B2B (must be B?)"
		}
		if s == "2B2" {
			err = "check value: 2B2 (must be 2?)"
		}
		if s == "233" {
			err = "check value: 233 (must be 3?)"
		}
	*/
	return err
}

func fix_OH(s []byte) []byte {
	new_s := make([]byte, len(s))
	j := 0
	for _, b := range s {
		if (b != '-') && (b != ' ') { // skip '-' and ' ' characters
			if b == '{' { // change '{' to '('
				b = '('
			} else if b == '}' { // change '}' to ')'
				b = ')'
			}
			new_s[j] = b
			j += 1
		}
	}
	new_s = new_s[0:j]
	return new_s
}

func TakeOutNum(RE_com []byte) (bas []byte, n int, sep byte, both bool) {
	bas = make([]byte, len(RE_com))
	for i, b := range RE_com {
		bas[i] = b
	}
	ln := len(bas)
	if ln > 0 {
		if bas[0] == '(' {
			sep = '('
			if bas[ln-1] == ')' {
				bas = bas[1 : ln-1]
				both = true
			} else {
				bas = bas[1:ln]
			}
		} else if bas[0] == '{' {
			sep = '{'
			if bas[ln-1] == '}' {
				bas = bas[1 : ln-1]
				both = true
			} else {
				bas = bas[1:ln]
			}
		}
		j := 0
		for _, b := range bas {
			if unicode.IsDigit(rune(b)) {
				if n == 0 {
					bas[j] = '%'
					j += 1
				}
				n = 10*n + int(b-'0')
			} else {
				bas[j] = b
				j += 1
			}
		}
		bas = bas[0:j]
	}
	return bas, n, sep, both
}

func check_Rank(strVal []byte) (err string) {
	s := string(strVal)
	if s == "9" {
		err = "check rank:" + s
	}
	if s == "9di" {
		err = "check rank:" + s
	}
	if s == "9f" {
		err = "check rank:" + s
	}
	if s == "9g" {
		err = "check rank:" + s
	}
	if s == "5e" {
		err = "check rank:" + s
	}
	if s == "98d" {
		err = "check rank:" + s
	}
	if s == "2D" {
		err = "check rank:" + s
	}
	if s == "7d ams" {
		err = "check rank:" + s
	}
	if s == "0d" {
		err = "check rank:" + s
	}
	if s == "2.5" {
		err = "check rank:" + s
	}
	if s == "NR" {
		err = "check rank:" + s
	}
	if s == "a5" {
		err = "check rank:" + s
	}
	if s == "a6" {
		err = "check rank:" + s
	}
	if s == "Holder" {
		err = "check rank:" + s
	}
	if s == " 7d" {
		err = "check rank:" + s
	}
	if s == "Wangwi" {
		err = "check rank:" + s
	}
	if s == "7d Ama" {
		err = "check rank:" + s
	}
	return err
}

func check_Name(strVal []byte) (err string) {
	s := string(strVal)
	if s == "artu" {
		err = "check name:" + s
	}
	if s == "jy23" {
		err = "check name:" + s
	}
	if s == "thug" {
		err = "check name:" + s
	}
	if s == "Yoshida" {
		err = "check name:" + s
	}
	if s == "Yi" {
		err = "check name:" + s
	}
	if s == "World" {
		err = "check name:" + s
	}
	if s == "Turtles" {
		err = "check name:" + s
	}
	if s == "Two shodans" {
		err = "check name:" + s
	}
	if s == "Storks" {
		err = "check name:" + s
	}
	if s == "Seo" {
		err = "check name:" + s
	}
	if s == "Old Lady of Black-horse Mountain" {
		err = "check name:" + s
	}
	if s == "NHK viewers, by internet poll" {
		err = "check name:" + s
	}
	if s == "NO1NO1" {
		err = "check name:" + s
	}
	if s == "MoGo Titan" {
		err = "check name:" + s
	}
	if s == "Miss Y." {
		err = "check name:" + s
	}
	if s == "Maeda" {
		err = "check name:" + s
	}
	if s == "Li Ang" {
		err = "check name:" + s
	}
	if s == "Kuwata" {
		err = "check name:" + s
	}
	if s == "Kuboniwa" {
		err = "check name:" + s
	}
	if s == "KCC Igo program" {
		err = "check name:" + s
	}
	if s == "Harada" {
		err = "check name:" + s
	}
	if s == "Goemate" {
		err = "check name:" + s
	}
	if s == "Go Professional III" {
		err = "check name:" + s
	}
	if s == "GO4++ program" {
		err = "check name:" + s
	}
	if s == "Fukuhara" {
		err = "check name:" + s
	}
	if s == "Fuji Hiroshi" {
		err = "check name:" + s
	}
	if s == "Fairy" {
		err = "check name:" + s
	}
	if s == "Anon." {
		err = "check name:" + s
	}
	if s == "An Immortal" {
		err = "check name:" + s
	}
	if s == "A Go Review subscriber" {
		err = "check name:" + s
	}
	if s == "99P" {
		err = "check name:" + s
	}
	if s == "Another Immortal" {
		err = "check name:" + s
	}
	return err
}

func check_RE(strVal []byte) (err string) {
	s := string(strVal)
	/* First round of fixes:
	if s == "B+2,5" {
		err = "check value: B+2,5 (must be B+2.5)"
	}
	if s == "W+0,5" {
		err = "check value: W+0,5 (must be W+0.5)"
	}
	if s == "W+40 zi" {
		err = "check value: W+40 zi (must be BW+40)"
	}
	if s == "W+34 zi" {
		err = "check value: W+34 zi (must be W+34)"
	}
	if s == "B+05." {
		err = "check value: B+05. (must be B+0.5)"
	}
	if s == "W4.5" {
		err = "check value: W4.5 (must be W+4.5)"
	}
	if s == "w+R" {
		err = "check value: w+R (must be W+R)"
	}
	if s == "B8.5" {
		err = "check value: B8.5 (must be B+8.5)"
	}
	if s == "B++" {
		err = "check value: B++ (must be B+)"
	}
	*/
	if s == "UF" {
		err = "check value: UF"
	}
	if s == "W+jigo" {
		err = "check value: W+jigo"
	}
	if s == "B+).5" {
		err = "check value: B+).5"
	}
	if s == "B+8.8" {
		err = "check value: B+8.8"
	}
	// This is O.K. (in 1988/1988-08-23c.sgf:)
	// There is a seki and rules are Ing Goe
	//
	//	if s == "W+6.6" {
	//		err = "check value: W+6.6"
	//	}
	if s == "W+3.35" {
		err = "check value: W+3.35"
	}
	// This is O.K. (in 1988/1988-08-23b.sgf:)
	// There is a seki and rules are Ing Goe
	//
	//	if s == "B+1.83" {
	//		err = "check value: B+1.83"
	//	}
	if s == "B+11.4" {
		err = "check value: B+11.4"
	}
	if s == "B+2.4" {
		err = "check value: B+2.4"
	}
	if s == "B+1.55" {
		err = "check value: B+1.55"
	}
	if s == "B+35." {
		err = "check value: B+35."
	}
	if s == "B+.75" {
		err = "check value: B+.75"
	}
	if s == "B+.5" {
		err = "check value: B+.5"
	}
	return err
}

func SplitRE(str []byte) (val []byte, com []byte) {
	iLP := bytes.IndexByte(str, '(')
	iLB := bytes.IndexByte(str, '{')
	if iLP > 0 {
		if iLB > 0 {
			if iLB < iLP {
				iLP = iLB
			}
		}
		ix := iLP - 1
		for str[ix] == ' ' {
			if ix == 0 {
				break
			}
			ix -= 1
		}
		val = str[0 : ix+1]
		com = str[iLP:]
	} else if iLB > 0 {
		ix := iLB - 1
		for str[ix] == ' ' {
			if ix == 0 {
				break
			}
			ix -= 1
		}
		val = str[0 : ix+1]
		com = str[iLB:]
	} else {
		val = str
	}
	if string(val) == "Left unfinished" {
		val = []byte("Void")
		if len(com) == 0 {
			com = []byte("(Left unfinished)")
		} else {
			var newCom []byte
			if com[0] == '(' {
				newCom = append([]byte("(Left unfinished. "), com[1:]...)
			} else if com[0] == '{' {
				newCom = append([]byte("{Left unfinished. "), com[1:]...)
			} else {
				newCom = append([]byte("(Left unfinished) "), com...)
			}
			com = newCom
		}
	}
	return val, com
}

func (p *Parser) processProperty(pv PropertyValue, nodd TreeNodeIdx) (ret TreeNodeIdx) {
	if p.trace {
		defer un(trace(p, "processProperty"))
	}
	ret = nodd
	// TODO: get rid of idx when all the "Not implemented:" messages are gone...
	idx := pv.PropType
	//	fmt.Println("Node", nodd, "idx", idx, "Description", GetProperty(idx).Description);
	//	os.Exit(998)
	if idx >= 0 {
		ID_Counts[idx] += 1
	} else {
		Unkn_Count += 1
	}
	switch idx {

	case AB_idx:
		if pv.ValType == Point {
			// Add point to Board
			mov, err := SGFPoint(pv.StrValue)
			if len(err) != 0 {
				p.errors.Add(p.pos, "Bad Point for AB: "+err.Error()+" B["+string(pv.StrValue)+"]")
			} else {
				err = p.DoAB(mov, p.play)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Error from DoAB: "+err.Error()+" B["+string(pv.StrValue)+"]")
				}
			}
			// Record the property:
			p.addProp(ret, pv)
		} else {
			str := pv.StrValue
			for len(str) > 0 { // compressed list of Points
				// Make a new pv
				npv := pv
				npv.StrValue = str[0:2]
				// Add point to Board
				mov, err := SGFPoint(npv.StrValue)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Bad Point List element for AB: "+err.Error()+" B["+string(pv.StrValue)+"]")
				} else {
					err = p.DoAB(mov, p.play)
					if len(err) != 0 {
						p.errors.Add(p.pos, "Error from DoAB:"+err.Error()+"in List element, B["+string(pv.StrValue)+"]")
					}
				}
				str = str[2:]
			}
			// Record the property: (only once)
			p.addProp(ret, pv)
		}

	case AE_idx:
		if pv.ValType == Point {
			// Add point to Board
			mov, err := SGFPoint(pv.StrValue)
			if len(err) != 0 {
				p.errors.Add(p.pos, "Bad Point for AE: "+err.Error()+": from "+string(pv.StrValue))
			} else {
				err = p.DoAE(mov, p.play)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Error from DoAE: "+err.Error()+": caused by "+string(pv.StrValue))
				}
			}
			// Record the property:
			p.addProp(ret, pv)
		} else {
			str := pv.StrValue
			for len(str) > 0 { // compressed list of Points
				// Make a new pv
				npv := pv
				npv.StrValue = str[0:2]
				// Add point to Board
				mov, err := SGFPoint(npv.StrValue)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Bad Point List element for AE: "+err.Error()+": from "+string(npv.StrValue))
				} else {
					err = p.DoAE(mov, p.play)
					if len(err) != 0 {
						p.errors.Add(p.pos, "Error from DoAE:"+err.Error()+" in List element: "+string(npv.StrValue))
					}
				}
				str = str[2:]
			}
			// Record the property: (only once)
			p.addProp(ret, pv)
		}

	case AN_idx:
		// set the board AN:
		p.SetAN(pv.StrValue)
		// Record the property:
		p.addProp(ret, pv)

	case AP_idx:
		// set the board AP:
		p.SetAP(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case AR_idx:
		// set the board AR:
		p.DoAR(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case AS_idx:
		// record the property:
		p.addProp(ret, pv)

	case AW_idx:
		if pv.ValType == Point {
			// Add point to Board
			mov, err := SGFPoint(pv.StrValue)
			if len(err) != 0 {
				p.errors.Add(p.pos, "Bad Point for AW: "+err.Error()+": from "+string(pv.StrValue))
			} else {
				err = p.DoAW(mov, p.play)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Error from DoAW: "+err.Error()+": caused by "+string(pv.StrValue))
				}
			}
			// Record the property:
			p.addProp(ret, pv)
		} else {
			str := pv.StrValue
			for len(str) > 0 { // compressed list of Points
				// Make a new pv
				npv := pv
				npv.StrValue = str[0:2]
				// Add point to Board
				mov, err := SGFPoint(npv.StrValue)
				if len(err) != 0 {
					p.errors.Add(p.pos, "Bad Point List element for AW: "+err.Error()+": from "+string(npv.StrValue))
				} else {
					err = p.DoAW(mov, p.play)
					if len(err) != 0 {
						p.errors.Add(p.pos, "Error from DoAW: "+err.Error()+" in List element: "+string(npv.StrValue))
					}
				}
				str = str[2:]
			}
			// Record the property: (once)
			p.addProp(ret, pv)
		}

	case B_idx:
		p.treeNodes[ret].TNodType = BlackMoveNode
		mov, err := SGFPoint(pv.StrValue)
		if len(err) != 0 {
			p.errors.Add(p.pos, err.Error()+" in SGFPoint, B["+string(pv.StrValue)+"]")
		} else {
			p.treeNodes[ret].propListOrNodeLoc = PropIdx(mov)
			movN, err := p.DoB(mov, p.play)
			if len(err) != 0 {
				p.warnings.Add(p.pos, err.Error()+" B["+string(pv.StrValue)+"]")
			}
			if (p.moveLimit > 0) && (movN >= p.moveLimit) {
				p.limitReached = true
				if p.trace {
					p.printTrace("limitReached set to true in B")

				}
			}
		}

	case BL_idx:
		// record the property:
		p.addProp(ret, pv)

	case BM_idx:
		// record the property:
		p.addProp(ret, pv)

	case BR_idx:
		// count the BR values:
		idx := string(pv.StrValue)
		n, _ := BWRank_map[idx]
		BWRank_map[idx] = n + 1
		// check the rank
		errStr := check_Rank(pv.StrValue)
		if errStr != "" {
			p.ReportException(BR_idx, pv.StrValue, errStr)
		}
		// set the board BR:
		p.SetBR(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case BT_idx:
		// set the board BT:
		p.SetBT(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case C_idx:
		// record the property:
		// TODO: need to process escape characters?
		if (p.mode & ParseComments) != 0 {
			p.addProp(ret, pv)
		}

	case CA_idx:
		// TODO: implement support for different char sets
		// record the property:
		p.addProp(ret, pv)

	case CP_idx:
		// record the property:
		p.addProp(ret, pv)

	case CR_idx:
		// record the property:
		p.addProp(ret, pv)

	case DD_idx:
		// record the property:
		p.addProp(ret, pv)

	case DM_idx:
		// record the property:
		p.addProp(ret, pv)

	case DO_idx:
		// record the property:
		p.addProp(ret, pv)

	case DT_idx:
		// set the board DT:
		p.SetDT(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case EV_idx:
		// set the board EV:
		p.SetEV(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case FF_idx:
		// set the board FF:
		p.SetFF(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case FG_idx:
		// record the property:
		p.addProp(ret, pv)

	case GB_idx:
		// record the property:
		p.addProp(ret, pv)

	case GC_idx:
		// set the board GC:
		p.SetGC(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case GM_idx:
		// Check the GM:
		i, _ := strconv.Atoi(string(pv.StrValue))
		if i != 1 {
			p.errors.Add(p.pos, "GM not 1: "+string(pv.StrValue))
		}
		// record the property:
		p.addProp(ret, pv)

	case GN_idx:
		// record the property:
		p.addProp(ret, pv)

	case GW_idx:
		// record the property:
		p.addProp(ret, pv)

	case HA_idx:
		// count the HA values:
		idx := string(pv.StrValue)
		n, _ := HA_map[idx]
		HA_map[idx] = n + 1
		// set the board HA:
		i, err := strconv.Atoi(idx)
		if err != nil {
			p.ReportException(HA_idx, pv.StrValue, err.Error())
		}
		p.SetHA(i)
		// record the property:
		p.addProp(ret, pv)

	case HO_idx:
		// record the property:
		p.addProp(ret, pv)

	case IP_idx:
		// record the property:
		p.addProp(ret, pv)

	case IT_idx:
		// record the property:
		p.addProp(ret, pv)

	case IY_idx:
		// record the property:
		p.addProp(ret, pv)

	case KM_idx:
		// set the board KM:
		f, err := strconv.ParseFloat(string(pv.StrValue), 64)
		if err == nil {
			p.SetKM(float32(f), true)
		} else {
			if (len(pv.StrValue) == 1) && (pv.StrValue[0] == '?') {
				p.SetKM(0.0, false)
			} else {
				p.ReportException(KM_idx, pv.StrValue, err.Error())
			}
		}
		// record the property
		p.addProp(ret, pv)

	case KO_idx:
		// TODO: process?
		// record the property:
		p.addProp(ret, pv)

	case LB_idx:
		// record the property:
		p.addProp(ret, pv)

	case LN_idx:
		// set the board LN:
		p.DoLN(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case MA_idx:
		// record the property:
		p.addProp(ret, pv)

	case MN_idx:
		// record the property:
		p.addProp(ret, pv)

	case N_idx:
		// record the property
		// TODO: allow lookup of nodes by name?
		p.addProp(ret, pv)

	case OB_idx:
		// record the property:
		p.addProp(ret, pv)

	case OH_idx:
		// count the OH values:
		strVal := fix_OH(pv.StrValue)
		//use this to validate OH values:
		errStr := check_OH(strVal)
		if errStr != "" {
			p.ReportException(OH_idx, pv.StrValue, errStr)
		}
		idx := string(strVal)
		n, _ := OH_map[idx]
		OH_map[idx] = n + 1
		// set the board OH:
		p.SetOH(strVal)
		// record the property:
		p.addProp(ret, pv)

	case ON_idx:
		// record the property:
		p.addProp(ret, pv)

	case OT_idx:
		// record the property:
		p.addProp(ret, pv)

	case OW_idx:
		// record the property:
		p.addProp(ret, pv)

	case PB_idx:
		// count the Player values:
		idx := string(pv.StrValue)
		n, _ := BWPlayer_map[idx]
		n.NGames += 1
		if n.FirstGame == "" {
			n.FirstGame = GameName(p.pos.Filename)
			p.GameTree.setFirstBRank = true
		}
		n.LastGame = GameName(p.pos.Filename)
		BWPlayer_map[idx] = n
		// check the name
		errStr := check_Name(pv.StrValue)
		if errStr != "" {
			p.ReportException(PB_idx, []byte(""), errStr)
		}
		// set the board PB:
		p.SetPB(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case PC_idx:
		// set the board PC:
		p.SetPC(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case PL_idx:
		// record the property:
		p.addProp(ret, pv)

	case PM_idx:
		// record the property:
		p.addProp(ret, pv)

	case PW_idx:
		// count the Player values:
		idx := string(pv.StrValue)
		n, _ := BWPlayer_map[idx]
		n.NGames += 1
		if n.FirstGame == "" {
			n.FirstGame = GameName(p.pos.Filename)
			p.GameTree.setFirstWRank = true
		}
		n.LastGame = GameName(p.pos.Filename)
		BWPlayer_map[idx] = n
		// check the name
		errStr := check_Name(pv.StrValue)
		if errStr != "" {
			p.ReportException(PW_idx, []byte(""), errStr)
		}
		// set the board PW:
		p.SetPW(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case RE_idx:
		// separate RE and RC (Result Comment)
		RE_val, RE_com := SplitRE(pv.StrValue)
		// count the RE values:
		idx := string(RE_val)
		n, _ := RE_map[idx]
		RE_map[idx] = n + 1
		errStr := check_RE(RE_val)
		if errStr != "" {
			p.ReportException(RE_idx, RE_val, errStr)
		}
		// count the RE comments:
		RE_bas, n, ch, both := TakeOutNum(RE_com)
		if len(RE_bas) > 0 {
			idx2 := string(RE_bas)
			n, _ := RC_map[idx2]
			RC_map[idx2] = n + 1
		}
		// set the board RE:
		p.SetRE(RE_val, RE_bas, n, ch, both)
		// record the property:
		p.addProp(ret, pv)

	case RO_idx:
		// set the board RO:
		p.SetRO(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case RU_idx:
		// count the RU values:
		idx := string(pv.StrValue)
		n, _ := RU_map[idx]
		RU_map[idx] = n + 1
		// set the board RU:
		p.SetRU(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case S_idx: // Represent a S property as 1 node with Black move + property list of other moves:
		p.treeNodes[ret].TNodType = SequenceNode
		mov, err := SGFPoint(pv.StrValue)
		p.treeNodes[ret].propListOrNodeLoc = PropIdx(mov)
		if len(err) != 0 {
			p.errors.Add(p.pos, err.Error()+": from "+string(pv.StrValue))
		}
		movColor := ah.Black
		for len(pv.StrValue) > 2 {
			pv.StrValue = pv.StrValue[2:]
			movColor = ah.OppositeColor(movColor)
			ret = p.addNode(ret, SequenceNode)
			mov, err = SGFPoint(pv.StrValue)
			p.treeNodes[ret].propListOrNodeLoc = PropIdx(mov)
			if len(err) != 0 {
				p.errors.Add(p.pos, err.Error()+": from "+string(pv.StrValue))
			}
		}
		// record the property:
		// TODO: figure out how to undo the expansion: (need a test or 2)
		//			p.addProp(ret, pv)	// don't record, got expanded into a sequence of Nodes...

	case SE_idx:
		// record the property:
		p.addProp(ret, pv)

	case SL_idx:
		// record the property:
		p.addProp(ret, pv)

	case SO_idx:
		// set the board SO:
		p.SetSO(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case SQ_idx:
		// record the property:
		p.addProp(ret, pv)

	case ST_idx:
		// set the board ST:
		p.SetST(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case SU_idx:
		// record the property:
		p.addProp(ret, pv)

	case SZ_idx:
		// set the board size:
		var col int
		var row int
		idx_colon := strings.Index(string(pv.StrValue), ":")
		if idx_colon > 0 {
			col, _ = strconv.Atoi(string(pv.StrValue[0:idx_colon]))
			row, _ = strconv.Atoi(string(pv.StrValue[idx_colon+1:]))
		} else {
			col, _ = strconv.Atoi(string(pv.StrValue))
			row = col
		}
		p.InitAbstHier(ah.ColSize(col), ah.RowSize(row), ah.StringLevel, p.play) // TODO: vary this?
		// record the property:
		p.addProp(ret, pv)

	case TB_idx:
		// record the property:
		p.addProp(ret, pv)

	case TE_idx:
		// record the property:
		p.addProp(ret, pv)

	case TM_idx:
		// set the board TM:
		// see if it is a real (float32) value (in seconds)
		f, err := strconv.ParseFloat(string(pv.StrValue), 32)
		if err != nil {
			str := pv.StrValue
			// see if it is hours:
			idx := bytes.IndexByte(str, 'h')
			if idx > 0 {
				h, err2 := strconv.ParseFloat(string(str[0:idx]), 32)
				if err2 == nil {
					f = 60 * 60 * h
					if idx < len(str)-1 {
						str = str[idx+1:]
						for (len(str) > 0) && (str[0] == ' ') {
							str = str[1:]
						}
						if len(str) > 0 {
							idx = bytes.IndexByte(str, 'm')
							if idx > 0 {
								m, err3 := strconv.ParseFloat(string(str[0:idx]), 32)
								if err3 == nil {
									f += 60 * m
									if idx < len(str)-1 {
										p.ReportException(TM_idx, pv.StrValue, "Extra h+m value: "+string(str[idx+1:]))
										//											p.errors.Add(p.pos, "Extra h+m value: " + str[idx+1 :] + ": from " + string(pv.StrValue))
									}
								} else {
									p.ReportException(TM_idx, pv.StrValue, err3.Error())
									//										p.errors.Add(p.pos, "Bad minute value: " + err3.Error() + ": from " + string(pv.StrValue))
								}
							} else {
								idx = bytes.Index(str, []byte("sudden death"))
								if idx < 0 {
									idx = bytes.Index(str, []byte("each"))
									if idx < 0 {
										p.ReportException(TM_idx, pv.StrValue, "Extra h value: "+string(str))
										//											p.errors.Add(p.pos, "Extra h value: " + str +  ": from " + string(pv.StrValue))
									}
								}
							}
						}
					}
				} else {
					p.ReportException(TM_idx, pv.StrValue, err2.Error())
					//						p.errors.Add(p.pos, "Bad hour value: " + err2.Error() + ": from " + string(pv.StrValue[0:idx]))
				}
			} else {
				idx := bytes.IndexByte(str, 'm')
				if idx > 0 {
					m, err3 := strconv.ParseFloat(string(str[0:idx]), 32)
					if err3 == nil {
						f = 60 * m
						if idx < len(str)-1 {
							str = str[idx+1:]
							for (len(str) > 0) && (str[0] == ' ') {
								str = str[1:]
							}
							if len(str) > 0 {
								idx = bytes.IndexByte(str, 's')
								if idx > 0 {
									s, err4 := strconv.ParseFloat(string(str[0:idx]), 32)
									if err4 == nil {
										f += s
										if idx < len(str)-1 {
											p.ReportException(TM_idx, pv.StrValue, "Extra m+s value: "+string(str[idx+1:]))
											//												p.errors.Add(p.pos, "Extra m+s value: " + str[idx+1 :] + ": from " + string(pv.StrValue))
										}
									} else {
										p.ReportException(TM_idx, pv.StrValue, "Extra m+s value: "+string(str[idx+1:]))
										//											p.errors.Add(p.pos, "Bad second value: " + err4.String() + ": from " + string(pv.StrValue))
									}
								} else {
									idx = bytes.Index(str, []byte("sudden death"))
									if idx < 0 {
										p.ReportException(TM_idx, pv.StrValue, "Extra m value: "+string(str))
										//											p.errors.Add(p.pos, "Extra m value: " + str + ": from " + string(pv.StrValue))
									}
								}
							}

						}
					} else {
						p.ReportException(TM_idx, pv.StrValue, err3.Error()+"Bad minute value: "+string(str))
						//							p.errors.Add(p.pos, "Bad minute value: " + err3.String() + ": from " + string(pv.StrValue))
					}
				} else {
					p.ReportException(TM_idx, pv.StrValue, err.Error())
					//						p.errors.Add(p.pos, "Bad Timelimit value: " + err.Error() + ": from " + str)
				}
			}
		}
		p.SetTM(float32(f))
		// record the property
		p.addProp(ret, pv)

	case TR_idx:
		// record the property:
		p.addProp(ret, pv)

	case TW_idx:
		// record the property:
		p.addProp(ret, pv)

	case UC_idx:
		// record the property:
		p.addProp(ret, pv)

	case US_idx:
		// set the board US:
		p.SetUS(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case V_idx:
		// record the property:
		p.addProp(ret, pv)

	case VW_idx:
		// record the property:
		p.addProp(ret, pv)

	case W_idx:
		p.treeNodes[ret].TNodType = WhiteMoveNode
		mov, err := SGFPoint(pv.StrValue)
		if len(err) != 0 {
			p.errors.Add(p.pos, err.Error()+" in SGFPoint, W["+string(pv.StrValue)+"]")
		} else {
			p.treeNodes[ret].propListOrNodeLoc = PropIdx(mov)
			movN, err := p.DoW(mov, p.play)
			if len(err) != 0 {
				p.warnings.Add(p.pos, err.Error()+" W["+string(pv.StrValue)+"]")
			}
			if (p.moveLimit > 0) && (movN >= p.moveLimit) {
				p.limitReached = true
				if p.trace {
					p.printTrace("limitReached set to true in W")

				}
			}
		}

	case WB_idx:
		// record the property:
		p.addProp(ret, pv)

	case WC_idx:
		// record the property:
		p.addProp(ret, pv)

	case WL_idx:
		// record the property:
		p.addProp(ret, pv)

	case WO_idx:
		// record the property:
		p.addProp(ret, pv)

	case WR_idx:
		// count the WR values:
		idx := string(pv.StrValue)
		n, _ := BWRank_map[idx]
		BWRank_map[idx] = n + 1
		// check the rank
		errStr := check_Rank(pv.StrValue)
		if errStr != "" {
			p.ReportException(WR_idx, pv.StrValue, errStr)
		}
		// set the board WR:
		p.SetWR(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case WT_idx:
		// set the board WT:
		p.SetWT(pv.StrValue)
		// record the property:
		p.addProp(ret, pv)

	case WW_idx:
		// record the property:
		p.addProp(ret, pv)

	case UnknownPropIdx:
		// for UnknownProperty, add composed two strings: first is name, second is value
		str := string(p.UnknownProperty.ID) + ":" + string(pv.StrValue)
		pv.StrValue = []byte(str)
		p.addProp(ret, pv)
		p.warnings.Add(p.pos, "Unknown SGF property: "+str)

	default:
		p.errors.Add(p.pos, "Not Implemented: "+"default:")
	}
	return ret
}

// SetPlayerRank
func (gam *GameTree) SetPlayerRank() {
	// set the rank for the black name
	bn := gam.GetPB()
	if bn != nil {
		br := gam.GetBR()
		if br != nil {
			ix := string(bn)
			np, _ := BWPlayer_map[ix]
			np.LastRank = string(br)
			if gam.setFirstBRank {
				np.FirstRank = string(br)
			}
			BWPlayer_map[ix] = np
		}
	}
	// set the rank for the white name
	wn := gam.GetPW()
	if wn != nil {
		wr := gam.GetWR()
		if wr != nil {
			ix := string(wn)
			np, _ := BWPlayer_map[ix]
			np.LastRank = string(wr)
			if gam.setFirstWRank {
				np.FirstRank = string(wr)
			}
			BWPlayer_map[ix] = np
		}
	}
}

// parseProperties parses all the property names starting in a given node.
// Note: some properties, such S[abcdefgh] are expanded into additional nodes.
func (p *Parser) parseProperties(inRoot bool, parentNode TreeNodeIdx) (returnNode TreeNodeIdx) {
	if p.trace {
		defer un(trace(p, "parseProperties"))
	}
	returnNode = parentNode
	//	for (p.tok == IDENT ) && (p.limitReached != true) {
	for p.tok == IDENT {
		var prop *Property
		IDidx := LookUp(p.lit)
		if IDidx == UnknownPropIdx {
			prop = &p.UnknownProperty

			p.UnknownProperty.Note = Unknown_SGF4
			p.UnknownProperty.ID = p.lit // record the id of this unknown property
			p.UnknownProperty.Description = "<Unknown Property>"
			p.UnknownProperty.FF4Type = 0   //  "--".
			p.UnknownProperty.Qualifier = 0 // none.
			p.UnknownProperty.Value = 0     // Unknown.

		} else {
			prop = GetProperty(IDidx)
			err := checkPropertyType(prop, inRoot)
			if err != nil {
				p.warningPropertyType(p.pos, err)
			}
		}
		p.next()
		// TODO: does this need to be a for loop? for more than one value?
		// TODO: Is more than one value handled in parsePropValue?
		p.expect(LBRACK)
		propVal := p.parsePropValue(prop.Value)
		propVal.PropType = IDidx
		returnNode = p.processProperty(propVal, returnNode)
	}
	return returnNode
}

func (p *Parser) parseNodeSequence(parentNode TreeNodeIdx) (returnNode TreeNodeIdx) {
	if p.trace {
		defer un(trace(p, "parseNodeSequence"))
	}

	// parse a Node
	p.expect(SEMICOLON)

	// add node
	returnNode = p.addNode(parentNode, InteriorNode)

	//	p.next()
	returnNode = p.parseProperties(false, returnNode)

	// parse additional Nodes
	for (p.tok != RPAREN) && (p.tok != EOF) && (p.limitReached != true) {
		switch p.tok {
		case SEMICOLON:
			returnNode = p.addNode(returnNode, InteriorNode)
			p.next()
			returnNode = p.parseProperties(false, returnNode)
		case LPAREN:
			p.next()
			newLeaf := p.parseNodeSequence(TreeNodeIdx(returnNode))
			if p.limitReached != true {
				// Backup the Board State to returnNode
				currentNode := newLeaf
				for currentNode != returnNode {
					// for loop allows for more than one move at a node, i.e. S[]
					for p.GetMovDepth() > p.treeNodes[currentNode].movDepth {
						p.UndoBoardMove(p.play)
					}
					currentNode = p.treeNodes[currentNode].Parent
				}
				p.expect(RPAREN)
			}
		default:
			p.expect2(RPAREN, SEMICOLON)
		}
	}

	return returnNode
}

func (p *Parser) parseGame(parentNode TreeNodeIdx) (returnNode TreeNodeIdx) {
	if p.trace {
		defer un(trace(p, "parseGame"))
	}

	p.expect(SEMICOLON)

	// add GameInfo node
	newGame := p.addNode(parentNode, GameInfoNode)

	// parse GameInfo properties
	returnNode = p.parseProperties(true, newGame)

	// parse interior move Nodes
	for (p.tok != RPAREN) && (p.tok != EOF) && (p.limitReached != true) {
		switch p.tok {
		case SEMICOLON:
			returnNode = p.addNode(returnNode, InteriorNode)
			p.next()
			returnNode = p.parseProperties(false, returnNode)
		case LPAREN:
			p.next()
			newLeaf := p.parseNodeSequence(TreeNodeIdx(returnNode))
			if p.limitReached != true {
				// Backup the Board State to returnNode
				currentNode := newLeaf
				for currentNode != returnNode {
					// for loop allows for more than one move at a node, i.e. S[]
					for p.GetMovDepth() > p.treeNodes[currentNode].movDepth {
						p.UndoBoardMove(p.play)
					}
					currentNode = p.treeNodes[currentNode].Parent
				}
				p.expect(RPAREN)
			}
		default:
			p.expect2(RPAREN, SEMICOLON)
		}
	}

	if p.tok == RPAREN {
		p.next()
	}

	return returnNode
}

func (p *Parser) parseFile() {
	if p.trace {
		defer un(trace(p, "parseFile"))
	}

	fileCollection := p.addNode(0, CollectionNode)

	for (p.tok != EOF) && (p.limitReached != true) {
		p.expect(LPAREN)

		p.parseGame(fileCollection)
	}

	if p.treeNodes[fileCollection].Children == nilTreeNodeIdx {
		p.errors.Add(p.pos, "file contains no games")
	}

	if p.errors.ErrorCount() > 0 {
		p.errors.RemoveMultiples()
		ah.PrintError(os.Stderr, p.errors)
		ah.PrintError(os.Stdout, p.errors)
	}

	if p.warnings.ErrorCount() > 0 {
		p.warnings.RemoveMultiples()
		ah.PrintError(os.Stderr, p.warnings)
		ah.PrintError(os.Stdout, p.warnings)
	}
	return
}
