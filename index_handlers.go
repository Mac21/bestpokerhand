package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	boardstr, _ := c.GetPostForm("board")
	pickedScore, _ := c.GetPostForm("chosenscore")
	hand1scorestr, _ := c.GetPostForm("hand1score")
	hand2scorestr, _ := c.GetPostForm("hand2score")
	hand1str, _ := c.GetPostForm("hand1")
	hand2str, _ := c.GetPostForm("hand2")

	hand1score, _ := strconv.Atoi(hand1scorestr)
	hand2score, _ := strconv.Atoi(hand2scorestr)

	fmt.Printf("Board: %s\n", boardstr)
	fmt.Printf("\tHand1: %s\n", hand1str)
	fmt.Printf("\tHand2: %s\n", hand2str)
	fmt.Printf("\tpicked score: %s, hand1score: %d, hand2score: %d, overallscore: %d\n", pickedScore, hand1score, hand2score, overallScore)

	switch pickedScore {
	case "hand1":
		if hand1score >= hand2score {
			overallScore++
		}
	case "hand2":
		if hand2score >= hand1score {
			overallScore++
		}
	}
	numHands++
	session.Set("numHands", numHands)
	session.Set("overallScore", overallScore)
	session.Save()
    c.Redirect(http.StatusFound, "/")
}
