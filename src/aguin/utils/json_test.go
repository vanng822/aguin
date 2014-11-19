package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBytes2jsonWrongFormat(t *testing.T) {
	data := []byte(`<xml><x>y</x><yy>boo</yy></xml>`)
	d, err := Bytes2json(data)
	assert.Error(t, err)
	assert.Empty(t, d)	
}

func TestBytes2jsonSyntaxError(t *testing.T) {
	data := []byte(`{"obj":{"c":b,"j":"F"}}`)
	d, err := Bytes2json(data)
	assert.Error(t, err)
	assert.Empty(t, d)	
}

func TestBytes2json(t *testing.T) {
	data := []byte(`{"t":"s","two":2,"obj":{"c":"b","j":"F"},"a":[1,2,3]}`)
	
	d, err := Bytes2json(data)
	d2 := d.(map[string]interface{})
	assert.Empty(t, err)
	assert.Equal(t, d2["t"].(string), "s")	
}

func TestJson2bytes(t *testing.T) {
	data := map[string]interface{}{"t":"s","two":2,"obj":map[string]interface{}{"c":"b","j":"F"},"a":[]interface{}{1,2,3}}
	
	d, err := Json2bytes(data)
	
	assert.Empty(t, err)
	rdata, _ := Bytes2json(d)
	d2 := rdata.(map[string]interface{})
	assert.Equal(t, d2["t"].(string), data["t"].(string))
}