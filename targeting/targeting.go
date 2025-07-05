package targeting

import "strings"

type Rule struct {
	Dimension string
	Type      string
	Value     string
}

func MatchCampaigns(rules []Rule, ctx map[string]string) bool {
	// Track inclusion logic for each dimension
	includeCheck := make(map[string]bool)
	hasInclude := make(map[string]bool)

	for _, r := range rules {
		value, ok := ctx[strings.ToLower(r.Dimension)]
		if !ok {
			value = "" // in case context is missing, treat as empty
		}

		switch r.Type {
		case "INCLUDE":
			hasInclude[r.Dimension] = true
			if strings.EqualFold(r.Value, value) {
				includeCheck[r.Dimension] = true
			}
		case "EXCLUDE":
			if strings.EqualFold(r.Value, value) {
				return false // excluded, reject immediately
			}
		}
	}

	for dimension := range hasInclude {
		if !includeCheck[dimension] {
			return false // dimension had INCLUDE but no match
		}
	}

	return true // passed all filters
}
