package targeting

import "strings"

type Rule struct {
	Dimension string
	Type      string
	Value     string
}

func MatchCampaigns(rules []Rule, ctx map[string]string) bool {
	includeCheck := make(map[string]bool)
	hasInclude := make(map[string]bool)

	for _, r := range rules {
		value, ok := ctx[strings.ToLower(r.Dimension)]
		if !ok {
			value = ""
		}

		switch r.Type {
		case "INCLUDE":
			hasInclude[r.Dimension] = true
			if strings.EqualFold(r.Value, value) {
				includeCheck[r.Dimension] = true
			}
		case "EXCLUDE":
			if strings.EqualFold(r.Value, value) {
				return false
			}
		}
	}

	for dimension := range hasInclude {
		if !includeCheck[dimension] {
			return false
		}
	}

	return true
}
