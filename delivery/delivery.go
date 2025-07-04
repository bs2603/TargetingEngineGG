package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"TargetingEngineGG/app"
	"TargetingEngineGG/campaign"
	"TargetingEngineGG/database"
	"TargetingEngineGG/targeting"

	"github.com/gin-gonic/gin"
)

func DeliverCampaigns(c *gin.Context) {
	start := time.Now()

	appName := c.Query("app")
	country := c.Query("country")
	os := c.Query("os")

	if appName == "" || country == "" || os == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		return
	}

	rows, err := database.DB.Query("SELECT id,image_url,cta FROM campaigns WHERE state='ACTIVE'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var campaigns []campaign.Campaign

	for rows.Next() {
		var camp campaign.Campaign
		if err := rows.Scan(&camp.ID, &camp.ImageURL, &camp.CTA); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error scanning DB": err})
		}

		ruleRows, err := database.DB.Query("SELECT dimension,type,value FROM targeting_rules WHERE campaign_id = ?", camp.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error querying DB": err})
			return
		}
		for ruleRows.Next() {
			var r targeting.Rule
			if err := ruleRows.Scan(&r.Dimension, &r.Type, &r.Value); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"DB scan error": err})
			}
			camp.Rules = append(camp.Rules, r)
		}
		ruleRows.Close()
		campaigns = append(campaigns, camp)

	}

	var matched []campaign.MatchedCampaigns

	for _, camp := range campaigns {
		if targeting.MatchCampaigns(camp.Rules, appName, country, os) {
			matched = append(matched, campaign.MatchedCampaigns{
				CampaignID: camp.ID,
				ImageURL:   camp.ImageURL,
				CTA:        camp.CTA,
			})
		}
	}

	// response := map[string]interface{}{
	// 	"app":     app,
	// 	"country": country,
	// 	"os":      os,
	// }

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
