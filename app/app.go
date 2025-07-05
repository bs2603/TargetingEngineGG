package app

import (
	"TargetingEngineGG/cache"
	"TargetingEngineGG/database"
	"io"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	cache.InitRedis()
	logDir := "logs"
	os.MkdirAll(logDir, os.ModePerm)

	logFile := &lumberjack.Logger{
		Filename:   logDir + "/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // gzip old logs
	}

	multiOut := io.MultiWriter(os.Stdout, logFile)

	Info = log.New(multiOut, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiOut, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestLatency)

	go func() {
		for {
			cache.RefreshCampaigns(database.DB)
			time.Sleep(3600 * time.Second)
		}
	}()

}
