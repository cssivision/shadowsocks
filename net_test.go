package shadowsocks

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialWithCipher(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		addr := ":3001"
		message := bytes.Repeat([]byte("test message"), 1<<10)
		cipher, err := NewCipher("aes-128-cfb", "password")
		require.Nil(t, err)

		go func() {
			sconn, err := DialWithCipher(addr, cipher)
			assert.Nil(t, err)
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
	})

	t.Run("dail with error", func(t *testing.T) {
		addr := "address not exists"
		cipher, err := NewCipher("aes-128-cfb", "password")
		require.Nil(t, err)

		_, err = DialWithCipher(addr, cipher)
		require.NotNil(t, err)
	})
}
