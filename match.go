package tournaments

import "fmt"

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
	plr1, err := trn.Server.getPlayerByID(m.Player1)
	if err != nil {
		fmt.Println(err)
		return
	}
	plr2, err := trn.Server.getPlayerByID(m.Player2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(plr1, plr2)
	if g.Winner == m.Player1 {
		m.P1Count++
		part1.GameWins++
		plr1.GamesWon++
	} else if g.Winner == m.Player2 {
		m.P2Count++
		part2.GameWins++
		plr2.GamesWon++
	}
	part1.Games++
	part2.Games++

	plr1.GamesPlayed++
	plr2.GamesPlayed++
	if m.P1Count < trn.GamesToWin && m.P2Count < trn.GamesToWin {
		m.Games = append(m.Games, Game{})
	} else {
		if m.P1Count == trn.GamesToWin {
			m.Winner = m.Player1
			part1.MatchWins++
			plr1.MatchesWon++
		} else if m.P2Count == trn.GamesToWin {
			m.Winner = m.Player2
			part2.MatchWins++
			plr2.MatchesWon++
		}
		if m.P1Count == m.P2Count {
			m.Draw = true
		}
		m.Ended = true

		part1.Matches++
		part2.Matches++

		// err = plr1.IncrementMatches()
		// if err != nil {
		// 	log.Println(err)
		// }
		// err = plr2.IncrementMatches()
		// if err != nil {
		// 	log.Println(err)
		// }
	}
}
