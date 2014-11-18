package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRandomHex(t *testing.T) {
	res := RandomHex(16)
	assert.Equal(t, 32, len(res))
}