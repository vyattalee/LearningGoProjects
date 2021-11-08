package Deck

// Suit represents the suit of the card (spade, heart, diamond, club)
type Suit int

// Rank represents the face of the card (ace, two...queen, king)
type Rank int

// Constants for Suit ♠♥♦♣
const (
	CLUB Suit = iota
	DIAMOND
	HEART
	SPADE
)

// Constants for Rank
const (
	ACE Rank = iota
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	JACK
	QUEEN
	KING
)

// Global Variables representing the default suits and Ranks in a deck of cards
var (
	SUITS = map[Suit]string{CLUB: "\u2663", DIAMOND: "\u2666", HEART: "\u2665", SPADE: "\u2660"}
	RANKS = map[Rank]string{ACE: "A", TWO: "2", THREE: "3", FOUR: "4", FIVE: "5", SIX: "6", SEVEN: "7", EIGHT: "8", NINE: "9", TEN: "10", JACK: "J", QUEEN: "Q", KING: "K"}
)
