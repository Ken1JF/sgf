package sgf_test

import (
	"fmt"
	"gitHub.com/Ken1JF/ah"
	. "gitHub.com/Ken1JF/sgf"
	"strconv"
)

// Eight boards of various sizes.
var brd_5, brd_7, brd_9, brd_11, brd_13, brd_15, brd_17, brd_19 *ah.AbstHier

// and an array to hold them.
var brds [8]*ah.AbstHier

// Transformation test data:
// For these tests, use char values instead of defined PointStatus values
var test_5 = []string{
	"1....",
	"2..x.",
	"3.+..",
	"4....",
	"5....",
}

var test_7 = []string{
	"1......",
	"2......",
	"3.+.+..",
	"4...x..",
	"5.+.+..",
	"6......",
	"7......",
}

var test_9 = []string{
	"1........",
	"2........",
	"3.+...+..",
	"4.....x..",
	"5........",
	"6........",
	"7.+...+..",
	"8........",
	"9........",
}

var test_11 = []string{
	"1..........",
	"2..........",
	"3.+.....+..",
	"4.......x..",
	"5..........",
	"6....+.....",
	"7..........",
	"8..........",
	"9.+.....+..",
	"A..........",
	"B..........",
}

var test_13 = []string{
	"1............",
	"2............",
	"3............",
	"4..+.....+X..",
	"5............",
	"6............",
	"7....+.......",
	"8............",
	"9............",
	"A..+.....+...",
	"B............",
	"C............",
	"D............",
}

var test_15 = []string{
	"1..............",
	"2..............",
	"3..............",
	"4..+.......+X..",
	"5..............",
	"6..............",
	"7..............",
	"8.....+........",
	"9..............",
	"A..............",
	"B..............",
	"C..+.......+...",
	"D..............",
	"E..............",
	"F..............",
}

var test_17 = []string{
	"1................",
	"2................",
	"3................",
	"4..+.........+X..",
	"5................",
	"6................",
	"7................",
	"8................",
	"9......+.........",
	"A................",
	"B................",
	"C................",
	"D................",
	"E..+.........+...",
	"F................",
	"G................",
	"H................",
}

var test_19 = []string{
	"1..................",
	"2..................",
	"3..................",
	"4..+.....+.....+X..",
	"5..................",
	"6..................",
	"7..................",
	"8..................",
	"9..................",
	"A..+.....+.....+...",
	"B..................",
	"C..................",
	"D..................",
	"E..................",
	"F..................",
	"G..+.....+.....+...",
	"H..................",
	"I..................",
	"J..................",
}

// printInitBoard prints the PointType values
// after a Board is initialized (via SetSize)
func printInitBoard(abhr *ah.AbstHier, title string) {

	//	Black_Occ_Pt:		"◉",
	//	White_Occ_Pt:		"◎",

	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := abhr.GetSize()
	fmt.Println(title, "Board", int(nCol), "by", int(nRow))
	for r = 0; ah.RowSize(r) < nRow; r++ {
		for c = 0; ah.ColSize(c) < nCol; c++ {
			bp := abhr.Graphs[ah.PointLevel].GetPoint(c, r)
			hs := bp.GetNodeHighState()
			if hs == uint16(ah.White) {
				fmt.Print("◎")
			} else if hs == uint16(ah.Black) {
				fmt.Print("◉")
			} else {
				fmt.Print(ah.PtTypeNames[bp.GetPointType()])
			}
		}
		fmt.Println()
	}
}

// printInitBoard2 is equivalent to printInitBoard
// but uses the iteration function ah.EachNode
// and a literal func.
func printInitBoard2(abhr *ah.AbstHier) {
	var row ah.RowValue = 0
	nCol, nRow := abhr.GetSize()
	fmt.Println("Board", int(nCol), "by", int(nRow))
	abhr.EachNode(ah.PointLevel,
		func(brd *ah.Graph, nl ah.NodeLoc) {
			_, r := brd.Nodes[nl].GetPointColRow()
			if r != row {
				fmt.Println()
				row = r
			}
			fmt.Print(ah.PtTypeNames[brd.Nodes[nl].GetPointType()])
		})
	fmt.Println()
}

// Print the boards, after transformation
func printBrds(msg string, brd *ah.AbstHier, newBrd *ah.AbstHier, tName string) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd.GetSize()
	fmt.Println("Board size", int(nCol), "by", int(nRow), "after", tName)
	for r = 0; ah.RowSize(r) < nRow; r++ {
		for c = 0; ah.ColSize(c) < nCol; c++ {
			bp := brd.Graphs[ah.PointLevel].GetPoint(c, r)
			ch := bp.GetNodeLowState()
			fmt.Printf("%c", byte(ch))
		}
		fmt.Print(" | ")
		for c = 0; ah.ColSize(c) < nCol; c++ {
			nbp := newBrd.Graphs[ah.PointLevel].GetPoint(c, r)
			ch := nbp.GetNodeLowState()
			fmt.Printf("%c", byte(ch))
		}
		fmt.Println()
	}
}

// SetUpTestBoard stores the test data (string characters)
// in the Board as PointStatus information.
func SetUpTestBoard(N int, brd *ah.AbstHier, data *[]string) {
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			brd.SetPoint(ah.MakeNodeLoc(ah.ColValue(c), ah.RowValue(r)), ah.PointStatus((*data)[r][c]))
		}
	}
}

