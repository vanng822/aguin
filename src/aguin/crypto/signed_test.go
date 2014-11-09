
package crypto
import (
	"testing"
	"fmt"
)

func TestDecrypt(t *testing.T) {
	t.Fail()
}

func TestEncrypt(t *testing.T) {
	res, _ := Encrypt(map[string]interface{}{"entity": "something2"}, []byte("123123123123123123123123213213"))
	fmt.Println(res)
	r, err := Decrypt(res, []byte("123123123123123123123123213213"))
	fmt.Println(r, err)
}


