/*
 *  File:		src/github.com/Ken1JF/sgf/parserdbstat.go
 *
 *  Created by Ken Friedenbach on Jan 27, 2014
 *  Copyright 2014 Ken Friedenbach. All rights reserved.
 *
 */

// The SGF parser records optional information useful when
// reading data bases, such as the GoGoD database.
package sgf

import (
	"bytes"
	"fmt"
	//    "github.com/Ken1JF/ah"
	//    "os"
	//    "strconv"
	//    "strings"
	//    "unicode"
)

// TODO:  1. these globals need to be collected into a struct.
// TODO:  2. type Parser needs to have a pointer to this struct.
// TODO:  3. Parser state includes a bool to indicate collecting or not.
// TODO:  4. if collecting, if pointer nil, allocate.
// TODO:  5. move all of this into another file: parsedb.go or something
var ID_Counts ID_CountArray
var Unkn_Count int

var HA_map map[string]int = make(map[string]int, 100)
var OH_map map[string]int = make(map[string]int, 100)

// Break RE into value and comment:
var RE_map map[string]int = make(map[string]int, 100)
var RC_map map[string]int = make(map[string]int, 100)

var RU_map map[string]int = make(map[string]int, 100)
var BWRank_map map[string]int = make(map[string]int, 100)

type PlayerInfo struct {
	NGames    int
	FirstGame string
	FirstRank string
	LastGame  string
	LastRank  string
}

var BWPlayer_map map[string]PlayerInfo = make(map[string]PlayerInfo, 100)

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

func reportIDCounts() {
	for i, c := range ID_Counts {
		if c > 0 {
			fmt.Printf("Property %s used %d times.\n", string(GetProperty(PropertyDefIdx(i)).ID), c)
		}
	}
	if Unkn_Count > 0 {
		fmt.Printf("Property ?Unkn? used %d times.\n", Unkn_Count)
	}
}

// TODO: sort by second field (last name) if present
func gtr(a []byte, b []byte) bool {
	idx := 0
	for (idx < len(a)) && (idx < len(b)) {
		if a[idx] > b[idx] {
			return true
		} else if a[idx] < b[idx] {
			return false
		}
		idx += 1
	}
	if len(a) > len(b) {
		return true
	}
	return false
}

func ReportSGFCounts() {
	reportIDCounts()
	// report the HA map
	sum := 0
	for s, n := range HA_map {
		fmt.Printf("Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Handicap games %d with %d different handicaps\n", sum, len(HA_map))

	// report the OH map
	sum = 0
	for s, n := range OH_map {
		fmt.Printf("Old Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Old Handicap games %d with %d different settings\n", sum, len(OH_map))

	// report the RE map
	sum = 0
	for s, n := range RE_map {
		fmt.Printf("Result %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Results %d among %d different settings\n", sum, len(RE_map))

	// report the RC (result comments)
	sum = 0
	for s, n := range RC_map {
		fmt.Printf("Result comment %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Result comments %d with %d different comments\n", sum, len(RC_map))

	// report the RU map
	sum = 0
	for s, n := range RU_map {
		fmt.Printf("Rules %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Rules %d with %d different settings\n", sum, len(RU_map))

	// report the BWRank map
	sum = 0
	for s, n := range BWRank_map {
		fmt.Printf("Rank %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total players with Ranks %d among %d different settings\n", sum, len(BWRank_map))

	// report the BWPlayer map
	//	sum = 0
	//	for s, n := range BWPlayer_map {
	//		fmt.Printf("Player %s occurred %d times, first %s, %s last %s, %s.\n", s, n.NGames, n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	//		sum += n.NGames
	//	}
	//	fmt.Printf("Total players %d with %d different names\n", sum, len(BWPlayer_map))

	// sort the Player names, with counts:
	nPlayers := len(BWPlayer_map)
	var playerNames [][]byte
	var playerCount []int
	playerNames = make([][]byte, nPlayers)
	playerCount = make([]int, nPlayers)
	idx := 0
	for s, n := range BWPlayer_map {
		playerNames[idx] = make([]byte, len(s))
		_ = copy(playerNames[idx], s)
		playerCount[idx] = n.NGames
		idx += 1
	}
	// Sort them alphabetically:
	for ix := 0; ix < nPlayers; ix++ {
		for iy := ix; iy < nPlayers; iy++ {
			if gtr(playerNames[ix], playerNames[iy]) {
				playerNames[ix], playerNames[iy] = playerNames[iy], playerNames[ix]
				playerCount[ix], playerCount[iy] = playerCount[iy], playerCount[ix]
			}
		}
	}
	for i, s := range playerNames {
		n, _ := BWPlayer_map[string(s)]
		fmt.Printf("Player %s: %d, first: %s, %s, last: %s, %s\n", s, playerCount[i], n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	}

	// Sort them numerically:
	for ix := 0; ix < nPlayers; ix++ {
		for iy := ix; iy < nPlayers; iy++ {
			if playerCount[ix] < playerCount[iy] {
				playerNames[ix], playerNames[iy] = playerNames[iy], playerNames[ix]
				playerCount[ix], playerCount[iy] = playerCount[iy], playerCount[ix]
			}
		}
	}
	for i, s := range playerCount {
		n, _ := BWPlayer_map[string(playerNames[i])]
		fmt.Printf(" %d : %s, first:  %s, %s, last: %s, %s\n", s, playerNames[i], n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	}
}
