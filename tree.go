/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/sgf/tree.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/08/09.
 *  Copyright 2009-2014 Ken Friedenbach. All rights reserved.
 *
 */

// This package declares the types used to represent Go games
// as SGF trees with attached SGF properties.
//
// Acyclic Directed Graphs are supported, to allow representation
// of transformations and transpositions leading to equivalent postions.
//
package sgf

import (
	"fmt"
	"gitHub.com/Ken1JF/ahgo/ah"
	"strconv"
)

// TreeNode Types are represented by unsigned ints
//
type TreeNodeType uint8

const (
	RootNode TreeNodeType = iota
	CollectionNode
	GameInfoNode
	InteriorNode
	BlackMoveNode
	WhiteMoveNode
	SequenceNode // for S[m1m2m3...] property, (first move is Black)
	TransferNode
)

var TreeNodeTypeNames = []string{
	"RootNode",
	"CollectionNode",
	"GameInfoNode",
	"InteriorNode",
	"BlackMoveNode",
	"WhiteMoveNode",
	"SequenceNode",
	"TransferNode",
}

// Instead of pointers, Nodes and properties are placed in dynamic arrays,
// and indexed. This limits the largest collection of Nodes and properties
// to 64K.
//
// TODO: add ExtensionNodes to allow larger collections.
// TODO: add DirNodes and FileNodes to store larger collections
// in structures of directories and files.
//
type TreeNodeIdx uint16
type PropIdx uint16

const (
	nilTreeNodeIdx TreeNodeIdx = 0xFFFF
	nilPropIdx     PropIdx     = 0xFFFF
	MAX_NODE_IDX   int         = 0xFFFE
	MAX_PROP_IDX   int         = 0xFFFE
)

// A TreeNode can access its Parent, its Children, and its Siblings via indices into []TreeNode
//
type TreeNode struct {
	Parent            TreeNodeIdx // Root (0) has nilTreeNodeIdx as Parent.
	Children          TreeNodeIdx // Tail of circular linked list. nilTreeNodeIdx => no children.
	NextSib           TreeNodeIdx // Next Sibling in circular list. self => no siblings.
	propListOrNodeLoc PropIdx     // index into []Property or NodeLoc
	movDepth          int16       // index into []movs
	TNodType          TreeNodeType
}

// Property Values are stored in a tail circular list
//
type PropertyValue struct {
	StrValue []byte
	NextProp PropIdx
	PropType PropertyDefIdx
	ValType  PropValueType
}

// A GameTree consists of two slices:
// the first holds the tree/ADG Nodes
// the second holds property values other than moves
//
type GameTree struct {
	ah.AbstHier
	treeNodes      []TreeNode
	propertyValues []PropertyValue // TODO: add an avail list for deleted properties
	// for now, count and report
	NumberOfDeletedProperties int
	kM                        Komi
	rU                        []byte  // Rules
	rE                        Result  // Result
	tM                        float32 // Timelimit
	// how the board was setup
	aB ah.NodeLocList // Add Black
	// not needed to check consistency:	aE ah.NodeLocList	// Add Empty
	aW ah.NodeLocList // Add White
	// about the players
	// TODO: make pB, bR, pW, wR arrays, for team play?
	pB []byte // Player Black
	bR []byte // Black Rank
	bT []byte // Black Team
	pW []byte // Player White
	wR []byte // White Rank
	wT []byte // White Team
	// TODO: can these two "helper" vaules, and parser logic now be removed?
	// helper values for the parser
	setFirstBRank bool
	setFirstWRank bool
	// info about the game record
	fF []byte // File Format
	sT []byte // Style
	dT []byte // Date
	pC []byte // Place
	gC []byte // Game Comment
	uS []byte // User
	eV []byte // Event
	rO []byte // Round
	sO []byte // Source
	aP []byte // Application
	aN []byte // Annotation
	cP []byte // Copyright
	oH []byte // Old Handicap
	// drawing info
	// TODO: implement SGF drawing?
	// TODO: or remove these (currently) unused arrays
	aR [][2]ah.NodeLoc // Arrows
	lN [][2]ah.NodeLoc // Lines
}

