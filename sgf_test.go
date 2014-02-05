package sgf_test

import (
	"fmt"
	. "github.com/Ken1JF/sgf"
)

func ExampleSGFTypeSizes() {
	PrintSGFTypeSizes()
	// Output:
	// Type Token size 1 alignment 1
	// Type ah.Position size 40 alignment 8
	// Type TreeNodeType size 1 alignment 1
	// Type TreeNodeIdx size 2 alignment 2
	// Type PropIdx size 2 alignment 2
	// Type TreeNode size 12 alignment 2
	// Type PropertyValue size 32 alignment 8
	// Type GameTree size 1520 alignment 8
	// Type Parser size 1840 alignment 8
	// Type PlayerInfo size 72 alignment 8
	// Type DBStatistics size 704 alignment 8
	// Type FF4Note size 1 alignment 1
	// Type SGFPropNodeType size 1 alignment 1
	// Type QualifierType size 1 alignment 1
	// Type PropValueType size 1 alignment 1
	// Type Property size 56 alignment 8
	// Type PropertyDefIdx size 1 alignment 1
	// Type ID_CountArray size 624 alignment 8
	// Type Scanner size 112 alignment 8
	// Type ErrorHandler size 8 alignment 8
	// Type ah.ErrorList size 24 alignment 8
	// Type Komi size 8 alignment 4
	// Type Result size 64 alignment 8
}

var defaultSpecFile = "../sgf/sgf_properties_spec.txt"

