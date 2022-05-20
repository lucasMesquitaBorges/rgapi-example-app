package httpcontrollers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"riot-developer-proxy/internal/domain/entities"
	"riot-developer-proxy/internal/domain/services"
	"strings"

	"github.com/labstack/echo/v4"
)

const INTERNAL_SERVER_ERROR_MESSAGE = "Internal Server Error"
const RSO_BASE_URI = "https://auth.riotgames.com"
const APP_BASE_URL = "http://local.exampleapp.com:3000"
const APP_CALLBACK_PATH = "/oauth-callback"
const AUTHORIZE_URL = RSO_BASE_URI + "/authorize"
const TOKEN_URL = RSO_BASE_URI + "/token"
const LOGOUT_URL = RSO_BASE_URI + "/logout"
const APP_CALLBACK_URL = APP_BASE_URL + APP_CALLBACK_PATH

var tokenResponse *TokenResponse

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AccountResponse struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
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

type HTTPController struct {
	SummonerService *services.SummonerService
}

func (hr *HTTPController) Login(c echo.Context) error {
	if tokenResponse == nil {
		return hr.loginPage(c)
	}

	return hr.loggedPage(c)
}

func (hr *HTTPController) OAUTHCallback(c echo.Context) error {
	code := c.QueryParam("code")
	data := make(url.Values)
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", APP_CALLBACK_URL)

	req, err := http.NewRequestWithContext(
		c.Request().Context(),
		"POST",
		TOKEN_URL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return internalServerError(c, err)
	}
	req.SetBasicAuth(os.Getenv("RSO_CLIENT_ID"), os.Getenv("RSO_CLIENT_SECRET"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return internalServerError(c, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return internalServerError(c, err)
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return internalServerError(c, err)
	}

	return c.HTML(http.StatusOK, `<script>window.location.href = "/";</script>`)
}

func (hr *HTTPController) SummonerProfileByName(c echo.Context) error {
	summoner, err := hr.SummonerService.GetSummonerProfile(
		c.Request().Context(),
		c.QueryParam("region"),
		c.QueryParam("name"),
	)
	if err != nil {
		return internalServerError(c, err)
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

func (hr *HTTPController) MyAccount(c echo.Context) error {
	if tokenResponse == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	account, err := hr.SummonerService.GetAccountByAccessToken(c.Request().Context(), tokenResponse.AccessToken)
	if err != nil {
		return internalServerError(c, err)
	}

	return c.JSON(http.StatusOK, AccountResponse(*account))
}

func (hr *HTTPController) Logout(c echo.Context) error {
	tokenResponse = nil
	return c.HTML(http.StatusOK, `<script>window.location.href = "`+LOGOUT_URL+`";</script>`)
}

func (hr *HTTPController) loginPage(c echo.Context) error {
	signInLink := AUTHORIZE_URL +
		"?redirect_uri=" + APP_CALLBACK_URL +
		"&client_id=" + "riot-example-app" +
		"&response_type=code" +
		"&scope=openid"

	return c.HTML(http.StatusOK, `<a href="`+signInLink+`">Sign In</a>`)
}

func (hr *HTTPController) loggedPage(c echo.Context) error {
	myAccountLink := APP_BASE_URL + "/accounts/me"
	logoutLink := APP_BASE_URL + "/logout"

	return c.HTML(http.StatusOK, `
		<a href="`+myAccountLink+`">My Account</a>
		&nbsp;
		<a href="`+logoutLink+`">Logout</a>
	`)
}

func buildLeagueEntriesFromSummoner(summoner *entities.Summoner) []LeagueEntryResponse {
	leagueEntriesResponse := make([]LeagueEntryResponse, len(summoner.LeagueEntries))
	for i, leagueEntry := range summoner.LeagueEntries {
		leagueEntriesResponse[i] = LeagueEntryResponse(leagueEntry)
	}

	return leagueEntriesResponse
}

func internalServerError(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(http.StatusInternalServerError, &Message{
		Message: INTERNAL_SERVER_ERROR_MESSAGE,
	})
}
