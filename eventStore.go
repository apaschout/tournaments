package tournaments

import "github.com/cognicraft/event"

type EventStore interface {
}

type SeasonEventRepository interface {
	CreateSeason(ID string) (event.Record, error)
	StartSeason() (event.Record, error)
	EndSeason() (event.Record, error)
	ChangeSeasonName(ID SeasonID, name string) (event.Record, error)
	AddPlayerToSeason(sID SeasonID, pID PlayerID) (event.Record, error)
	RemovePlayerFromSeason(sID SeasonID, pID PlayerID) (event.Record, error)
}

type ES struct {
	*event.Store
}
