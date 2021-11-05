package postgres_utils

//import (
//"github.com/Moriartii/url-shortner-api/utils/errors"
//"github.com/lib/pq"
//"strings"
//)

const (
	ErrorNoRows    = "no rows in result set"
	ErrorDuplicate = "duplicate key value violates unique"
)

// func ParseError(err error) *errors.RestErr {
// 	sqlErr, ok := err.(*pq.Error)
// 	if !ok {
// 		if strings.Contains(err.Error(), ErrorNoRows) {
// 			return errors.NewNotFoundError("Не найдено совпадений в БД")
// 		} else if strings.Contains(err.Error(), ErrorDuplicate) {
// 			return errors.NewNotFoundError("Уже есть в БД такой PK (hash)")
// 		}
// 		return errors.NewInternalServerError("error parsing database response")
// 	}

// 	switch sqlErr.Code.Name() {
// 	case "unique_violation":
// 		return errors.NewBadRequestError("Invalid data (уникальное значение уже есть в БД?)")
// 	case "short_urls_pkey":
// 		return errors.NewBadRequestError("short_urls_pkey (уникальное значение уже есть в БД?)")
// 	}
// 	return errors.NewInternalServerError("error processing request")
// }
