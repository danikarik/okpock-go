package pkpass

import (
	"crypto/sha1"
	"encoding/hex"
)

// HashFile generates SHA1 hash of file content.
func HashFile(content []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(content)
	if err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	hash := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(hash, sum)
	return string(hash), nil
}
