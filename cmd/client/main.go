package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cssivision/shadowsocks"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()
	var config *shadowsocks.Config
	var err error
	config, err = shadowsocks.ParseConfig(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	configLog := fmt.Sprintf(
		"LocalAddr:  %v\nServerAddr: %v\nPassword:   %v\nMethod:     %v",
		config.LocalAddr,
		config.ServerAddr,
		config.Password,
		config.Method,
	)
	fmt.Println(configLog)
}
