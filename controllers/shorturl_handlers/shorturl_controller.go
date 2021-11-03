package handlers

import (
	"github.com/Moriartii/url-shortner-api/domain/shorturl"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/Moriartii/url-shortner-api/services"
	"github.com/Moriartii/url-shortner-api/utils/encode"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var (
	ShortUrlsController shortUrlsControllerInterface = &shortUrlsController{}
	log                 *zap.SugaredLogger
)

func init() {
	log = logger.GetLogger().Named("handlers (shorturl_controller.go)").Sugar()
}

type shortUrlsControllerInterface interface {
	Create(*gin.Context)
	Information(*gin.Context)
	Redirect(*gin.Context)
	GetAll(*gin.Context)
}

type shortUrlsController struct{}

func (i *shortUrlsController) Create(c *gin.Context) {
	var body []byte
	var tempErr error
	body, c.Request.Body, tempErr = encode.RequestBodyForLogger(c.Request.Body)
	if tempErr != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: encode.RequestBodyForLogger(): %+v", tempErr)
		c.JSON(http.StatusBadRequest, tempErr)
		return
	}
	var shortUrlRequest shorturl.ShortUrlRequest

	log.Infof("[HANDLER] %s request from %s to endpoint %s with body: %s", c.Request.Method, c.ClientIP(), c.Request.URL, body)

	if err := c.ShouldBindJSON(&shortUrlRequest); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		log.Errorf("[HANDLER] ERROR when trying to call: c.ShouldBindJSON(): %+v", err)
		c.JSON(restErr.Status, restErr)
		return
	}
	shortedUrl, err := services.ShortUrlService.CreateShortUrl(shortUrlRequest.Url)
	if err != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: services.ShortUrlService.CreateShortUrl(): %+v", err)
		c.JSON(err.Status, err)
		return
	}
	log.Infof("[HANDLER] SUCCESS response code: %d with body: %+v", http.StatusCreated, shortedUrl)
	c.JSON(http.StatusCreated, shortedUrl)
}

func (i *shortUrlsController) Information(c *gin.Context) {
	log.Infof("[HANDLER] %s request from: %s to endpoint: %s with params: %s", c.Request.Method, c.ClientIP(), c.Request.URL, c.Request.URL.Query())
	shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	redirectUrl, err := services.ShortUrlService.GetShortUrlByShortPath(shortUrlRequest.ShortBase32)
	if err != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: services.ShortUrlService.GetShortUrlByShortPath(): %+v", err)
		c.JSON(err.Status, err)
		return
	}
	log.Infof("[HANDLER] SUCCESS response code: %d with body: %+v", http.StatusCreated, redirectUrl)
	c.JSON(http.StatusOK, redirectUrl)
}

func (i *shortUrlsController) Redirect(c *gin.Context) {
	log.Infof("[HANDLER] %s request from: %s to endpoint: %s with params: %s", c.Request.Method, c.ClientIP(), c.Request.URL, c.Request.URL.Query())
	shortUrlRequest := shorturl.ShortUrlRequest{ShortBase32: c.Param("short_path")}
	redirectUrl, err := services.ShortUrlService.GetShortUrlByShortPath(shortUrlRequest.ShortBase32)
	if err != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: services.ShortUrlService.GetShortUrlByShortPath(): %+v", err)
		c.JSON(err.Status, err)
		return
	}
	errIncremCount := services.ShortUrlService.IncrementShortUrlCount(shortUrlRequest.ShortBase32)
	if errIncremCount != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: services.ShortUrlService.IncrementShortUrlCount(): %+v", err)
		c.JSON(err.Status, err)
		return
	}
	log.Infof("[HANDLER] SUCCESS response code: %d with body: %+v", http.StatusPermanentRedirect, redirectUrl)
	c.Redirect(http.StatusPermanentRedirect, redirectUrl.Url)
}

func (i *shortUrlsController) GetAll(c *gin.Context) {
	log.Infof("[HANDLER] %s request from: %s to endpoint: %s with params: %s", c.Request.Method, c.ClientIP(), c.Request.URL, c.Request.URL.Query())
	allUrls, err := services.ShortUrlService.GetAllUrls()
	if err != nil {
		log.Errorf("[HANDLER] ERROR when trying to call: services.ShortUrlService.GetAllUrls(): %+v", err)
		c.JSON(err.Status, err)
		return
	}
	log.Infof("[HANDLER] SUCCESS response code: %d with body: %+v", http.StatusOK, allUrls)
	c.JSON(http.StatusOK, allUrls)
}
