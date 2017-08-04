package shadowsocks

import (
	"net"
)

// DialWithCipher ...
func DialWithCipher(address string, cipher *Cipher) (*Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewConn(conn, cipher), nil
}
