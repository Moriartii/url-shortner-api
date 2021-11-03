package shorturl

import (
	"fmt"
	"github.com/Moriartii/url-shortner-api/db/postgres"
	"github.com/Moriartii/url-shortner-api/logger"
	"github.com/Moriartii/url-shortner-api/utils/errors"
	"github.com/Moriartii/url-shortner-api/utils/postgres_utils"
	"go.uber.org/zap"
	"strings"
)

const (
	queryGetAllUrls             = "SELECT id, url, hash FROM short_urls ORDER BY id ;"
	queryGetShortUrlByHash      = "SELECT id, url, short_base32 FROM short_urls WHERE hash = $1 ;"
	queryGetShortUrlByBase32    = "SELECT id, url, hash, short_base32_inc, url_http_status, last_check_time, redirect_count FROM short_urls WHERE short_base32 = $1 ;"
	queryCreateShortUrl         = "INSERT INTO short_urls(url, hash, short_base32, short_base32_inc) VALUES($1, $2, $3, $4) returning id ;"
	queryIncrementRedirectCount = "UPDATE short_urls SET redirect_count = redirect_count+1 WHERE short_base32 = $1 ;"
)

var (
	log *zap.SugaredLogger
)

func init() {
	log = logger.GetLogger().Named("shorturl (shorturl_dao.go)").Sugar()
}

type ShortUrlInterface interface {
	CreateShortUrl() *errors.RestErr
	GetShortUrlByHash() *errors.RestErr
	GetShortUrlByShortBase32() *errors.RestErr
	IncrementRedirectCount() *errors.RestErr
}

func GetAllUrls(a []ShortUrlRequestWithId) ([]ShortUrlRequestWithId, *errors.RestErr) {
	log.Info("[DOMAIN] BEGIN CALL: GetAllUrls()")
	stmt, err := postgres.Client.Prepare(queryGetAllUrls)
	if err != nil {
		log.Error("Ошибка подготовки запроса в БД", err)
		return nil, errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Error("Ошибка: ", err)
		return nil, errors.NewInternalServerError("ошибка при запросе списка url с БД")
	}

	for rows.Next() {
		var temp ShortUrlRequestWithId
		err = rows.Scan(&temp.Id, &temp.Url, &temp.UrlHash)
		if err != nil {
			log.Error("Ошибка: ", err)
			return nil, errors.NewInternalServerError("ошибка при row.Scan в структуры")
		}
		a = append(a, temp)
	}
	log.Infof("[DOMAIN] SUCCESS called GetAllUrls() and returned: %+v", a)
	return a, nil
}

func (s *ShortUrl) IncrementRedirectCount() *errors.RestErr {
	log.Infof("[DOMAIN] BEGIN CALL: IncrementRedirectCount() with DAO: %+v", s)
	stmt, err := postgres.Client.Prepare(queryIncrementRedirectCount)
	if err != nil {
		log.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	_, sqlErr := stmt.Exec(s.ShortBase32)
	if sqlErr != nil {
		log.Error("Ошибка: ", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	log.Info("[DOMAIN] SUCCESS called IncrementRedirectCount()")
	return nil
}

func (s *ShortUrl) CreateShortUrl() *errors.RestErr {
	log.Infof("[DOMAIN] BEGIN CALL: CreateShortUrl() with DAO: %+v", s)
	stmt, err := postgres.Client.Prepare(queryCreateShortUrl)
	if err != nil {
		log.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.Url, s.UrlHash, s.ShortBase32, s.ShortBase32Inc)
	err = result.Scan(&s.Id)

	if err != nil {
		fmt.Println(err)
		log.Error("Ошибка получения информации о ShortUrl из БД", err)
		if strings.Contains(err.Error(), postgres_utils.ErrorDuplicate) {
			log.Errorf("Ошибка_1: %+v", err)
			return errors.NewAlreadyExistError("ErrorDuplicate")
		} else {
			log.Errorf("Ошибка_2: %+v", err)
			return errors.NewInternalServerError("ошибка при работе с БД")
		}
	}
	log.Infof("[DOMAIN] SUCCESS called CreateShortUrl() and returned: %+v", s)
	return nil
}

func (s *ShortUrl) GetShortUrlByHash() *errors.RestErr {
	log.Infof("[DOMAIN] BEGIN CALL: GetShortUrlByHash() with DAO: %+v", s)
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByHash)
	if err != nil {
		log.Error("Ошибка подготовки запроса в БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.UrlHash)

	err = result.Scan(
		&s.Id,
		&s.Url,
		&s.ShortBase32)

	if err != nil {
		log.Errorf("[DOMAIN] ERROR when trying to call: dao.GetShortUrlByHash(): %+v", err)
		if strings.Contains(err.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		//log.Error("Ошибка получения информации о ShortUrl из БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	log.Infof("[DOMAIN] SUCCESS called CreateShortUrl() and returned: %+v", s)
	return nil
}

func (s *ShortUrl) GetUrlByShortBase32() *errors.RestErr {
	log.Infof("[DOMAIN] BEGIN CALL: GetUrlByShortBase32() with DAO: %+v", s)
	stmt, err := postgres.Client.Prepare(queryGetShortUrlByBase32)
	if err != nil {
		log.Error("Ошибка подготовки запроса в БД", err)
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
		log.Error("Ошибка: ", err)
		if strings.Contains(err.Error(), postgres_utils.ErrorNoRows) {
			return errors.NewNotFoundError("no value")
		}
		log.Error("Ошибка получения информации о ShortUrl из БД", err)
		return errors.NewInternalServerError("ошибка при работе с БД")
	}
	log.Infof("[DOMAIN] SUCCESS called GetUrlByShortBase32() and returned: %+v", s)
	return nil
}