// initGameTree needs to be called before the GameTree can be used
//
func (gT *GameTree) initGameTree() {
	// TODO: remove this? nilTreeNodeIdx has been changed to 0xFFFF, so 0 can be used...
	// add the RootNode
	gT.AddChild(nilTreeNodeIdx, RootNode, 0)
	// TODO: Remove this restriction: add a dummy property (can't use 0 location)
	//	_ := gT.AddAProp(0, pv)
}

// TODO: add an avail list for deleted properties
// for now count and report
func (gT *GameTree) ReportDeletedProperties() {
	fmt.Println("The number of deleted Properties = ", gT.NumberOfDeletedProperties)
}

func (gT *GameTree) AddToAvailProps(pidx PropIdx) {
	// TODO: add an avail list for deleted properties
	gT.NumberOfDeletedProperties += 1
	// for now, take off any list and reset values
	gT.propertyValues[pidx].NextProp = nilPropIdx
	gT.propertyValues[pidx].StrValue = nil
	gT.propertyValues[pidx].PropType = UnknownPropIdx
	gT.propertyValues[pidx].ValType = None
}

type TreeTraverseVisitFunc func(*GameTree, TreeNodeIdx)

var NumberOfAddedLabels = 0

func DoAddLabels(gamT *GameTree, nodIdx TreeNodeIdx) {
	lastCh := gamT.treeNodes[nodIdx].Children
	if lastCh != nilTreeNodeIdx {
		ch := gamT.treeNodes[lastCh].NextSib
		i := 1
		for ch != lastCh {
			nodLoc, _, _ := gamT.GetMove(gamT.treeNodes[ch])
			if nodLoc != ah.PassNodeLoc {
				i += 1
			}
			ch = gamT.treeNodes[ch].NextSib
		}
		if i > 1 {
			Labels := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			//Labels := [26]byte{'A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z'}
			var sep = make([]byte, 2)
			var colon = make([]byte, 1)
			sep[0] = ']'
			sep[1] = '['
			colon[0] = ':'
			// fmt.Println(" i = ", i)
			needsLabel := false
			if gamT.treeNodes[nodIdx].TNodType == BlackMoveNode {
				// change to InteriorNode and Add B_idx property
				nodLoc := ah.NodeLoc(gamT.treeNodes[nodIdx].propListOrNodeLoc)
				if nodLoc != ah.PassNodeLoc {
					pv := PropertyValue{StrValue: SGFCoords(nodLoc, gamT.IsFF4()), NextProp: nilPropIdx, PropType: B_idx, ValType: Move}
					gamT.treeNodes[nodIdx].TNodType = InteriorNode
					gamT.treeNodes[nodIdx].propListOrNodeLoc = nilPropIdx
					_ = gamT.addProperty(pv, nodIdx)
					needsLabel = true
				}
			} else if gamT.treeNodes[nodIdx].TNodType == WhiteMoveNode {
				// change to InteriorNode and Add W_idx property
				nodLoc := ah.NodeLoc(gamT.treeNodes[nodIdx].propListOrNodeLoc)
				if nodLoc != ah.PassNodeLoc {
					pv := PropertyValue{StrValue: SGFCoords(nodLoc, gamT.IsFF4()), NextProp: nilPropIdx, PropType: W_idx, ValType: Move}
					gamT.treeNodes[nodIdx].TNodType = InteriorNode
					gamT.treeNodes[nodIdx].propListOrNodeLoc = nilPropIdx
					_ = gamT.addProperty(pv, nodIdx)
					needsLabel = true
				}
			} else {
				needsLabel = true
			}
			if needsLabel {
				// Build the LB property
				pv := PropertyValue{StrValue: nil, NextProp: nilPropIdx, PropType: LB_idx, ValType: ListOfCompPoint_simpTest}
				lastCh = gamT.treeNodes[nodIdx].Children
				ch = gamT.treeNodes[lastCh].NextSib
				j := 0
				for ch != lastCh {
					nodLoc, _, _ := gamT.GetMove(gamT.treeNodes[ch])
					if nodLoc != ah.PassNodeLoc {
						// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
						pv.StrValue = append(pv.StrValue, SGFCoords(nodLoc, gamT.IsFF4())...)
						// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
						pv.StrValue = append(pv.StrValue, colon[0])
						// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
						pv.StrValue = append(pv.StrValue, Labels[j])
						// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
						j += 1
						if j < i {
							pv.StrValue = append(pv.StrValue, sep...)
							// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
						}
					}
					ch = gamT.treeNodes[ch].NextSib
				}
				nodLoc, _, _ := gamT.GetMove(gamT.treeNodes[ch])
				if nodLoc != ah.PassNodeLoc {
					pv.StrValue = append(pv.StrValue,
						SGFCoords(nodLoc, gamT.IsFF4())...)
					pv.StrValue = append(pv.StrValue, colon[0])
					pv.StrValue = append(pv.StrValue, Labels[j])
					// fmt.Printf(" StrValue = %s\n", string(pv.StrValue))
				}
				//add the LB property
				_ = gamT.addProperty(pv, nodIdx)
				NumberOfAddedLabels += 1
			}
		}
	}
}

