package main

import (
	"github.com/bs2603/TargetingEngineGG/delivery"

	"github.com/bs2603/TargetingEngineGG/app"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	app.Init()
	router := gin.Default()

	router.GET("/health", getHealth)
	router.GET("/v1/delivery", delivery.DeliverCampaigns)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app.Info.Println("Server started on :8888")
	router.Run("localhost:8888")
}

func getHealth(c *gin.Context) {
	c.String(200, "ok")
}
