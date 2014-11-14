package validator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestValidateTrue(t *testing.T) {
	d := map[string]interface{}{"entity": "testing", "data": map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true}}
	r := ValidateEntity(d)
	assert.Equal(t, true, r.Validated)
	assert.Equal(t, d["data"], r.Data)
	assert.Equal(t, d["entity"], r.Entity)
}

func TestValidateDateTrue(t *testing.T) {
	d := map[string]interface{}{"entity": "testing", "data": map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true}}
	xd := map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true}
	r := ValidateEntity(d)
	assert.Equal(t, true, r.Validated)
	assert.Equal(t, xd, r.Data)
	assert.Equal(t, d["entity"], r.Entity)
}

func TestValidateTrueArray(t *testing.T) {
	d := map[string]interface{}{"entity": "testing", "data": map[string]interface{}{"field1": 1, "field2": 3.0, "field3": []interface{}{1,2,3.5,4}}}
	r := ValidateEntity(d)
	assert.Equal(t, true, r.Validated)
	assert.Equal(t, d["data"], r.Data)
	assert.Equal(t, d["entity"], r.Entity)
}

func TestValidateFalse(t *testing.T) {
	r := ValidateEntity(map[string]interface{}{"entity": "testing", "data": map[string]interface{}{"field1": 1, "field2": "sdfsdf", "field3": []interface{}{1,2,3.5,4}}})
	assert.Equal(t, false, r.Validated)
	assert.Equal(t, map[string]interface{}{}, r.Data)
	assert.Equal(t, "testing", r.Entity)
}