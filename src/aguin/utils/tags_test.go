package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTagsFull(t *testing.T) {
	type T struct {
		F string `json:"hej" bson:"f"`
		J string `json:"you" bson:""`
	}
	
	tags := GetFieldsTag(T{}, "bson")
	assert.Equal(t, "f", tags.Get("F"))
	assert.Panics(t, func() {
		tags.Get("J")
	})
	
	assert.Panics(t, func() {
		tags.Get("None")
	})
}