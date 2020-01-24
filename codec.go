package tournaments

import "github.com/cognicraft/event"

func Codec() (*event.Codec, error) {
	c := event.NewCodec()
	err = c.Register("season:created", SeasonCreated{})
	if err != nil {
		return nil, err
	}
	err = c.Register("season:name-changed", SeasonNameChanged{})
	if err != nil {
		return nil, err
	}
	err = c.Register("season:started", SeasonStarted{})
	if err != nil {
		return nil, err
	}
	err = c.Register("season:ended", SeasonEnded{})
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
