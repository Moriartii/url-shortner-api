package shorturl_handlers

import (
	"github.com/Moriartii/url-shortner-api/domain/shorturl"
	"github.com/Moriartii/url-shortner-api/services"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ShortUrlsController shortUrlsControllerInterface = &shortUrlsController{}
)

type shortUrlsControllerInterface interface {
	Create(*gin.Context)
	Information(*gin.Context)
	Redirect(*gin.Context)
	GetAll(*gin.Context)
}

type shortUrlsController struct{}

func (i *shortUrlsController) Create(c *gin.Context) {
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

func (i *shortUrlsController) Information(c *gin.Context) {
	shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	redirectUrl, err := services.ShortUrlService.GetShortUrlByShortPath(shortUrlRequest.ShortBase32)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, redirectUrl)
}

func (i *shortUrlsController) Redirect(c *gin.Context) {
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

func (i *shortUrlsController) GetAll(c *gin.Context) {
	allUrls, err := services.ShortUrlService.GetAllUrls()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, allUrls)
}