// Test the transformation logic
func ExampleTestTrans() {
	// Set up the test data boards.
	var col ah.ColSize
	var row ah.RowSize
	for size := 5; size <= 19; size += 2 {
		switch size {
		case 5:
			col = 5
			row = 5
			brd_5 = brd_5.InitAbstHier(col, row, ah.StringLevel, true)
			//				ah.SetAHTrace(false)
			printInitBoard(brd_5, "Initial 5x5 Board")
			brd_5.PrintAbstHier("Initial 5x5 Board", true)
			SetUpTestBoard(size, brd_5, &test_5)
			brds[0] = brd_5
		case 7:
			col = 7
			row = 7
			//				brd_7 = new(ah.AbstHier)
			//				brd_7.SetSize(col, row)
			brd_7 = brd_7.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_7)
			SetUpTestBoard(size, brd_7, &test_7)
			brds[1] = brd_7
		case 9:
			col = 9
			row = 9
			//				brd_9 = new(ah.AbstHier)
			//				brd_9.SetSize(col, row)
			brd_9 = brd_9.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_9, "Initial 9x9 Board")
			SetUpTestBoard(size, brd_9, &test_9)
			brds[2] = brd_9
		case 11:
			col = 11
			row = 11
			//				brd_11 = new(ah.AbstHier)
			//				brd_11.SetSize(col, row)
			brd_11 = brd_11.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_11)
			SetUpTestBoard(size, brd_11, &test_11)
			brds[3] = brd_11
		case 13:
			col = 13
			row = 13
			brd_13 = brd_13.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_13, "Initial 13x13 Board")
			SetUpTestBoard(size, brd_13, &test_13)
			brds[4] = brd_13
		case 15:
			col = 15
			row = 15
			brd_15 = brd_15.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_15)
			SetUpTestBoard(size, brd_15, &test_15)
			brds[5] = brd_15
		case 17:
			col = 17
			row = 17
			brd_17 = brd_17.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard(brd_17, "Initial 17x17 Board")
			SetUpTestBoard(size, brd_17, &test_17)
			brds[6] = brd_17
		case 19:
			col = 19
			row = 19
			brd_19 = brd_19.InitAbstHier(col, row, ah.StringLevel, true)
			printInitBoard2(brd_19)
			SetUpTestBoard(size, brd_19, &test_19)
			brds[7] = brd_19
		}
	}
	// Print each board, after applying one of the transformations,
	// and print it (for visual verification)
	//	ah.SetAHTrace(true) // trace first one
	for i, brd := range brds {
		fmt.Println("Checking brds[", i, "]")
		if brd == nil {
			fmt.Println("Error in setup: brd == nil")
		} else {
			newBrd := brd.TransBoard(ah.BoardTrans(i))
			printBrds("Visual Check", brd, newBrd, ah.TransName[i])
		}
		ah.SetAHTrace(false) // turn off after first one
	}
	// Verify that the inverse transformations produce the original
	for i, brd := range brds {
		t := ah.BoardTrans(i)
		inv := ah.InverseTrans[t]
		fmt.Println("Checking", ah.TransName[i], "and its inverse:", ah.TransName[inv])
		newBrd := brd.TransBoard(t)
		newBrdInv := newBrd.TransBoard(inv)
		if differBrds(brd, newBrdInv) {
			printBrds("Error: inverse differs", brd, newBrdInv, ah.TransName[i])
		}
	}
	// Verify the transformation composition table
	nxtBrd := 0 // used to pick the next board
	for A := ah.T_FIRST; A <= ah.T_LAST; A++ {
		for B := ah.T_FIRST; B <= ah.T_LAST; B++ {
			C := ah.ComposeTrans[A][B]
			fmt.Println("Checking", ah.TransName[C], "=", ah.TransName[A], "*", ah.TransName[B])
			brd := brds[nxtBrd]
			nxtBrd++
			if nxtBrd >= 8 {
				nxtBrd = 0
			}
			brdA := brd.TransBoard(A)
			brdAB := brdA.TransBoard(B)
			brdC := brd.TransBoard(C)
			if differBrds(brdAB, brdC) {
				printBrds("Error: "+ah.TransName[ah.ComposeTrans[A][B]], brdAB, brdC,
					"not equal"+ah.TransName[A]+"*"+ah.TransName[B])
			}
		}
	}
	// Output:
	// Initial 5x5 Board Board 5 by 5
	// ┏┯┯┯┓
	// ┠╬┼╬┨
	// ┠┼◘┼┨
	// ┠╬┼╬┨
	// ┗┷┷┷┛
	// Abstraction Hierarchy: Initial 5x5 Board
	// Level 1
	// Black nodes
	// Total 0 nodes, with 0 members
	// White nodes
	// Total 0 nodes, with 0 members
	// Unocc nodes
	// 0:202,3 1-mem:(E,5):202:202,adj:16(1),8(1),
	// 1:778,3 1-mem:(A,1):778:778,adj:4(1),2(1),
	// 2:906,3 3-mem:(D,1):906:906,(C,1),(B,1),adj:1(1),7(1),6(1),5(1),3(1),
	// 3:394,3 1-mem:(E,1):394:394,adj:2(1),8(1),
	// 4:842,3 3-mem:(A,4):842:842,(A,3),(A,2),adj:1(1),15(1),12(1),9(1),5(1),
	// 5:3018,3 1-mem:(B,2):3018:3018,adj:4(1),2(1),9(1),6(1),
	// 6:4042,3 1-mem:(C,2):4042:4042,adj:5(1),2(1),10(1),7(1),
	// 7:3018,3 1-mem:(D,2):3018:3018,adj:6(1),2(1),11(1),8(1),
	// 8:458,3 3-mem:(E,4):458:458,(E,3),(E,2),adj:0(1),7(1),3(1),14(1),11(1),
	// 9:4042,3 1-mem:(B,3):4042:4042,adj:4(1),5(1),12(1),10(1),
	// 10:1994,3 1-mem:(C,3):1994:1994,adj:9(1),6(1),13(1),11(1),
	// 11:4042,3 1-mem:(D,3):4042:4042,adj:8(1),10(1),7(1),14(1),
	// 12:3018,3 1-mem:(B,4):3018:3018,adj:4(1),9(1),16(1),13(1),
	// 13:4042,3 1-mem:(C,4):4042:4042,adj:12(1),10(1),16(1),14(1),
	// 14:3018,3 1-mem:(D,4):3018:3018,adj:8(1),13(1),11(1),16(1),
	// 15:586,3 1-mem:(A,5):586:586,adj:4(1),16(1),
	// 16:714,3 3-mem:(D,5):714:714,(C,5),(B,5),adj:0(1),14(1),13(1),15(1),12(1),
	// Total 17 nodes, with 25 members
	// Board 7 by 7
	// ┏┯┯┯┯┯┓
	// ┠╬┼┼┼╬┨
	// ┠┼◘┼◘┼┨
	// ┠┼┼╋┼┼┨
	// ┠┼◘┼◘┼┨
	// ┠╬┼┼┼╬┨
	// ┗┷┷┷┷┷┛
	// Initial 9x9 Board Board 9 by 9
	// ┏┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼╬┨
	// ┠┼◘┼┼┼◘┼┨
	// ┠┼┼╬┼╬┼┼┨
	// ┠┼┼┼◘┼┼┼┨
	// ┠┼┼╬┼╬┼┼┨
	// ┠┼◘┼┼┼◘┼┨
	// ┠╬┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┛
	// Board 11 by 11
	// ┏┯┯┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼┼┼╬┨
	// ┠┼◘┼┼┼┼┼◘┼┨
	// ┠┼┼╬┼┼┼╬┼┼┨
	// ┠┼┼┼╬┼╬┼┼┼┨
	// ┠┼┼┼┼◘┼┼┼┼┨
	// ┠┼┼┼╬┼╬┼┼┼┨
	// ┠┼┼╬┼┼┼╬┼┼┨
	// ┠┼◘┼┼┼┼┼◘┼┨
	// ┠╬┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┛
	// Initial 13x13 Board Board 13 by 13
	// ┏┯┯┯┯┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼┼┼┼┼╬┨
	// ┠┼╬┼┼┼┼┼┼┼╬┼┨
	// ┠┼┼◘┼┼┼┼┼◘┼┼┨
	// ┠┼┼┼╬┼┼┼╬┼┼┼┨
	// ┠┼┼┼┼╬┼╬┼┼┼┼┨
	// ┠┼┼┼┼┼◘┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼╬┼┼┼┨
	// ┠┼┼◘┼┼┼┼┼◘┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┛
	// Board 15 by 15
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠┼┼◘┼┼┼◘┼┼┼◘┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼╬┼┼┼┼┼┨
	// ┠┼┼◘┼┼┼◘┼┼┼◘┼┼┨
	// ┠┼┼┼┼┼╬┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◘┼┼┼◘┼┼┼◘┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	// Initial 17x17 Board Board 17 by 17
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠┼┼◘┼┼┼┼◘┼┼┼┼◘┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼◘┼┼┼╋◘╋┼┼┼◘┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◘┼┼┼┼◘┼┼┼┼◘┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	// Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┓
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠┼┼◘┼┼┼┼┼◘┼┼┼┼┼◘┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼◘┼┼┼╋╋◘╋╋┼┼┼◘┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◘┼┼┼┼┼◘┼┼┼┼┼◘┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	// Checking brds[ 0 ]
	// Board size 5 by 5 after T_IDENTITY
	// 1.... | 1....
	// 2..x. | 2..x.
	// 3.+.. | 3.+..
	// 4.... | 4....
	// 5.... | 5....
	// Checking brds[ 1 ]
	// Board size 7 by 7 after T_ROTA_090
	// 1...... | .......
	// 2...... | .......
	// 3.+.+.. | ..+x+..
	// 4...x.. | .......
	// 5.+.+.. | ..+.+..
	// 6...... | .......
	// 7...... | 1234567
	// Checking brds[ 2 ]
	// Board size 9 by 9 after T_ROTA_180
	// 1........ | ........9
	// 2........ | ........8
	// 3.+...+.. | ..+...+.7
	// 4.....x.. | ........6
	// 5........ | ........5
	// 6........ | ..x.....4
	// 7.+...+.. | ..+...+.3
	// 8........ | ........2
	// 9........ | ........1
	// Checking brds[ 3 ]
	// Board size 11 by 11 after T_ROTA_270
	// 1.......... | BA987654321
	// 2.......... | ...........
	// 3.+.....+.. | ..+.....+..
	// 4.......x.. | ...........
	// 5.......... | ...........
	// 6....+..... | .....+.....
	// 7.......... | ...........
	// 8.......... | ...........
	// 9.+.....+.. | ..+....x+..
	// A.......... | ...........
	// B.......... | ...........
	// Checking brds[ 4 ]
	// Board size 13 by 13 after T_FLP_SLAS
	// 1............ | .............
	// 2............ | .............
	// 3............ | .........X...
	// 4..+.....+X.. | ...+.....+...
	// 5............ | .............
	// 6............ | .............
	// 7....+....... | .............
	// 8............ | ......+......
	// 9............ | .............
	// A..+.....+... | ...+.....+...
	// B............ | .............
	// C............ | .............
	// D............ | DCBA987654321
	// Checking brds[ 5 ]
	// Board size 15 by 15 after T_FLP_VERT
	// 1.............. | ..............1
	// 2.............. | ..............2
	// 3.............. | ..............3
	// 4..+.......+X.. | ..X+.......+..4
	// 5.............. | ..............5
	// 6.............. | ..............6
	// 7.............. | ..............7
	// 8.....+........ | ........+.....8
	// 9.............. | ..............9
	// A.............. | ..............A
	// B.............. | ..............B
	// C..+.......+... | ...+.......+..C
	// D.............. | ..............D
	// E.............. | ..............E
	// F.............. | ..............F
	// Checking brds[ 6 ]
	// Board size 17 by 17 after T_FLP_BACK
	// 1................ | 123456789ABCDEFGH
	// 2................ | .................
	// 3................ | .................
	// 4..+.........+X.. | ...+.........+...
	// 5................ | .................
	// 6................ | .................
	// 7................ | .................
	// 8................ | ........+........
	// 9......+......... | .................
	// A................ | .................
	// B................ | .................
	// C................ | .................
	// D................ | .................
	// E..+.........+... | ...+.........+...
	// F................ | ...X.............
	// G................ | .................
	// H................ | .................
	// Checking brds[ 7 ]
	// Board size 19 by 19 after T_FLP_HORI
	// 1.................. | J..................
	// 2.................. | I..................
	// 3.................. | H..................
	// 4..+.....+.....+X.. | G..+.....+.....+...
	// 5.................. | F..................
	// 6.................. | E..................
	// 7.................. | D..................
	// 8.................. | C..................
	// 9.................. | B..................
	// A..+.....+.....+... | A..+.....+.....+...
	// B.................. | 9..................
	// C.................. | 8..................
	// D.................. | 7..................
	// E.................. | 6..................
	// F.................. | 5..................
	// G..+.....+.....+... | 4..+.....+.....+X..
	// H.................. | 3..................
	// I.................. | 2..................
	// J.................. | 1..................
	// Checking T_IDENTITY and its inverse: T_IDENTITY
	// Checking T_ROTA_090 and its inverse: T_ROTA_270
	// Checking T_ROTA_180 and its inverse: T_ROTA_180
	// Checking T_ROTA_270 and its inverse: T_ROTA_090
	// Checking T_FLP_SLAS and its inverse: T_FLP_SLAS
	// Checking T_FLP_VERT and its inverse: T_FLP_VERT
	// Checking T_FLP_BACK and its inverse: T_FLP_BACK
	// Checking T_FLP_HORI and its inverse: T_FLP_HORI
	// Checking T_IDENTITY = T_IDENTITY * T_IDENTITY
	// Checking T_ROTA_090 = T_IDENTITY * T_ROTA_090
	// Checking T_ROTA_180 = T_IDENTITY * T_ROTA_180
	// Checking T_ROTA_270 = T_IDENTITY * T_ROTA_270
	// Checking T_FLP_SLAS = T_IDENTITY * T_FLP_SLAS
	// Checking T_FLP_VERT = T_IDENTITY * T_FLP_VERT
	// Checking T_FLP_BACK = T_IDENTITY * T_FLP_BACK
	// Checking T_FLP_HORI = T_IDENTITY * T_FLP_HORI
	// Checking T_ROTA_090 = T_ROTA_090 * T_IDENTITY
	// Checking T_ROTA_180 = T_ROTA_090 * T_ROTA_090
	// Checking T_ROTA_270 = T_ROTA_090 * T_ROTA_180
	// Checking T_IDENTITY = T_ROTA_090 * T_ROTA_270
	// Checking T_FLP_HORI = T_ROTA_090 * T_FLP_SLAS
	// Checking T_FLP_SLAS = T_ROTA_090 * T_FLP_VERT
	// Checking T_FLP_VERT = T_ROTA_090 * T_FLP_BACK
	// Checking T_FLP_BACK = T_ROTA_090 * T_FLP_HORI
	// Checking T_ROTA_180 = T_ROTA_180 * T_IDENTITY
	// Checking T_ROTA_270 = T_ROTA_180 * T_ROTA_090
	// Checking T_IDENTITY = T_ROTA_180 * T_ROTA_180
	// Checking T_ROTA_090 = T_ROTA_180 * T_ROTA_270
	// Checking T_FLP_BACK = T_ROTA_180 * T_FLP_SLAS
	// Checking T_FLP_HORI = T_ROTA_180 * T_FLP_VERT
	// Checking T_FLP_SLAS = T_ROTA_180 * T_FLP_BACK
	// Checking T_FLP_VERT = T_ROTA_180 * T_FLP_HORI
	// Checking T_ROTA_270 = T_ROTA_270 * T_IDENTITY
	// Checking T_IDENTITY = T_ROTA_270 * T_ROTA_090
	// Checking T_ROTA_090 = T_ROTA_270 * T_ROTA_180
	// Checking T_ROTA_180 = T_ROTA_270 * T_ROTA_270
	// Checking T_FLP_VERT = T_ROTA_270 * T_FLP_SLAS
	// Checking T_FLP_BACK = T_ROTA_270 * T_FLP_VERT
	// Checking T_FLP_HORI = T_ROTA_270 * T_FLP_BACK
	// Checking T_FLP_SLAS = T_ROTA_270 * T_FLP_HORI
	// Checking T_FLP_SLAS = T_FLP_SLAS * T_IDENTITY
	// Checking T_FLP_VERT = T_FLP_SLAS * T_ROTA_090
	// Checking T_FLP_BACK = T_FLP_SLAS * T_ROTA_180
	// Checking T_FLP_HORI = T_FLP_SLAS * T_ROTA_270
	// Checking T_IDENTITY = T_FLP_SLAS * T_FLP_SLAS
	// Checking T_ROTA_090 = T_FLP_SLAS * T_FLP_VERT
	// Checking T_ROTA_180 = T_FLP_SLAS * T_FLP_BACK
	// Checking T_ROTA_270 = T_FLP_SLAS * T_FLP_HORI
	// Checking T_FLP_VERT = T_FLP_VERT * T_IDENTITY
	// Checking T_FLP_BACK = T_FLP_VERT * T_ROTA_090
	// Checking T_FLP_HORI = T_FLP_VERT * T_ROTA_180
	// Checking T_FLP_SLAS = T_FLP_VERT * T_ROTA_270
	// Checking T_ROTA_270 = T_FLP_VERT * T_FLP_SLAS
	// Checking T_IDENTITY = T_FLP_VERT * T_FLP_VERT
	// Checking T_ROTA_090 = T_FLP_VERT * T_FLP_BACK
	// Checking T_ROTA_180 = T_FLP_VERT * T_FLP_HORI
	// Checking T_FLP_BACK = T_FLP_BACK * T_IDENTITY
	// Checking T_FLP_HORI = T_FLP_BACK * T_ROTA_090
	// Checking T_FLP_SLAS = T_FLP_BACK * T_ROTA_180
	// Checking T_FLP_VERT = T_FLP_BACK * T_ROTA_270
	// Checking T_ROTA_180 = T_FLP_BACK * T_FLP_SLAS
	// Checking T_ROTA_270 = T_FLP_BACK * T_FLP_VERT
	// Checking T_IDENTITY = T_FLP_BACK * T_FLP_BACK
	// Checking T_ROTA_090 = T_FLP_BACK * T_FLP_HORI
	// Checking T_FLP_HORI = T_FLP_HORI * T_IDENTITY
	// Checking T_FLP_SLAS = T_FLP_HORI * T_ROTA_090
	// Checking T_FLP_VERT = T_FLP_HORI * T_ROTA_180
	// Checking T_FLP_BACK = T_FLP_HORI * T_ROTA_270
	// Checking T_ROTA_090 = T_FLP_HORI * T_FLP_SLAS
	// Checking T_ROTA_180 = T_FLP_HORI * T_FLP_VERT
	// Checking T_ROTA_270 = T_FLP_HORI * T_FLP_BACK
	// Checking T_IDENTITY = T_FLP_HORI * T_FLP_HORI
}

