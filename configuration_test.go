package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlingConfig(t *testing.T) {
	Init("development", "test", INFO)

	assert.Equal(t, "dev_table", GetStringConfig("database.settingsCollection"))
	assert.Equal(t, 8888, GetIntegerConfig("server.port"))
	assert.Equal(t, "cool", GetStringConfig("nested.more.extrem.super.extrem.fancy"))
	assert.False(t, GetBooleanConfig("auth.enabled"))
}
