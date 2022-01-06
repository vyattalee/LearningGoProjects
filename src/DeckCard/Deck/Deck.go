package Deck

import (
	. "github.com/google/uuid"
	"math/rand"
	"time"
)

type Deck struct {
	ID    UUID
	Cards []Card
}

func init() {
	println("Deck.go init")
	//ranks := list("23456789JQKA")
	//_suit := '"\u2660", "\u2665", "\u2666", "\u2663"'.split()
	//
	//Card := namedtuple.New("Card", ["_rank", "_suit"])
}

func NewDeck() Deck {
	deck := Deck{}

	deck.ID = New()

	println("A:", ACE, "  CLUB:", CLUB)
	//for range generate the random
	for suit, svalue := range SUITS {
		if suit < Joker {
			for rank, rvalue := range RANKS {
				card := NewCard(rank, suit)
				println(card.name, suit, rank, svalue, rvalue)
				deck.Cards = append(deck.Cards, card)
			}
		} else if suit == Joker {
			white_joker := NewCard(52, suit)
			println(white_joker.name, suit, 52)
			deck.Cards = append(deck.Cards, white_joker)
		} else {
			black_joker := NewCard(53, suit)
			println(black_joker.name, suit, 53)
			deck.Cards = append(deck.Cards, black_joker)
		}

	}

	deck = deck.Shuffle()
	return deck
}

//Shuffle is a method that will shuffle any deck given using the random unix of the machine it's running on
func (d Deck) Shuffle() Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })

	return d
}

// ShufflePerm uses rand.Perm instead of the many calls to rand.Intn.
//  When compared to the current implementation:
//
//  benchmark                        old ns/op     new ns/op     delta
//  BenchmarkTinyDeckShuffle-8       524           537           +2.48%
//  BenchmarkSmallDeckShuffle-8      1119          1070          -4.38%
//  BenchmarkMediumDeckShuffle-8     1611          1626          +0.93%
//  BenchmarkDeckShuffle-8           2115          2194          +3.74%
//  BenchmarkLargeDeckShuffle-8      21301         21408         +0.50%
//
//  Conclusion: Not Recommended
func (d *Deck) ShufflePerm() {
	N := len(d.Cards)
	perm := rand.Perm(N)
	for i := 0; i < N; i++ {
		d.Cards[perm[i]], d.Cards[i] = d.Cards[i], d.Cards[perm[i]]
	}
}

// NumberOfCards is a utility function that tells you how many cards are left in the deck
func (d *Deck) NumberOfCards() int {
	return len(d.Cards)
}

// Deal distributes cards to other decks/hands
func (d *Deck) Deal(cards int, hands ...*Deck) {
	for i := 0; i < cards; i++ {
		for _, hand := range hands {
			card := d.Cards[0]
			d.Cards = d.Cards[1:]
			hand.Cards = append(hand.Cards, card)
		}
	}
}
func (d *Deck) DumpCards() {
	for i, card := range d.Cards {
		println("第", i+1, "张牌:", card.name)

	}
}

// DefaultCompare is the default comparison function
// Currently not used in any games.
func (d *Deck) DefaultCompare(i, j Card) CompareResult {
	if i.Rank() > j.Rank() {
		return 1
	}

	if i.Rank() < j.Rank() {
		return -1
	}

	if i.Suit() > j.Suit() {
		return 1
	}

	if i.Suit() < j.Suit() {
		return -1
	}

	return 0
}

// DrawCards card from deck
func (d *Deck) DrawCards(numberOfCards int) Deck {

	newD := NewDeck()
	var draw = (*d).Cards[len((*d).Cards)-numberOfCards : len((*d).Cards)]
	(*d).Cards = (*d).Cards[:len((*d).Cards)-numberOfCards]

	newD.Cards = draw

	return newD
}

// AddCard : add card to top of deck
func (d *Deck) AddCard(card Card) {
	(*d).Cards = append((*d).Cards, card)
}
