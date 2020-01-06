package tournaments

type Datastore interface {
	CreateSeason(seas *Season) error
	FindSeasonByID(id string) (*Season, error)
	FindAllSeasons() (*[]Season, error)

	FindAllPlayers() (*[]Player, error)
	FindPlayerByID(id string) (*Player, error)
	CreatePlayer(plr *Player) error

	FindAllDecks() (*[]Deck, error)
	FindDeckByID(id string) (*Deck, error)
	CreateDeck(deck *Deck) error
}
