# shadowsocks
minimalist port of shadowsocks, only reserve basic feature for personal usage.

# Installation
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