// differBrds checks the LowStates of the Nodes
// only suitable for special set boards
func differBrds(brd1, brd2 *ah.AbstHier) (ret bool) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd1.GetSize()
	nCol2, nRow2 := brd2.GetSize()
	if (nCol != nCol2) || (nRow != nRow2) {
		ret = true
	} else {
		for r = 0; ah.RowSize(r) < nRow; r++ {
			for c = 0; ah.ColSize(c) < nCol; c++ {
				nl := ah.MakeNodeLoc(c, r)
				bp1 := &brd1.Graphs[ah.PointLevel].Nodes[nl]
				bp2 := &brd2.Graphs[ah.PointLevel].Nodes[nl]
				if bp1.GetNodeLowState() != bp2.GetNodeLowState() {
					ret = true
					break
				}
			}
		}
	}
	return ret
}

// checkHandicapBrds checks the LowStates of the Nodes
// only suitable for special set boards
func checkHandicapBrds(brd1, brd2 *ah.AbstHier) (ret bool) {
	var c ah.ColValue
	var r ah.RowValue
	nCol, nRow := brd1.GetSize()
	nCol2, nRow2 := brd2.GetSize()
	if (nCol != nCol2) || (nRow != nRow2) {
		ret = true
	} else {
		for r = 0; ah.RowSize(r) < nRow; r++ {
			for c = 0; ah.ColSize(c) < nCol; c++ {
				nl := ah.MakeNodeLoc(c, r)
				bp1 := &brd1.Graphs[ah.PointLevel].Nodes[nl]
				bp2 := &brd2.Graphs[ah.PointLevel].Nodes[nl]
				low1 := bp1.GetNodeLowState()
				low2 := bp2.GetNodeLowState()
				// check that both are occupied or unoccupied
				if ah.IsOccupied(ah.PointStatus(low1)) != ah.IsOccupied(ah.PointStatus(low2)) {
					ret = true
					break
				}
			}
		}
	}
	return ret
}

