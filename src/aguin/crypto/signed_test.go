
package crypto
import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	authKey := []byte("545e0716f2fea0c7a9c46c74fec46c71")
	expected := map[string]interface{}{"entity": "something2"}
	res, _ := Encrypt(expected, authKey, authKey)
	fmt.Println(res)
	r, err := Decrypt(res, authKey, authKey)
	fmt.Println(r, err)
	assert.Equal(t, expected, r)
}


