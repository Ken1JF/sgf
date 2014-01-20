/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/sgf/printer.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/08/09.
 *  Copyright 2009-2014 Ken Friedenbach. All rights reserved.
 *
 *	This package implements writing of SGF game trees.
 */

package sgf

import (
	"bufio"
	"errors"
	"fmt"
	"gitHub.com/Ken1JF/ahgo/ah"
	"os"
	"strconv"
	"strings"
)

var indent int = 0 // TODO: make this a per parser variable?

func u(s string) {
	if ah.TraceAH {
		indent -= 1
		for i := indent; i > 0; i-- {
			fmt.Print(". ")
		}
		fmt.Println("Leaving", s)
	}
}

func tr(s string) string {
	if ah.TraceAH {
		for i := indent; i > 0; i-- {
			fmt.Print(". ")
		}
		fmt.Println("Entering", s)
		indent += 1
	}
	return s
}

func (pv *PropertyValue) writeProperty(w *bufio.Writer, FF4 bool) (err error) {
	defer u(tr("writeProperty"))
	pt := pv.PropType
	prop := GetProperty(pt)
	if prop == nil { // either error or UnknownProperty
		if pt == UnknownPropIdx {
			idx := strings.Index(string(pv.StrValue), ":")
			if idx > 0 {
				_, err = w.Write(pv.StrValue[0:idx])
				if err == nil {
					err = w.WriteByte('[')
					str := pv.StrValue[idx+1:]
					if err == nil {
						_, err = w.Write(str)
						if err == nil {
							err = w.WriteByte(']')
						}
					}
				}
			} else {
				fmt.Println("*** no \":\" in UnknownProperty value in writeProperty")
				err = errors.New("writeProperty: no \":\" in UnknownProperty value.")
			}
		} else {
			fmt.Println("*** nil property pointer in writeProperty")
			return errors.New("writeProperty: BAD PropertyDefIdx " + strconv.FormatInt(int64(pt), 10))
		}
	} else {
		_, err = w.Write(prop.ID)
		if err == nil {
			err = w.WriteByte('[')
			str := pv.StrValue
			if err == nil {
				if (pt == AB_idx) || (pt == AE_idx) || (pt == AW_idx) || (pt == S_idx) || (pt == TB_idx) || (pt == TR_idx) || (pt == TW_idx) { // split into pairs
					for len(str) > 2 {
						_, err = w.Write(str[0:2])
						_, err = w.WriteString("][")
						str = str[2:]
					}
				}
				if ((len(str) == 2) && (str[0] == 't') && (str[1] == 't')) &&
					(((pt == B_idx) || (pt == W_idx)) && (FF4 == true)) {
					// replace with empty string
				} else {
					_, err = w.Write(str)
				}
				if err == nil {
					err = w.WriteByte(']')
				}
			}
		}
	}
	return err
}

func (p *GameTree) writeProperties(w *bufio.Writer, n TreeNodeIdx, onePer bool) (err error) {
	defer u(tr("writeProperties"))
	lastProp := p.treeNodes[n].propListOrNodeLoc
	if lastProp != nilPropIdx {
		prop := p.propertyValues[lastProp].NextProp
		err = p.propertyValues[prop].writeProperty(w, p.IsFF4())
		if err == nil {
			for (prop != lastProp) && (err == nil) {
				prop = p.propertyValues[prop].NextProp
				err = p.propertyValues[prop].writeProperty(w, p.IsFF4())
				if err == nil {
					if onePer {
						err = w.WriteByte('\n')
					}
				}
			}
		}
	}
	return err
}

