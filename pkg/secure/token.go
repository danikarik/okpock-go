package secure

import (
	"crypto/sha256"
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

// Token generates a new secure token.
func Token() string {
	id := uuid.NewV4().String()
	hasher := sha256.New()
	hasher.Write([]byte(id))
	return hex.EncodeToString(hasher.Sum(nil))
}
