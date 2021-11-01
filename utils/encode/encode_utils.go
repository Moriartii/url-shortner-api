package encode

import (
	"encoding/base32"
)

func HashToBase32(s string) string {
	b32 := base32.StdEncoding.EncodeToString([]byte(s))
	return b32
}