// checkHandicapCanonical
func checkHandicapCanonical() {
	// Verify that the handicap patterns are preserved by transformaions,
	for i, brd := range brds {
		t := ah.BoardTrans(i)
		inv := ah.InverseTrans[t]
		fmt.Println("Checking", ah.TransName[i], "and its inverse:", ah.TransName[inv])
		if brd != nil {
			newBrd := brd.TransBoard(t)
			newBrdInv := newBrd.TransBoard(inv)
			if differBrds(brd, newBrdInv) {
				printBrds("Error: inverse differs", brd, newBrdInv, ah.TransName[i])
			}
		} else {
			fmt.Println("Error: brds [", i, "] has not been initialized.")
		}
	}
	// Output: look around
	// its all around you
}

// ExampleCannonicalHandicap points]
func ExampleCannonicalHandicap() {
	for ha := 0; ha <= 9; ha++ {
		if ha != 1 {
			var gam *GameTree = new(GameTree)
			gam.InitAbstHier(19, 19, ah.StringLevel, true)
			gam.SetHandicap(ha)
			gam.PlaceHandicap(ha, 19)
			for r := 0; r < 19; r++ {
				for c := 0; c < 19; c++ {
					nl := ah.MakeNodeLoc(ah.ColValue(c), ah.RowValue(r))
					bp := &gam.Graphs[ah.PointLevel].Nodes[nl]
					if bp.GetNodeLowState() != uint16(ah.Black) {
						if gam.IsCanonical(nl, ah.BoardHandicapSymmetry[ha]) {
							bp.SetNodeHighState(uint16(ah.White))
							//gam.SetPoint(nl, ah.White)
						}
					}
				}
			}
			str := "Handicap pattern " + strconv.Itoa(ha)
			printInitBoard(&gam.AbstHier, str)
			for trans := ah.T_FIRST; trans <= ah.T_LAST; trans += 1 {
				newBrd := gam.AbstHier.TransBoard(trans)
				if checkHandicapBrds(&gam.AbstHier, newBrd) {
					fmt.Print(" false, /* ", ah.TransName[trans], " */")
					// TODO: replace?					printBrds("Error: inverse differs", &gam.AbstHier, newBrd, ah.TransName[trans])
				} else {
					fmt.Print(" true, /* ", ah.TransName[trans], " */")
				}
			}
			fmt.Println()
		}
	}
	checkHandicapCanonical()
	// Output:
	// Handicap pattern 0 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◘┼┼┼┼┼◘┼┼┼┼┼◎◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◘┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◘┼┼┼┼┼◘┼┼┼┼┼◘┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ true, /* T_ROTA_090 */ true, /* T_ROTA_180 */ true, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ true, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 2 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◘┼┼┼┼┼◘┼┼┼┼┼◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◘┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◎◎◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷◎
	//  true, /* T_IDENTITY */ false, /* T_ROTA_090 */ true, /* T_ROTA_180 */ false, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ false, /* T_FLP_VERT */ true, /* T_FLP_BACK */ false, /* T_FLP_HORI */
	// Handicap pattern 3 Board 19 by 19
	// ◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◎◎◎◎◎◎◎◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼◎◎◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼◎◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋◎◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◘┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉◎◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷◎
	//  true, /* T_IDENTITY */ false, /* T_ROTA_090 */ false, /* T_ROTA_180 */ false, /* T_ROTA_270 */ false, /* T_FLP_SLAS */ false, /* T_FLP_VERT */ true, /* T_FLP_BACK */ false, /* T_FLP_HORI */
	// Handicap pattern 4 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◘┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ true, /* T_ROTA_090 */ true, /* T_ROTA_180 */ true, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ true, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 5 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◘┼┼┼╋╋◉◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ true, /* T_ROTA_090 */ true, /* T_ROTA_180 */ true, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ true, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 6 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯◎◎◎◎◎◎◎◎◎◎
	// ┠╬┼┼┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼╬┼┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼┼┼◎◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼╋╋◎◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ false, /* T_ROTA_090 */ true, /* T_ROTA_180 */ false, /* T_ROTA_270 */ false, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ false, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 7 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯◎◎◎◎◎◎◎◎◎◎
	// ┠╬┼┼┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼╬┼┼┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼┼┼◎◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋◎◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼╋╋◉◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◘┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ false, /* T_ROTA_090 */ true, /* T_ROTA_180 */ false, /* T_ROTA_270 */ false, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ false, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 8 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◉┼┼┼┼┼◉┼┼┼┼┼◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼╋╋◎◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◉┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ true, /* T_ROTA_090 */ true, /* T_ROTA_180 */ true, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ true, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Handicap pattern 9 Board 19 by 19
	// ┏┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯┯◎
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼◎◎◎
	// ┠┼┼◉┼┼┼┼┼◉┼┼┼┼┼◉◎◎◎
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼◎◎◎◎◎
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼◎◎◎◎◎◎
	// ┠┼┼┼┼┼╬┼┼┼┼┼◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋◎◎◎◎◎◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋◎◎◎◎◎◎◎◎◎
	// ┠┼┼◉┼┼┼╋╋◉◎◎◎◎◎◉◎◎◎
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼┼╋╋╋╋╋┼┼┼┼┼┼┨
	// ┠┼┼┼┼┼╬┼┼┼┼┼╬┼┼┼┼┼┨
	// ┠┼┼┼┼╬┼┼┼┼┼┼┼╬┼┼┼┼┨
	// ┠┼┼┼╬┼┼┼┼┼┼┼┼┼╬┼┼┼┨
	// ┠┼┼◉┼┼┼┼┼◉┼┼┼┼┼◉┼┼┨
	// ┠┼╬┼┼┼┼┼┼┼┼┼┼┼┼┼╬┼┨
	// ┠╬┼┼┼┼┼┼┼┼┼┼┼┼┼┼┼╬┨
	// ┗┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┷┛
	//  true, /* T_IDENTITY */ true, /* T_ROTA_090 */ true, /* T_ROTA_180 */ true, /* T_ROTA_270 */ true, /* T_FLP_SLAS */ true, /* T_FLP_VERT */ true, /* T_FLP_BACK */ true, /* T_FLP_HORI */
	// Checking T_IDENTITY and its inverse: T_IDENTITY
	// Checking T_ROTA_090 and its inverse: T_ROTA_270
	// Checking T_ROTA_180 and its inverse: T_ROTA_180
	// Checking T_ROTA_270 and its inverse: T_ROTA_090
	// Checking T_FLP_SLAS and its inverse: T_FLP_SLAS
	// Checking T_FLP_VERT and its inverse: T_FLP_VERT
	// Checking T_FLP_BACK and its inverse: T_FLP_BACK
	// Checking T_FLP_HORI and its inverse: T_FLP_HORI
}

