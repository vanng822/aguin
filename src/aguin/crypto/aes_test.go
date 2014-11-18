package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/base64"
)

func TestAESEncryption(t *testing.T) {
	key := []byte("545e0716f2fea0c7a9c46c74fec46c71")
	expected := []byte("testing string")
	encrypted := AESEncrypt(expected, key)
	assert.Equal(t, expected, AESDecrypt(encrypted, key))
}

func TestAESDecrypt(t *testing.T) {
	key := []byte("545e0716f2fea0c7a9c46c74fec46c71")
	input, _ := base64.StdEncoding.DecodeString("JZMk6FZloYh5BL0K7dHGSyTqB4lTgd9annrFEgLTELnxR3bHweL2")
	expected := []byte("{\"entity\":\"something2\"}")
	assert.Equal(t, expected, AESDecrypt(input, key))
}