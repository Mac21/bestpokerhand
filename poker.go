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
	for _, c := range d {
		s += c.String()
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

	runningCount := hand[0].Score() % 14
	for i := 0; i < n; i++ {
		c := hand[i]
		if (runningCount % 14) != (c.Score() % 14) {
			return false, runningCount
		}

		if c.Score() == 14 {
			runningCount += 2
		} else {
			runningCount++
		}
	}

	return true, runningCount
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

	if isOutsideStraight, score := isStraight(aceHigh[1:]); isOutsideStraight {
		return isOutsideStraight, score
	}

	if isHighStraight, score := isStraight(aceHigh[2:]); isHighStraight {
		return isHighStraight, score
	}

	if isLowStraight, score := isStraight(aceLow[:aceLow.Len()-2]); isLowStraight {
		return isLowStraight, score
	}

	return false, -1
}

func (d Deck) IsFlush(hand Deck) (bool, bool) {
	var cards Deck
	cards = append(cards, d...)
	cards = append(cards, hand...)

	flushgroups := make(map[byte]Deck)
	for _, c := range cards {
		flushgroups[c.suit] = append(flushgroups[c.suit], c)
	}

	for _, groups := range flushgroups {
		if groups.Len() > 4 {
			isStraight, _ := isStraight(groups)
			return true, isStraight
		}
	}
	return false, false
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
	flush, isStraightFlush := d.IsFlush(hand)
	straight, straightScore := d.IsStraight(hand)

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
		case isStraightFlush:
			return 9 * (handScore * handScore)
		case flush:
			return 6 * handScore
		case straight:
			return 5 * straightScore
		default:
			card1Group := groups[card1.Score()]
			g1len := card1Group.Len()
			card2Group := groups[card2.Score()]
			g2len := card2Group.Len()

			// Trips
			if g1len == 3 {
				return 4 * card1.Score()
			} else if g2len == 3 {
				return 4 * card2.Score()
			}

			// Two pair
			if g1len == 2 && g2len == 2 && card1.Score() != card2.Score() {
				return 3 * handScore
			} else if g1len == 2 {
				return 2 * card1.Score()
			} else if g2len == 2 {
				return 2 * card2.Score()
			}

			return handScore - 13
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
