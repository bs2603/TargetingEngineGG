package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"TargetingEngineGG/campaign"
	"TargetingEngineGG/targeting"
)

func RefreshCampaigns(db *sql.DB) {
	rows, err := db.Query("SELECT id, image_url, cta FROM campaigns WHERE state='ACTIVE'")
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var camp campaign.Campaign
		if err := rows.Scan(&camp.ID, &camp.ImageURL, &camp.CTA); err != nil {
			fmt.Println("Scan error:", err)
			continue
		}

		ruleRows, err := db.Query("SELECT dimension, type, value FROM targeting_rules WHERE campaign_id = ?", camp.ID)
		if err != nil {
			fmt.Println("Rule query error:", err)
			continue
		}

		for ruleRows.Next() {
			var r targeting.Rule
			if err := ruleRows.Scan(&r.Dimension, &r.Type, &r.Value); err != nil {
				fmt.Println("Rule scan error:", err)
				continue
			}
			camp.Rules = append(camp.Rules, r)
		}
		ruleRows.Close()

		data, _ := json.Marshal(camp)

		key := "campaign:" + camp.ID
		RDB.Set(Ctx, key, data, time.Hour)
	}
}
