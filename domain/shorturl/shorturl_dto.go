package shorturl

import (
	"gopkg.in/guregu/null.v4"
)

type ShortUrl struct {
	ShortUrlRequest
	Id             int64       `json:"id"`
	ShortBase32Inc null.Int    `json:"short_base32_inc"`
	UrlStatus      null.Int    `json:"url_status"`
	LastCheckTime  null.String `json:"last_check_time"`
	RedirectCount  int64       `json:"redirect_count"`
}

type ShortUrlRequest struct {
	Url         string `json:"url"`
	UrlHash     string `json:"hash"`
	ShortBase32 string `json:"short_base32,omitempty"`
}

type ShortUrlRequestWithId struct {
	Id int64 `json:"id"`
	ShortUrlRequest
}