// test EachNode and EachAdjNode
func ExamplePrintBoard() {
	var count int

	printPoint := func(nl ah.NodeLoc) {
		pp := &brds[3].Graphs[ah.PointLevel].Nodes[nl]
		c, r := pp.GetPointColRow()
		if count > 0 {
			fmt.Print(", ")
		}
		fmt.Print("[", c, r, "]")
		count++
	}

	brds[3].EachNode(ah.PointLevel,
		func(brd *ah.Graph, nl ah.NodeLoc) {
			count = 0
			printPoint(nl)
			fmt.Print(": ")
			count = 0
			brds[3].EachAdjNode(ah.PointLevel, nl, printPoint)
			fmt.Println()
		})
	fmt.Println()
	// Output:
	// [A 1]: [A 2], [B 1]
	// [B 1]: [A 1], [B 2], [C 1]
	// [C 1]: [B 1], [C 2], [D 1]
	// [D 1]: [C 1], [D 2], [E 1]
	// [E 1]: [D 1], [E 2], [F 1]
	// [F 1]: [E 1], [F 2], [G 1]
	// [G 1]: [F 1], [G 2], [H 1]
	// [H 1]: [G 1], [H 2], [J 1]
	// [J 1]: [H 1], [J 2], [K 1]
	// [K 1]: [J 1], [K 2], [L 1]
	// [L 1]: [K 1], [L 2]
	// [A 2]: [A 1], [A 3], [B 2]
	// [B 2]: [B 1], [A 2], [B 3], [C 2]
	// [C 2]: [C 1], [B 2], [C 3], [D 2]
	// [D 2]: [D 1], [C 2], [D 3], [E 2]
	// [E 2]: [E 1], [D 2], [E 3], [F 2]
	// [F 2]: [F 1], [E 2], [F 3], [G 2]
	// [G 2]: [G 1], [F 2], [G 3], [H 2]
	// [H 2]: [H 1], [G 2], [H 3], [J 2]
	// [J 2]: [J 1], [H 2], [J 3], [K 2]
	// [K 2]: [K 1], [J 2], [K 3], [L 2]
	// [L 2]: [L 1], [K 2], [L 3]
	// [A 3]: [A 2], [A 4], [B 3]
	// [B 3]: [B 2], [A 3], [B 4], [C 3]
	// [C 3]: [C 2], [B 3], [C 4], [D 3]
	// [D 3]: [D 2], [C 3], [D 4], [E 3]
	// [E 3]: [E 2], [D 3], [E 4], [F 3]
	// [F 3]: [F 2], [E 3], [F 4], [G 3]
	// [G 3]: [G 2], [F 3], [G 4], [H 3]
	// [H 3]: [H 2], [G 3], [H 4], [J 3]
	// [J 3]: [J 2], [H 3], [J 4], [K 3]
	// [K 3]: [K 2], [J 3], [K 4], [L 3]
	// [L 3]: [L 2], [K 3], [L 4]
	// [A 4]: [A 3], [A 5], [B 4]
	// [B 4]: [B 3], [A 4], [B 5], [C 4]
	// [C 4]: [C 3], [B 4], [C 5], [D 4]
	// [D 4]: [D 3], [C 4], [D 5], [E 4]
	// [E 4]: [E 3], [D 4], [E 5], [F 4]
	// [F 4]: [F 3], [E 4], [F 5], [G 4]
	// [G 4]: [G 3], [F 4], [G 5], [H 4]
	// [H 4]: [H 3], [G 4], [H 5], [J 4]
	// [J 4]: [J 3], [H 4], [J 5], [K 4]
	// [K 4]: [K 3], [J 4], [K 5], [L 4]
	// [L 4]: [L 3], [K 4], [L 5]
	// [A 5]: [A 4], [A 6], [B 5]
	// [B 5]: [B 4], [A 5], [B 6], [C 5]
	// [C 5]: [C 4], [B 5], [C 6], [D 5]
	// [D 5]: [D 4], [C 5], [D 6], [E 5]
	// [E 5]: [E 4], [D 5], [E 6], [F 5]
	// [F 5]: [F 4], [E 5], [F 6], [G 5]
	// [G 5]: [G 4], [F 5], [G 6], [H 5]
	// [H 5]: [H 4], [G 5], [H 6], [J 5]
	// [J 5]: [J 4], [H 5], [J 6], [K 5]
	// [K 5]: [K 4], [J 5], [K 6], [L 5]
	// [L 5]: [L 4], [K 5], [L 6]
	// [A 6]: [A 5], [A 7], [B 6]
	// [B 6]: [B 5], [A 6], [B 7], [C 6]
	// [C 6]: [C 5], [B 6], [C 7], [D 6]
	// [D 6]: [D 5], [C 6], [D 7], [E 6]
	// [E 6]: [E 5], [D 6], [E 7], [F 6]
	// [F 6]: [F 5], [E 6], [F 7], [G 6]
	// [G 6]: [G 5], [F 6], [G 7], [H 6]
	// [H 6]: [H 5], [G 6], [H 7], [J 6]
	// [J 6]: [J 5], [H 6], [J 7], [K 6]
	// [K 6]: [K 5], [J 6], [K 7], [L 6]
	// [L 6]: [L 5], [K 6], [L 7]
	// [A 7]: [A 6], [A 8], [B 7]
	// [B 7]: [B 6], [A 7], [B 8], [C 7]
	// [C 7]: [C 6], [B 7], [C 8], [D 7]
	// [D 7]: [D 6], [C 7], [D 8], [E 7]
	// [E 7]: [E 6], [D 7], [E 8], [F 7]
	// [F 7]: [F 6], [E 7], [F 8], [G 7]
	// [G 7]: [G 6], [F 7], [G 8], [H 7]
	// [H 7]: [H 6], [G 7], [H 8], [J 7]
	// [J 7]: [J 6], [H 7], [J 8], [K 7]
	// [K 7]: [K 6], [J 7], [K 8], [L 7]
	// [L 7]: [L 6], [K 7], [L 8]
	// [A 8]: [A 7], [A 9], [B 8]
	// [B 8]: [B 7], [A 8], [B 9], [C 8]
	// [C 8]: [C 7], [B 8], [C 9], [D 8]
	// [D 8]: [D 7], [C 8], [D 9], [E 8]
	// [E 8]: [E 7], [D 8], [E 9], [F 8]
	// [F 8]: [F 7], [E 8], [F 9], [G 8]
	// [G 8]: [G 7], [F 8], [G 9], [H 8]
	// [H 8]: [H 7], [G 8], [H 9], [J 8]
	// [J 8]: [J 7], [H 8], [J 9], [K 8]
	// [K 8]: [K 7], [J 8], [K 9], [L 8]
	// [L 8]: [L 7], [K 8], [L 9]
	// [A 9]: [A 8], [A 10], [B 9]
	// [B 9]: [B 8], [A 9], [B 10], [C 9]
	// [C 9]: [C 8], [B 9], [C 10], [D 9]
	// [D 9]: [D 8], [C 9], [D 10], [E 9]
	// [E 9]: [E 8], [D 9], [E 10], [F 9]
	// [F 9]: [F 8], [E 9], [F 10], [G 9]
	// [G 9]: [G 8], [F 9], [G 10], [H 9]
	// [H 9]: [H 8], [G 9], [H 10], [J 9]
	// [J 9]: [J 8], [H 9], [J 10], [K 9]
	// [K 9]: [K 8], [J 9], [K 10], [L 9]
	// [L 9]: [L 8], [K 9], [L 10]
	// [A 10]: [A 9], [A 11], [B 10]
	// [B 10]: [B 9], [A 10], [B 11], [C 10]
	// [C 10]: [C 9], [B 10], [C 11], [D 10]
	// [D 10]: [D 9], [C 10], [D 11], [E 10]
	// [E 10]: [E 9], [D 10], [E 11], [F 10]
	// [F 10]: [F 9], [E 10], [F 11], [G 10]
	// [G 10]: [G 9], [F 10], [G 11], [H 10]
	// [H 10]: [H 9], [G 10], [H 11], [J 10]
	// [J 10]: [J 9], [H 10], [J 11], [K 10]
	// [K 10]: [K 9], [J 10], [K 11], [L 10]
	// [L 10]: [L 9], [K 10], [L 11]
	// [A 11]: [A 10], [B 11]
	// [B 11]: [B 10], [A 11], [C 11]
	// [C 11]: [C 10], [B 11], [D 11]
	// [D 11]: [D 10], [C 11], [E 11]
	// [E 11]: [E 10], [D 11], [F 11]
	// [F 11]: [F 10], [E 11], [G 11]
	// [G 11]: [G 10], [F 11], [H 11]
	// [H 11]: [H 10], [G 11], [J 11]
	// [J 11]: [J 10], [H 11], [K 11]
	// [K 11]: [K 10], [J 11], [L 11]
	// [L 11]: [L 10], [K 11]
}
