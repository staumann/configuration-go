package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlingConfig(t *testing.T) {
	Init("development", "test", DEBUG)

	assert.Equal(t, "dev_table", GetString("database.settingsCollection"))
	assert.Equal(t, 8888, GetInteger("server.port"))
	assert.Equal(t, "cool", GetString("nested.more.extrem.super.extrem.fancy"))
	assert.False(t, GetBoolean("auth.enabled"))
	assert.True(t, GetBooleanWithDefaultValue("auth.enabled.new", true))

	AddMapToConfig(map[interface{}]interface{}{
		"foobar": "toll",
		"lorem": map[interface{}]interface{}{
			"ipsum": true,
		},
	})

	assert.Equal(t, "toll", GetString("foobar"))
	assert.True(t, GetBooleanWithDefaultValue("lorem.ipsum", false))
}
