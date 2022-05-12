package services

import (
	"context"
	"riot-developer-proxy/internal/domain/entities"
	"riot-developer-proxy/rgapis"
	"time"
)

type SummonerService struct {
	rgapi *rgapis.RGAPIWrapper
}

func NewSummonerService(rgapi *rgapis.RGAPIWrapper) *SummonerService {
	return &SummonerService{
		rgapi: rgapi,
	}
}

func (rs *SummonerService) GetSummonerProfile(ctx context.Context, region string, summonerName string) (*entities.Summoner, error) {
	summoner, err := rs.getSummonerAccount(ctx, region, summonerName)
	if err != nil {
		return nil, err
	}

	leagueEntries, err := rs.getSummonerLeagueEntries(ctx, region, summoner)
	if err != nil {
		return nil, err
	}

	for _, leagueEntry := range leagueEntries {
		summoner.AddLeagueEntry(&leagueEntry)
	}

	return summoner, nil
}

func (rs *SummonerService) GetAccountByAccessToken(ctx context.Context, accessToken string) (*entities.Account, error) {
	accountDTO, err := rs.rgapi.GetAccountByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	account := entities.Account(*accountDTO)
	return &account, nil
}

func (rs *SummonerService) getSummonerLeagueEntries(ctx context.Context, region string, summoner *entities.Summoner) ([]entities.LeagueEntry, error) {
	leagueEntriesDTO, err := rs.rgapi.GetSummonerLeagueEntries(ctx, region, summoner)
	if err != nil {
		return nil, err
	}

	if len(leagueEntriesDTO) == 0 {
		return nil, nil
	}

	var leagueEntries []entities.LeagueEntry
	for _, leagueEntryDTO := range leagueEntriesDTO {
		leagueEntry := entities.LeagueEntry(leagueEntryDTO)
		leagueEntries = append(leagueEntries, leagueEntry)
	}

	return leagueEntries, nil
}

func (rs *SummonerService) getSummonerAccount(ctx context.Context, region string, summonerName string) (*entities.Summoner, error) {
	summonerDTO, err := rs.rgapi.GetAccountBySummonerName(ctx, region, summonerName)
	if err != nil {
		return nil, err
	}

	return &entities.Summoner{
		ID:            summonerDTO.ID,
		PUUID:         summonerDTO.PUUID,
		AccountID:     summonerDTO.AccountID,
		Name:          summonerDTO.Name,
		ProfileIconID: summonerDTO.ProfileIconID,
		Level:         summonerDTO.Level,
		RevisionDate:  time.Unix(0, summonerDTO.RevisionDate*int64(time.Millisecond)),
	}, nil
}
