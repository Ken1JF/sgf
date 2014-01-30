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
	"sort"
	//    "github.com/Ken1JF/ah"
	//    "os"
	//    "strconv"
	//    "strings"
	//    "unicode"
)

// Tyoe PlayerInfo is used to record statistics about players in a data base
// of games in .sgf format. It records the number of games, the name of the
// first game encountered, the rank of the player in the first game, the name
// of the last game, and rank at the last game. Note: this assumes the games
// are retrieved in chronological order, from oldest to newest.
type PlayerInfo struct {
	NGames    int
	FirstGame string
	FirstRank string
	LastGame  string
	LastRank  string
}

// The type indexedCount is used to hold the contents of
// a mapStringInt so that it can be sorted.
type indexedCount struct {
	idx string
	cnt int
}

type mapStringInt map[string]int

// Type ByCount implements the sort.Interface based on
// the count, with ties broken by the index string.
type ByCount []indexedCount

func (bc ByCount) Len() int      { return len(bc) }
func (bc ByCount) Swap(i, j int) { bc[i], bc[j] = bc[j], bc[i] }
func (bc ByCount) Less(i, j int) bool {
	if bc[i].cnt > bc[j].cnt {
		return true
	} else {
		if bc[i].cnt < bc[j].cnt {
			return false
		}
	}
	return bc[i].idx < bc[j].idx
}

// The struct DBStatistics is used to record statistics about a data base
// of games in .sgf format.
type DBStatistics struct {
	ID_Counts  ID_CountArray // count of occurances of SGF IDs
	Unkn_Count int           // count of unknown SGF IDs

	FirstBRankNotSet int // count of times setFirstBRank was not set
	FirstWRankNotSet int // count of times setFirstWRank was not set

	HA_map mapStringInt // count the occurances of handicap values
	OH_map mapStringInt // count the occurances of old handicap values

	// Break RE into value and comment:
	RE_map mapStringInt // count the occurances of result values
	RC_map mapStringInt // count the occurances of result comment values

	RU_map     mapStringInt // count the occurances of rules values
	BWRank_map mapStringInt // count the occurances of rank values

	BWPlayer_map map[string]PlayerInfo // accumulate player information
}

// The variable theDBStatistics is used to share statistics among many
// parsers.
var theDBStatistics *DBStatistics

// The function initStats must be called before using a DBStatistics struct.
func (dbStat *DBStatistics) initStats() {
	dbStat.HA_map = make(map[string]int, 100)
	dbStat.OH_map = make(map[string]int, 100)

	// Break RE into value and optional comment:
	dbStat.RE_map = make(map[string]int, 100)
	dbStat.RC_map = make(map[string]int, 100)

	dbStat.RU_map = make(map[string]int, 100)
	dbStat.BWRank_map = make(map[string]int, 100)

	dbStat.BWPlayer_map = make(map[string]PlayerInfo, 100)
}

// The function GameName returns the last portion of a filename for
// use as the game name.
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

