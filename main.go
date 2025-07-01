package main

import (
	"github.com/gin-gonic/gin"
)

// step 1 register v1/delivery
func main() {

	router := gin.Default()

	router.GET("/health", getHealth)
	router.GET("/v1/delivery", getCampaigns)
	router.Run("localhost:8888")
}

func getHealth(c *gin.Context) {
	c.String(200, "ok")
}

func getCampaigns(c *gin.Context) {

	App := c.Query("app")
	Country := c.Query("country")
	Os := c.Query("os")

	response := map[string]interface{}{
		"app":     App,
		"country": Country,
		"os":      Os,
	}

	c.JSON(200, response)

}
