package app

import (
	"github.com/Moriartii/url-shortner-api/controllers/ping"
	"github.com/Moriartii/url-shortner-api/controllers/shorturl_handlers"
	"github.com/gin-gonic/gin"
)

func mapUrls(router *gin.Engine) {
	router.GET("/ping", ping.Ping)
	router.GET("/all_short_urls", handlers.ShortUrlsController.GetAll)
	router.GET("/:short_path", handlers.ShortUrlsController.Redirect)
	router.GET("/:short_path/info", handlers.ShortUrlsController.Information)
	router.POST("/", handlers.ShortUrlsController.Create)

}
