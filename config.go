package shadowsocks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
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

// ParseConfig parse config both from file or env variable
func ParseConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	if config.ServerAddr == "" {
		serverAddr := os.Getenv("SHADOWSOCKS_SERVER_ADDR")
		if serverAddr == "" {
			serverAddr = defaultServerAddr
		}

		config.ServerAddr = serverAddr
	}

	if config.LocalAddr == "" {
		localAddr := os.Getenv("SHADOWSOCKS_LOCAL_ADDR")
		if localAddr == "" {
			localAddr = defaultLocalAddr
		}
		config.LocalAddr = localAddr
	}

	if config.Password == "" {
		password := os.Getenv("SHADOWSOCKS_PASSWORD")
		if password == "" {
			password = defaultPassword
		}
		config.Password = password
	}

	if config.Timeout == 0 {
		var timeout int64
		var err error
		timeoutStr := os.Getenv("SHADOWSOCKS_TIMEOUT")
		if timeoutStr == "" {
			timeout = int64(defaultTimeout)
		} else {
			timeout, err = strconv.ParseInt(timeoutStr, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		config.Timeout = int(timeout)
	}

	if config.Method == "" {
		method := os.Getenv("SHADOWSOCKS_METHOD")
		if method == "" {
			method = defaultMethod
		}
		config.Method = method
	}

	return config, nil
}

// GetConfig ...
func GetConfig() *Config {
	return config
}
