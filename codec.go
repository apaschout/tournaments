package tournaments

import "github.com/cognicraft/event"

func Codec() (*event.Codec, error) {
	c := event.NewCodec()
	err = c.Register("tournament:created", TournamentCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:deleted", TournamentDeleted{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:name-changed", TournamentNameChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:gamestowin-changed", TournamentGamesToWinChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:phase-changed", TournamentPhaseChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:format-changed", TournamentFormatChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:player-registered", TournamentPlayerRegistered{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:player-dropped", TournamentPlayerDropped{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:started", TournamentStarted{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:ended", TournamentEnded{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:matches-created", TournamentMatchesCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:game-ended", TournamentGameEnded{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:match-ended", TournamentMatchEnded{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:created", PlayerCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:name-changed", PlayerNameChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:matches-played", PlayerMatchPlayed{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:match-won", PlayerMatchWon{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:tournament-registered", PlayerTournamentRegistered{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:games-played", PlayerGamePlayed{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:game-won", PlayerGameWon{})
	if err != nil {
		return nil, err
	}
	return c, nil
}
