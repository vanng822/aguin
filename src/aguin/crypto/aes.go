package crypto
import(
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)
func AESDecrypt(message, key []byte) []byte {
	if len(message) < aes.BlockSize {
		panic("Message too short")
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	
	iv := message[:aes.BlockSize]
	ciphertext := message[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func AESEncrypt(message, key []byte) []byte {
	block, err := aes.NewCipher(key)
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
