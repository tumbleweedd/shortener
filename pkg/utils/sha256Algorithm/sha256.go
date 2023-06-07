package sha256Algorithm

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHash(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
