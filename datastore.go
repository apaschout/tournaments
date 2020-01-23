package tournaments

type Datastore interface {
	SeasonRepository
	PlayerRepository
	DeckRepository
}

type SeasonRepository interface {
	FindAllSeasons() ([]Season, error)
	FindSeasonByID(id SeasonID) (Season, error)
	FindPlayersInSeason(seasID SeasonID) ([]PlayerID, error)
	SavePlayersToSeason(seasID SeasonID, plrs []PlayerID) error
	SaveSeason(seas Season) error
	SeasonNameAvailable(name string) (bool, error)
}

type PlayerRepository interface {
	FindAllPlayers() ([]Player, error)
	FindPlayerByID(id PlayerID) (Player, error)
	SavePlayer(plr Player) error
	PlayerNameAvailable(name string) (bool, error)
	PlayerExists(id PlayerID) (bool, error)
}

type DeckRepository interface {
	FindAllDecks() ([]Deck, error)
	FindDeckByID(id DeckID) (Deck, error)
	SaveDeck(deck Deck) error
	DeckNameAvailable(name string) (bool, error)
}
