package shorturl

import (
	"fmt"
	"strings"

	"github.com/Moriartii/url-shortner-api/db/postgres"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/Moriartii/url-shortner-api/utils/postgres_utils"
)

const (
	queryGetShortUrlByHash      = "SELECT id, url, short_base32 FROM short_urls WHERE hash = $1 ;"
	queryGetShortUrlByBase32    = "SELECT id, url, hash, short_base32_inc, url_http_status, last_check_time, redirect_count FROM short_urls WHERE short_base32 = $1 ;"
	queryCreateShortUrl         = "INSERT INTO short_urls(url, hash, short_base32, short_base32_inc) VALUES($1, $2, $3, $4) returning id ;"
	queryIncrementRedirectCount = "UPDATE short_urls SET redirect_count = redirect_count+1 WHERE short_base32 = $1 ;"
)

type ShortUrlInterface interface {
	CreateShortUrl() *errors.RestErr
	GetShortUrlByHash() *errors.RestErr
	GetShortUrlByShortBase32() *errors.RestErr
	IncrementRedirectCount() *errors.RestErr
}

func (s *ShortUrl) IncrementRedirectCount() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryIncrementRedirectCount)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	_, sqlErr := stmt.Exec(s.ShortBase32)
	if sqlErr != nil {
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	return nil
}

func (s *ShortUrl) CreateShortUrl() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryCreateShortUrl)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.Url, s.UrlHash, s.ShortBase32, s.ShortBase32Inc)
	err = result.Scan(&s.Id)

	if err != nil {
		fmt.Println(err)
		logger.Error("Ошибка получения информации о ShortUrl из БД", err)
		if strings.Contains(err.Error(), postgres_utils.ErrorDuplicate) {
			return errors.NewAlreadyExistError("ErrorDuplicate")
		} else {
			return errors.NewInternalServerError("ошибка при работе с БД")
		}
	}
	return nil
}

func (s *ShortUrl) GetShortUrlByHash() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByHash)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.UrlHash)

	err = result.Scan(
		&s.Id,
		&s.Url,
		&s.ShortBase32)

	if err != nil {
		if strings.Contains(err.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		logger.Error("Ошибка получения информации о ShortUrl из БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	return nil
}

func (s *ShortUrl) GetUrlByShortBase32() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByBase32)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.ShortBase32)

	err = result.Scan(
		&s.Id,
		&s.Url,
		&s.UrlHash,
		&s.ShortBase32Inc,
		&s.UrlStatus,
		&s.LastCheckTime,
		&s.RedirectCount)

	if err != nil {
		if strings.Contains(err.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		logger.Error("Ошибка получения информации о ShortUrl из БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	return nil
}
