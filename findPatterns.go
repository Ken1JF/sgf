/*
 *  File:		src/github.com/Ken1JF/sgf/findPatterns.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 4/17/2011.
 *  Copyright 2011-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file walks SGF game trees, and records patterns, etc.
 */

package sgf

import (
	"github.com/Ken1JF/ah"
	"strconv"
)

type traversePoint struct {
	cGam, cPat TreeNodeIdx
	pDep       int
	mkGd       bool
}

// AddTeachingPattern adds one or more patterns from a GameTree, to a global pattern tree (DAG)
//	szCol, sizRow is the board size
//	ha is the handicap
//	firstMove is the color of the first move, based on handicap (0 => Black first, 2 or more => White first)
//		But, some games have a starting pattern placed on the board with AB, AW, etc.
//			at least one game has first move by White (skip it (them?)).
//	pattTree is the GameTree which needs the pattern to be added, based on the handicap.
//		Note: may be nil on first use.
//	pattType is the type of Pattern being stored in pattTree
//		WHOLE_BOARD_PATTERN, etc.
// TODO: add ohter types of patterns
//	moveLimit is the maximum move number to place in the pattTree
//		Does not include handicap or other pre-placed stones.
//
// returns an Error if one is detected, a translation that takes the first move into a canonical location, and the updated pattTree
func (gamT *GameTree) AddTeachingPattern(szCol ah.ColSize, szRow ah.RowSize, ha int, pattTree *GameTree,
	pattType ah.PatternType, moveLimit int, patternLimit int, skipFiles int) (err ah.ErrorList, trans ah.BoardTrans, upPattTree *GameTree) {
	var collPatt, gInfoPatt TreeNodeIdx
	var curGam, curPatt TreeNodeIdx
	var pv PropertyValue
	var nodColr ah.PointStatus = ah.White
	var traverseStack []traversePoint
	var onMain bool = true
	var markGood bool = false
	var nodLoc ah.NodeLoc
	var newNodLoc ah.NodeLoc
	var patternDepth int
	var limitReached bool = false

	// compute color of firstMoveColor
	firstMoveColor := ah.White
	if ha == 0 {
		firstMoveColor = ah.Black
	}
	// count the moves:
	nMoves := 0
	found := false
	findOrAdd := func(newNL ah.NodeLoc) (idx TreeNodeIdx) { // find a child of curPatt or add one
		patternDepth += 1
		if (patternLimit > 0) && (patternDepth >= patternLimit) {
			limitReached = true
		}
		//		str := strconv.Itoa(int(nodColr))
		idx = pattTree.FindChild(curPatt, newNL)
		if idx == nilTreeNodeIdx { // not found, add it
			moveType := WhiteMoveNode
			if nodColr == ah.Black {
				moveType = BlackMoveNode
			}
			idxx, err := pattTree.AddChild(curPatt, moveType, 0)
			if len(err) != 0 {
				return
			}
			//			fmt.Printf("findOrAdd: added %d\n", idxx)
			pattTree.treeNodes[idxx].propListOrNodeLoc = PropIdx(newNL)
			idx = idxx
		}
		return idx
	}
	// if pattTree doesn't exist, create it, and initialize it
	if pattTree == nil {
		pattTree = new(GameTree)
		pattTree.initGameTree()
		collPatt, err = pattTree.AddChild(0, CollectionNode, 0)
		if len(err) != 0 {
			return err, trans, upPattTree
		}
		gInfoPatt, err = pattTree.AddChild(collPatt, GameInfoNode, 0)
		if len(err) != 0 {
			return err, trans, upPattTree
		}
		// set the curPatt node in the pattTree
		curPatt = gInfoPatt
		// TODO: are these needed? FF, GM, CA, AP, ST?
		// ADD FF
		pv.StrValue = []byte("4")
		pv.PropType = FF_idx
		pv.ValType = Num_1_4
		pattTree.AddAProp(gInfoPatt, pv)
		// ADD GM
		pv.StrValue = []byte("1")
		pv.PropType = GM_idx
		pv.ValType = Num_1_5_or_7_16
		pattTree.AddAProp(gInfoPatt, pv)
		// ADD CA
		pv.StrValue = []byte("UTF-8")
		pv.PropType = CA_idx
		pv.ValType = SimpText
		pattTree.AddAProp(gInfoPatt, pv)
		// ADD AP
		// TODO: make it vary with releases?
		pv.StrValue = []byte("ahgo:0.8")
		pv.PropType = AP_idx
		pv.ValType = CompSimpText_simpText
		pattTree.AddAProp(gInfoPatt, pv)
		// ADD ST
		pv.StrValue = []byte("1")
		pv.PropType = ST_idx
		pv.ValType = Num_0_3
		pattTree.AddAProp(gInfoPatt, pv)

		// TODO: support n x m boards. Add SZ
		pv.StrValue = []byte(strconv.Itoa(int(szCol)))
		pv.PropType = SZ_idx
		pattTree.AddAProp(gInfoPatt, pv)
		// Add HA
		pv.StrValue = []byte(strconv.Itoa(ha))
		pv.PropType = HA_idx
		pattTree.AddAProp(gInfoPatt, pv)
		pattTree.InitAbstHier(szCol, szRow, ah.StringLevel, true)
		pattTree.SetHandicap(ha)
		pv.StrValue = pattTree.PlaceHandicap(ha, int(szCol))

		if pv.StrValue != nil {
			// Add the AB for handicap points
			pv.PropType = AB_idx
			pv.ValType = ListOfStone
			pattTree.AddAProp(gInfoPatt, pv)
		}
		//		fmt.Printf("Created pattTree: collPatt %d gInfoPatt %d curPatt %d\n", collPatt, gInfoPatt, curPatt)
	} else {
		collPatt = 1  // CollectionNode is child of RootNode
		gInfoPatt = 2 // GameInfoNode is child of CollectionNode
		curPatt = gInfoPatt
		//		fmt.Printf("Reuse pattTree: collPatt %d gInfoPatt %d curPatt %d\n", collPatt, gInfoPatt, curPatt)
	}

	// traverse gamT
	// first visit the main line of play, via children (tail of circular linked list)
	// find the first move:
findFirst:
	for i, nod := range gamT.treeNodes {
		switch nod.TNodType {
		case RootNode:
		case CollectionNode:
			if i != 1 {
				err.Add(ah.NoPos, "AddTeachingPattern: CollectionNode is not correct "+strconv.Itoa(i))
				return
			}
			//				collGam = TreeNodeIdx(i)
		case GameInfoNode:
			if i != 2 {
				err.Add(ah.NoPos, "AddTeachingPattern: GameInfoNode is not correct "+strconv.Itoa(i))
				return
			}
			//				gInfoGam = TreeNodeIdx(i)
		case InteriorNode, BlackMoveNode, WhiteMoveNode:
			nMoves += 1
			curGam = TreeNodeIdx(i)
			nodLoc, nodColr, err = gamT.GetMove(nod)
			if len(err) != 0 {
				return err, trans, upPattTree
			}
			if nMoves > gamT.Board.GetNMoves() {
				//					fmt.Printf("Need to re-DoBoardMove: %d %d \n", nMoves, gamT.Board.GetNMoves())
				_, err = gamT.AbstHier.DoBoardMove(nodLoc, nodColr, true)

			} else {
				//					fmt.Printf("Don't need to re-DoBoardMove: %d %d \n", nMoves, gamT.Board.GetNMoves())
			}
			if len(err) != 0 {
				return err, trans, upPattTree
			}
			if (nMoves == 1) && (firstMoveColor == nodColr) {
				//					str := strconv.Itoa(int(nodColr))
				found = true
				newNodLoc, trans = gamT.FindCanonicalRep(nodLoc, ah.BoardHandicapSymmetry[ha])
				// check that nod has no siblings
				//					if nod.NextSib != nilTreeNodeIdx {
				//						err.Add(ah.NoPos, "AddTeachingPattern: unsupported sibling of first node")
				//						return
				//					}
				break findFirst
			}
		case SequenceNode:
			err.Add(ah.NoPos, "AddTeachingPattern: SequenceNode not supported "+strconv.Itoa(i))
		case TransferNode:
			err.Add(ah.NoPos, "AddTeachingPattern: TransferNode not supported "+strconv.Itoa(i))
		}
	}
again:
	if found {
		// traverse the gamTree and put in the pattTree
		for (curGam != nilTreeNodeIdx) && (limitReached == false) {
			// traverse tree via children links
			markBad := gamT.treeNodes[curGam].NextSib != curGam
			//			str := strconv.Itoa(int(nodColr))
			curPatt = findOrAdd(newNodLoc)
			if onMain && markBad {
				var pv PropertyValue
				// BM Bad Move
				pv.StrValue = []byte("1")
				pv.NextProp = nilPropIdx
				pv.PropType = BM_idx
				pv.ValType = Double
				// TODO: something missing here. pv not saved before overwritten.
				// TR Triangle
				pv.StrValue = SGFCoords(newNodLoc, gamT.IsFF4())
				pv.NextProp = nilPropIdx
				pv.PropType = TR_idx
				pv.ValType = ListOfPoint
				if newNodLoc != ah.PassNodeLoc {
					pattTree.AddAProp(curPatt, pv)
				}
			}
			if markGood {
				var pv PropertyValue
				// GB Good for Black or GW Good for White
				pv.StrValue = []byte("1")
				pv.NextProp = nilPropIdx
				if nodColr == ah.Black {
					pv.PropType = GB_idx
				} else {
					pv.PropType = GW_idx
				}
				pv.ValType = Double
				// TODO: something missing here. pv not saved before overwritten.
				// SQ Square
				pv.StrValue = SGFCoords(newNodLoc, gamT.IsFF4())
				pv.NextProp = nilPropIdx
				pv.PropType = SQ_idx
				pv.ValType = ListOfPoint
				if newNodLoc != ah.PassNodeLoc {
					pattTree.AddAProp(curPatt, pv)
				}
				markGood = false
			}
			// if curGam is the firstChild of parent:
			//	push siblings of curGam
			parent := gamT.treeNodes[curGam].Parent
			if parent != nilTreeNodeIdx {
				lastCh := gamT.treeNodes[parent].Children
				if lastCh != nilTreeNodeIdx {
					firstCh := gamT.treeNodes[lastCh].NextSib
					if curGam == firstCh {
						// push the Siblings
						for firstCh != lastCh {
							sib := gamT.treeNodes[firstCh].NextSib
							if sib != curGam {
								var newTraverseRec traversePoint
								newTraverseRec.cGam = sib
								newTraverseRec.cPat = curPatt
								newTraverseRec.mkGd = onMain
								newTraverseRec.pDep = patternDepth
								//								fmt.Printf("  pushing sib: %d curPatt: %d curGam: %d\n", sib, curPatt, curGam)
								traverseStack = append(traverseStack, newTraverseRec)
							}
							firstCh = sib
						}
					}
				}
			}

			// move down to next generation
			curGam = gamT.treeNodes[curGam].Children
			if curGam != nilTreeNodeIdx {
				curGam = gamT.treeNodes[curGam].NextSib // move to first child
				nodLoc, nodColr, err = gamT.GetMove(gamT.treeNodes[curGam])
				if len(err) != 0 {
					return
				}
				c, r := ah.GetColRow(nodLoc)
				newNodLoc = gamT.TransNodeLoc(trans, c, r)
			}
		}
	}

	// see if any stacked nodes to visit:
	if len(traverseStack) > 0 {
		var nxtTraverseRec = traverseStack[0]
		//		fmt.Printf("  popping curGam: %d curPatt: %d\n", curGam, curPatt)
		traverseStack = traverseStack[1:]
		curGam = nxtTraverseRec.cGam
		curPatt = nxtTraverseRec.cPat
		curPatt = pattTree.treeNodes[curPatt].Parent // move to parent
		markGood = nxtTraverseRec.mkGd
		patternDepth = nxtTraverseRec.pDep
		patternDepth -= 1 // decrement, due to move to parent
		limitReached = false
		onMain = false
		nodLoc, nodColr, err = gamT.GetMove(gamT.treeNodes[curGam])
		if len(err) != 0 {
			return
		}
		c, r := ah.GetColRow(nodLoc)
		newNodLoc = gamT.TransNodeLoc(trans, c, r)
		found = true
		goto again
	}

	upPattTree = pattTree
	return err, trans, upPattTree
}
