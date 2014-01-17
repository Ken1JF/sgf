/*
 *  File:		src/gitHub.com/Ken1JF/ahgo/sgf/sgf.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 11/25/09.
 *  Copyright 2009-2014 Ken Friedenbach. All rights reserved.
 *
 *	This package provides external types, constants, and functions to support
 *	the reading and interpretation of SGF FF4 property definition files.
 *
 *	If a property ID is encountered which is unknown, the index value will be -1,
 *	and the id will be saved in the global variable:
 *		UnknownProperty - a pointer to the most recently encountered unknown property.
 */

package sgf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gitHub.com/Ken1JF/ahgo/ah"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const TRACE_SGF bool = true

func tra(s string) string {
	if TRACE_SGF && ah.TraceAH {
		fmt.Println("Entering: ", s)
	}
	return s
}

func unt(s string) {
	if TRACE_SGF && ah.TraceAH {
		fmt.Println("Leaving: ", s)
	}
}

// FF4 note codes: (optional character before ID)
//
type FF4Note uint8

const ( // code	meaning
	Std_SGF4     FF4Note = iota // none "Std"
	New_SGF4                    // '*'	"New"
	Changed_SGF4                // '!'	"Changed"
	Non_std_SGF4                // '#'	"Non_std"
	Unknown_SGF4                // ?	"Unknown"
)

// FF4NoteNames is an array containing strings used
// when printing the SGF properties with FF4 notes.
// TODO: cleaner to use a map???
//
var FF4NoteNames = [...]string{
	"Std",
	"New",
	"Changed",
	"Non_std",
	"Unknown",
}

// SGFPropNodeType is used to provide named constants for
// node type descriptions used in the SGF Specification file
//
type SGFPropNodeType uint8

const (
	NoType SGFPropNodeType = iota
	RootProp
	SetupProp
	GameInfoProp
	MoveProp
)

// SGFPropNodeTypeNames is an array containing the strings used
// in the SGF Specification file to describe classes of properties.
// Note: the original Specification file used "-" for none. This was changed
// to "--" so it could be found with Index.
// SGFPropNodeTypeNames must be kept in sync with the SGFPropNodeType constants.
// Note: unknown properties will have NoType for property SGFPropNodeType.
// TODO: cleaner to use a map?
//
var SGFPropNodeTypeNames = [...]string{
	"--",
	"root",
	"setup",
	"game-info",
	"move",
}

// checkPropertyType checks that node types appear where they
// should, as specified by the FF4 Specification.
// TODO: allow some flexibility for trees of move patterns
//
func checkPropertyType(prop *Property, inGameRoot bool) (ret error) {
	if prop != nil {
		if inGameRoot {
			switch prop.FF4Type {
			case NoType, RootProp, SetupProp, GameInfoProp:
				// do nothing
			case MoveProp:
				ret = errors.New("move property in root node")
			default:
				ret = errors.New("unknown SGFPropNodeType")
			}
		} else {
			switch prop.FF4Type {
			case RootProp:
				ret = errors.New("RootProp not in root node")
			case SetupProp:
				ret = errors.New("SetupProp not in root node")
			case GameInfoProp:
				ret = errors.New("GameInfoProp not in root node")
			case NoType, MoveProp:
				// do nothing
			default:
				ret = errors.New("unknown SGFPropNodeType")
			}
		}
	} else {
		ret = errors.New("checkPropertyType: nil prop pointer")
	}
	return ret
}

// QualifierType is used to provide named constants for
// qualifier value descriptions used in the SGF Specification file
//
type QualifierType uint8

const (
	NoQualifier QualifierType = iota // ""
	Inherit                          // "(inherit)"
	LOA                              // "(LOA)"
	GoGO                             // "(Go)"
	SGC                              // "(SGC)"
)

// QualifierNames is an array containing the strings used
// in the SGF Specification file to describe qualifiers.
// QualifierNames must be kept in sync with the QualifierType constants.
// Note: unknown properties will have NoQualifier for property QualifierType.
// TODO: cleaner to use a map?
//
var QualifierNames = [...]string{
	"",
	"(inherit)",
	"(LOA)",
	"(Go)",
	"(SGC)",
}