func (p *GameTree) writeLabel(w *bufio.Writer, n ah.NodeLoc, LabelIdx int) (err error) {
	Labels := [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	err = w.WriteByte('[')
	if err == nil {
		_, err = w.Write(SGFCoords(n, false))
		if err == nil {
			err = w.WriteByte(':')
			if err == nil {
				err = w.WriteByte(Labels[LabelIdx])
				if err == nil {
					err = w.WriteByte(']')
				}
			}
		}
	}
	return err
}

//	writeTree writes a .sgf tree from the treeNodes array
//		w is a buffered I/O writer
//		n is the TreeNodeIdx of the root of this tree
//		needs is a bool that is true when a \n is needed
//		nMov keeps a count of moves per line.
//	writeTree first writes one node, then recursively calls writeTree
//	writeTree is only called from writeGame, which has one active call, and one that is never reached (&& false)
func (p *GameTree) writeTree(w *bufio.Writer, n TreeNodeIdx, needs bool, nMov int, nMovPerLine int) (err error) {
	defer u(tr("writeTree"))
	if needs == true {
		if nMov > 0 {
			err = w.WriteByte('\n')
			nMov = 0
		}
		err = w.WriteByte('(')
	}
	if err == nil {
		if nMov == nMovPerLine {
			err = w.WriteByte('\n')
			nMov = 0
		}
		err = w.WriteByte(';')
		// write the node
		typ := p.treeNodes[n].TNodType
		switch typ {
		case GameInfoNode:
			//           fmt.Println("writing GameInfoNode\n")
			err = p.writeProperties(w, n, true)
		case InteriorNode:
			//           fmt.Println("writing InteriorNode\n")
			err = p.writeProperties(w, n, false)
		case BlackMoveNode:
			_, err = w.WriteString("B[")
			_, err = w.Write(SGFCoords(ah.NodeLoc(p.treeNodes[n].propListOrNodeLoc), p.IsFF4()))
			err = w.WriteByte(']')
			nMov += 1
		case WhiteMoveNode:
			_, err = w.WriteString("W[")
			_, err = w.Write(SGFCoords(ah.NodeLoc(p.treeNodes[n].propListOrNodeLoc), p.IsFF4()))
			err = w.WriteByte(']')
			nMov += 1
		default:
			fmt.Println("*** unsupported TreeNodeType in writeTree")
			err = errors.New("writeTree: unsupported TreeNodeType" + strconv.FormatInt(int64(typ), 10))
			return err
		}
		if err == nil {
			// write the children
			lastCh := p.treeNodes[n].Children
			if lastCh != nilTreeNodeIdx && err == nil {
				ch := p.treeNodes[lastCh].NextSib
				chNeeds := (lastCh != ch)
				err = p.writeTree(w, ch, chNeeds, nMov, nMovPerLine)
				for ch != lastCh && err == nil {
					ch = p.treeNodes[ch].NextSib
					//					nMov += 1
					err = p.writeTree(w, ch, chNeeds, nMov, nMovPerLine)
				}
			}
			if (err == nil) && (needs == true) {
				err = w.WriteByte(')')
			}
		}
	}
	return err
}

//	writeGame is called to write a .sgf game tree from writeCollection
//		w is a buffered I/O writer
//		n is the TreeNodeIdx where the game begins
//	writeGame returns an Error (nil if no error encounterd)
//	writeGame writes the initial "(", then calls writeTree.
//	there is logic, which is forced to fail (&& false) for writing siblings of n (doesn't writeTree do this)
//	if writeTree does not return an error, writeGame writes the terminating ")" with newlines before and after.
func (p *GameTree) writeGame(w *bufio.Writer, n TreeNodeIdx, nMovPerLine int) (err error) {
	defer u(tr("writeGame"))
	err = w.WriteByte('(')
	if err == nil {
		// TODO: allow more than one game in a file
		//		hasSibs := p.HasSiblings(n)
		err := p.writeTree(w, n, false, 0, nMovPerLine)
		if err == nil && false /* && hasSibs */ {
			nxtSib := p.treeNodes[n].NextSib
			for nxtSib != n {
				err = p.writeTree(w, nxtSib, false, 0, nMovPerLine)
				nxtSib = p.treeNodes[nxtSib].NextSib
			}
		}
		if err == nil {
			_, err = w.WriteString("\n)\n")
		}
	}
	return err
}

//	writeCollection is only called from writeParseTree
//		w is a buffered I/O writer
//		coll is the TreeNodeIdx of the collection
//	returns an Error if one is encountered
//	writeCollection verifies that coll has TNodType == CollectionNode
//	TODO: could crash if coll is out of range...
//	then checks if coll has children, if so gets first child, and calls writeGame
//	if there are additional siblings, it calls writeGame for each sibling
//	TODO: does not appear to have logic to stop after first err is returned.
//	TODO: need a short sgf test for multiple games in a collection?
//	TODO: this level of logic seems to satisfy TODO statements above, and forced to fail logic.
func (p *GameTree) writeCollection(w *bufio.Writer, coll TreeNodeIdx, nMovPerLine int) (err error) {
	defer u(tr("writeCollection"))
	typ := p.treeNodes[coll].TNodType
	if typ == CollectionNode {
		lastCh := p.treeNodes[coll].Children
		if lastCh != nilTreeNodeIdx {
			ch := p.treeNodes[lastCh].NextSib // get first child
			err = p.writeGame(w, ch, nMovPerLine)
			for ch != lastCh {
				ch = p.treeNodes[ch].NextSib
				err = p.writeGame(w, ch, nMovPerLine)
			}
		} else {
			return errors.New("writeCollection, no Games.")
		}
	} else {
		return errors.New("writeCollection, no CollectionNode: " + strconv.FormatInt(int64(typ), 10))
	}
	return err
}

// writeParseTree is only called from WriteFile
//		w is a buffered I/O writer, which has been successfully opened for writing
// return an Error if one is encountered
// writeParseTrre checks that treeNodes[0] has type RootNode,
// then calls writeCollection for the Children node
// TODO: could crash if 0 is out of range? OR does init function in parser.go prevent this?
// TODO: could crash if RootNode (0) has no Children
// TODO: does not check if RootNode has more than one child.
func (p *GameTree) writeParseTree(w *bufio.Writer, nMovPerLine int) (err error) {
	defer u(tr("writeParseTree"))
	typ := p.treeNodes[0].TNodType
	if typ == RootNode {
		coll := p.treeNodes[0].Children
		err = p.writeCollection(w, coll, nMovPerLine)
	} else {
		return errors.New("writeParseTree, no RootNode: " + strconv.FormatInt(int64(typ), 10))
	}
	return err
}

const filePERM uint32 = 0644 // owner RW, group R, others R

//	WriteFile is used to write a .sgf file from a tree contained in Parser structure
//		fileName is a string that contains the full path name of the file to write
//	WriteFile calls os.Open for write only, create if necessary, truncate if exists,
//	and sets the file permisions to Owner RW, Group R, Others R
//	TODO: the file open options and permisions are currently hard coded?
//  Do these need to be varied, say by program options?
//	WriteFile checks for errors on open, and returns immediately on an error.
//	WriteFile sets up a deferred Close on the

func (tree *GameTree) WriteFile(fileName string, nMovPerLine int) (err error) {
	defer u(tr("WriteFile"))
	// old parms to Open(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePERM)
	f, err := os.Create(fileName)
	if err != nil {
		return errors.New("OpenFile:" + fileName + " " + err.Error())
	}
	defer f.Close() // TODO: should this be conditional on not being closed?
	w := bufio.NewWriter(f)
	if w == nil {
		return errors.New("nil from NewWriter:" + fileName + " " + err.Error())
	}
	err = tree.writeParseTree(w, nMovPerLine)
	if err != nil {
		return errors.New("Error:" + fileName + " " + err.Error())
	}
	err = w.Flush()
	err = f.Close()
	return err
}
