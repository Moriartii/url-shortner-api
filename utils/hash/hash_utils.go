package hash

import (
	"crypto/sha256"
	"github.com/Moriartii/url-shortner-api/utils/errors"
)

func UrlToHash(url string) ([]byte, *errors.RestErr) {
	h := sha256.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		return nil, errors.NewInternalServerError("Error when trying to hash")
	}
	return h.Sum(nil), nil
}
