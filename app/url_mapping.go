package app

import (
	"github.com/Moriartii/url-shortner-api/controllers/ping"
	"github.com/Moriartii/url-shortner-api/controllers/shorturl_handlers"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	router.GET("/all_short_urls", shorturl_handlers.GetAll)
	router.GET("/:short_path", shorturl_handlers.Redirect)
	router.GET("/:short_path/info", shorturl_handlers.Information)
	router.POST("/", shorturl_handlers.Create)

}
