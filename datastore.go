package tournaments

type Datastore interface {
	SeasonRepository
	PlayerRepository
	DeckRepository
}

type SeasonRepository interface {
	FindAllSeasons() ([]Season, error)
	FindSeasonByID(id string) (Season, error)
	SaveSeason(seas Season) error
	UpdateSeason(seas Season) error
}

type PlayerRepository interface {
	FindAllPlayers() ([]Player, error)
	FindPlayerByID(id string) (Player, error)
	SavePlayer(plr Player) error
	UpdatePlayer(plr Player) error
}

type DeckRepository interface {
	FindAllDecks() ([]Deck, error)
	FindDeckByID(id string) (Deck, error)
	SaveDeck(deck Deck) error
}