// The function SetPlayerRank is called after parsing the first move.
func (p *Parser) SetPlayerRank() {
	// set the rank for the black name
	bn := p.GameTree.GetPB()
	if bn != nil {
		br := p.GameTree.GetBR()
		if br != nil {
			ix := string(bn)
			np, _ := p.DBStats.BWPlayer_map[ix]
			np.LastRank = string(br)
			if p.GameTree.setFirstBRank {
				np.FirstRank = string(br)
			} else {
				p.DBStats.FirstBRankNotSet += 1
			}
			p.DBStats.BWPlayer_map[ix] = np
		} else {
			// fmt.Println("Error, Black Player:",bn,"has nil reank.")
		}
	} else {
		// TODO: this rare. do we need an option to report this?
		// fmt.Println("Error: PB name is nil.")
	}
	// set the rank for the white name
	wn := p.GameTree.GetPW()
	if wn != nil {
		wr := p.GameTree.GetWR()
		if wr != nil {
			ix := string(wn)
			np, _ := p.DBStats.BWPlayer_map[ix]
			np.LastRank = string(wr)
			if p.GameTree.setFirstWRank {
				np.FirstRank = string(wr)
			} else {
				p.DBStats.FirstWRankNotSet += 1
			}
			p.DBStats.BWPlayer_map[ix] = np
		} else {
			// fmt.Println("Error, white Player:",wn,"has nil reank.")
		}
	} else {
		// TODO: this rare. do we need an option to report this?
		// fmt.Println("Error: PW name is nil.")
	}
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

func (dbstat *DBStatistics) reportIDCounts() {
	for i, c := range dbstat.ID_Counts {
		if c > 0 {
			fmt.Printf("Property %s used %d times.\n", string(GetProperty(PropertyDefIdx(i)).ID), c)
		}
	}
	if dbstat.Unkn_Count > 0 {
		fmt.Printf("Property ?Unkn? used %d times.\n", dbstat.Unkn_Count)
	}
	if dbstat.FirstBRankNotSet > 0 {
		fmt.Printf("FirstBRankNotSet %d times.\n", dbstat.FirstBRankNotSet)
	}
	if dbstat.FirstWRankNotSet > 0 {
		fmt.Printf("FirstWRankNotSet %d times.\n", dbstat.FirstWRankNotSet)
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

func (aMap mapStringInt) reportSortedCounts(what string) {
	var bc ByCount
	bc = make([]indexedCount, len(aMap))
	// move map to array
	sum := 0
	i := 0
	for s, n := range aMap {
		bc[i].idx = s
		bc[i].cnt = n
		i += 1
		sum += n
	}
	// sort array
	sort.Sort(bc)
	// report
	for _, n := range bc {
		fmt.Println(what, n.idx, "occurred", n.cnt, "times.")
	}
	fmt.Printf("Total %s games %d with %d different %ss\n", what, sum, len(aMap), what)
}

func (dbstat *DBStatistics) reportHACounts() {
	dbstat.HA_map.reportSortedCounts("Handicap")
}

func (dbstat *DBStatistics) reportOHCounts() {
	dbstat.OH_map.reportSortedCounts("Old Handicap")
}

func (dbstat *DBStatistics) reportRECounts() {
	dbstat.RE_map.reportSortedCounts("Result")
}

func (dbstat *DBStatistics) reportRCCounts() {
	dbstat.RC_map.reportSortedCounts("Result comment")
}

func (dbstat *DBStatistics) reportRUCounts() {
	dbstat.RU_map.reportSortedCounts("Rules")
}

func (dbstat *DBStatistics) reportRankCounts() {
	dbstat.BWRank_map.reportSortedCounts("Rank")
}

func (dbstat *DBStatistics) reportPlayers() {
	// report the BWPlayer map
	//	sum = 0
	//	for s, n := range BWPlayer_map {
	//		fmt.Printf("Player %s occurred %d times, first %s, %s last %s, %s.\n", s, n.NGames, n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	//		sum += n.NGames
	//	}
	//	fmt.Printf("Total players %d with %d different names\n", sum, len(BWPlayer_map))

	// sort the Player names, with counts:
	nPlayers := len(dbstat.BWPlayer_map)
	var playerNames [][]byte
	var playerCount []int
	playerNames = make([][]byte, nPlayers)
	playerCount = make([]int, nPlayers)
	idx := 0
	for s, n := range dbstat.BWPlayer_map {
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
		n, _ := dbstat.BWPlayer_map[string(s)]
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
		n, _ := dbstat.BWPlayer_map[string(playerNames[i])]
		fmt.Printf(" %d : %s, first:  %s, %s, last: %s, %s\n", s, playerNames[i], n.FirstGame, n.FirstRank, n.LastGame, n.LastRank)
	}
}

func ReportSGFCounts() {
	if theDBStatistics != nil {
		theDBStatistics.reportIDCounts()

		theDBStatistics.reportHACounts()
		theDBStatistics.reportOHCounts()
		theDBStatistics.reportRECounts()
		theDBStatistics.reportRCCounts()
		theDBStatistics.reportRUCounts()
		theDBStatistics.reportRankCounts()

		theDBStatistics.reportPlayers()
		theDBStatistics = nil
	}
}
