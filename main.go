package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Campaign struct {
	ID       string
	ImageURL string
	CTA      string
	State    string
	Rules    []Rule
}

type Rule struct {
	Dimension string
	Type      string
	Value     string
}

// step 1 register v1/delivery
func main() {

	router := gin.Default()

	router.GET("/health", getHealth)
	router.GET("/v1/delivery", getCampaigns)
	router.Run("localhost:8888")

	// //mysql
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/targeting")
	// if err != nil {
	// 	log.Fatal("DB connection error: ", err)
	// }
	// defer db.Close()
}

func getHealth(c *gin.Context) {
	c.String(200, "ok")
}

func getCampaigns(c *gin.Context) {

	//mysql
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/GGtargeting")
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}
	defer db.Close()

	app := c.Query("app")
	country := c.Query("country")
	os := c.Query("os")

	if app == "" || country == "" || os == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameters"})
		return
	}

	rows, err := db.Query("SELECT id,image_url,cta FROM campaigns WHERE state='ACTIVE'")
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

		ruleRows, err := db.Query("SELECT dimension,type,value FROM targeting_rules WHERE campaign_id = ?", camp.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error querying DB": err})
			return
		}
		for ruleRows.Next() {
			var r Rule
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
		if matches(camp.Rules, app, country, os) {
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

func matches(rules []Rule, app, country, os string) bool {
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
