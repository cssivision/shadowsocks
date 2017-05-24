package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/cssivision/shadowsocks"
)

var (
	config *shadowsocks.Config
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("socks connect from %s\n", conn.RemoteAddr().String())
	rawaddr, host, err := shadowsocks.Socks5Negotiate(conn)
	if err != nil {
		log.Printf("socks negotiate with %v error: %v", host, err)
	}

	// Sending connection established message immediately to client.
	// This some round trip time for creating socks connection with the client.
	// But if connection failed, the client will get connection reset error.
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		log.Println("send connection confirmation: ", err)
		return
	}

	cipher, err := shadowsocks.NewCipher(config.Method, config.Password)
	if err != nil {
		log.Println("init cipher method error: ", err)
	}

	serverConn, err := shadowsocks.DialWithCipher(config.ServerAddr, cipher)
	if err != nil {
		log.Println("connect to server error: ", err)
	}
	defer serverConn.Close()

	_, err = serverConn.Write(rawaddr)
	if err != nil {

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
		}

		go handleConnection(conn)
	}
}
