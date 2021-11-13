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

func (c *Card) SuitLevel() Suit {
	return c._suit
}

func (c *Card) RankLevel() Rank {
	return c._rank
}

func (c *Card) Suit() string {
	return SUITS[c._suit]
}

func (c *Card) Rank() string {
	return RANKS[c._rank]
}

// Info : Show human readable card info
func (c Card) Info() (string, string) {
	// translate card to human readable info

	return c.Rank(), c.Suit()
}

// IsCardHigher : Compare 2 card, return true if card 1 number and symbol is higher
func IsCardHigher(c1 Card, c2 Card) bool {
	if c1.RankLevel() >= c2.RankLevel() && c1.SuitLevel() > c2.SuitLevel() {
		return true
	}
	return false
}
