package main

import (
	"DeckCard/Deck"
	//_"DeckCard/Deck"
)

func main() {
	fdeck := Deck.NewDeck()
	println(&fdeck)
	fdeck.DumpCards()
	//println("\u2660", "\u2665", "\u2666", "\u2663")
}