// PropValueType is used to provide named constants for
// property value descriptions used in the SGF Specification file
//
type PropValueType uint8

const (
	Unknown                  PropValueType = iota // "<?Unknown?>"
	None_OR_compNum_simpText                      // "none | composed number \":\" simpletext"
	None                                          // "none"
	CompNum_simpText                              // "composed number \":\" simpletext"
	CompressedListOfPoint                         // "compressed list of point"
	ListOfCompPoint_simpTest                      // "list of composed point ':' simpletext"
	ListOfCompPoint_Point                         // "list of composed point ':' point"
	EListOfPoint                                  // "elist of point"
	ListOfPoint                                   // "list of point"
	Point                                         // "point"
	CompSimpText_simpText                         // "composed simpletext ':' simpletext"
	SimpText                                      // "simpletext"
	Text                                          // "text"
	Num_OR_compNum_num                            // "(number | composed number ':' number)"
	Num_0_3                                       // "number (range: 0-3)"
	Num_1_4                                       // "number (range: 1-4)"
	Num_1_5_or_7_16                               // "number (range: 1-5,7-16)"
	Num                                           // "number"
	Move                                          // "move"
	ListOfStone                                   // "list of stone"
	Stone                                         // "stone"
	Real                                          // "real"
	Double                                        // "double"
	Color                                         // "color"
)

// ValueNames is an array containing the strings used
// in the SGF Specification file to describe value syntax.
// ValueNames must be kept in sync with the PropValueType constants.
// Note: unknown properties will have Unknown for property PropValueType.
// TODO: cleaner to use a map?
//
var ValueNames = [...]string{
	"<?Unknown?>",
	"none | composed number \":\" simpletext",
	"none",
	"composed number \":\" simpletext",
	"compressed list of point",
	"list of composed point ':' simpletext",
	"list of composed point ':' point",
	"elist of point",
	"list of point",
	"point",
	"composed simpletext ':' simpletext",
	"simpletext",
	"text",
	"(number | composed number ':' number)",
	"number (range: 0-3)",
	"number (range: 1-4)",
	"number (range: 1-5,7-16)",
	"number",
	"move",
	"list of stone",
	"stone",
	"real",
	"double",
	"color",
}

// FF4 Property specifications:
//
type Property struct {
	Note        FF4Note // see FF4 note codes, above
	ID          []byte  // one or two bytes
	Description string
	FF4Type     SGFPropNodeType // see Property types, above
	Qualifier   QualifierType   // see Qualifiers, above
	Value       PropValueType   // see Value types, above
}

// theProperties is an internal array containing the properties
// It is const after initialization, i.e. thread-safe to share.
//
var theProperties []Property

// PropertyDefIdx is an index into theProperties
//
type PropertyDefIdx int8 // index into theProperties

