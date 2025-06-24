package util

import (
	"crypto/md5"
	"fmt"
)

func MD5Hash(data []byte) string {
	md5h := md5.New()
	md5h.Write(data)
	return fmt.Sprintf("%x", md5h.Sum(nil))
}
