package app

import (
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	log *zap.SugaredLogger
)

func init() {
	log = logger.GetLogger().Named("application (application.go)").Sugar()
}

func StartApplication(srv *http.Server, router *gin.Engine) {
	mapUrls(router)
	log.Info("[MAIN-APP] Starting server now")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorf("[MAIN-APP] ERROR when trying to START server: %s", err)
		panic(err)
	}
}
