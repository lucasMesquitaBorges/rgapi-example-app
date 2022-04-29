package http

import (
	"log"
	"net/http"
	"riot-developer-proxy/internal/services"

	"github.com/labstack/echo/v4"
)

type HTTPController struct {
	RiotClient *services.RiotClient
}

const INTERNAL_SERVER_ERROR_MESSAGE = "Internal Server Error"

func (hr *HTTPController) ProxyToRIOTApi(c echo.Context) error {
	proxyResponse, err := hr.RiotClient.DoReq(
		c.Request().Context(),
		c.Request().Method,
		c.Request().URL.Path,
	)

	if err != nil {
		log.Println("err when retrieving data from api", err)
		return c.JSON(http.StatusInternalServerError, &Message{
			Message: INTERNAL_SERVER_ERROR_MESSAGE,
		})
	}

	return c.JSONBlob(proxyResponse.StatusCode, proxyResponse.Body)
}
