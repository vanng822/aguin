package crypto

import (
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"aguin/utils"
)

func Test() {
	fmt.Println("Crypto Testing")
}

func Decrypt(data string, key []byte) (map[string]interface{}, error) {
	parts := strings.Split(data, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("")
	}
	messageMAC := []byte(parts[0])
	message, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("")
	}
	
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	if hmac.Equal(messageMAC, expectedMAC) == true {
		return utils.Bytes2json(&message)
	}
	return nil, fmt.Errorf("")
}

func Encrypt(data map[string]interface{}, key []byte) {
	
}