package rgapis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"riot-developer-proxy/internal/domain/entities"
)

const RGAPI_DEFAULT_REGION = "americas"

type RGAPIWrapper struct {
	rgapiClient *rgapiClient
}

type SummonerDTO struct {
	ID            string `json:"id"`
	PUUID         string `json:"puuid"`
	AccountID     string `json:"accountId"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	Level         int    `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

type LeagueEntryDTO struct {
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

type AccountDTO struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

func NewRGAPIWrapper(rgapiClient *rgapiClient) *RGAPIWrapper {
	return &RGAPIWrapper{
		rgapiClient: rgapiClient,
	}
}

func (rr *RGAPIWrapper) GetAccountBySummonerName(ctx context.Context, region string, summonerName string) (*SummonerDTO, error) {
	resp, err := rr.rgapiClient.DoReq(ctx, "GET", region, "/lol/summoner/v4/summoners/by-name/"+url.PathEscape(summonerName), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("err when getting account by summoner name Status: %d Response: %s", resp.StatusCode, string(body))
	}

	var summonerDTO SummonerDTO
	if err := json.Unmarshal(body, &summonerDTO); err != nil {
		return nil, err
	}

	return &summonerDTO, nil
}

func (rr *RGAPIWrapper) GetSummonerLeagueEntries(ctx context.Context, region string, summoner *entities.Summoner) ([]LeagueEntryDTO, error) {
	resp, err := rr.rgapiClient.DoReq(ctx, "GET", region, "/lol/league/v4/entries/by-summoner/"+url.PathEscape(summoner.ID), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("err when getting summoner league Status: %d Response: %s", resp.StatusCode, string(body))
	}

	var leagueEntryDTO []LeagueEntryDTO
	if err := json.Unmarshal(body, &leagueEntryDTO); err != nil {
		return nil, err
	}

	return leagueEntryDTO, nil
}

func (rr *RGAPIWrapper) GetAccountByAccessToken(ctx context.Context, accessToken string) (*AccountDTO, error) {
	resp, err := rr.rgapiClient.DoReq(ctx, "GET", RGAPI_DEFAULT_REGION, "/riot/account/v1/accounts/me", &accessToken)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("err when getting account Status: %d Response: %s", resp.StatusCode, string(body))
	}

	var accountDTO AccountDTO
	if err := json.Unmarshal(body, &accountDTO); err != nil {
		return nil, err
	}

	return &accountDTO, nil
}
