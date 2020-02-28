package tournaments

type Match struct {
	Player1 PlayerID `json:"player1"`
	Player2 PlayerID `json:"player2"`
	Winner  PlayerID `json:"winner"`
	Games   []Game   `json:"game"`
	Draw    bool     `json:"draw"`
	Ended   bool     `json:"ended"`
}

type Game struct {
	Winner PlayerID `json:"winner"`
	Draw   bool     `json:"draw"`
	Ended  bool     `json:"ended"`
}

func (trn *Tournament) MakeMatches() {
	plrs := []PlayerID{}
	for _, par := range trn.Players {
		plrs = append(plrs, par.Player)
	}
	if len(plrs)%2 != 0 {
		// create dummy player, later delete matches with dummy
		plrs = append(plrs, "x")
	}

	numRounds := len(plrs) - 1
	halfSize := len(plrs) / 2

	withoutFirst := plrs[1:]
	plrsLen := len(withoutFirst)
	matches := []Match{}

	for round := 0; round < numRounds; round++ {
		plrIdx := round % plrsLen
		matches = append(matches, Match{Player1: plrs[0], Player2: withoutFirst[plrIdx]})
		for i := 1; i < halfSize; i++ {
			plr1 := (round + plrsLen - i) % plrsLen
			plr2 := (round + i) % plrsLen
			matches = append(matches, Match{Player1: withoutFirst[plr1], Player2: withoutFirst[plr2]})
		}
	}
	trn.Matches = deleteDummyMatches(matches)
}

func deleteDummyMatches(matches []Match) []Match {
	for i := 0; i < len(matches); i++ {
		if matches[i].Player1 == "x" || matches[i].Player2 == "x" {
			matches = append(matches[:i], matches[i+1:]...)
			i--
		}
	}
	return matches
}

//call on TournamentMatchesCreated + on TournamentGameEnded
func (trn *Tournament) handleGames() {
	for i, m := range trn.Matches {
		p1Count := 0
		p2Count := 0
		for _, g := range m.Games {
			if g.Winner == m.Player1 {
				p1Count++
			} else if g.Winner == m.Player2 {
				p2Count++
			}
		}
		if p1Count < trn.GamesToWin || p2Count < trn.GamesToWin {
			m.Games = append(m.Games, Game{})
		} else {
			var wnr PlayerID
			var draw bool
			if p1Count == trn.GamesToWin {
				wnr = m.Player1
			} else if p2Count == trn.GamesToWin {
				wnr = m.Player2
			}
			if p1Count == p2Count {
				draw = true
			}
			trn.EndMatch(i, wnr, draw)
		}
	}
}
