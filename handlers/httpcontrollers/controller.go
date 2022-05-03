package httpcontrollers

import (
	"log"
	"net/http"
	"os"
	"riot-developer-proxy/internal/domain/entities"
	"riot-developer-proxy/internal/domain/services"
	"riot-developer-proxy/repositories"
	"time"

	"github.com/labstack/echo/v4"
)

const INTERNAL_SERVER_ERROR_MESSAGE = "Internal Server Error"

type HTTPController struct {
	SummonerService *services.SummonerService
}

type SummonerResponse struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	ProfileIconID int                   `json:"profileIconId"`
	Level         int                   `json:"summonerLevel"`
	RevisionDate  string                `json:"revisionDate"`
	LeagueEntries []LeagueEntryResponse `json:"leagueEntries"`
}

type LeagueEntryResponse struct {
	LeagueID     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
}

func (hr *HTTPController) SummonerProfileByName(c echo.Context) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	riotClient := repositories.NewRiotClient(client, os.Getenv("RIOT_API_BASE_URI"))
	riotClient.WithToken(os.Getenv("RIOT_API_TOKEN"))

	repo := repositories.NewRiotRepository(riotClient)

	svc := services.NewSummonerService(repo)
	summoner, err := svc.GetSummonerProfile(
		c.Request().Context(),
		c.QueryParam("region"),
		c.QueryParam("name"),
	)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, &Message{
			Message: INTERNAL_SERVER_ERROR_MESSAGE,
		})
	}

	summonerResponse := &SummonerResponse{
		ID:            summoner.ID,
		Name:          summoner.Name,
		ProfileIconID: summoner.ProfileIconID,
		Level:         summoner.Level,
		LeagueEntries: buildLeagueEntriesFromSummoner(summoner),
		RevisionDate:  summoner.GetRevisionDateAsISOFormat(),
	}

	return c.JSON(http.StatusOK, summonerResponse)
}

func buildLeagueEntriesFromSummoner(summoner *entities.Summoner) []LeagueEntryResponse {
	leagueEntriesResponse := make([]LeagueEntryResponse, len(summoner.LeagueEntries))
	for i, leagueEntry := range summoner.LeagueEntries {
		leagueEntriesResponse[i] = LeagueEntryResponse(leagueEntry)
	}

	return leagueEntriesResponse
}
