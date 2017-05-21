package shadowsocks

import (
	"net"
)

type Conn struct {
	net.Conn
	*Cipher
}

func NewConn() *Conn {

}

func (c *Conn) Read(b []byte) (int, error) {

}

func (c *Conn) Write(b []byte) (int, error) {

}
