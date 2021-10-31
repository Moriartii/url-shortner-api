package shorturl

import (
	"github.com/Moriartii/url-shortner-api/db/postgres"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/Moriartii/url-shortner-api/utils/postgres_utils"
	"github.com/Moriartii/url-shortner-api/utils/shufflestring"
	"strings"
)

const (
	queryGetShortUrlByHash      = "SELECT id, url, hash, hash_in_base32, short_base32 FROM short_urls WHERE hash = $1 ;"
	queryGetShortUrlByBase32    = "SELECT id, url, hash, hash_in_base32, short_base32, short_base32_shuffle, url_status, last_check_time, redirect_count FROM short_urls WHERE short_base32 = $1 ;"
	queryCreateShortUrl         = "INSERT INTO short_urls(url, hash, hash_in_base32, short_base32, short_base32_shuffle, redirect_count) VALUES($1, $2, $3, $4, $5, $6) returning id, url, hash, hash_in_base32, short_base32, short_base32_shuffle ;"
	queryIncrementRedirectCount = "UPDATE short_urls SET redirect_count = redirect_count+1 WHERE short_base32 = $1 ;"
)

var shuffleCount int

type ShortUrlInterface interface {
	CreateShortUrl() *errors.RestErr
	GetShortUrlByHash() *errors.RestErr
	GetShortUrlByShortBase32() *errors.RestErr
	IncrementRedirectCount() *errors.RestErr
}

func (short_url *ShortUrl) IncrementRedirectCount() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryIncrementRedirectCount)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	_, sqlErr := stmt.Exec(short_url.ShortBase32)
	if sqlErr != nil {
		return errors.NewInternalServerError("ошибка при работе 1 с БД")
	}
	return nil
}

func (short_url *ShortUrl) CreateShortUrl() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryCreateShortUrl)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(short_url.Url, short_url.UrlHash, short_url.HashInBase32, short_url.ShortBase32, short_url.ShortBase32Shuffle, 0)
	retErr := result.Scan(&short_url.Id, &short_url.Url, &short_url.UrlHash, &short_url.HashInBase32, &short_url.ShortBase32, &short_url.ShortBase32Shuffle)

	if retErr != nil {
		logger.Error("Ошибка получения информации о ShortUrl из БД", retErr)
		if strings.Contains(retErr.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no hash in db")
		} else if strings.Contains(retErr.Error(), postgres_utils.ErrorDuplicate) {
			short_url.ShortBase32 = shufflestring.Shuffle(short_url.ShortBase32)
			short_url.ShortBase32Shuffle = true
			shuffleCount++
			if shuffleCount == 10 {
				return errors.NewInternalServerError("Ошибка при работе с БД (все варианты short_base32 использованы?)")
			}
			short_url.CreateShortUrl()
		}
	}
	return nil
}

func (short_url *ShortUrl) GetShortUrlByHash() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByHash)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(short_url.UrlHash)

	getErr := result.Scan(&short_url.Id, &short_url.Url, &short_url.UrlHash, &short_url.HashInBase32, &short_url.ShortBase32)
	if getErr != nil {
		if strings.Contains(getErr.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		logger.Error("Ошибка получения информации о ShortUrl из БД", getErr)
		return errors.NewInternalServerError("Ошибка при работе с БД")
	}
	return nil
}

func (short_url *ShortUrl) GetShortUrlByShortBase32() *errors.RestErr {
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByBase32)
	if err != nil {
		logger.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(short_url.ShortBase32)
	getErr := result.Scan(
		&short_url.Id,
		&short_url.Url,
		&short_url.UrlHash,
		&short_url.HashInBase32,
		&short_url.ShortBase32,
		&short_url.ShortBase32Shuffle,
		&short_url.UrlStatus,
		&short_url.LastCheckTime,
		&short_url.RedirectCount)
	if !strings.HasPrefix(short_url.Url, "http") {
		short_url.Url = "http://" + short_url.Url
	}

	if getErr != nil {
		if strings.Contains(getErr.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		logger.Error("Ошибка получения информации о ShortUrl из БД", getErr)
		return errors.NewInternalServerError("Ошибка при работе с БД")
	}
	return nil
}
