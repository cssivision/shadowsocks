package shadowsocks

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConn(t *testing.T) {
	addr := ":3000"
	message := bytes.Repeat([]byte("test message"), 1<<10)
	cipher, err := NewCipher("aes-128-cfb", "password")
	assert.Nil(t, err)

	go func() {
		conn, err := net.Dial("tcp", addr)
		assert.Nil(t, err)
		sconn := NewConn(conn, cipher)
		defer sconn.Close()

		_, err = sconn.Write(message)
		assert.Nil(t, err)
	}()

	lis, err := net.Listen("tcp", addr)
	assert.Nil(t, err)

	for {
		conn, err := lis.Accept()
		assert.Nil(t, err)
		sconn := NewConn(conn, cipher)
		defer sconn.Close()
		buf, err := ioutil.ReadAll(sconn)
		assert.Nil(t, err)
		assert.Equal(t, message, buf)
		return
	}
}
