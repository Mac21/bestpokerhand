package main

import (
	"testing"
)

func buildBoardAndHand(input string) Deck {
	deck := Deck{}
	i := 0
	for i < len(input)-1 {
		deck = append(deck, &Card{
			face:     input[i : i+2][0],
			suit:     input[i : i+2][1],
			priority: 0,
		})
		i += 2
	}
	return deck
}

func validateTest(deck, hand Deck, expected bool, t *testing.T) {
	if deck.IsStraight(hand) != expected {
		t.Fatalf("deck %v with hand %v failed expected %v", deck, hand, expected)
	}
}

func Test67Straight(t *testing.T) {
	deck := buildBoardAndHand("as2c3d4h5h6h7h")
	board := deck.DealCards(5)
	hand := deck
	validateTest(board, hand, true, t)
}

func TestWheelStraight(t *testing.T) {
	deck := buildBoardAndHand("as2s3s4s9s5sqs")
	board := deck.DealCards(5)
	hand := deck
	validateTest(board, hand, true, t)
}

func TestRoyalFlush(t *testing.T) {
    deck := buildBoardAndHand("asksqsjs2hts2c")
    board := deck.DealCards(5)
    hand := deck
	validateTest(board, hand, true, t)
}

func TestNoStraight(t *testing.T) {
    deck := buildBoardAndHand("3s7cthasqd3h3c")
    board := deck.DealCards(5)
    hand := deck
	validateTest(board, hand, false, t)
}

func TestOutsideStraight(t *testing.T) {
    deck := buildBoardAndHand("9c2ctcqdjhkd7d")
    board := deck.DealCards(5)
    hand := deck
	validateTest(board, hand, true, t)
}
