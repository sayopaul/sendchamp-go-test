package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256HashString(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