// These constants are generated by the "verbose" option
// of the Setup function.
// If the SGF_Properties_Spec.txt file is changed, be sure
// to regenerate these:
//
const (
	AB_idx PropertyDefIdx = 0
	AE_idx PropertyDefIdx = 1
	AN_idx PropertyDefIdx = 2
	AP_idx PropertyDefIdx = 3
	AR_idx PropertyDefIdx = 4
	AS_idx PropertyDefIdx = 5
	AW_idx PropertyDefIdx = 6
	B_idx  PropertyDefIdx = 7
	BL_idx PropertyDefIdx = 8
	BM_idx PropertyDefIdx = 9
	BR_idx PropertyDefIdx = 10
	BT_idx PropertyDefIdx = 11
	C_idx  PropertyDefIdx = 12
	CA_idx PropertyDefIdx = 13
	CP_idx PropertyDefIdx = 14
	CR_idx PropertyDefIdx = 15
	DD_idx PropertyDefIdx = 16
	DM_idx PropertyDefIdx = 17
	DO_idx PropertyDefIdx = 18
	DT_idx PropertyDefIdx = 19
	EV_idx PropertyDefIdx = 20
	FF_idx PropertyDefIdx = 21
	FG_idx PropertyDefIdx = 22
	GB_idx PropertyDefIdx = 23
	GC_idx PropertyDefIdx = 24
	GM_idx PropertyDefIdx = 25
	GN_idx PropertyDefIdx = 26
	GW_idx PropertyDefIdx = 27
	HA_idx PropertyDefIdx = 28
	HO_idx PropertyDefIdx = 29
	IP_idx PropertyDefIdx = 30
	IT_idx PropertyDefIdx = 31
	IY_idx PropertyDefIdx = 32
	KM_idx PropertyDefIdx = 33
	KO_idx PropertyDefIdx = 34
	LB_idx PropertyDefIdx = 35
	LN_idx PropertyDefIdx = 36
	MA_idx PropertyDefIdx = 37
	MN_idx PropertyDefIdx = 38
	N_idx  PropertyDefIdx = 39
	OB_idx PropertyDefIdx = 40
	OH_idx PropertyDefIdx = 41
	ON_idx PropertyDefIdx = 42
	OT_idx PropertyDefIdx = 43
	OW_idx PropertyDefIdx = 44
	PB_idx PropertyDefIdx = 45
	PC_idx PropertyDefIdx = 46
	PL_idx PropertyDefIdx = 47
	PM_idx PropertyDefIdx = 48
	PW_idx PropertyDefIdx = 49
	RE_idx PropertyDefIdx = 50
	RO_idx PropertyDefIdx = 51
	RU_idx PropertyDefIdx = 52
	S_idx  PropertyDefIdx = 53
	SE_idx PropertyDefIdx = 54
	SL_idx PropertyDefIdx = 55
	SO_idx PropertyDefIdx = 56
	SQ_idx PropertyDefIdx = 57
	ST_idx PropertyDefIdx = 58
	SU_idx PropertyDefIdx = 59
	SZ_idx PropertyDefIdx = 60
	TB_idx PropertyDefIdx = 61
	TE_idx PropertyDefIdx = 62
	TM_idx PropertyDefIdx = 63
	TR_idx PropertyDefIdx = 64
	TW_idx PropertyDefIdx = 65
	UC_idx PropertyDefIdx = 66
	US_idx PropertyDefIdx = 67
	V_idx  PropertyDefIdx = 68
	VW_idx PropertyDefIdx = 69
	W_idx  PropertyDefIdx = 70
	WB_idx PropertyDefIdx = 71
	WC_idx PropertyDefIdx = 72
	WL_idx PropertyDefIdx = 73
	WO_idx PropertyDefIdx = 74
	WR_idx PropertyDefIdx = 75
	WT_idx PropertyDefIdx = 76
	WW_idx PropertyDefIdx = 77
)

type ID_CountArray [78]int

// TODO: should this be 0 (initially) => len(theProperties) after initialization?
// to allow more then 127 properties?
//
const UnknownPropIdx PropertyDefIdx = PropertyDefIdx(-1)

// GetProperty is an accessor function, returning the Property strut associated
// with a PropertyDefIdx value
//
func GetProperty(idx PropertyDefIdx) (p *Property) {
	if idx >= 0 && idx < PropertyDefIdx(len(theProperties)) {
		p = &((theProperties)[idx])
	}
	return p
}

//TODO: make the above global variables part of an SGFReader strut

