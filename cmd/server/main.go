package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"syscall"
	"time"

	"github.com/cssivision/shadowsocks"
)

var (
	config *shadowsocks.Config
)

func init() {
	rand.Seed(time.Now().Unix())
}

const (
	idType  = 0 // address type index
	idIP0   = 1 // ip addres start index
	idDmLen = 1 // domain address length index
	idDm0   = 2 // domain address start index

	typeIPv4 = 1 // type is ipv4 address
	typeDm   = 3 // type is domain address
	typeIPv6 = 4 // type is ipv6 address

	lenIPv4     = net.IPv4len + 2 // ipv4 + 2port
	lenIPv6     = net.IPv6len + 2 // ipv6 + 2port
	lenDmBase   = 2               // 1addrLen + 2port, plus addrLen
	lenHmacSha1 = 10

	addrMask byte = 0xf
)

func getAddrInfo(conn *shadowsocks.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 255(max length address) + 2(port) + 10(hmac-sha1)
	buf := make([]byte, 269)
	if _, err := io.ReadFull(conn, buf[:idType+1]); err != nil {
		return "", err
	}

	var reqStart, reqEnd int
	addrType := buf[idType]
	switch addrType & addrMask {
	case typeIPv4:
		reqStart, reqEnd = idIP0, idIP0+lenIPv4
	case typeIPv6:
		reqStart, reqEnd = idIP0, idIP0+lenIPv6
	case typeDm:
		if _, err := io.ReadFull(conn, buf[idType+1:idDmLen+1]); err != nil {
			return "", err
		}
		reqStart, reqEnd = idDm0, idDm0+int(buf[idDmLen])+lenDmBase
	default:
		return "", fmt.Errorf("addr type %d not supported", addrType&addrMask)
	}

	if _, err := io.ReadFull(conn, buf[reqStart:reqEnd]); err != nil {
		return "", err
	}

	var host string
	switch addrType & addrMask {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+int(buf[idDmLen])])
	}

	port := binary.BigEndian.Uint16(buf[reqEnd-2 : reqEnd])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return host, nil
}

func handleConnection(conn *shadowsocks.Conn) {
	host, err := getAddrInfo(conn)
	if err != nil {
		log.Printf("get host error: %v", err)
		return
	}

	// log.Printf("connection to host: %v", host)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			log.Printf("dial error: %v", err)
		} else {
			log.Printf("connecting to %v error: %v", host, err)
		}
		return
	}

	go func() {
		defer conn.Close()
		shadowsocks.CopyBuffer(conn, remote)
	}()

	shadowsocks.CopyBuffer(remote, conn)
	remote.Close()
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

	lis, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		panic(err)
	}

	cipher, err := shadowsocks.NewCipher(config.Method, config.Password)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(shadowsocks.NewConn(conn, cipher.Clone()))
	}
}
