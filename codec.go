package tournaments

import "github.com/cognicraft/event"

func Codec() (*event.Codec, error) {
	c := event.NewCodec()
	err = c.Register("tournament:created", TournamentCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:name-changed", TournamentNameChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("tournament:player-registered", TournamentPlayerRegistered{})
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

	err = c.Register("player:created", PlayerCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("player:name-changed", PlayerNameChanged{})
	if err != nil {
		return nil, err
	}
	return c, nil
}
