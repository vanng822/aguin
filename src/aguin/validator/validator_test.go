package validator

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
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

func TestValidateSearchDefaultDates(t *testing.T) {
	now := time.Now()
	r := ValidateSearch(map[string]interface{}{"entity": "testing"})
	assert.Equal(t, r.Entity, "testing")
	assert.Equal(t, r.EndDate.Year(), now.Year())
	assert.Equal(t, r.EndDate.YearDay(), now.YearDay())
	assert.Equal(t, r.EndDate.Hour(), 23)
	assert.Equal(t, r.EndDate.Minute(), 59)
	assert.Equal(t, r.EndDate.Second(), 59)
	
	end := now.AddDate(0, 0, -30)
	assert.Equal(t, r.StartDate.Year(), end.Year())
	assert.Equal(t, r.StartDate.YearDay(), end.YearDay())
	assert.Equal(t, r.StartDate.Hour(), 0)
	assert.Equal(t, r.StartDate.Minute(), 0)
	assert.Equal(t, r.StartDate.Second(), 0)
}

