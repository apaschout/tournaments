package tournaments

import "github.com/cognicraft/event"

type Projection interface {
	TournamentRepository
	PlayerRepository
	DeckRepository
	On(rec event.Record)
	GetVersion() uint64
}

type TournamentRepository interface {
	FindAllTournaments() ([]Tournament, error)
	FindTournamentByID(id TournamentID) (Tournament, error)
	IsTournamentNameAvailable(name string) (bool, error)
}

type PlayerRepository interface {
	FindAllPlayers() ([]Player, error)
	FindPlayerByID(id PlayerID) (Player, error)
	IsPlayerNameAvailable(name string) (bool, error)
	PlayerExists(id PlayerID) (bool, error)
}

type DeckRepository interface {
	FindAllDecks() ([]Deck, error)
	FindDeckByID(id DeckID) (Deck, error)
	IsDeckNameAvailable(name string) (bool, error)
}
