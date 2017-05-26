package shadowsocks

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	t.Run("config file not exists", func(t *testing.T) {
		config, err := ParseConfig("file not exists")
		assert.Nil(t, config)
		assert.NotNil(t, err)
	})

	t.Run("parse config", func(t *testing.T) {
		configPath := "./testdata/config.json"
		data, err := ioutil.ReadFile(configPath)
		assert.Nil(t, err)
		sourceConfig := new(Config)
		err = json.Unmarshal(data, sourceConfig)
		assert.Nil(t, err)

		config, err := ParseConfig(configPath)
		assert.Nil(t, err)
		assert.Equal(t, sourceConfig.LocalAddr, config.LocalAddr)
		assert.Equal(t, sourceConfig.ServerAddr, config.ServerAddr)
		assert.Equal(t, sourceConfig.Method, config.Method)
		assert.Equal(t, sourceConfig.Password, config.Password)
		assert.Equal(t, sourceConfig.Timeout, config.Timeout)
	})
}
