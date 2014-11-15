package crypto

import (
	"encoding/hex"
	"io"
	"crypto/rand"
)

func RandomHex(byteLenth int) string {
	result := make([]byte, byteLenth)
	if _, err := io.ReadFull(rand.Reader, result); err != nil {
		panic(err)
	}
	return hex.EncodeToString(result)
}