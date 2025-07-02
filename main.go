package main

import (
	"TargetingEngineGG/campaigns"

	"github.com/gin-gonic/gin"
)

// step 1 register v1/delivery
func main() {

	router := gin.Default()

	router.GET("/health", getHealth)
	router.GET("/v1/delivery", campaigns.GetCampaigns)
	router.Run("localhost:8888")
}

func getHealth(c *gin.Context) {
	c.String(200, "ok")
}
