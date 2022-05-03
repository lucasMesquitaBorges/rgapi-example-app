package entities

import "time"

type Summoner struct {
	ID            string
	PUUID         string
	AccountID     string
	Name          string
	ProfileIconID int
	Level         int
	RevisionDate  time.Time
	LeagueEntries []LeagueEntry
}

func (s *Summoner) AddLeagueEntry(leagueEntry *LeagueEntry) {
	s.LeagueEntries = append(s.LeagueEntries, *leagueEntry)
}

func (s *Summoner) GetRevisionDateAsISOFormat() string {
	return s.RevisionDate.Format(time.RFC3339)
}
