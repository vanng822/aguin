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

func TestValidateSearchInvalid(t *testing.T) {
	r := ValidateSearch(map[string]interface{}{"entity": ""})
	assert.Equal(t, r.Entity, "")
	assert.Equal(t, r.Validated, false)
	assert.Equal(t, r.EndDate.IsZero(), true)
	assert.Equal(t, r.StartDate.IsZero(), true)
}

func TestValidateSearchDefaultDates(t *testing.T) {
	now := time.Now()
	r := ValidateSearch(map[string]interface{}{"entity": "testing"})
	assert.Equal(t, r.Entity, "testing")
	assert.Equal(t, r.Validated, true)
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

func TestValidateSearchValidDateRange(t *testing.T) {
	r := ValidateSearch(map[string]interface{}{"entity": "testing", "startDate": "2014-10-11", "endDate": "2014-11-10"})
	assert.Equal(t, r.Entity, "testing")
	assert.Equal(t, r.EndDate.Year(), 2014)
	assert.Equal(t, r.EndDate.Month(), time.November)
	assert.Equal(t, r.EndDate.Day(), 10)
	assert.Equal(t, r.EndDate.Hour(), 23)
	assert.Equal(t, r.EndDate.Minute(), 59)
	assert.Equal(t, r.EndDate.Second(), 59)
	
	assert.Equal(t, r.StartDate.Year(), 2014)
	assert.Equal(t, r.StartDate.Month(), time.October)
	assert.Equal(t, r.StartDate.Day(), 11)
	assert.Equal(t, r.StartDate.Hour(), 0)
	assert.Equal(t, r.StartDate.Minute(), 0)
	assert.Equal(t, r.StartDate.Second(), 0)
}

func TestValidateSearchLongDateRange(t *testing.T) {
	r := ValidateSearch(map[string]interface{}{"entity": "testing", "startDate": "2011-10-11", "endDate": "2014-11-10"})
	assert.Equal(t, r.Entity, "testing")
	assert.Equal(t, r.EndDate.Year(), 2014)
	assert.Equal(t, r.EndDate.Month(), time.November)
	assert.Equal(t, r.EndDate.Day(), 10)
	assert.Equal(t, r.EndDate.Hour(), 23)
	assert.Equal(t, r.EndDate.Minute(), 59)
	assert.Equal(t, r.EndDate.Second(), 59)
	// end date minus 30 days
	assert.Equal(t, r.StartDate.Year(), 2014)
	assert.Equal(t, r.StartDate.Month(), time.August)
	assert.Equal(t, r.StartDate.Day(), 12)
	assert.Equal(t, r.StartDate.Hour(), 0)
	assert.Equal(t, r.StartDate.Minute(), 0)
	assert.Equal(t, r.StartDate.Second(), 0)
}
