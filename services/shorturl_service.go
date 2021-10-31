package services

import (
	"fmt"
	"github.com/Moriartii/url-shortner-api/domain/shorturl"
	"github.com/Moriartii/url-shortner-api/utils/encode"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/Moriartii/url-shortner-api/utils/hash"
)

type shortUrlService struct {
}

type shortUrlServiceInterface interface {
	GetShortUrl(string) (*shorturl.ShortUrl, *errors.RestErr)
	CreateShortUrl(string) (*shorturl.ShortUrl, *errors.RestErr)
	GetShortUrlByShortPath(string) (*shorturl.ShortUrl, *errors.RestErr)
	IncrementShortUrlCount(string) *errors.RestErr
}

var (
	ShortUrlService shortUrlServiceInterface = &shortUrlService{}
)

func (s *shortUrlService) IncrementShortUrlCount(short_base32 string) *errors.RestErr {
	dao := &shorturl.ShortUrl{
		ShortBase32: short_base32,
	}
	if err := dao.IncrementRedirectCount(); err != nil {
		return err
	}
	return nil
}

func (s *shortUrlService) CreateShortUrl(url string) (*shorturl.ShortUrl, *errors.RestErr) {
	dao := &shorturl.ShortUrl{
		Url: url,
	}
	hashed, err := hash.UrlToHash(url)
	if err != nil {
		return nil, err
	}
	dao.UrlHash = fmt.Sprintf("%x", string(hashed)) // Теперь в струтуре заполнены URL и HASH

	// Проверяем по хэшу, есть ли уже в БД такая запись. Если есть, возвращаем её, если нет, то создаем новую.
	if err := dao.GetShortUrlByHash(); err != nil {
		if err.Message == "no value" {
			// Дополняем струтурку: base32, short_base32 и делаем новую запись в БД.
			dao.HashInBase32 = encode.HashToB32(dao.UrlHash)
			dao.ShortBase32 = dao.HashInBase32[0:6]
			err := dao.CreateShortUrl()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return dao, nil
}

func (s *shortUrlService) GetShortUrl(url string) (*shorturl.ShortUrl, *errors.RestErr) {
	dao := &shorturl.ShortUrl{
		Url: url,
	}
	hashed, err := hash.UrlToHash(url)
	if err != nil {
		return nil, err
	}
	dao.UrlHash = string(hashed)

	if err := dao.GetShortUrlByHash(); err != nil {
		return nil, err
	}
	return dao, nil
}

func (s *shortUrlService) GetShortUrlByShortPath(short_path string) (*shorturl.ShortUrl, *errors.RestErr) {
	dao := &shorturl.ShortUrl{
		ShortBase32: short_path,
	}
	if err := dao.GetShortUrlByShortBase32(); err != nil {
		return nil, err
	}
	return dao, nil
}
