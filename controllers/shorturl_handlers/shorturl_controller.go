package shorturl_handlers

import (
	"github.com/Moriartii/url-shortner-api/domain/shorturl"
	"github.com/Moriartii/url-shortner-api/services"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(c *gin.Context) {
	var shortUrlRequest shorturl.ShortUrlRequest
	if err := c.ShouldBindJSON(&shortUrlRequest); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	shortedUrl, err := services.ShortUrlService.CreateShortUrl(shortUrlRequest.Url)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, shortedUrl)
}

func Information(c *gin.Context) {
	shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	redirectUrl, err := services.ShortUrlService.GetShortUrlByShortPath(shortUrlRequest.ShortBase32)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, redirectUrl)
}

func Redirect(c *gin.Context) {
	shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	redirectUrl, err := services.ShortUrlService.GetShortUrlByShortPath(shortUrlRequest.ShortBase32)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	errIncremCount := services.ShortUrlService.IncrementShortUrlCount(shortUrlRequest.ShortBase32)
	if errIncremCount != nil {
		c.JSON(err.Status, err)
		return
	}
	c.Redirect(http.StatusPermanentRedirect, redirectUrl.Url)
}

func GetAll(c *gin.Context) {
	//shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	allUrls, err := services.ShortUrlService.GetAllUrls()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, allUrls)
}
