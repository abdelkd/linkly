package util

import (
	"crypto/sha1"
	"fmt"
)

func HashString(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
