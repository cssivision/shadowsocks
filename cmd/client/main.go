package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"

	"github.com/cssivision/shadowsocks"
)

var (
	config *shadowsocks.Config
)

func init() {
	rand.Seed(time.Now().Unix())
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("socks connect from %s\n", conn.RemoteAddr().String())
	rawaddr, host, err := shadowsocks.Socks5Negotiate(conn)
	if err != nil {
		log.Printf("socks negotiate with %v error: %v", host, err)
	}

	cipher, err := shadowsocks.NewCipher(config.Method, config.Password)
	if err != nil {
		log.Println("init cipher method error: ", err)
	}

	serverConn, err := shadowsocks.DialWithCipher(config.ServerAddr, cipher.Reset())
	if err != nil {
		log.Printf("connect to server error: %v", err)
		return
	}
	defer serverConn.Close()

	_, err = serverConn.Write(rawaddr)
	if err != nil {
		log.Printf("write data to server error: %v", err)
		return
	}

	go io.Copy(conn, serverConn)
	io.Copy(serverConn, conn)
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
			continue
		}

		go handleConnection(conn)
	}
}
