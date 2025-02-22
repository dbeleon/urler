package tiny

import (
	"crypto/md5"

	"github.com/dromara/dongle"
)

func Get(url string) string {
	hash := md5.Sum([]byte(url))
	return dongle.Encode.FromBytes(hash[:]).ByBase62().ToString()
}
