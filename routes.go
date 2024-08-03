package main

import (
	"github.com/gin-gonic/gin"
)

func initRoutes(rtr *gin.Engine) {
	rtr.GET("/", indexTemplateHandler)
	rtr.POST("/current", runningGameHandler)
}
