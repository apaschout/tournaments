package tournaments

type Projection interface {
	SeasonRepository
	PlayerRepository
	DeckRepository
}

type SeasonRepository interface {
	FindAllSeasons() ([]Season, error)
	FindSeasonByID(id SeasonID) (Season, error)
	IsSeasonNameAvailable(name string) (bool, error)
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