func DoRemoveLabels(gamT *GameTree, nodIdx TreeNodeIdx) {
	nod := gamT.treeNodes[nodIdx]
	switch nod.TNodType {
	case RootNode, CollectionNode, BlackMoveNode, WhiteMoveNode:
		{
			// nothing to do, no properties
		}
	case GameInfoNode, InteriorNode:
		{
			lastProp := gamT.treeNodes[nodIdx].propListOrNodeLoc
			prevProp := nilPropIdx
			process := func(prop PropIdx) {
				if gamT.propertyValues[prop].PropType == LB_idx {
					if prop == lastProp { // deleting last prop
						if prevProp == nilPropIdx { // and only one
							gamT.treeNodes[nodIdx].propListOrNodeLoc = nilPropIdx
						} else { // there are others
							// set the new last
							gamT.treeNodes[nodIdx].propListOrNodeLoc = prevProp
							// set the first value
							gamT.propertyValues[prevProp].NextProp = gamT.propertyValues[lastProp].NextProp
						}
					} else { // not deleting the last one
						if prevProp == nilPropIdx { // deleting the first
							gamT.propertyValues[lastProp].NextProp = gamT.propertyValues[prop].NextProp
						} else {
							gamT.propertyValues[prevProp].NextProp = gamT.propertyValues[prop].NextProp
						}
					}
					// add deleted property to avail list
					gamT.AddToAvailProps(prop)
				}
			}
			if lastProp != nilPropIdx {
				prop := gamT.propertyValues[lastProp].NextProp
				nextProp := gamT.propertyValues[prop].NextProp // peek ahead, in case deleted
				process(prop)
				for prop != lastProp {
					prevProp = prop
					prop = nextProp
					nextProp = gamT.propertyValues[prop].NextProp // peek ahead, in case deleted
					process(prop)
				}
			}
		}
	case SequenceNode:
	case TransferNode:
	default:
		fmt.Println("Unknown NodeType, nod = ", nod, " TNodType = ", nod.TNodType)
	}
	/*
	   for i := 0; i <= int(nod.movDepth)+1; i +=1 {
	       fmt.Print(".")
	   }
	   switch nod.TNodType {
	       case RootNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (RootNode)")
	       case CollectionNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (CollectionNode)")
	       case GameInfoNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (GameInfoNode)")
	       case InteriorNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (InteriorNode)")
	       case BlackMoveNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (BlackMoveNode)")
	       case WhiteMoveNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (WhiteMoveNode)")
	       case SequenceNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (SequenceNode)")
	       case TransferNode:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (TransferNode)")
	       default:
	           fmt.Println("Node ", nod,
	                       " TNodType ", nod.TNodType, " (default)")
	   }
	*/

}

/*
func printQueue(nodQueue []TreeNodeIdx) {
    fmt.Println(" node queue, len = ", len(nodQueue), " cap = ", cap(nodQueue))
    for i:=0; i<len(nodQueue); i+=1 {
        fmt.Print(" node # ", i, " = ", nodQueue[i])
    }
    fmt.Println(".")
}
*/

