package tournaments

type Match struct {
	Player1 PlayerID `json:"player1"`
	Player2 PlayerID `json:"player2"`
	Winner  PlayerID `json:"winner"`
	P1Count int      `json:"p1Count"`
	P2Count int      `json:"p2Count"`
	Games   []Game   `json:"games"`
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
		matches = append(matches, Match{Player1: plrs[0], Player2: withoutFirst[plrIdx], Games: []Game{Game{}}})
		for i := 1; i < halfSize; i++ {
			plr1 := (round + plrsLen - i) % plrsLen
			plr2 := (round + i) % plrsLen
			matches = append(matches, Match{Player1: withoutFirst[plr1], Player2: withoutFirst[plr2], Games: []Game{Game{}}})
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

//call on TournamentGameEnded
func (trn *Tournament) manageGameWins(match int, game int) {
	m := &trn.Matches[match]
	g := &m.Games[game]
	part1 := trn.getParticipantByID(m.Player1)
	part2 := trn.getParticipantByID(m.Player2)
	if g.Winner == m.Player1 {
		m.P1Count++
		part1.GameWins++
	} else if g.Winner == m.Player2 {
		m.P2Count++
		part2.GameWins++
	}
	part1.Games++
	part2.Games++
	if m.P1Count < trn.GamesToWin && m.P2Count < trn.GamesToWin {
		m.Games = append(m.Games, Game{})
	}
	//  else {
	// 	if m.P1Count == trn.GamesToWin {
	// 		m.Winner = m.Player1
	// 		part1.MatchWins++
	// 	} else if m.P2Count == trn.GamesToWin {
	// 		m.Winner = m.Player2
	// 		part2.MatchWins++
	// 	}
	// 	if m.P1Count == m.P2Count {
	// 		m.Draw = true
	// 	}
	// 	m.Ended = true

	// 	part1.Matches++
	// 	part2.Matches++

	// }
}

func (trn *Tournament) manageMatchWin(match int) {
	m := &trn.Matches[match]
	part1 := trn.getParticipantByID(m.Player1)
	part2 := trn.getParticipantByID(m.Player2)
	part1.Matches++
	part2.Matches++
	if m.Winner == part1.Player {
		part1.MatchWins++
	} else if m.Winner == part2.Player {
		part2.MatchWins++
	}
}
