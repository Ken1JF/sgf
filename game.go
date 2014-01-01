/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/sgf/game.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 2/11/2010.
 *  Copyright 2010-2014 Ken Friedenbach. All rights reserved.
 *
 *	This file implements the data structures for storing a game,
 *	in a format suitable for reading/writing in .sgf format.
 */

package sgf

import (
    "strconv"
	"gitHub.com/Ken1JF/ahgo/ah"
)

// This Komi structure allows the notation KM[?] which indicates that a komi
// was given, but the value is unknown. The default is set == false. 
// If KM[xxx] appears in an SGF file, set == true. 
// If known == true, then val is the komi value. Otherwise, val is 0.0.
// KM[0] is used to indicate that no komi was given.
//
type Komi struct {
	val	float32
	set	bool
	known bool
}

// This Result structure supports the addition of comments in the RE[] field.
// Many of the comments differ only in the separators used, 
// or the number of moves known. So this structure allows the specification of
// a number that occurs in the comment, and a left separator: either "{" or "(".
// The boolean "both" indicates that a matching right separator was found.
// 
type Result struct {
	val []byte
	com []byte
	n int
	sep byte
	both bool
}

// Set functions for SGF properties:
//
func (gam* GameTree) SetFF(f []byte) {
	gam.fF = f
}

func (gam* GameTree) IsFF4() (ret bool) {
    if gam.fF == nil {
        ret = false
    } else {
        i, _ := strconv.Atoi(string(gam.fF))
        ret = (i == 4)
    }
	return ret
}

func (gam* GameTree) SetST(s []byte) {
	gam.sT = s
}

func (gam* GameTree) SetPB(p []byte) {
	gam.pB = p
}

func (gam* GameTree) GetPB() ( []byte) {
	return gam.pB
}

func (gam* GameTree) SetBR(p []byte) {
	gam.bR = p
}

func (gam* GameTree) GetBR() ([]byte) {
	return gam.bR
}

func (gam* GameTree) SetBT(p []byte) {
	gam.bT = p
}

func (gam* GameTree) GetPW() ([]byte) {
	return gam.pW
}

func (gam* GameTree) SetPW(p []byte) {
	gam.pW = p
}

func (gam* GameTree) SetWR(p []byte) {
	gam.wR = p
}

func (gam* GameTree) GetWR() ([]byte) {
	return gam.wR
}

func (gam* GameTree) SetWT(p []byte) {
	gam.wT = p
}

func (gam* GameTree) SetDT(p []byte) {
	gam.dT = p
}

func (gam* GameTree) SetPC(p []byte) {
	gam.pC = p
}

func (gam* GameTree) SetRU(p []byte) {
	gam.rU = p
}

func (gam* GameTree) SetRE(v []byte, b [] byte, num int, ch byte, two bool) {
	// TODO: set Board value
	// set SGF details
	gam.rE.val = v
	gam.rE.com = b
	gam.rE.n = num
	gam.rE.sep = ch
	gam.rE.both = two
}

func (gam* GameTree) SetGC(p []byte) {
	gam.gC = p
}

func (gam* GameTree) SetEV(p []byte) {
	gam.eV = p
}

func (gam* GameTree) SetRO(p []byte) {
	gam.rO = p
}

func (gam* GameTree) SetAP(p []byte) {
	gam.aP = p
}

func (gam* GameTree) SetAN(p []byte) {
	gam.aN = p
}

func (gam* GameTree) SetCP(p []byte) {
	gam.cP = p
}

func (gam* GameTree) SetSO(p []byte) {
	gam.sO = p
}

func (gam* GameTree) SetUS(p []byte) {
	gam.uS = p
}

func (gam* GameTree) SetKM(p float32, k bool) {
	// Set the Board value
	if k {
		gam.SetKomi(p)
	}
	// Set the SGF details
	gam.kM.set = true
	gam.kM.known = k
	if k {
		gam.kM.val = p
	}
}

func (gam* GameTree) SetTM(p float32) {
	gam.tM = p
}

func (gam* GameTree) DoAB(p ah.NodeLoc, doPlay bool) (err ah.ErrorList) {
	movType := ah.Unocc
	gam.aB = append(gam.aB, p)
	bp := &gam.Graphs[ah.PointLevel].Nodes[p]
	cur := bp.GetNodeLowState()
	if ah.IsOccupied(ah.PointStatus(cur)) {
		// TODO: check for white and add AB_B?
		movType = ah.AB_W
	} else {
		movType = ah.AB_U
	}
	if doPlay {
		gam.ChangeNodeState(ah.PointLevel, p, ah.NodeStatus(ah.Black), true)
	}
	_ = gam.AddMove(p, movType, 0, ah.NilNodeLoc)
	if doPlay {
		gam.EachAdjNode(ah.PointLevel, p, 
			func(adjNl ah.NodeLoc) {
				newSt := gam.Graphs[ah.PointLevel].CompHigh(&(gam.AbstHier), ah.PointLevel, adjNl, uint16(ah.Unocc))
				gam.ChangeNodeState(ah.PointLevel, adjNl, ah.NodeStatus(newSt), true)
			})
	}
	return err
}

