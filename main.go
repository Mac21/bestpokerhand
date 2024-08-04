package main

import (
	"errors"
	"net"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	// APIHost is the address the api listens on
	APIHost = "localhost"
	// APIPort is the port the api listens on
	APIPort      = "9090"
	CookieSecret = "123"
)

func init() {
	// Change to gin.ReleaseMode when running in production
	gin.SetMode(gin.DebugMode)

	// Gitlab-ci defines PORT in the env which is what the kube service exposed port will be
	APIPort = getEnvDefault("PORT", APIPort)
	APIHost = getEnvDefault("HOST", APIHost)
	CookieSecret = getEnvDefault("COOKIE_SECRET", CookieSecret)
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

func getSession(ctx *gin.Context) sessions.Session {
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		Path: "/",
	})

	return session
}

func main() {
	rtr := gin.Default()
	rtr.SetTrustedProxies([]string{
		"127.0.0.1",
		"localhost",
	})

	store := cookie.NewStore([]byte(CookieSecret))
	store.Options(sessions.Options{
		Secure: true,
	})

	rtr.Use(
		sessions.Sessions("session", store),
	)

	rtr.HTMLRender = loadTemplates("./templates")
	rtr.Static("/static", "./static")
	initRoutes(rtr)
	rtr.Run(net.JoinHostPort(APIHost, APIPort))
}
