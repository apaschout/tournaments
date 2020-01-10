package tournaments

type Datastore interface {
	SeasonRepository
	PlayerRepository
	DeckRepository
	CheckForDuplicateName(table string, name string) error
}

type SeasonRepository interface {
	FindAllSeasons() ([]Season, error)
	FindSeasonByID(id string) (Season, error)
	FindPlayersInSeason(seasID string) ([]Player, error)
	SaveSeason(seas Season) error
}

type PlayerRepository interface {
	FindAllPlayers() ([]Player, error)
	FindPlayerByID(id string) (Player, error)
	SavePlayer(plr Player) error
}

type DeckRepository interface {
	FindAllDecks() ([]Deck, error)
	FindDeckByID(id string) (Deck, error)
	SaveDeck(deck Deck) error
}