func (gam* GameTree) DoAE(p ah.NodeLoc, doPlay bool) (err ah.ErrorList) {
	movType := ah.Unocc
//	gam.aE = append(gam.aE, p)
	bp := &gam.Graphs[ah.PointLevel].Nodes[p]
	cur := bp.GetNodeLowState()
	if ah.PointStatus(cur) == ah.Black {
		movType = ah.AE_B
	} else {
		// TODO: check for white and add AE_E?
		movType = ah.AE_W
	}
	if doPlay {
		newSt := gam.Graphs[ah.PointLevel].CompHigh(&(gam.AbstHier), ah.PointLevel, p, uint16(ah.Unocc))
		gam.ChangeNodeState(ah.PointLevel, p, ah.NodeStatus(newSt), true)
	}
	_ = gam.AddMove(p, movType, 0, ah.NilNodeLoc)
	return err
}

func (gam* GameTree) DoAW(p ah.NodeLoc, doPlay bool) (err ah.ErrorList) {
	movType := ah.Unocc
	gam.aW = append(gam.aW, p)
	bp := &gam.Graphs[ah.PointLevel].Nodes[p]
	cur := bp.GetNodeLowState()
	if ah.IsOccupied(ah.PointStatus(cur)) {
		// TODO: check for black and add AW_W?
		movType = ah.AW_B
	} else {
		movType = ah.AW_U
	}
	if doPlay {
		gam.ChangeNodeState(ah.PointLevel, p, ah.NodeStatus(ah.White), true)
	}
	_ = gam.AddMove(p, movType, 0, ah.NilNodeLoc)
	if doPlay {
		gam.EachAdjNode(ah.PointLevel, p, 
			func(adjNl ah.NodeLoc) {
				newSt := gam.Graphs[ah.PointLevel].CompHigh(&(gam.AbstHier), ah.PointLevel, adjNl, uint16(ah.Unocc))
				gam.ChangeNodeState(ah.PointLevel, adjNl, ah.NodeStatus(newSt), true)
			})
	}
	return err
}

func (gam* GameTree) DoB(nl ah.NodeLoc, doPlay bool) (movN int, err ah.ErrorList) {
	movN, err = gam.DoBoardMove(nl, ah.Black, doPlay)
	if movN == 1 {
		gam.SetPlayerRank()
	}
	return movN, err
}

func (gam* GameTree) DoW(nl ah.NodeLoc, doPlay bool) (movN int, err ah.ErrorList) {
	movN, err = gam.DoBoardMove(nl, ah.White, doPlay)
	if movN == 1 {
		gam.SetPlayerRank()
	}
	return movN, err
}

func (gam* GameTree) SetOH(p []byte) {
	gam.oH = p
}

func (gam* GameTree) GetOH() ( []byte) {
	return gam.oH
}

func (gam* GameTree) SetHA(p int) {
	gam.SetHandicap(p)
}

func (gam* GameTree) GetHA() ( int) {
	return gam.GetHandicap()
}

// PlaceHandicap sets the handicap stones, and returns the list of points
//
func (gam *GameTree) PlaceHandicap(n int, siz int) (pts []uint8) {
	var lin, mid int
//	var nl ah.NodeLoc
//	var play bool = true
	place := func(c int, r int) {
		nl := ah.MakeNodeLoc(ah.ColValue(c), ah.RowValue(r))
		gam.DoAB(nl, true)
		pts = append(pts, SGFCoords(nl, gam.IsFF4())...)
	}
	lin = 3
	if siz < 13 {
		lin = 2
	}
	mid = (siz-1)/2
	switch n {
		case 0:
		case 2: 
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
		case 3:
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
		case 4:
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
		case 5: 
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
			place(mid, mid)					// mid point
		case 6: 
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
			place(lin, mid)					// mid point Left side
			place(siz-(lin+1), mid)			// mid point Right side
		case 7: 
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
			place(lin, mid)					// mid point Left side
			place(siz-(lin+1), mid)			// mid point Right side
			place(mid, mid)					// mid point
		case 8:
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
			place(lin, mid)					// mid point Left side
			place(siz-(lin+1), mid)			// mid point Right side
			place(mid, lin)					// mid point Top side
			place(mid, siz-(lin+1))			// mid point Bottom side
		case 9:
			place(siz-(lin+1), lin)			// UpperRight
			place(lin, siz-(lin+1))			// LowerLeft
			place(siz-(lin+1), siz-(lin+1))	// LowerRight
			place(lin, lin)					// UpperLeft
			place(lin, mid)					// mid point Left side
			place(siz-(lin+1), mid)			// mid point Right side
			place(mid, lin)					// mid point Top side
			place(mid, siz-(lin+1))			// mid point Bottom side
			place(mid, mid)					// mid point
	}
	return pts
}


