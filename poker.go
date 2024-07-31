package main

import (
	"cmp"
	"crypto/rand"
	"math/big"
	"slices"
	"strings"
)

const (
	faces = "23456789tjqka"
	suits = "shdc"
)

var (
	suitToUnicode = map[byte]rune{
		's': 0x1F0A0,
		'h': 0x1F0B0,
		'd': 0x1F0C0,
		'c': 0x1F0D0,
	}
	faceToUnicode = map[byte]rune{
		'a': 0x00001,
		'2': 0x00002,
		'3': 0x00003,
		'4': 0x00004,
		'5': 0x00005,
		'6': 0x00006,
		'7': 0x00007,
		'8': 0x00008,
		'9': 0x00009,
		't': 0x0000A,
		'j': 0x0000B,
		'q': 0x0000D,
		'k': 0x0000E,
	}
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
	for i, c := range d {
		if i == d.Len()-1 {
			s += c.String()
		} else {
			s += c.String() + ", "
		}
	}
	return s
}

func (d Deck) UnicodeString() string {
	var b strings.Builder
	for _, c := range d {
		b.WriteString(c.UnicodeString())
	}
	return b.String()
}

func (d Deck) Len() int {
	return len(d)
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
	hand := (*d)[:n]
	*d = (*d)[n:]
	return hand
}

func isStraight(hand Deck) bool {
	for i := 0; i < hand.Len()-1; i++ {
		a := hand[i]
		b := hand[i+1]
		if b.Score()-a.Score() != 1 {
			return b.Score() == 13
		}
	}
	return true
}

func (d Deck) IsStraight(hand Deck) bool {
	var sorted Deck
	sorted = append(sorted, d...)
	sorted = append(sorted, hand...)

	slices.SortStableFunc(sorted, func(a, b *Card) int {
		return cmp.Compare(a.Score(), b.Score())
	})

	insideIsStraight := isStraight(sorted[1:6])
	highIsStraight := isStraight(sorted[2:])
	lowIsStraight := isStraight(append(sorted[:4], sorted[sorted.Len()-1]))

	return lowIsStraight || insideIsStraight || highIsStraight
}

func (d Deck) IsFlush(hand Deck) bool {
	numSuitsMatched := 0
	for i := 0; i < d.Len(); i++ {
		if d[i].suit == hand[0].suit {
			numSuitsMatched++
		}
		if d[i].suit == hand[1].suit {
			numSuitsMatched++
		}
	}
	return numSuitsMatched > 5
}

func (d Deck) AnalyzeHand(hand Deck) int {
	var cards Deck
	cards = append(cards, d...)
	cards = append(cards, hand...)

	groups := make(map[int]Deck)
	for _, c := range cards {
		groups[c.Score()] = append(groups[c.Score()], c)
	}

	card1 := hand[0]
	card2 := hand[1]
	handScore := max(card1.Score(), card2.Score())
	flush := d.IsFlush(hand)
	straight := d.IsStraight(hand)

	switch len(groups) {
	case 4:
		numPairs := 0
		for _, group := range groups {
			if len(group) == 4 {
				// quads
				return 8 * handScore
			}

			if len(group) == 2 {
				numPairs++
			}
		}

		// If the board has paired itself and the player has two pair then it's not a full house
		if numPairs < 3 {
			// full house
			return 7 * handScore
		}
		fallthrough
	default:
		switch {
		case flush && straight:
			return 9 * handScore
		case flush:
			return 6 * handScore
		case straight:
			return 5 * handScore
		default:
			card1Group := groups[card1.Score()]
			g1len := card1Group.Len()
			card2Group := groups[card2.Score()]
			g2len := card2Group.Len()

			// Trips
			if g1len == 3 || g2len == 3 {
				return 4 * handScore
			}

			// Two pair
			if g1len == 2 && g2len == 2 {
				return 3 * (card1.Score() + card2.Score())
			} else if g1len == 2 {
				return 2 * card1.Score()
			} else if g2len == 2 {
				return 2 * card2.Score()
			}

			// Ace high card
			if handScore == 13 {
				return 0
			}
			// High card
			return -handScore
		}
	}
}

type Card struct {
	face     byte
	suit     byte
	priority int
}

func (c Card) UnicodeString() string {
	return string(suitToUnicode[c.suit] + faceToUnicode[c.face])
}

func (c Card) String() string {
	return string(c.face) + string(c.suit)
}

func (c Card) Color() string {
	switch c.suit {
	case 'c':
		return "green"
	case 'd':
		return "blue"
	case 'h':
		return "red"
	default:
		return "black"
	}
}

func (c Card) Score() int {
	switch c.face {
	case 'a':
		return 13
	case 'k':
		return 12
	case 'q':
		return 11
	case 'j':
		return 10
	case 't':
		return 9
	case '9':
		return 8
	case '8':
		return 7
	case '7':
		return 6
	case '6':
		return 5
	case '5':
		return 4
	case '4':
		return 3
	case '3':
		return 2
	case '2':
		return 1
	default:
		return 0
	}
}