// BreadthFirstTraverse
//
func (gamT *GameTree) BreadthFirstTraverse(preVisit bool, Visit TreeTraverseVisitFunc) {
	var nodQueue []TreeNodeIdx
	// fmt.Println("Parent Children NextSib propListOrNodeLoc movDepth TNodType")
	nodQueue = append(nodQueue, 0)
	// printQueue(nodQueue)
	j := 0
	for len(nodQueue) > 0 {
		nod := nodQueue[0]
		nodQueue = nodQueue[1:]
		if preVisit {
			Visit(gamT, nod)
		}
		lastCh := gamT.treeNodes[nod].Children
		if lastCh != nilTreeNodeIdx {
			ch := gamT.treeNodes[lastCh].NextSib
			nodQueue = append(nodQueue, ch)
			// printQueue(nodQueue)
			i := 0
			for ch != lastCh {
				// fmt.Println("ch = ", ch, " lastCh = ", lastCh)
				ch = gamT.treeNodes[ch].NextSib
				nodQueue = append(nodQueue, ch)
				// printQueue(nodQueue)
				i += 1
				if i > 20 {
					j += 1
					fmt.Println("  more than 20 children? ")
					break
				}
			}
		}
		if !preVisit {
			Visit(gamT, nod)
		}
		// fmt.Println("  queue len = ", len(nodQueue))
		if j > 5 {
			break
		}
	}
}

// DepthFirstTraverse - nodes are visited in either Pre-order, before children
//  or Post-order, after all children
//
func (gamT *GameTree) DepthFirstTraverse(preVisit bool, Visit TreeTraverseVisitFunc) {
	type dftElement struct {
		nod_tIdx TreeNodeIdx
		cur_ch   TreeNodeIdx
		last_ch  TreeNodeIdx
	}
	var nodStack []dftElement
	var nod dftElement
	var rootNode TreeNodeIdx = 0

	// fmt.Println("Parent Children NextSib propListOrNodeLoc movDepth TNodType")

	lastChild := gamT.treeNodes[rootNode].Children
	firstChild := nilTreeNodeIdx
	nodStack = append(nodStack, dftElement{nod_tIdx: rootNode, cur_ch: firstChild, last_ch: lastChild})
	// fmt.Println("  stack len = ", len(nodStack), " should be 1 (push rootnode)")
	i := 0
	for len(nodStack) > 0 {
		// take an element off the stack
		nod = nodStack[len(nodStack)-1]
		nodStack = nodStack[0 : len(nodStack)-1]
		// fmt.Println("  stack len = ", len(nodStack), " should be 1 less (pop a node)")

		if preVisit && (nod.cur_ch == nilTreeNodeIdx) {
			// fmt.Print("PreVisit: ")
			Visit(gamT, nod.nod_tIdx)
		}
		// see if there are children
		if nod.last_ch != nilTreeNodeIdx { // there are children

			if nod.cur_ch == nilTreeNodeIdx { // this is the first child, set cur_ch
				nod.cur_ch = gamT.treeNodes[nod.last_ch].NextSib
				// and put this node back, for later children
				nodStack = append(nodStack, nod)
				// fmt.Println("  stack len = ", len(nodStack), " should be 1 more (put back node, on first child)")
				// build a child element
				var ch_nod dftElement
				ch_nod.nod_tIdx = nod.cur_ch
				ch_nod.last_ch = gamT.treeNodes[nod.cur_ch].Children
				ch_nod.cur_ch = nilTreeNodeIdx
				// and put on the stack
				nodStack = append(nodStack, ch_nod)
				// fmt.Println("  stack len = ", len(nodStack), " should be 1 more (put first child)")
			} else { // this is second, etc. if any
				if nod.cur_ch != nod.last_ch { // more
					nod.cur_ch = gamT.treeNodes[nod.cur_ch].NextSib
					// and put this node back, for later children
					nodStack = append(nodStack, nod)
					// fmt.Println("  stack len = ", len(nodStack), " should be 1 more (put back next child)")
					// build a child element
					var ch_nod dftElement
					ch_nod.nod_tIdx = nod.cur_ch
					ch_nod.last_ch = gamT.treeNodes[nod.cur_ch].Children
					ch_nod.cur_ch = nilTreeNodeIdx
					// and put on the stack
					nodStack = append(nodStack, ch_nod)
					// fmt.Println("  stack len = ", len(nodStack), " should be 1 more(put next child)")
				} else { // done with all children
					if !preVisit {
						// fmt.Print("PostVisit: ")
						Visit(gamT, nod.nod_tIdx)
					}
				}
			}
		} else { // there are no children
			if !preVisit {
				// fmt.Print("PostVisit: ")
				Visit(gamT, nod.nod_tIdx)
			}
		}
		i += 1
		// fmt.Println(" end loop # ", i, " stack len = ", len(nodStack))
	}
}

