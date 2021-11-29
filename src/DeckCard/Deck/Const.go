package Deck

// Suit represents the _suit of the card (spade, heart, diamond, club)
type Suit int

// Rank represents the face of the card (ace, two...queen, king)
type Rank int

// CompareResult is the custom type returned when comparing cards
type CompareResult int

// Constants for Suit ♠♥♦♣
const (
	CLUB Suit = iota
	DIAMOND
	HEART
	SPADE
	Joker
	JOKER
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
	SUITS = map[Suit]string{CLUB: "\u2663", DIAMOND: "\u2666", HEART: "\u2665", SPADE: "\u2660", Joker: "\u265A", JOKER: "\u265B"}
	RANKS = map[Rank]string{ACE: "A", TWO: "2", THREE: "3", FOUR: "4", FIVE: "5", SIX: "6", SEVEN: "7", EIGHT: "8", NINE: "9", TEN: "10", JACK: "J", QUEEN: "Q", KING: "K"}
	//REV_SUITS = map[Suit]string{CLUB: "\u2663", DIAMOND: "\u2666", HEART: "\u2665", SPADE: "\u2660", Joker: "U+1F0DF", JOKER: "U+1F0BF"}
	//REV_RANKS = map[Rank]string{ACE: "A", TWO: "2", THREE: "3", FOUR: "4", FIVE: "5", SIX: "6", SEVEN: "7", EIGHT: "8", NINE: "9", TEN: "10", JACK: "J", QUEEN: "Q", KING: "K"}
	//SUITS     = map[string]Suit{"\u2663": CLUB, "\u2666": DIAMOND, "\u2665": HEART, "\u2660": SPADE, "U+1F0DF":Joker, "U+1F0BF": JOKER}
	//RANKS     = map[string]Rank{"A": ACE, "2": TWO, "3": THREE, "4": FOUR, "5": FIVE, "6": SIX, "7": SEVEN, "8": EIGHT, "9": NINE, "10": TEN, "J": JACK, "Q": QUEEN, "K": KING}
)
