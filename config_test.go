package shadowsocks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	t.Run("config file", func(t *testing.T) {
		config = new(Config)
		configPath := "./testdata/config.json"
		data, err := ioutil.ReadFile(configPath)
		assert.Nil(t, err)
		sourceConfig := new(Config)
		err = json.Unmarshal(data, sourceConfig)
		assert.Nil(t, err)

		config, err := ParseConfig(configPath)
		require.Nil(t, err)
		assert.Equal(t, sourceConfig.LocalAddr, config.LocalAddr)
		assert.Equal(t, sourceConfig.ServerAddr, config.ServerAddr)
		assert.Equal(t, sourceConfig.Method, config.Method)
		assert.Equal(t, sourceConfig.Password, config.Password)
		assert.Equal(t, sourceConfig.Timeout, config.Timeout)
	})

	t.Run("env variables", func(t *testing.T) {
		config = new(Config)
		require.Nil(t, os.Setenv("SHADOWSOCKS_SERVER_ADDR", ":8087"))
		require.Nil(t, os.Setenv("SHADOWSOCKS_LOCAL_ADDR", ":1090"))
		require.Nil(t, os.Setenv("SHADOWSOCKS_PASSWORD", "password"))
		require.Nil(t, os.Setenv("SHADOWSOCKS_TIMEOUT", "400"))
		require.Nil(t, os.Setenv("SHADOWSOCKS_METHOD", "aes-256-cfb"))
		config, err := ParseConfig("")
		require.Nil(t, err)
		assert.Equal(t, config.ServerAddr, ":8087")
		assert.Equal(t, config.LocalAddr, ":1090")
		assert.Equal(t, config.Password, "password")
		assert.Equal(t, config.Timeout, 400)
		assert.Equal(t, config.Method, "aes-256-cfb")
	})

	t.Run("default values", func(t *testing.T) {
		require.Nil(t, os.Unsetenv("SHADOWSOCKS_SERVER_ADDR"))
		require.Nil(t, os.Unsetenv("SHADOWSOCKS_LOCAL_ADDR"))
		require.Nil(t, os.Unsetenv("SHADOWSOCKS_PASSWORD"))
		require.Nil(t, os.Unsetenv("SHADOWSOCKS_TIMEOUT"))
		require.Nil(t, os.Unsetenv("SHADOWSOCKS_METHOD"))
		config = new(Config)
		config, err := ParseConfig("")
		require.Nil(t, err)
		assert.Equal(t, config.ServerAddr, defaultServerAddr)
		assert.Equal(t, config.LocalAddr, defaultLocalAddr)
		assert.Equal(t, config.Password, defaultPassword)
		assert.Equal(t, config.Timeout, defaultTimeout)
		assert.Equal(t, config.Method, defaultMethod)
	})
}
