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

type PlayerInfo struct {
	NGames    int
	FirstGame string
	FirstRank string
	LastGame  string
	LastRank  string
}

// TODO:  1. these globals need to be collected into a struct.
// TODO:  2. type Parser needs to have a pointer to this struct.
// TODO:  3. Parser state includes a bool to indicate collecting or not.
// TODO:  4. if collecting, if pointer nil, allocate.
type DBStatistics struct {
	ID_Counts  ID_CountArray
	Unkn_Count int

	FirstBRankNotSet int
	FirstWRankNotSet int

	HA_map map[string]int
	OH_map map[string]int

	// Break RE into value and comment:
	RE_map map[string]int
	RC_map map[string]int

	RU_map     map[string]int
	BWRank_map map[string]int

	BWPlayer_map map[string]PlayerInfo
}

var theDBStatistics *DBStatistics

func (dbStat *DBStatistics) initStats() {
	dbStat.HA_map = make(map[string]int, 100)
	dbStat.OH_map = make(map[string]int, 100)

	// Break RE into value and comment:
	dbStat.RE_map = make(map[string]int, 100)
	dbStat.RC_map = make(map[string]int, 100)

	dbStat.RU_map = make(map[string]int, 100)
	dbStat.BWRank_map = make(map[string]int, 100)

	dbStat.BWPlayer_map = make(map[string]PlayerInfo, 100)
}

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

// SetPlayerRank
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
		fmt.Println("Error: PB name is nil.")
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
		fmt.Println("Error: PW name is nil.")
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

func (dbstat *DBStatistics) reportHACounts() {
	// TODO: sort for canonical output
	// report the HA map
	sum := 0
	for s, n := range dbstat.HA_map {
		fmt.Printf("Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Handicap games %d with %d different handicaps\n", sum, len(dbstat.HA_map))
}

func (dbstat *DBStatistics) reportOHCounts() {
	// TODO: sort for canonical output
	// report the OH map
	sum := 0
	for s, n := range dbstat.OH_map {
		fmt.Printf("Old Handicap %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Old Handicap games %d with %d different settings\n", sum, len(dbstat.OH_map))
}

func (dbstat *DBStatistics) reportRECounts() {
	// TODO: sort for canonical output
	// report the RE map
	sum := 0
	for s, n := range dbstat.RE_map {
		fmt.Printf("Result %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Results %d among %d different settings\n", sum, len(dbstat.RE_map))
}

func (dbstat *DBStatistics) reportRCCounts() {
	// TODO: sort for canonical output
	// report the RC (result comments)
	sum := 0
	for s, n := range dbstat.RC_map {
		fmt.Printf("Result comment %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total Result comments %d with %d different comments\n", sum, len(dbstat.RC_map))
}

func (dbstat *DBStatistics) reportRUCounts() {
	// TODO: sort for canonical output
	// report the RU map
	sum := 0
	for s, n := range dbstat.RU_map {
		fmt.Printf("Rules %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total games with Rules %d with %d different settings\n", sum, len(dbstat.RU_map))
}

func (dbstat *DBStatistics) reportRankCounts() {
	// TODO: sort for canonical output
	// report the BWRank map
	sum := 0
	for s, n := range dbstat.BWRank_map {
		fmt.Printf("Rank %s occurred %d times.\n", s, n)
		sum += n
	}
	fmt.Printf("Total players with Ranks %d among %d different settings\n", sum, len(dbstat.BWRank_map))
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
