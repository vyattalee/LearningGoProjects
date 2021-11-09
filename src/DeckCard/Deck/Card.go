package Deck

type Card struct {
	name  string
	_rank Rank
	_suit Suit
}

func NewCard(rank Rank, suit Suit) Card {
	return Card{
		_rank: rank,
		_suit: suit,
		name:  SUITS[suit] + RANKS[rank]}
}

func (c *Card) Suit() string {
	return SUITS[c._suit]
}

func (c *Card) Rank() string {
	return RANKS[c._rank]
}
