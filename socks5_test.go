package shadowsocks

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
)

func TestSocks5Negotiate(t *testing.T) {
	addr := ":3002"

	go func() {
		forward := proxy.FromEnvironment()
		dailer, err := proxy.SOCKS5("tcp", addr, nil, forward)
		assert.Nil(t, err)

		conn, err := dailer.Dial("tcp", "looli.xyz:80")
		assert.Nil(t, err)
		defer conn.Close()
	}()

	lis, err := net.Listen("tcp", addr)
	assert.Nil(t, err)

	for {
		conn, err := lis.Accept()
		assert.Nil(t, err)

		_, host, err := Socks5Negotiate(conn)
		assert.Nil(t, err)
		assert.Equal(t, "looli.xyz:80", host)
		// assert.Equal(t, []byte{0x3, 0x9, 0x6c, 0x6f, 0x6f, 0x6c, 0x69, 0x2e, 0x78, 0x79, 0x7a, 0x0, 0x50}, rawaddr)
		return
	}
}
