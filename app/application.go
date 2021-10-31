package app

import (
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.Info("*** Begin start appication ***")
	router.Run(":8081")
}
