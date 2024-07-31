package main

import (
	"errors"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	// APIHost is the address the api listens on
	APIHost = "localhost"
	// APIPort is the port the api listens on
	APIPort = "9090"
	// APIKey used for auth on rest end points
)

func init() {
	// Change to gin.ReleaseMode when running in production
	gin.SetMode(gin.DebugMode)

	// Gitlab-ci defines PORT in the env which is what the kube service exposed port will be
	if port := os.Getenv("PORT"); port != "" {
		APIPort = port
	}
}

func errorMissingEnv(env string) {
	panic(errors.New(env + " env var is missing"))
}

func getEnvDefault(key, def string) string {
	tempEnvVar := os.Getenv(key)
	if tempEnvVar == "" {
		return def
	}
	return tempEnvVar
}

func getEnv(key string) string {
	tempEnvVar := os.Getenv(key)
	if tempEnvVar == "" {
		errorMissingEnv(key)
	}
	return tempEnvVar
}

func main() {
	rtr := gin.Default()
	rtr.SetTrustedProxies([]string{
		"127.0.0.1",
		"localhost",
	})

	rtr.HTMLRender = loadTemplates("./templates")
	rtr.Static("/static", "./static")
	initRoutes(rtr)
	rtr.Run(net.JoinHostPort(APIHost, APIPort))
}