// HasSiblings returns true if the node has siblings
//
func (gamT *GameTree) HasSiblings(nd TreeNodeIdx) (ret bool) {
	if nd > 0 { // RootNode can't have siblings
		sib := gamT.treeNodes[nd].NextSib
		if sib != nd {
			ret = true
		}
	}
	return ret
}

// addProperty appends the new property, and maintains a circular linked list
//
func (gamT *GameTree) addProperty(pv PropertyValue, nd TreeNodeIdx) (err ah.ErrorList) {
	cur_l := len(gamT.propertyValues)
	if cur_l <= MAX_PROP_IDX {
		gamT.propertyValues = append(gamT.propertyValues, pv)
		if gamT.treeNodes[nd].propListOrNodeLoc == nilPropIdx { // first property
			gamT.treeNodes[nd].propListOrNodeLoc = PropIdx(cur_l)
			gamT.propertyValues[cur_l].NextProp = PropIdx(cur_l) // circular tail list
		} else {
			head := gamT.propertyValues[gamT.treeNodes[nd].propListOrNodeLoc].NextProp // get head of list
			gamT.propertyValues[cur_l].NextProp = head
			gamT.propertyValues[gamT.treeNodes[nd].propListOrNodeLoc].NextProp = PropIdx(cur_l)
			gamT.treeNodes[nd].propListOrNodeLoc = PropIdx(cur_l)
		}
	} else {
		// TODO: add logic to split the GameTree propertyValues
		err.Add(ah.NoPos, "addProperty: too many properties "+strconv.Itoa(MAX_PROP_IDX))
	}
	return err
}

// AddAProp adds a property to a node.
// AddAProp changes a BlackMoveNode or a WhiteMoveNode into an InteriorNode,
// when a property is added, making the B or W property the first in the list.
//
func (gamT *GameTree) AddAProp(n TreeNodeIdx, pv PropertyValue) (err ah.ErrorList) {
	mov := gamT.treeNodes[n].propListOrNodeLoc
	if gamT.treeNodes[n].TNodType == BlackMoveNode {
		var movPV *PropertyValue = new(PropertyValue)
		movPV.StrValue = SGFCoords(ah.NodeLoc(mov), gamT.IsFF4())
		movPV.NextProp = nilPropIdx
		movPV.PropType = B_idx
		movPV.ValType = Move
		gamT.treeNodes[n].TNodType = InteriorNode
		gamT.treeNodes[n].propListOrNodeLoc = nilPropIdx
		err := gamT.addProperty(*movPV, n)
		if len(err) != 0 {
			err.Add(ah.NoPos, "adding B property "+err.Error())
		}
	} else if gamT.treeNodes[n].TNodType == WhiteMoveNode {
		var movPV *PropertyValue = new(PropertyValue)
		movPV.StrValue = SGFCoords(ah.NodeLoc(mov), gamT.IsFF4())
		movPV.NextProp = nilPropIdx
		movPV.PropType = W_idx
		movPV.ValType = Move
		gamT.treeNodes[n].TNodType = InteriorNode
		gamT.treeNodes[n].propListOrNodeLoc = nilPropIdx
		err := gamT.addProperty(*movPV, n)
		if len(err) != 0 {
			err.Add(ah.NoPos, "adding W property "+err.Error())
		}
	}
	err = gamT.addProperty(pv, n)
	if len(err) != 0 {
		err.Add(ah.NoPos, "adding property "+err.Error())
	}
	return err
}

