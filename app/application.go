package app

import (
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func init() {
	log = logger.GetLogger().Named("application (application.go)").Sugar()
}

func StartApplication() {
	defer log.Sync()

	router := gin.Default()

	mapUrls(router)

	log.Infof("[MAIN-APP] starting server")

	router.Run("192.168.0.4:8081")
}
