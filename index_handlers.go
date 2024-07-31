package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var overallScore = 0

func indexTemplateHandler(c *gin.Context) {
	overallScore = 0
	deck := NewDeck()
	deck.Shuffle()
	board := deck.DealCards(5)
	hand1 := deck.DealCards(2)
	hand2 := deck.DealCards(2)
	renderHTML(c, http.StatusOK, "index.tmpl", gin.H{
		"overallScore": overallScore,
		"board":        board,
		"hand1": gin.H{
			"cards": hand1,
			"score": board.AnalyzeHand(hand1),
		},
		"hand2": gin.H{
			"cards": hand2,
			"score": board.AnalyzeHand(hand2),
		},
	})
}

func runningGameHandler(c *gin.Context) {
	pickedScore, _ := c.GetPostForm("chosenscore")
	hand1score, _ := c.GetPostForm("hand1score")
	hand2score, _ := c.GetPostForm("hand2score")
	boardstr, _ := c.GetPostForm("board")
	hand1str, _ := c.GetPostForm("hand1")
	hand2str, _ := c.GetPostForm("hand2")

    fmt.Printf("Board: %s\n", boardstr)
    fmt.Printf("\tHand1: %s\n", hand1str)
    fmt.Printf("\tHand2: %s\n", hand2str)
    fmt.Printf("\tpicked score: %s, hand1score: %s, hand2score: %s, overallscore: %s\n", pickedScore, hand1score, hand2score, overallScore)

	switch pickedScore {
	case "hand1":
		if hand1score > hand2score {
			overallScore++
		}
	case "hand2":
		if hand2score > hand1score {
			overallScore++
		}
	}

	deck := NewDeck()
	deck.Shuffle()
	board := deck.DealCards(5)
	hand1 := deck.DealCards(2)
	hand2 := deck.DealCards(2)
	renderHTML(c, http.StatusOK, "index.tmpl", gin.H{
		"overallScore": overallScore,
		"board":        board,
		"hand1": gin.H{
			"cards": hand1,
			"score": board.AnalyzeHand(hand1),
		},
		"hand2": gin.H{
			"cards": hand2,
			"score": board.AnalyzeHand(hand2),
		},
	})
}
