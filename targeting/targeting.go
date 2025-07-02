package targeting

import "strings"

type Rule struct {
	Dimension string
	Type      string
	Value     string
}

func MatchCampaigns(rules []Rule, app, country, os string) bool {
	includeApp := true
	includeCountry := true
	includeOS := true

	hasIncludeApp := false
	hasIncludeCountry := false
	hasIncludeOS := false

	for _, r := range rules {
		switch r.Dimension {
		case "APP":
			if r.Type == "INCLUDE" {
				hasIncludeApp = true
				if r.Value == app {
					includeApp = true
				}
			}
			if r.Type == "EXCLUDE" && strings.EqualFold(r.Value, app) {
				return false
			}
		case "COUNTRY":
			if r.Type == "INCLUDE" {
				hasIncludeCountry = true
				if r.Value == country {
					includeCountry = true
				}
			}
			if r.Type == "EXCLUDE" && strings.EqualFold(r.Value, country) {
				return false
			}

		case "OS":
			if r.Type == "INCLUDE" {
				hasIncludeOS = true
				if r.Value == os {
					includeOS = true
				}
			}
			if r.Type == "EXCLUDE" && strings.EqualFold(r.Value, os) {
				return false
			}
		}
	}
	if hasIncludeApp && !includeApp {
		return false
	}
	if hasIncludeCountry && !includeCountry {
		return false
	}
	if hasIncludeOS && !includeOS {
		return false
	}
	return true
}
