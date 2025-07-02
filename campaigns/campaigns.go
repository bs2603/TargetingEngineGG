package campaigns

import (
	"TargetingEngineGG/database"
	"TargetingEngineGG/targeting"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Campaign struct {
	ID       string
	ImageURL string
	CTA      string
	State    string
	Rules    []targeting.Rule
}

func GetCampaigns(c *gin.Context) {

	//mysql

	app := c.Query("app")
	country := c.Query("country")
	os := c.Query("os")

	if app == "" || country == "" || os == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		return
	}

	rows, err := database.DB.Query("SELECT id,image_url,cta FROM campaigns WHERE state='ACTIVE'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer rows.Close()

	var campaigns []Campaign

	for rows.Next() {
		var camp Campaign
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

	var matched []gin.H

	for _, camp := range campaigns {
		if targeting.MatchCampaigns(camp.Rules, app, country, os) {
			matched = append(matched, gin.H{
				"cid": camp.ID,
				"img": camp.ImageURL,
				"cta": camp.CTA,
			})
		}
	}

	if len(matched) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// response := map[string]interface{}{
	// 	"app":     app,
	// 	"country": country,
	// 	"os":      os,
	// }

	c.JSON(http.StatusOK, matched)

}