func ExampleTheProperties_tt() {
	err := SetupSGFProperties(defaultSpecFile, true, true)
	if err == 0 {

	} else {
		fmt.Println("Can't read Specification file:", defaultSpecFile)
	}
	// Output:
	// Read: 78 lines, 3914 bytes, 78 properties.
	//  0: AB:       Add Black:     setup:          :     Std: 19:list of stone
	//  1: AE:       Add Empty:     setup:          :     Std:  8:list of point
	//  2: AN:      Annotation: game-info:          :     Std: 11:simpletext
	//  3: AP:     Application:      root:          :     New: 10:composed simpletext ':' simpletext
	//  4: AR:           Arrow:        --:          :     New:  6:list of composed point ':' point
	//  5: AS: Who adds stones:        --:     (LOA):     New: 11:simpletext
	//  6: AW:       Add White:     setup:          :     Std: 19:list of stone
	//  7:  B:           Black:      move:          :     Std: 18:move
	//  8: BL: Black time left:      move:          :     Std: 21:real
	//  9: BM:        Bad move:      move:          :     Std: 22:double
	// 10: BR:      Black rank: game-info:          :     Std: 11:simpletext
	// 11: BT:      Black team: game-info:          :     Std: 11:simpletext
	// 12:  C:         Comment:        --:          :     Std: 12:text
	// 13: CA:         Charset:      root:          :     New: 11:simpletext
	// 14: CP:       Copyright: game-info:          :     Std: 11:simpletext
	// 15: CR:          Circle:        --:          :     Std:  8:list of point
	// 16: DD:      Dim points:        --: (inherit):     New:  7:elist of point
	// 17: DM:   Even position:        --:          :     Std: 22:double
	// 18: DO:        Doubtful:      move:          :     Std:  2:none
	// 19: DT:            Date: game-info:          : Changed: 11:simpletext
	// 20: EV:           Event: game-info:          :     Std: 11:simpletext
	// 21: FF:      Fileformat:      root:          :     Std: 15:number (range: 1-4)
	// 22: FG:          Figure:        --:          : Changed:  1:none | composed number ":" simpletext
	// 23: GB:  Good for Black:        --:          :     Std: 22:double
	// 24: GC:    Game comment: game-info:          :     Std: 12:text
	// 25: GM:            Game:      root:          :     Std: 16:number (range: 1-5,7-16)
	// 26: GN:       Game name: game-info:          :     Std: 11:simpletext
	// 27: GW:  Good for White:        --:          :     Std: 22:double
	// 28: HA:        Handicap: game-info:      (Go):     Std: 17:number
	// 29: HO:         Hotspot:        --:          :     Std: 22:double
	// 30: IP:    Initial pos.: game-info:     (LOA):     New: 11:simpletext
	// 31: IT:     Interesting:      move:          :     Std:  2:none
	// 32: IY:   Invert Y-axis: game-info:     (LOA):     New: 11:simpletext
	// 33: KM:            Komi: game-info:      (Go):     Std: 21:real
	// 34: KO:              Ko:      move:          :     Std:  2:none
	// 35: LB:           Label:        --:          : Changed:  5:list of composed point ':' simpletext
	// 36: LN:            Line:        --:          :     New:  6:list of composed point ':' point
	// 37: MA:            Mark:        --:          :     Std:  8:list of point
	// 38: MN: set move number:      move:          :     Std: 17:number
	// 39:  N:        Nodename:        --:          :     Std: 11:simpletext
	// 40: OB:  OtStones Black:      move:          :     Std: 17:number
	// 41: OH:    Old Handicap: game-info:          : Non_std: 12:text
	// 42: ON:         Opening: game-info:          :     Std: 12:text
	// 43: OT:        Overtime: game-info:          :     New: 11:simpletext
	// 44: OW:  OtStones White:      move:          :     Std: 17:number
	// 45: PB:    Player Black: game-info:          :     Std: 11:simpletext
	// 46: PC:           Place: game-info:          :     Std: 11:simpletext
	// 47: PL:  Player to play:     setup:          :     Std: 23:color
	// 48: PM: Print move mode:        --: (inherit):     New: 17:number
	// 49: PW:    Player White: game-info:          :     Std: 11:simpletext
	// 50: RE:          Result: game-info:          : Changed: 11:simpletext
	// 51: RO:           Round: game-info:          :     Std: 11:simpletext
	// 52: RU:           Rules: game-info:          : Changed: 11:simpletext
	// 53:  S:        Sequence:      move:     (SGC): Non_std:  4:compressed list of point
	// 54: SE:          Markup:        --:     (LOA):     New:  9:point
	// 55: SL:        Selected:        --:          :     Std:  8:list of point
	// 56: SO:          Source: game-info:          :     Std: 11:simpletext
	// 57: SQ:          Square:        --:          :     New:  8:list of point
	// 58: ST:           Style:      root:          :     New: 14:number (range: 0-3)
	// 59: SU:      Setup type: game-info:     (LOA):     New: 11:simpletext
	// 60: SZ:            Size:      root:          : Changed: 13:(number | composed number ':' number)
	// 61: TB: Territory Black:        --:      (Go):     Std:  7:elist of point
	// 62: TE:          Tesuji:      move:          :     Std: 22:double
	// 63: TM:       Timelimit: game-info:          :     Std: 21:real
	// 64: TR:        Triangle:        --:          :     Std:  8:list of point
	// 65: TW: Territory White:        --:      (Go):     Std:  7:elist of point
	// 66: UC:     Unclear pos:        --:          :     Std: 22:double
	// 67: US:            User: game-info:          :     Std: 11:simpletext
	// 68:  V:           Value:        --:          :     Std: 21:real
	// 69: VW:            View:        --: (inherit):     New:  7:elist of point
	// 70:  W:           White:      move:          :     Std: 18:move
	// 71: WB:      Wins Black:      move:          : Non_std: 17:number
	// 72: WC:    Win Continue:      move:          : Non_std: 10:composed simpletext ':' simpletext
	// 73: WL: White time left:      move:          :     Std: 21:real
	// 74: WO:      Wins Other:      move:          : Non_std: 17:number
	// 75: WR:      White rank: game-info:          :     Std: 11:simpletext
	// 76: WT:      White team: game-info:          :     Std: 11:simpletext
	// 77: WW:      Wins White:      move:          : Non_std: 17:number
	// AB_idx PropertyDefIdx = 0
	// AE_idx PropertyDefIdx = 1
	// AN_idx PropertyDefIdx = 2
	// AP_idx PropertyDefIdx = 3
	// AR_idx PropertyDefIdx = 4
	// AS_idx PropertyDefIdx = 5
	// AW_idx PropertyDefIdx = 6
	// B_idx PropertyDefIdx = 7
	// BL_idx PropertyDefIdx = 8
	// BM_idx PropertyDefIdx = 9
	// BR_idx PropertyDefIdx = 10
	// BT_idx PropertyDefIdx = 11
	// C_idx PropertyDefIdx = 12
	// CA_idx PropertyDefIdx = 13
	// CP_idx PropertyDefIdx = 14
	// CR_idx PropertyDefIdx = 15
	// DD_idx PropertyDefIdx = 16
	// DM_idx PropertyDefIdx = 17
	// DO_idx PropertyDefIdx = 18
	// DT_idx PropertyDefIdx = 19
	// EV_idx PropertyDefIdx = 20
	// FF_idx PropertyDefIdx = 21
	// FG_idx PropertyDefIdx = 22
	// GB_idx PropertyDefIdx = 23
	// GC_idx PropertyDefIdx = 24
	// GM_idx PropertyDefIdx = 25
	// GN_idx PropertyDefIdx = 26
	// GW_idx PropertyDefIdx = 27
	// HA_idx PropertyDefIdx = 28
	// HO_idx PropertyDefIdx = 29
	// IP_idx PropertyDefIdx = 30
	// IT_idx PropertyDefIdx = 31
	// IY_idx PropertyDefIdx = 32
	// KM_idx PropertyDefIdx = 33
	// KO_idx PropertyDefIdx = 34
	// LB_idx PropertyDefIdx = 35
	// LN_idx PropertyDefIdx = 36
	// MA_idx PropertyDefIdx = 37
	// MN_idx PropertyDefIdx = 38
	// N_idx PropertyDefIdx = 39
	// OB_idx PropertyDefIdx = 40
	// OH_idx PropertyDefIdx = 41
	// ON_idx PropertyDefIdx = 42
	// OT_idx PropertyDefIdx = 43
	// OW_idx PropertyDefIdx = 44
	// PB_idx PropertyDefIdx = 45
	// PC_idx PropertyDefIdx = 46
	// PL_idx PropertyDefIdx = 47
	// PM_idx PropertyDefIdx = 48
	// PW_idx PropertyDefIdx = 49
	// RE_idx PropertyDefIdx = 50
	// RO_idx PropertyDefIdx = 51
	// RU_idx PropertyDefIdx = 52
	// S_idx PropertyDefIdx = 53
	// SE_idx PropertyDefIdx = 54
	// SL_idx PropertyDefIdx = 55
	// SO_idx PropertyDefIdx = 56
	// SQ_idx PropertyDefIdx = 57
	// ST_idx PropertyDefIdx = 58
	// SU_idx PropertyDefIdx = 59
	// SZ_idx PropertyDefIdx = 60
	// TB_idx PropertyDefIdx = 61
	// TE_idx PropertyDefIdx = 62
	// TM_idx PropertyDefIdx = 63
	// TR_idx PropertyDefIdx = 64
	// TW_idx PropertyDefIdx = 65
	// UC_idx PropertyDefIdx = 66
	// US_idx PropertyDefIdx = 67
	// V_idx PropertyDefIdx = 68
	// VW_idx PropertyDefIdx = 69
	// W_idx PropertyDefIdx = 70
	// WB_idx PropertyDefIdx = 71
	// WC_idx PropertyDefIdx = 72
	// WL_idx PropertyDefIdx = 73
	// WO_idx PropertyDefIdx = 74
	// WR_idx PropertyDefIdx = 75
	// WT_idx PropertyDefIdx = 76
	// WW_idx PropertyDefIdx = 77
	// { 0, "AB", "Add Black",  2, 0, 19 },
	// { 0, "AE", "Add Empty",  2, 0, 8 },
	// { 0, "AN", "Annotation",  3, 0, 11 },
	// { 1, "AP", "Application",  1, 0, 10 },
	// { 1, "AR", "Arrow",  0, 0, 6 },
	// { 1, "AS", "Who adds stones",  0, 2, 11 },
	// { 0, "AW", "Add White",  2, 0, 19 },
	// { 0, "B", "Black",  4, 0, 18 },
	// { 0, "BL", "Black time left",  4, 0, 21 },
	// { 0, "BM", "Bad move",  4, 0, 22 },
	// { 0, "BR", "Black rank",  3, 0, 11 },
	// { 0, "BT", "Black team",  3, 0, 11 },
	// { 0, "C", "Comment",  0, 0, 12 },
	// { 1, "CA", "Charset",  1, 0, 11 },
	// { 0, "CP", "Copyright",  3, 0, 11 },
	// { 0, "CR", "Circle",  0, 0, 8 },
	// { 1, "DD", "Dim points",  0, 1, 7 },
	// { 0, "DM", "Even position",  0, 0, 22 },
	// { 0, "DO", "Doubtful",  4, 0, 2 },
	// { 2, "DT", "Date",  3, 0, 11 },
	// { 0, "EV", "Event",  3, 0, 11 },
	// { 0, "FF", "Fileformat",  1, 0, 15 },
	// { 2, "FG", "Figure",  0, 0, 1 },
	// { 0, "GB", "Good for Black",  0, 0, 22 },
	// { 0, "GC", "Game comment",  3, 0, 12 },
	// { 0, "GM", "Game",  1, 0, 16 },
	// { 0, "GN", "Game name",  3, 0, 11 },
	// { 0, "GW", "Good for White",  0, 0, 22 },
	// { 0, "HA", "Handicap",  3, 3, 17 },
	// { 0, "HO", "Hotspot",  0, 0, 22 },
	// { 1, "IP", "Initial pos.",  3, 2, 11 },
	// { 0, "IT", "Interesting",  4, 0, 2 },
	// { 1, "IY", "Invert Y-axis",  3, 2, 11 },
	// { 0, "KM", "Komi",  3, 3, 21 },
	// { 0, "KO", "Ko",  4, 0, 2 },
	// { 2, "LB", "Label",  0, 0, 5 },
	// { 1, "LN", "Line",  0, 0, 6 },
	// { 0, "MA", "Mark",  0, 0, 8 },
	// { 0, "MN", "set move number",  4, 0, 17 },
	// { 0, "N", "Nodename",  0, 0, 11 },
	// { 0, "OB", "OtStones Black",  4, 0, 17 },
	// { 3, "OH", "Old Handicap",  3, 0, 12 },
	// { 0, "ON", "Opening",  3, 0, 12 },
	// { 1, "OT", "Overtime",  3, 0, 11 },
	// { 0, "OW", "OtStones White",  4, 0, 17 },
	// { 0, "PB", "Player Black",  3, 0, 11 },
	// { 0, "PC", "Place",  3, 0, 11 },
	// { 0, "PL", "Player to play",  2, 0, 23 },
	// { 1, "PM", "Print move mode",  0, 1, 17 },
	// { 0, "PW", "Player White",  3, 0, 11 },
	// { 2, "RE", "Result",  3, 0, 11 },
	// { 0, "RO", "Round",  3, 0, 11 },
	// { 2, "RU", "Rules",  3, 0, 11 },
	// { 3, "S", "Sequence",  4, 4, 4 },
	// { 1, "SE", "Markup",  0, 2, 9 },
	// { 0, "SL", "Selected",  0, 0, 8 },
	// { 0, "SO", "Source",  3, 0, 11 },
	// { 1, "SQ", "Square",  0, 0, 8 },
	// { 1, "ST", "Style",  1, 0, 14 },
	// { 1, "SU", "Setup type",  3, 2, 11 },
	// { 2, "SZ", "Size",  1, 0, 13 },
	// { 0, "TB", "Territory Black",  0, 3, 7 },
	// { 0, "TE", "Tesuji",  4, 0, 22 },
	// { 0, "TM", "Timelimit",  3, 0, 21 },
	// { 0, "TR", "Triangle",  0, 0, 8 },
	// { 0, "TW", "Territory White",  0, 3, 7 },
	// { 0, "UC", "Unclear pos",  0, 0, 22 },
	// { 0, "US", "User",  3, 0, 11 },
	// { 0, "V", "Value",  0, 0, 21 },
	// { 1, "VW", "View",  0, 1, 7 },
	// { 0, "W", "White",  4, 0, 18 },
	// { 3, "WB", "Wins Black",  4, 0, 17 },
	// { 3, "WC", "Win Continue",  4, 0, 10 },
	// { 0, "WL", "White time left",  4, 0, 21 },
	// { 3, "WO", "Wins Other",  4, 0, 17 },
	// { 0, "WR", "White rank",  3, 0, 11 },
	// { 0, "WT", "White team",  3, 0, 11 },
	// { 3, "WW", "Wins White",  4, 0, 17 },
}

func ExampleTheProperties_tf() {
	err := SetupSGFProperties(defaultSpecFile, true, false)
	if err == 0 {

	} else {
		fmt.Println("Can't read Specification file:", defaultSpecFile)
	}
	// Output:
	//
}
func ExampleTheProperties_ff() {
	err := SetupSGFProperties(defaultSpecFile, false, false)
	if err == 0 {

	} else {
		fmt.Println("Can't read Specification file:", defaultSpecFile)
	}
	// Output:
	//
}
