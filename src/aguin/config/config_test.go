package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetConfig(t *testing.T) {
	config := AppConf()
	config.EncryptionEnabled = true
	assert.True(t, config.EncryptionEnabled)
	config2 := AppConf()
	assert.True(t, config2.EncryptionEnabled)

	conf := ServerConf()
	conf.Host = "0.0.0.0"
	assert.Equal(t, "0.0.0.0", conf.Host)
	conf.PidFile = "/var/run/aguin.pid"
	conf.Port = 8090
	
	conf2 := ServerConf()
	assert.Equal(t, "0.0.0.0", conf2.Host)
	assert.Equal(t, 8090, conf2.Port)
	assert.Equal(t, "/var/run/aguin.pid", conf2.PidFile)
}
