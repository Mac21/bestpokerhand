package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexTemplateHandler(c *gin.Context) {
	session := getSession(c)
	overallScore, hasOverallScore := session.Get("overallScore").(int)
	if !hasOverallScore {
		session.Set("overallScore", 0)
		overallScore = 0
	}

	numHands, hasNumHands := session.Get("numHands").(int)
	if !hasNumHands {
		session.Set("numHands", 0)
		numHands = 0
	}
	session.Save()

	deck := NewDeck()
	deck.Shuffle()
	board := deck.DealCards(5)
	hand1 := deck.DealCards(2)
	hand2 := deck.DealCards(2)
	renderHTML(c, http.StatusOK, "index.tmpl", gin.H{
		"overallScore": overallScore,
		"numHands":     numHands,
		"board":        board,
		"hand1":        hand1,
		"hand2":        hand2,
	})
}

func newGameHandler(c *gin.Context) {
	session := getSession(c)
	session.Set("overallScore", 0)
	session.Set("numHands", 0)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func runningGameHandler(c *gin.Context) {
	session := getSession(c)
	overallScore := session.Get("overallScore").(int)
	numHands := session.Get("numHands").(int)
	bs, _ := c.GetPostForm("board")
	pickedScore, _ := c.GetPostForm("chosenscore")
	h1s, _ := c.GetPostForm("hand1")
	h2s, _ := c.GetPostForm("hand2")

	board := NewDeckFromString(bs)
	h1 := NewDeckFromString(h1s)
	h2 := NewDeckFromString(h2s)
	hand1Score := board.AnalyzeHand(*h1)
	hand2Score := board.AnalyzeHand(*h2)

	fmt.Printf("Board: %s\n", board)
	fmt.Printf("\tHand1: %s\n", h1s)
	fmt.Printf("\tHand2: %s\n", h2s)
	fmt.Printf("\tpicked score: %s, hand1Score: %d, hand2Score: %d, overallscore: %d\n", pickedScore, hand1Score, hand2Score, overallScore)

	switch pickedScore {
	case "hand1":
		if hand1Score >= hand2Score {
			overallScore++
		}
	case "hand2":
		if hand2Score >= hand1Score {
			overallScore++
		}
	}
	numHands++
	session.Set("numHands", numHands)
	session.Set("overallScore", overallScore)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