func (gam* GameTree) DoAR(p []byte) {
	// TODO: implement Arrows
	// string format is:
	//	t1:h1t2:h2 ... tN:hN
	// where ti is the tail of arrow i
	// and hi is the head of arrow i
}

func (gam* GameTree) DoLN(p []byte) {
	// TODO: implement Lines
	// string format is:
	//	t1:h1t2:h2 ... tN:hN
	// where ti is the start of line i
	// and hi is the end of line i
}

// CheckProperties is called after parsing an SGF file.
//
func (gam* GameTree) CheckProperties(gogod bool) (errstr string) {

	// Check FF[]
	if (gam.fF == nil) && gogod {	// treat as an error for GoGoD feedback
		errstr = "No FF "
	}

	// Check SZ[]
	cSz, _ := gam.GetSize()
	if (cSz == 0) && gogod {	// treat as an error for GoGoD feedback
		errstr = errstr + "No SZ"
	}
	
	// TODO: Check GM[]

	// Check HA[] with AB[], AW[], and mov1.
	// Assumes B is taking handicap, and W playing first.
	// (This logic is complicated by at least one GoGoD game,
	// where B gets HA, and makes first move, on another handicap point. 
	// Check how many, and consider "correcting" the game record(s).)
	// TODO: check PL?
	handi := gam.GetHandicap()
	mov1, set := gam.GetMov1()
	if handi > 0 {
		if set {
			if mov1 == ah.Black {
				if handi != (len(gam.aW) - len(gam.aB)) &&
					(handi != (len(gam.aB) - len(gam.aW) +1)) {
					errstr = errstr + "HA not equal AW - AB (1.B)"
				}
			} else { // mov1 == ah.White
				if handi != (len(gam.aB) - len(gam.aW)) {
					errstr = errstr + "HA not equal AB - AW (1.W)"
				}
			}
		} else { // assume Black to play first
			if handi != (len(gam.aB) - len(gam.aW)) {
				errstr = errstr + "HA not equal AB - AW (1.B?)"
			}
			// else ok
		}
	} else {  // hA == 0
		if len(gam.aB) != 0 {
			if len(gam.aW) != len(gam.aB) {
				if set {
					if mov1 == ah.Black {
						if len(gam.aB)+1 != len(gam.aW) {
							errstr = errstr + "AB+1 not equal AW"
						}
						// else ok
					} else { // mov1 == ah.White
						if len(gam.aB) != len(gam.aW) +1 {
							errstr = errstr + "AB not equal AW + 1"
						}
						// else ok
					}
				} else {
					errstr = errstr + "AB not equal AW (mov1 not set)"
				}
			}
		} else if len(gam.aW) != 0 {
			errstr = errstr + "AW not zero"
		}
	}
	
	// TODO: check OH (GoGoD specific property) for consistency with HA, and ranks.
	
	// TODO: check RE with evaluation of final position.
	
	return errstr
}

// GetMove returns the move at a node
//
func (gamT *GameTree)GetMove(n TreeNode) (nl ah.NodeLoc, c ah.PointStatus, err ah.ErrorList) {
	if n.TNodType == BlackMoveNode {
		nl, c = ah.NodeLoc(n.propListOrNodeLoc), ah.Black
	} else if n.TNodType == WhiteMoveNode {
		nl, c = ah.NodeLoc(n.propListOrNodeLoc), ah.White
	} else if n.TNodType == InteriorNode {
		OK := false
		lastProp := n.propListOrNodeLoc
		if lastProp != nilPropIdx {
			pl := gamT.propertyValues[lastProp].NextProp
			prop := gamT.propertyValues[pl]
			process := func (prop PropertyValue) {
				if prop.PropType == B_idx {
					OK = true
					c = ah.Black
					nl, err = SGFPoint(prop.StrValue)
				} else if prop.PropType == W_idx {
					OK = true
					c = ah.White
					nl, err =  SGFPoint(prop.StrValue)
				} else {
					pl = prop.NextProp
				}
			}
			// pprocess first prop.
			process(prop)
			if !OK {
				for (pl != lastProp) && (err == nil) {
					pl := gamT.propertyValues[pl].NextProp
					prop := gamT.propertyValues[pl]
					// pprocess next prop.
					process(prop)
					if OK {
						break
					}
				}
			}
		}
		if !OK {
			err.Add(ah.NoPos, "sgf/GetMove: move property not found in interior node")
		}
	} else {
		// TODO: support moves with other properties
		nl = ah.IllegalNodeLoc
		c = ah.Unocc
		err.Add(ah.NoPos, "sgf/GetMove: not a move node or interior node")
	}
	return nl, c, err
}
