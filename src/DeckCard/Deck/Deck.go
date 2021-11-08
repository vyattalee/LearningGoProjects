package Deck

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Deck struct {
	ID    uuid.UUID
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

	deck.ID = uuid.New()

	println("A:", ACE, "  CLUB:", CLUB)

	for _, suit := range SUITS {
		for _, rank := range RANKS {
			card := NewCard(rank, suit)
			println(card.name)
			deck.Cards = append(deck.Cards, card)
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

func (d *Deck) DumpCards() {
	for card := range d.Cards {
		println(card)

	}
}
