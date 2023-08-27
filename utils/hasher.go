package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func GenerateHash(url string) string {
	data := []byte(url + strconv.FormatInt(time.Now().UTC().Unix(), 10))
	return fmt.Sprintf("%x", md5.Sum(data))[:8]
}