// parseProperty parses a line from the SGF Specification file,
// and builds a property. It returns the property,
// together with a string indicating any errors that occurred.
//
func parseProperty(b []byte) (ret Property, err string) {
	// process the note (or first char of ID)
	if len(b) < 1 {
		return ret, "empty string"
	}
	ch := b[0]
	b = b[1:]
	switch ch {
	case '*':
		ret.Note = New_SGF4
	case '!':
		ret.Note = Changed_SGF4
	case '#':
		ret.Note = Non_std_SGF4
	default:
		ret.Note = Std_SGF4
		ret.ID = append(ret.ID, ch)
	}
	// find the (rest of?) the ID (one char of Std_SGF4 properties handled above)
	if len(b) < 1 {
		return ret, "no ID"
	}
	ch = b[0]
	b = b[1:]
	for !unicode.IsSpace(rune(ch)) {
		ret.ID = append(ret.ID, ch)
		if len(b) < 1 {
			return ret, "no property"
		}
		ch = b[0]
		b = b[1:]
	}
	b = bytes.TrimSpace(b)
	// set the Description
	ret.Description = string(bytes.TrimSpace(b[0:15]))
	b = bytes.TrimSpace(b[15:])
	// find the type:
	typeIdx := -1
	for ti, ts := range SGFPropNodeTypeNames {
		typeIdx = bytes.Index(b, []byte(ts))
		if typeIdx >= 0 {
			ret.FF4Type = SGFPropNodeType(ti)
			b = bytes.TrimSpace(b[typeIdx+len(ts):])
			break
		}
	}
	if typeIdx == -1 {
		return ret, "type not found:" + string(b)
	}
	// find the qualifier, if any
	qualIdx := -1
	for qi, qs := range QualifierNames {
		// skip the empty string
		if qi > 0 {
			qualIdx = bytes.Index(b, []byte(qs))
			if qualIdx >= 0 {
				ret.Qualifier = QualifierType(qi)
				b = bytes.TrimSpace(b[qualIdx+len(qs):])
				break
			}
		}
	}
	if qualIdx == -1 {
		ret.Qualifier = NoQualifier
	}
	// find the property type
	switch string(b) {
	case "none | composed number \":\" simpletext":
		ret.Value = None_OR_compNum_simpText
	case "none":
		ret.Value = None
	case "compressed list of point":
		ret.Value = CompressedListOfPoint
	case "list of composed point ':' simpletext":
		ret.Value = ListOfCompPoint_simpTest
	case "list of composed point ':' point":
		ret.Value = ListOfCompPoint_Point
	case "elist of point":
		ret.Value = EListOfPoint
	case "list of point":
		ret.Value = ListOfPoint
	case "point":
		ret.Value = Point
	case "composed simpletext ':' simpletext":
		ret.Value = CompSimpText_simpText
	case "simpletext":
		ret.Value = SimpText
	case "text":
		ret.Value = Text
	case "(number | composed number ':' number)":
		ret.Value = Num_OR_compNum_num
	case "number (range: 0-3)":
		ret.Value = Num_0_3
	case "number (range: 1-4)":
		ret.Value = Num_1_4
	case "number (range: 1-5,7-16)":
		ret.Value = Num_1_5_or_7_16
	case "number":
		ret.Value = Num
	case "move":
		ret.Value = Move
	case "list of stone":
		ret.Value = ListOfStone
	case "stone":
		ret.Value = Stone
	case "real":
		ret.Value = Real
	case "double":
		ret.Value = Double
	case "color":
		ret.Value = Color
	default:
		ret.Value = Unknown
		fmt.Printf("Warning, unknown Value type: %s\n", b)
	}
	return ret, ""
}

// addPropertyDef appends a property to theProperties, which
// is grown dynamically to hold all the properties
//
func addPropertyDef(sp []Property, p Property) []Property {
	cur_l := len(sp)
	cur_c := cap(sp)
	if cur_l+1 > cur_c { // reallocate
		// Allocate double what's needed for future growth...
		newLen := (cur_c + 1) * 2
		// ... but avoid small startup sizes:
		if newLen < 16 {
			newLen = 16
		}
		newSlice := make([]Property, newLen)
		// Copy:
		for i, cur_p := range sp {
			newSlice[i] = cur_p
		}
		sp = newSlice
	}
	sp = sp[0 : cur_l+1]
	sp[cur_l] = p
	return sp
}

// readSpecFile reads the SGF Specification file,
// and stores the properties read in theProperties.
// if an error occurs, it returns an Error other than io.EOF
//
func readSpecFile(fn string, verbose bool) (err error) {
	var (
		line                               []byte
		line_count, byte_count, prop_count int
	)
	fd, err := os.Open(fn) // Old parms to Open(fn, int(os.O_RDONLY), uint32(0))
	if err != nil {
		fmt.Printf("Error while opening: \"%s\", %s.\n", fn, err.Error())
		return err
	}
	defer fd.Close()

	bf := bufio.NewReader(fd)
	for {
		line, err = bf.ReadBytes('\n')
		if err != nil {
			if (err != io.EOF) || (len(line) == 0) {
				fmt.Printf("Read Error, while reading: \"%s\", %s\n", fn, err.Error())
				return err
			}
		}
		//		fmt.Printf("Read %d bytes: %s", line)
		byte_count += len(line)
		line_count++
		p, e := parseProperty(bytes.TrimSpace(line))
		if e == "" {
			theProperties = addPropertyDef(theProperties, p)
			prop_count++
		} else if e == "empty string" {
			break
		} else {
			fmt.Printf("Error reading property: %s, %s\n", e, line)
		}
		if err == io.EOF {
			break
		}
	}
	if verbose {
		fmt.Printf("Read: %d lines, %d bytes, %d properties.\n", line_count, byte_count, prop_count)
	}
	return err
}

