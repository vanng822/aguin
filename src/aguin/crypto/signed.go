package crypto

import (
	"aguin/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)


func Decrypt(data string, key []byte) (interface{}, error) {
	parts := strings.Split(data, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("")
	}
	messageMAC := []byte(Urldecode(parts[0]))
	message, err := base64.StdEncoding.DecodeString(Urldecode(parts[1]))
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

func Encrypt(data interface{}, key []byte) (string, error) {
	jdata, err := utils.Json2bytes(data)
	if err != nil {
		return "", nil
	}
	
	mac := hmac.New(sha256.New, key)
	message := base64.StdEncoding.EncodeToString(jdata)
	
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	expectedMACBase64 := base64.StdEncoding.EncodeToString(expectedMAC)
	return fmt.Sprintf("%s.%s", Urlencode(expectedMACBase64), Urlencode(message)), nil
}


func Urlencode(base64String string) string {
	base64String = strings.Replace(base64String, "+", "-", -1)
	return strings.Replace(base64String, "/", "_", -1)
}

func Urldecode(base64String string) string {
	base64String = strings.Replace(base64String, "-", "+", -1)
	return strings.Replace(base64String, "_", "/", -1)
}

