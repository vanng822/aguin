package aguin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var userAgent = "Aguin/0.1"

type Aguin struct {
	apiKey  string
	authKey []byte
	aesKey  []byte
	url     string
}

func New(apiKey, authKey, aesKey, url string) *Aguin {
	url := strings.TrimRight(url, "/")
	return &Aguin{apiKey, []byte(authKey), []byte(aesKey), url}
}

func (a *Aguin) Decrypt(data string) (interface{}, error) {
	parts := strings.Split(data, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Signed data must contain base64_sig.base64_data")
	}
	expectedMAC, err := base64.StdEncoding.DecodeString(a.Urldecode(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("Invalid data")
	}
	message, err := base64.StdEncoding.DecodeString(a.Urldecode(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("Invalid data")
	}

	mac := hmac.New(sha256.New, a.authKey)
	mac.Write(message)
	messageMAC := mac.Sum(nil)
	if hmac.Equal(messageMAC, expectedMAC) {
		return bytes2json(a.AESDecrypt(message))
	}
	return nil, fmt.Errorf("Invalid signed request")
}

func (a *Aguin) Encrypt(data interface{}) (string, error) {
	message, err := json2bytes(data)
	if err != nil {
		return "", fmt.Errorf("Invalid data")
	}
	message = a.AESEncrypt(message)
	mac := hmac.New(sha256.New, a.authKey)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedMACBase64 := base64.StdEncoding.EncodeToString(expectedMAC)
	messageBase64 := base64.StdEncoding.EncodeToString(message)
	return fmt.Sprintf("%s.%s", a.Urlencode(expectedMACBase64), a.Urlencode(messageBase64)), nil
}

func (a *Aguin) AESDecrypt(message []byte) []byte {
	if len(message) < aes.BlockSize {
		panic("Message too short")
	}

	block, err := aes.NewCipher(a.aesKey)
	if err != nil {
		panic(err)
	}

	iv := message[:aes.BlockSize]
	ciphertext := message[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func (a *Aguin) AESEncrypt(message []byte) []byte {
	block, err := aes.NewCipher(a.aesKey)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], message)

	return ciphertext
}

func (a *Aguin) Urlencode(base64String string) string {
	return strings.NewReplacer("+", "-", "/", "_").Replace(base64String)
}

func (a *Aguin) Urldecode(base64String string) string {
	return strings.NewReplacer("-", "+", "_", "/").Replace(base64String)
}

func (a *Aguin) Get(entity string, criteria map[string]interface{}) (map[string]interface{}, error) {
	client := &http.Client{}
	message := url.Values{}
	if criteria == nil {
		criteria = map[string]interface{}{}
	}
	criteria["entity"] = entity

	cryptedMessage, err := a.Encrypt(criteria)

	message.Add("message", cryptedMessage)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/?%s", a.url, message.Encode()), nil)
	if err != nil {
		return nil, err
	}

	return a.makeRequest(client, req)
}

func (a *Aguin) Status() (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status", a.url), nil)
	if err != nil {
		return nil, err
	}

	return a.makeRequest(client, req)
}

func (a *Aguin) Post(entity string, data map[string]interface{}) (map[string]interface{}, error) {
	client := &http.Client{}
	cryptedMessage, err := a.Encrypt(map[string]interface{}{"entity": entity, "data": data})
	if err != nil {
		return nil, err
	}
	message := url.Values{}
	message.Add("message", cryptedMessage)

	req, err := http.NewRequest("POST", a.url, bytes.NewBufferString(message.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	return a.makeRequest(client, req)
}

func (a *Aguin) makeRequest(client *http.Client, req *http.Request) (map[string]interface{}, error) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-AGUIN-API-KEY", a.apiKey)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res, err := bytes2json(body)
	if err != nil {
		return nil, err
	}
	respJson := res.(map[string]interface{})

	encrypted, ok := respJson["encrypted"].(bool)
	if !ok {
		return nil, fmt.Errorf("Got unexpected response %v", respJson)
	}
	
	if encrypted == true {
		result, ok := respJson["result"].(string)
		if !ok {
			return nil, fmt.Errorf("Got unexpected response %v", respJson)
		}
		resultJson, err := a.Decrypt(result)
		if err != nil {
			return nil, err
		}
		respJson["result"] = resultJson
	}

	return respJson, nil
}

func bytes2json(data []byte) (interface{}, error) {
	var jsonData interface{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func json2bytes(data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}
