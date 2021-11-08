package Deck

import (
	"fmt"
	"math/rand"
	"time"
)

type Deck struct {
	//ID    uuid.UUID
	Cards []Card
}

func init(){
	//ranks := list("23456789JQKA")
	//suit := '"\u2660", "\u2665", "\u2666", "\u2663"'.split()
	//
	//Card := namedtuple.New("Card", ["rank", "suit"])
}


func (d *Deck) NewDeck() Deck{
	desk := Deck{}

	//desk.ID = uuid.New()
	for index := 0; index < 4; index++ {

		switch index {
		case 0:
			for value := 1; value < 14; value++ {
				stringvalue := fmt.Sprint(value)
				card1 := Card{ rank: RANKS[index], "h"}
				desk.Cards = append(desk.Cards, card1)
			}
		case 1:
			for value := 1; value < 14; value++ {
				stringvalue := fmt.Sprint(value)
				card1 := Card{(stringvalue + "desk"), value, "d"}
				desk.Cards = append(desk.Cards, card1)
			}
		case 2:
			for value := 1; value < 14; value++ {
				stringvalue := fmt.Sprint(value)
				card1 := Card{(stringvalue + "s"), value, "s"}
				desk.Cards = append(desk.Cards, card1)
			}
		case 3:
			for value := 1; value < 14; value++ {
				stringvalue := fmt.Sprint(value)
				card1 := Card{(stringvalue + "c"), value, "c"}
				desk.Cards = append(desk.Cards, card1)
			}
		}

	}
	desk = desk.Shuffle()
	return desk
}

//Shuffle is a method that will shuffle any deck given using the random unix of the machine it's running on
func (d Deck) Shuffle() Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })

	return d
}