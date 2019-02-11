package hashutils

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"log"
	"strings"
)

func MakeHash(b []byte) string {
	f := md5.New()
	for _, s := range strings.Split(string(b), "\n") {
		_, _ = io.WriteString(f, s)
	}
	s := base64.URLEncoding.EncodeToString(f.Sum(nil))
	log.Println(s)
	return s
}
