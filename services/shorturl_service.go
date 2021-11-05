package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Moriartii/url-shortner-api/domain/shorturl"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/Moriartii/url-shortner-api/utils/encode"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/Moriartii/url-shortner-api/utils/hash"
	"github.com/Moriartii/url-shortner-api/utils/shufflestring"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func init() {
	log = logger.GetLogger().Named("services (shorturl_service.go)").Sugar()
}

type shortUrlService struct {
}

type shortUrlServiceInterface interface {
	GetShortUrl(string) (*shorturl.ShortUrl, *errors.RestErr)
	CreateShortUrl(string) (*shorturl.ShortUrl, *errors.RestErr)
	GetShortUrlByShortPath(string) (*shorturl.ShortUrl, *errors.RestErr)
	IncrementShortUrlCount(string) *errors.RestErr
	GetAllUrls() ([]shorturl.ShortUrlRequestWithId, *errors.RestErr)
}

var (
	ShortUrlService shortUrlServiceInterface = &shortUrlService{}
)

func initUrlAndHashFromRequest(url string) (*shorturl.ShortUrl, *errors.RestErr) {
	log.Infof("[SERVICE] BEGIN CALL: initUrlAndHashFromRequest() with args: %s", url)
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	dao := &shorturl.ShortUrl{}
	dao.Url = url
	hashed, err := hash.UrlToHash(url)
	if err != nil {
		log.Errorf("[SERVICE] ERROR when trying to call: hash.UrlToHash(): %+v", err)
		return nil, err
	}
	dao.UrlHash = fmt.Sprintf("%x", string(hashed)) // Теперь в струтуре заполнены URL и HASH, есть префикс http
	log.Infof("[SERVICE] SUCCESS called initUrlAndHashFromRequest() and returned: %+v", dao)
	return dao, nil
}

func (s *shortUrlService) IncrementShortUrlCount(base32 string) *errors.RestErr {
	dao := &shorturl.ShortUrl{}
	dao.ShortBase32 = base32
	if err := dao.IncrementRedirectCount(); err != nil {
		return err
	}
	return nil
}

func (s *shortUrlService) CreateShortUrl(url string) (*shorturl.ShortUrl, *errors.RestErr) {
	log.Infof("[SERVICE] BEGIN CALL: CreateShortUrl() with args: %s", url)
	dao, err := initUrlAndHashFromRequest(url)
	if err != nil {
		return nil, err
	}

	// Проверяем по хэшу, есть ли уже в БД такая запись. Если есть, возвращаем её, если нет, то создаем новую.
	if err := dao.GetShortUrlByHash(); err != nil {
		log.Errorf("[SERVICE] ERROR when trying to call: dao.GetShortUrlByHash(): %+v", err)
		switch err.Message {
		case "no value":
			dao.ShortBase32 = encode.HashToBase32(dao.UrlHash)[0:6] // Дополняем струтурку полем short_base32 и делаем новую запись в БД.
			if daoErr := dao.CreateShortUrl(); daoErr != nil && daoErr.Is(errors.NewAlreadyExistError("ErrorDuplicate")) {
				log.Errorf("[SERVICE] ERROR when trying to call: dao.CreateShortUrl(): %+v", daoErr)
				for daoErr != nil {
					if dao.ShortBase32Inc.Int64 == 10 {
						return nil, daoErr
					}
					dao.ShortBase32Inc.Int64++
					dao.ShortBase32Inc.SetValid(dao.ShortBase32Inc.Int64)
					dao.ShortBase32 = shufflestring.Shuffle(dao.ShortBase32) + strconv.Itoa(int(dao.ShortBase32Inc.Int64)) // Shuffle на случай если в базе уже есть совпадение по ShortUrl
					daoErr = dao.CreateShortUrl()
				}
				log.Infof("[SERVICE] SUCCESS called CreateShortUrl() and returned: %+v", dao)
				return dao, nil
			}
			log.Infof("[SERVICE] SUCCESS called CreateShortUrl() and returned: %+v", dao)
			return dao, nil
		default:
			return nil, err
		}
	}
	log.Infof("[SERVICE] SUCCESS called CreateShortUrl() and returned: %+v", dao)
	return dao, nil
}

func (s *shortUrlService) GetShortUrl(url string) (*shorturl.ShortUrl, *errors.RestErr) {
	log.Infof("[SERVICE] BEGIN CALL: GetShortUrl() with args: %s", url)
	dao, err := initUrlAndHashFromRequest(url)
	if err != nil {
		return nil, err
	}
	if err := dao.GetShortUrlByHash(); err != nil {
		return nil, err
	}
	log.Infof("[SERVICE] SUCCESS called GetShortUrl() and returned: %+v", dao)
	return dao, nil
}

func (s *shortUrlService) GetShortUrlByShortPath(short_path string) (*shorturl.ShortUrl, *errors.RestErr) {
	log.Infof("[SERVICE] BEGIN CALL: GetShortUrl() with args: %s", short_path)
	dao := &shorturl.ShortUrl{}
	dao.ShortBase32 = short_path
	if err := dao.GetUrlByShortBase32(); err != nil {
		return nil, err
	}
	log.Infof("[SERVICE] SUCCESS called GetShortUrlByShortPath() and returned: %+v", dao)
	return dao, nil
}

func (s *shortUrlService) GetAllUrls() ([]shorturl.ShortUrlRequestWithId, *errors.RestErr) {
	log.Infof("[SERVICE] BEGIN CALL: GetAllUrls() with args: ")
	var allUrls []shorturl.ShortUrlRequestWithId
	allUrls, err := shorturl.GetAllUrls(allUrls)
	if err != nil {
		return nil, err
	}
	log.Infof("[SERVICE] SUCCESS called GetAllUrls() and returned: %+v", allUrls)
	return allUrls, nil
}