// AddChild appends a new node, and maintains a circular linked list of siblings
//
func (gamT *GameTree) AddChild(par TreeNodeIdx, ndty TreeNodeType, mDep int16) (idx TreeNodeIdx, err ah.ErrorList) {
	if ah.TraceAH {
		fmt.Println("AddChild:", par, "type:", TreeNodeTypeNames[ndty])
	}
	var newTn TreeNode
	cur_l := len(gamT.treeNodes)
	if cur_l <= MAX_NODE_IDX {
		gamT.treeNodes = append(gamT.treeNodes, newTn)
		gamT.treeNodes[cur_l].TNodType = ndty
		gamT.treeNodes[cur_l].propListOrNodeLoc = nilPropIdx
		gamT.treeNodes[cur_l].Parent = par
		gamT.treeNodes[cur_l].Children = nilTreeNodeIdx    // when added, no children
		gamT.treeNodes[cur_l].NextSib = TreeNodeIdx(cur_l) // when added, circular list points to self
		gamT.treeNodes[cur_l].movDepth = mDep
		if par != nilTreeNodeIdx { // nilTreeNodeIdx indicates no parent
			tail := gamT.treeNodes[par].Children
			if tail == nilTreeNodeIdx { // first member
				// adding cur_l as first child, already points to self
				gamT.treeNodes[par].Children = TreeNodeIdx(cur_l)
			} else {
				head := gamT.treeNodes[tail].NextSib
				// adding cur_l as a new tail element
				// cur_l will be new tail, point to  head
				gamT.treeNodes[cur_l].NextSib = head
				// old tail will point to cur_l
				gamT.treeNodes[tail].NextSib = TreeNodeIdx(cur_l)
				// parent points to new tail
				gamT.treeNodes[par].Children = TreeNodeIdx(cur_l)
			}
		}
	} else {
		// TODO: add logic to split the GameTree treeNodes
		err.Add(ah.NoPos, "AddChild: too many nodes "+strconv.Itoa(MAX_NODE_IDX))
	}
	return TreeNodeIdx(cur_l), err
}

// FindChild returns the index of a child of with a move at mov
//
func (gamT *GameTree) FindChild(par TreeNodeIdx, mov ah.NodeLoc) (found TreeNodeIdx) {
	found = nilTreeNodeIdx
	var ch TreeNodeIdx = nilTreeNodeIdx
	var p_idx PropIdx = nilPropIdx
	checkMov := func() {
		if gamT.propertyValues[p_idx].PropType == B_idx || gamT.propertyValues[p_idx].PropType == W_idx {
			mv, err := SGFPoint(gamT.propertyValues[p_idx].StrValue)
			if len(err) != 0 {
				return
			}
			if mv == mov {
				found = ch // set found to the child index
			}
		}
	}
	lookFor := func() {
		typ := gamT.treeNodes[ch].TNodType
		switch typ {
		case InteriorNode:
			var tail_p PropIdx = gamT.treeNodes[ch].propListOrNodeLoc
			if tail_p != nilPropIdx {
				p_idx = tail_p
				// check for mov
				checkMov()
				for p_idx != tail_p && found == nilTreeNodeIdx {
					p_idx = gamT.propertyValues[p_idx].NextProp
					// check for mov
					checkMov()
				}
			}
		case BlackMoveNode, WhiteMoveNode:
			// TODO: need to check the mov color? currently, no
			if ah.NodeLoc(gamT.treeNodes[ch].propListOrNodeLoc) == mov {
				found = ch
			}
		default:
		}
	}
	tail := gamT.treeNodes[par].Children
	if tail != nilTreeNodeIdx { // check if any children
		ch = gamT.treeNodes[tail].NextSib // get the first child
		// look for a move at mov
		lookFor()
		for ch != tail && found == nilTreeNodeIdx {
			ch = gamT.treeNodes[ch].NextSib // get the next child
			// look for a move at mov
			lookFor()
		}
	}
	return found
}
