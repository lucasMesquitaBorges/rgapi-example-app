package main

import (
	"os"
	"riot-developer-proxy/handlers/httpcontrollers"
	"riot-developer-proxy/internal/domain/services"
	"riot-developer-proxy/rgapis"

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
	hc := initHTTPController()

	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/summoners/overview", hc.SummonerProfileByName)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}

func initHTTPController() *httpcontrollers.HTTPController {
	rgapiClient := rgapis.NewRGAPIClient(os.Getenv("RIOT_API_TOKEN"))
	rgapi := rgapis.NewRGAPIWrapper(rgapiClient)
	svc := services.NewSummonerService(rgapi)

	return &httpcontrollers.HTTPController{
		SummonerService: svc,
	}
}
