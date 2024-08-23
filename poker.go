package main

import (
	"cmp"
	"crypto/rand"
	"math/big"
	"slices"
)

const (
	faces = "23456789tjqka"
	suits = "shdc"
)

type Score int

const (
    Highcard Score = (1000 * iota)
    Pair 
    TwoPair
    Trips
    Straight
    Flush
    FullHouse
    Quads
    StraightFlush 
    RoyalFlush
)

type Deck []*Card

func NewDeck() *Deck {
	deck := &Deck{}
	for _, suit := range suits {
		for _, face := range faces {
			*deck = append(*deck, &Card{
				face:     byte(face),
				suit:     byte(suit),
				priority: 0,
			})
		}
	}
	return deck
}

func NewDeckFromString(input string) *Deck {
	deck := &Deck{}
	i := 0
	for i < len(input)-1 {
		*deck = append(*deck, &Card{
			face:     input[i : i+2][0],
			suit:     input[i : i+2][1],
			priority: 0,
		})
		i += 2
	}
	return deck
}

func (d Deck) String() string {
	s := ""
	for _, c := range d {
		s += c.String()
	}
	return s
}

func (d Deck) Len() int {
	return len(d)
}

func (d Deck) Empty() bool {
	return len(d) == 0
}

func (d Deck) Shuffle() Deck {
	for _, c := range d {
		v, _ := rand.Int(rand.Reader, big.NewInt(64))
		c.priority = int(v.Int64())
	}
	slices.SortFunc(d, func(a, b *Card) int {
		return cmp.Compare(a.priority, b.priority)
	})
	return d
}

func (d *Deck) DealCards(n int) Deck {
	if d.Empty() {
		return nil
	}

	hand := (*d)[:n]
	*d = (*d)[n:]
	return hand
}

func aceHighSort(a, b *Card) int {
	return a.Score() - b.Score()
}

func aceLowSort(a, b *Card) int {
	return (a.Score() % 14) - (b.Score() % 14)
}

func isStraight(hand Deck) (bool, int) {
	n := hand.Len()
	if n < 5 {
		return false, -1
	}

	runningCount := hand[0].Score()
	for i := 0; i < n; i++ {
		c := hand[i]
		if (runningCount % 14) != (c.Score() % 14) {
			return false, -1
		}

		if c.Score() == 14 {
			runningCount += 2
		} else {
			runningCount++
		}
	}

	return true, runningCount % 14
}

func (d Deck) IsStraight(hand Deck) (bool, int) {
	var aceHigh Deck
	aceHigh = append(aceHigh, d...)
	aceHigh = append(aceHigh, hand...)
	slices.SortFunc(aceHigh, aceHighSort)

	aceHigh = slices.Clip(slices.CompactFunc(aceHigh, func(a, b *Card) bool {
		return a.Score() == b.Score()
	}))
	aceLow := slices.Clone(aceHigh)
	slices.SortFunc(aceLow, aceLowSort)

	for i := 1; i < 3; i++ {
		if yes, score := isStraight(aceHigh[i:]); yes {
			return yes, score
		}

		if yes, score := isStraight(aceHigh[:aceHigh.Len()-i]); yes {
			return yes, score
		}

		if yes, score := isStraight(aceLow[i:]); yes {
			return yes, score
		}

		if yes, score := isStraight(aceLow[:aceLow.Len()-i]); yes {
			return yes, score
		}
	}

	return false, -1
}

func (d Deck) IsFlush(hand Deck) (bool, bool, int) {
	var cards Deck
	cards = append(cards, d...)
	cards = append(cards, hand...)

	flushgroups := make(map[byte]Deck)
	for _, c := range cards {
		flushgroups[c.suit] = append(flushgroups[c.suit], c)
	}

	for _, groups := range flushgroups {
		if groups.Len() > 4 {
			slices.SortFunc(groups, aceHighSort)
			isStraight, _ := isStraight(groups)
			return true, isStraight, groups[groups.Len()-1].Score()
		}
	}
	return false, false, -1
}

func (d Deck) getPairs() map[int]Deck {
	groups := make(map[int]Deck)
	for _, c := range d {
		groups[c.Score()] = append(groups[c.Score()], c)
	}

	pairs := make(map[int]Deck)
	for _, g := range groups {
		pairs[g.Len()] = append(pairs[g.Len()], g[0])
	}

	for _, set := range pairs {
		slices.SortFunc(set, aceHighSort)
	}

	return pairs
}

func (d Deck) AnalyzeHand(hand Deck) int {
	var cards Deck
	cards = append(cards, d...)
	cards = append(cards, hand...)

	handScore := hand[0].Score() + hand[1].Score()
	flush, isStraightFlush, flushScore := d.IsFlush(hand)
	straight, straightScore := d.IsStraight(hand)

	groups := cards.getPairs()
	hasQuads := groups[4]
	hasTrips := groups[3]
	hasPair := groups[2]

	switch {
	case isStraightFlush:
		return int(StraightFlush) + handScore
	case !hasQuads.Empty():
		return int(Quads) + hasQuads[hasQuads.Len()-1].Score()
	case !hasTrips.Empty() && !hasPair.Empty():
		return int(FullHouse) + hasTrips[hasTrips.Len()-1].Score()
	case flush:
		return int(Flush) + flushScore
	case straight:
		return int(Straight) + straightScore
	case !hasTrips.Empty():
		return int(Trips) + hasTrips[hasTrips.Len()-1].Score() + handScore
	case !hasPair.Empty():
		// Two pair or pair + board pair
		if hasPair.Len() > 1 {
			return int(TwoPair) + 7*hasPair[hasPair.Len()-2].Score() + 14*hasPair[hasPair.Len()-1].Score() + handScore
		}
		return int(Pair) + hasPair[hasPair.Len()-1].Score() + handScore
	default:
		return int(Highcard) + handScore
	}
}

type Card struct {
	face     byte
	suit     byte
	priority int
}

func (c Card) String() string {
	return string(c.face) + string(c.suit)
}

// Title returns a string of "card face fullname of card suit" to be used in an HTML tags title attribute.
// E.g. Ace of Spades, Ten of Diamonds
func (c Card) Title() string {
    return c.faceName() + " of " + c.suitName()
}

func (c Card) suitName() string {
	switch c.suit {
	case 'c':
		return "Clubs"
	case 's':
		return "Spades"
	case 'd':
		return "Diamonds"
	case 'h':
		return "Hearts"
    default:
        return "Invalid"
	}
}

func (c Card) faceName() string {
	switch c.face {
	case 'a':
		return "Ace"
	case 'k':
		return "King"
	case 'q':
		return "Queen"
	case 'j':
		return "Jack"
	case 't':
		return "Ten"
	case '9':
		return "Nine"
	case '8':
		return "Eight"
	case '7':
		return "Seven"
	case '6':
		return "Six"
	case '5':
		return "Five"
	case '4':
		return "Four"
	case '3':
		return "Three"
	case '2':
		return "Two"
	default:
		return "Invalid"
	}
}

func (c Card) Score() int {
	switch c.face {
	case 'a':
		return 14
	case 'k':
		return 13
	case 'q':
		return 12
	case 'j':
		return 11
	case 't':
		return 10
	case '9':
		return 9
	case '8':
		return 8
	case '7':
		return 7
	case '6':
		return 6
	case '5':
		return 5
	case '4':
		return 4
	case '3':
		return 3
	case '2':
		return 2
	default:
		return 0
	}
}
