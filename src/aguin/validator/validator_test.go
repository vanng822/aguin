package validator

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestValidateTrue(t *testing.T) {
	d := map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true}
	vd, validated := Validate(d)
	assert.Equal(t, true, validated)
	assert.Equal(t, d, vd)
}

func TestValidateDateTrue(t *testing.T) {
	d := map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true, "field4": "2013-02-03 02:01:01"}
	xtime, _ := time.Parse("2006-01-02 03:04:01", "2013-02-03 02:01:01")
	xd := map[string]interface{}{"field1": 1, "field2": 3.0, "field3": true, "field4": xtime}
	vd, validated := Validate(d)
	assert.Equal(t, true, validated)
	assert.Equal(t, xd, vd)
}

func TestValidateTrueArray(t *testing.T) {
	d := map[string]interface{}{"field1": 1, "field2": 3.0, "field3": []interface{}{1,2,3.5,4}}
	vd, validated := Validate(d)
	assert.Equal(t, true, validated)
	assert.Equal(t, d, vd)
}

func TestValidateFalse(t *testing.T) {
	vd, validated := Validate(map[string]interface{}{"field1": 1, "field2": "sdfsdf", "field3": []interface{}{1,2,3.5,4}})
	assert.Equal(t, false, validated)
	assert.Equal(t, map[string]interface{}{}, vd)
}