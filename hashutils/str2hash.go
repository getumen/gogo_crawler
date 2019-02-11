package hashutils

import (
	"crypto/sha512"
	"encoding/base64"
)

func MakeHash(b []byte) string {
	f := sha512.New()
	s := base64.URLEncoding.EncodeToString(f.Sum(b))
	if len(s) > 256 {
		return s[:256]
	} else {
		return s
	}
}
