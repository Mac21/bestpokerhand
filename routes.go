package main

import (
	"github.com/gin-gonic/gin"
)

func initRoutes(rtr *gin.Engine) {
	rtr.GET("/", indexTemplateHandler)
	rtr.GET("/newgame", newGameHandler)
	rtr.POST("/score", runningGameHandler)
}
