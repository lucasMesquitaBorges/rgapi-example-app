package main

import (
	net_http "net/http"
	"os"
	"riot-developer-proxy/handlers/http"
	"riot-developer-proxy/internal/services"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	hc := &http.HTTPController{
		RiotClient: initRiotClient(),
	}

	e := echo.New()
	e.Use(middleware.Recover())

	e.Any("/*", hc.ProxyToRIOTApi)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}

func initRiotClient() *services.RiotClient {
	client := &net_http.Client{
		Timeout: time.Second * 10,
	}
	riotClient := services.NewRiotClient(*client, os.Getenv("RIOT_API_BASE_URI"))
	riotClient.WithToken(os.Getenv("RIOT_API_TOKEN"))

	return riotClient
}
