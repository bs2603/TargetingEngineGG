package campaign

import (
	"TargetingEngineGG/targeting"
)

type Campaign struct {
	ID       string
	ImageURL string
	CTA      string
	State    string
	Rules    []targeting.Rule
}

type MatchedCampaigns struct {
	CampaignID string `json:"cid"`
	ImageURL   string `json:"img"`
	CTA        string `json:"cta"`
}
