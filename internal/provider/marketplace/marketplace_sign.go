package marketplace

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(key string, base string) string {

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(base))

	return hex.EncodeToString(h.Sum(nil))
}
