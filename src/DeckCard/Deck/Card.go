package Deck

//const SUITS_STR ("\u2660", "\u2665", "\u2666", "\u2663")

type Card struct {
	name string
	rank Rank
	suit Suit
}

func (c *Card) NewCard(rank Rank, suit Suit) *Card {
	return &Card{
		rank: rank,
		suit: suit,
		name: SUITS[c.suit] + RANKS[c.rank]}
}

func (c *Card) Suit() string {
	return SUITS[c.suit]
}

func (c *Card) Rank() string {
	return RANKS[c.rank]
}