// Setup reads the SGF Specification file, builds theProperties array, and returns:
//	0 if all properties are in order
//	-1 if SGF Specification file cannot be read,
//	n if n properties are out of order
//
func Setup(specFile string, verifyOrder bool, verbose bool) (ret int) {
	err := readSpecFile(specFile, verbose)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading SGF Spec File: %s, %s\n", specFile, err)
		return -1
	}
	if verbose {
		for i, p := range theProperties {
			fmt.Printf("%3d:%3s:%16s:%10s:%10s:%8s:%3d:%s\n",
				i, p.ID, p.Description, SGFPropNodeTypeNames[p.FF4Type],
				QualifierNames[p.Qualifier], FF4NoteNames[p.Note],
				p.Value, ValueNames[p.Value])
		}
		for i, p := range theProperties {
			fmt.Printf("\t%s_idx PropertyDefIdx = %d\n", p.ID, i)
		}
		for i, p := range theProperties {
			fmt.Printf("\t { /* %d */ %d", i, p.Note)
			fmt.Printf(", \"%s\"", p.ID)
			fmt.Printf(", \"%s\", ", p.Description)
			fmt.Printf(" %d", p.FF4Type)
			fmt.Printf(", %d", p.Qualifier)
			fmt.Printf(", %d }, \n", p.Value)
		}
	}
	var prev_p Property
	if verifyOrder {
		for i, p := range theProperties {
			if i > 0 { // skip first one, no previous one to compare to
				if bytes.Compare(prev_p.ID, p.ID) >= 0 {
					fmt.Printf("Error, properties out of order: \"%s\" >= \"%s\"\n", prev_p.ID, p.ID)
					ret++
				}
			}
			prev_p = p
		}
	}
	return ret
}

// LookUp does a binary search on theProperties,
// and returns the PropertyDefIdx corresponding to the id.
// If not found, returns UnknownPropIdx
//
func LookUp(id []byte) (prop PropertyDefIdx) {
	var (
		LEN      int = len(theProperties)
		min, max int = 0, LEN - 1
		mid      int
	)
	prop = UnknownPropIdx
Loop:
	for max >= min {
		mid = (min + max) >> 1
		if mid >= LEN {
			break
		} // not found
		switch bytes.Compare(id, (theProperties)[mid].ID) {
		case -1:
			max = mid - 1
		case 0:
			prop = PropertyDefIdx(mid)
			break Loop
		case 1:
			min = mid + 1
		}
	}
	return prop
}

// SGF FF[4] suuports up to 52x52 board sizes:
//
const sgf_coords = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Given a NodeLoc, convert to SGF coordinates:
// isFF4 handles the FF[3] vs FF[4] variation in representation of "pass"
//
func SGFCoords(n ah.NodeLoc, isFF4 bool) (ret []byte) {
	ret = make([]byte, 2)
	if n != ah.PassNodeLoc {
		c, r := ah.GetColRow(n)
		ret[0] = sgf_coords[c]
		ret[1] = sgf_coords[r]
	} else if isFF4 {
		ret = make([]byte, 0)
	} else {
		ret[0] = 't'
		ret[1] = 't'
	}
	return ret
}

