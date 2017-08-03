# shadowsocks
[![Build Status](https://img.shields.io/travis/cssivision/shadowsocks.svg?style=flat-square)](https://travis-ci.org/cssivision/shadowsocks)
[![Coverage Status](http://img.shields.io/coveralls/cssivision/shadowsocks.svg?style=flat-square)](https://coveralls.io/github/cssivision/shadowsocks?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/cssivision/shadowsocks/blob/master/LICENSE)

minimalist port of shadowsocks, only reserve basic feature for personal usage.

# Installation
#### Fetch the [latest release](https://github.com/cssivision/shadowsocks/releases).
### linux 
```sh
https://github.com/cssivision/shadowsocks/releases/download/v0.1.5/ssserver.linux
```
### macos
```sh 
https://github.com/cssivision/shadowsocks/releases/download/v0.1.5/ssserver.macos
```
#### use go command
#### server
```sh
go get github.com/cssivision/shadowsocks/cmd/ssserver
```

#### client 
```sh
go get github.com/cssivision/shadowsocks/cmd/sslocal
```
# Configuration
config.json
```json
{
	"server_addr": ":8089",
	"password": "password",
	"local_addr": ":6009",
	"method": "aes-128-cfb",
	"timeout": 300
}
```

# Usage 
#### server
```sh
ssserver -c config.json
```

install a [client](https://shadowsocks.org/en/download/clients.html), connect to your server using your configuration, Done!

# Licenses

All source code is licensed under the [MIT License](https://github.com/cssivision/shadowsocks/blob/master/LICENSE).
