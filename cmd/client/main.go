package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/cssivision/shadowsocks"
)

var (
	config *shadowsocks.Config
)

func handleConnection(conn net.Conn) {

}

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()
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

	lis, err := net.Listen("tcp", config.LocalAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)
	}
}
