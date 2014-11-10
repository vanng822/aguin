package crypto

import (
	"aguin/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func Decrypt(data string, authKey, aesKey []byte) (interface{}, error) {
	parts := strings.Split(data, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Signed data must contain base64_sig.base64_data")
	}
	expectedMAC, err := base64.StdEncoding.DecodeString(Urldecode(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("Invalid data")
	}
	message, err := base64.StdEncoding.DecodeString(Urldecode(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("Invalid data")
	}

	mac := hmac.New(sha256.New, authKey)
	mac.Write(message)
	messageMAC := mac.Sum(nil)
	if hmac.Equal(messageMAC, expectedMAC) {
		return utils.Bytes2json(AESDecrypt(message, aesKey))
	}
	return nil, fmt.Errorf("Invalid signed request")
}

/*
* Data from server to client shouldn't be a problem with + and /
* but we encode it anyways so we can use this logic for both server and client side
 */
func Encrypt(data interface{}, authKey, aesKey []byte) (string, error) {
	message, err := utils.Json2bytes(data)
	if err != nil {
		return "", fmt.Errorf("Invalid data")
	}
	message = AESEncrypt(message, aesKey)
	mac := hmac.New(sha256.New, authKey)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedMACBase64 := base64.StdEncoding.EncodeToString(expectedMAC)
	messageBase64 := base64.StdEncoding.EncodeToString(message)
	return fmt.Sprintf("%s.%s", Urlencode(expectedMACBase64), Urlencode(messageBase64)), nil
}

func Urlencode(base64String string) string {
	return strings.NewReplacer("+", "-", "/", "_").Replace(base64String)
}

func Urldecode(base64String string) string {
	return strings.NewReplacer("-", "+", "_", "/").Replace(base64String)
}
