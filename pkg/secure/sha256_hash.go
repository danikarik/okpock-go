package secure

import (
	"crypto/sha256"
	"encoding/base64"
)

// Hash creates SHA256 hash.
func Hash(data []byte) (string, error) {
	h := sha256.New()

	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum), nil
}
