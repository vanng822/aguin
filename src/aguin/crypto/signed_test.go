
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
	r, _ := Decrypt(res, authKey, authKey)
	assert.Equal(t, expected, r)
}

func TestDecrypt(t *testing.T) {
	authKey := []byte("545e0716f2fea0c7a9c46c74fec46c71")
	input := "WmtlMxl4d_VTfUYnl-A0Uycpr2e3VswKDwoPd03XtoY=.JZMk6FZloYh5BL0K7dHGSyTqB4lTgd9annrFEgLTELnxR3bHweL2"
	expected := map[string]interface{}{"entity": "something2"}
	r, _ := Decrypt(input, authKey, authKey)
	assert.Equal(t, expected, r)
	
}
