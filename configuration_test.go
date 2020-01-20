package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestHandlingConfig(t *testing.T) {
	Init("development", "test", DEBUG)

	assert.Equal(t, "dev_table", GetString("database.settingsCollection"))
	assert.Equal(t, 8888, GetInteger("server.port"))
	assert.Equal(t, "cool", GetString("nested.more.extrem.super.extrem.fancy"))
	assert.False(t, GetBoolean("auth.enabled"))
	assert.True(t, GetBooleanWithDefaultValue("auth.enabled.new", true))

	AddMapToConfig("", map[string]interface{}{
		"foobar": "toll",
		"lorem": map[string]interface{}{
			"ipsum": true,
		},
	})

	assert.Equal(t, "toll", GetString("foobar"))
	assert.True(t, GetBooleanWithDefaultValue("lorem.ipsum", false))

	tenant := "tenant1"

	cfg := GetSubConfig("tenants." + tenant)

	assert.Equal(t, "t1", cfg.GetString("id"))
	keys := cfg.GetKeys()

	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "id")

	firstLevelKeys := GetFirstLevelKeys()

	assert.Equal(t, 9, len(firstLevelKeys))
	log.Printf("%v", firstLevelKeys)

}
