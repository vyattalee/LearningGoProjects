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
			println(black_joker.name, suit, 52)
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

func (d *Deck) DumpCards() {
	for i, card := range d.Cards {
		println("第", i+1, "张牌:", card.name)

	}
}
