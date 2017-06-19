package shadowsocks

import (
	"encoding/json"
	"io/ioutil"
)

var (
	defaultServerAddr = "0.0.0.0:8089"
	defaultLocalAddr  = "0.0.0.0:6009"
	defaultTimeout    = 300
	defaultMethod     = "aes-128-cfb"
	defaultPassword   = "shadowsocks-secret-key"

	config = new(Config)
)

type Config struct {
	// local addr, default is ":6009"
	LocalAddr string `json:"local_addr"`
	// server addr, default is ":7008"
	ServerAddr string `json:"server_addr"`
	// password, default is "shadowsocks-secret-key"
	Password string `json:"password"`
	// cipher method, default is "ase-128-cfb"
	Method string `json:"method"`
	// connection timeout, default is 300
	Timeout int `json:"timeout"`
}

func ParseConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	if config.ServerAddr == "" {
		config.ServerAddr = defaultServerAddr
	}

	if config.LocalAddr == "" {
		config.LocalAddr = defaultLocalAddr
	}

	if config.Password == "" {
		config.Password = defaultPassword
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	if config.Method == "" {
		config.Method = defaultMethod
	}

	return config, nil
}
