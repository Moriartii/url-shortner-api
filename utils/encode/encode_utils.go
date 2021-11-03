package encode

import (
	"encoding/base32"
	//"github.com/gin-gonic/gin"
	"bytes"
	"io"
	"io/ioutil"
)

func HashToBase32(s string) string {
	b32 := base32.StdEncoding.EncodeToString([]byte(s))
	return b32
}

func RequestBodyForLogger(req io.ReadCloser) ([]byte, io.ReadCloser, error) {
	body, err := ioutil.ReadAll(req)
	if err != nil {
		return nil, nil, err
	}
	req = ioutil.NopCloser(bytes.NewReader(body))
	return body, req, nil
}