// Convert the first 2 characters of a string into an SGF point represented as a NodeLoc
//
func SGFPoint(s []byte) (n ah.NodeLoc, err ah.ErrorList) {
	defer unt(tra("SGFPoint"))
	if len(s) < 2 {
		if len(s) == 0 {
			// TODO: check that SGF FF == 4 ???
			n = ah.PassNodeLoc
		} else {
			err.Add(ah.NoPos, "SGFPoint, len = "+strconv.Itoa(len(s))+" not 0 or 2 chars, "+string(s))
		}
	} else {
		col := strings.Index(sgf_coords, string(s[0:1]))
		if col < 0 {
			err.Add(ah.NoPos, "SGFPoint: bad col "+string(s[0:1]))
		}
		row := strings.Index(sgf_coords, string(s[1:2]))
		if row < 0 {
			err.Add(ah.NoPos, "SGFPoint: bad row "+string(s[1:2]))
		}
		if err == nil {
			n = ah.MakeNodeLoc(ah.ColValue(col), ah.RowValue(row))
		}
	}
	return n, err
}

// Some functions to print the size and alignment of types:
// TODO: delete when no longer needed
//
//func printSizeAlign(s string, sz int, al int) {
func printSizeAlign(s string, sz uintptr, al uintptr) {
	fmt.Println("Type", s, "size", sz, "alignment", al)
}

func PrintSGFStructSizes() {
	// token.go
	var t Token
	var pos ah.Position
	printSizeAlign("Token", unsafe.Sizeof(t), unsafe.Alignof(t))
	printSizeAlign("ah.Position", unsafe.Sizeof(pos), unsafe.Alignof(pos))
	// tree.go
	var tnt TreeNodeType
	var tni TreeNodeIdx
	var pi PropIdx
	var tn TreeNode
	var pv PropertyValue
	printSizeAlign("TreeNodeType", unsafe.Sizeof(tnt), unsafe.Alignof(tnt))
	printSizeAlign("TreeNodeIdx", unsafe.Sizeof(tni), unsafe.Alignof(tni))
	printSizeAlign("PropIdx", unsafe.Sizeof(pi), unsafe.Alignof(pi))
	printSizeAlign("TreeNode", unsafe.Sizeof(tn), unsafe.Alignof(tn))
	printSizeAlign("PropertyValue", unsafe.Sizeof(pv), unsafe.Alignof(pv))
	// parser.go
	var pt GameTree
	var pr Parser
	var pyi PlayerInfo
	printSizeAlign("GameTree", unsafe.Sizeof(pt), unsafe.Alignof(pt))
	printSizeAlign("Parser", unsafe.Sizeof(pr), unsafe.Alignof(pr))
	printSizeAlign("PlayerInfo", unsafe.Sizeof(pyi), unsafe.Alignof(pyi))
	// sgf.go
	var ff4n FF4Note
	var sgfpnt SGFPropNodeType
	var qt QualifierType
	var vt PropValueType
	var prop Property
	var pdi PropertyDefIdx
	var idca ID_CountArray
	printSizeAlign("FF4Note", unsafe.Sizeof(ff4n), unsafe.Alignof(ff4n))
	printSizeAlign("SGFPropNodeType", unsafe.Sizeof(sgfpnt), unsafe.Alignof(sgfpnt))
	printSizeAlign("QualifierType", unsafe.Sizeof(qt), unsafe.Alignof(qt))
	printSizeAlign("PropValueType", unsafe.Sizeof(vt), unsafe.Alignof(vt))
	printSizeAlign("Property", unsafe.Sizeof(prop), unsafe.Alignof(prop))
	printSizeAlign("PropertyDefIdx", unsafe.Sizeof(pdi), unsafe.Alignof(pdi))
	printSizeAlign("ID_CountArray", unsafe.Sizeof(idca), unsafe.Alignof(idca))
	// scanner.go
	var s Scanner
	printSizeAlign("Scanner", unsafe.Sizeof(s), unsafe.Alignof(s))
	// errors.go
	var eh ErrorHandler
	var el ah.ErrorList
	printSizeAlign("ErrorHandler", unsafe.Sizeof(eh), unsafe.Alignof(eh))
	printSizeAlign("ah.ErrorList", unsafe.Sizeof(el), unsafe.Alignof(el))
	// game.go
	var k Komi
	var r Result
	printSizeAlign("Komi", unsafe.Sizeof(k), unsafe.Alignof(k))
	printSizeAlign("Result", unsafe.Sizeof(r), unsafe.Alignof(r))
}
