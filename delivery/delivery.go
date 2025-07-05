package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"TargetingEngineGG/app"
	"TargetingEngineGG/cache"
	"TargetingEngineGG/campaign"
	"TargetingEngineGG/database"
	"TargetingEngineGG/targeting"

	"github.com/gin-gonic/gin"
)

func DeliverCampaigns(c *gin.Context) {
	start := time.Now()
	ctx := make(map[string]string)
	missing := []string{}

	for _, dim := range RequiredDimensions {
		val := c.Query(dim)
		if val == "" {
			missing = append(missing, dim)
		}
		ctx[dim] = val
	}
	for _, dim := range OptionalDimensions {
		val := c.Query(dim)
		if val != "" {
			ctx[dim] = val
		}
	}

	if len(missing) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Missing query parameters: %v", missing)})
		return
	}

	var matched []campaign.MatchedCampaigns
	var campaigns []campaign.Campaign

	keys, err := cache.RDB.Keys(cache.Ctx, "campaign:*").Result()
	if err == nil && len(keys) > 0 {
		for _, key := range keys {
			data, err := cache.RDB.Get(cache.Ctx, key).Result()
			if err != nil {
				continue
			}
			var camp campaign.Campaign
			if err := json.Unmarshal([]byte(data), &camp); err != nil {
				continue
			}
			campaigns = append(campaigns, camp)
		}
	} else {
		rows, err := database.DB.Query("SELECT id,image_url,cta FROM campaigns WHERE state='ACTIVE'")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var camp campaign.Campaign
			if err := rows.Scan(&camp.ID, &camp.ImageURL, &camp.CTA); err != nil {
				continue
			}

			ruleRows, err := database.DB.Query("SELECT dimension,type,value FROM targeting_rules WHERE campaign_id = ?", camp.ID)
			if err != nil {
				continue
			}
			for ruleRows.Next() {
				var r targeting.Rule
				if err := ruleRows.Scan(&r.Dimension, &r.Type, &r.Value); err != nil {
					continue
				}
				camp.Rules = append(camp.Rules, r)
			}
			ruleRows.Close()

			campaigns = append(campaigns, camp)

			data, _ := json.Marshal(camp)
			cache.RDB.Set(cache.Ctx, "campaign:"+camp.ID, data, time.Hour)
		}
	}

	for _, camp := range campaigns {
		if targeting.MatchCampaigns(camp.Rules, ctx) {
			matched = append(matched, campaign.MatchedCampaigns{
				CampaignID: camp.ID,
				ImageURL:   camp.ImageURL,
				CTA:        camp.CTA,
			})
		}
	}
	access := app.AccessLog{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		StatusCode: http.StatusOK,
		Error:      "",
		Request:    c.Request.URL.RawQuery,
		Response:   matched,
		DurationMS: time.Since(start).Milliseconds(),
	}
	accessJSON, _ := json.Marshal(access)
	app.Info.Println(string(accessJSON))

	app.RequestsTotal.WithLabelValues("/v1/delivery", "GET", "200").Inc()

	if len(matched) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, matched)
}
