package toolbox

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func HashSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	hashSum := hash.Sum(nil)
	return hashSum
}

func NewRandomBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func NewRandomString(length int) string {
	len := int(float64(length)/float64(1.333333336)) + 1
	bytes := NewRandomBytes(len)
	return base64.URLEncoding.EncodeToString(bytes)[0:length]
}
