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
		name:  REV_SUITS[suit] + REV_RANKS[rank]}
}

func (c *Card) Suit() string {
	return REV_SUITS[c._suit]
}

func (c *Card) Rank() string {
	return REV_RANKS[c._rank]
}
