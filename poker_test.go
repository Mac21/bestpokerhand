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

func validateStraightTest(deck, hand Deck, expected bool, t *testing.T) {
    if yes, _ := deck.IsStraight(hand); yes != expected {
        t.Fatalf("Board: %v, Hand: %v IsStraight expected %v", deck, hand, expected)
	}
}

func validateHandStrength(board, winning, losing Deck, t *testing.T) {
    winningScore := board.AnalyzeHand(winning)
    losingScore := board.AnalyzeHand(losing)
    winningLess :=  winningScore < losingScore
    if winningLess {
        t.Fatalf("Board: %v, Winning: %v(%d), Losing: %v(%d) expected winning > losing got winning < losing", board, winning, winningScore, losing, losingScore)
    }
}

func Test67Straight(t *testing.T) {
	deck := buildBoardAndHand("as2c3d4h5h6h7h")
	board := deck.DealCards(5)
	hand := deck
	validateStraightTest(board, hand, true, t)
}

func TestWheelStraight(t *testing.T) {
	deck := buildBoardAndHand("as2s3s4s9s5sqs")
	board := deck.DealCards(5)
	hand := deck
	validateStraightTest(board, hand, true, t)
}

func TestRoyalFlush(t *testing.T) {
    deck := buildBoardAndHand("asksqsjs2hts2c")
    board := deck.DealCards(5)
    hand := deck
	validateStraightTest(board, hand, true, t)
}

func TestNoStraight(t *testing.T) {
    deck := buildBoardAndHand("3s7cthasqd3h3c")
    board := deck.DealCards(5)
    hand := deck
	validateStraightTest(board, hand, false, t)
}

func TestOutsideStraight(t *testing.T) {
    deck := buildBoardAndHand("9c2ctcqdjhkd7d")
    board := deck.DealCards(5)
    hand := deck
	validateStraightTest(board, hand, true, t)
}

func TestTripsBeatsTwoPair(t *testing.T) {
    deck := buildBoardAndHand("asacqsjs2had4hqhjh")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestQuadsBeatsFullhouse(t *testing.T) {
    deck := buildBoardAndHand("asacqsjs2hadahqhqc")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTripsBeatsTrips(t *testing.T) {
    deck := buildBoardAndHand("asacqsjs2had4hah4c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestHighcardBeatsHighcard(t *testing.T) {
    deck := buildBoardAndHand("jc7s9sqcad4s3d4c2c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestPairBeatsHighcard(t *testing.T) {
    deck := buildBoardAndHand("kh5d2s3s8h3d9hjsas")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestPairBeatsPair(t *testing.T) {
    deck := buildBoardAndHand("kh5d2s3c8h8d9h3das")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTwoPairBeatsPair(t *testing.T) {
    deck := buildBoardAndHand("kh5d2s3c8h8d2h3das")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTwoPairBeatsTwoPair(t *testing.T) {
    deck := buildBoardAndHand("kh5d2s3c8h3dks8d2h")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTwoPairBeatsPairBug(t *testing.T) {
    deck := buildBoardAndHand("5c7h3h7sah3d4c8dks")
    board := deck.DealCards(5)
    // Two pair
    winning := deck.DealCards(2)
    // Pair
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTripsBeatsTwoPairBug(t *testing.T) {
    deck := buildBoardAndHand("ac9cjdah4cas2s7h9s")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestTwopairBeatsTwoPairBug(t *testing.T) {
    deck := buildBoardAndHand("kh4h4c5cks7c7d9cqs")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestAceHighBeatsHighcard(t *testing.T) {
    deck := buildBoardAndHand("tsqc3s7h2dasks6h5d")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestFlushBeatsPair(t *testing.T) {
    deck := buildBoardAndHand("8s2s7sks4s8caskd9h")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestFlushBeatsPocketPair(t *testing.T) {
    deck := buildBoardAndHand("2h8cjc4h3cqckcasah")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestPairBeatsHighcardStraightBug(t *testing.T) {
    deck := buildBoardAndHand("8h5htdas7s4h8d6dkh")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestStraightBeatsStraight(t *testing.T) {
    deck := buildBoardAndHand("2c3s4djsah5s6d5hkc")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestFlushBeatsFlush(t *testing.T) {
    deck := buildBoardAndHand("2c3c4c5hthqc5c6c7c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestFlushBeatsStraight(t *testing.T) {
    deck := buildBoardAndHand("2c3c4c5hthqc5c6h7h")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestStraightFlushBeatsFlush(t *testing.T) {
    deck := buildBoardAndHand("2c3c4c5cth6c7cqc9c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestQuadsBeatsFlush(t *testing.T) {
    deck := buildBoardAndHand("2c3c4c5c2h2s2dac9c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}

func TestStraightFlushBeatsStraightFlush(t *testing.T) {
    deck := buildBoardAndHand("2c3c4c5cth6c7cac9c")
    board := deck.DealCards(5)
    winning := deck.DealCards(2)
    losing := deck
	validateHandStrength(board, winning, losing,  t)
}
