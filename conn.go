package shadowsocks

import (
	"errors"
	"io"
	"net"
)

type Conn struct {
	net.Conn
	*Cipher
}

func NewConn(c net.Conn, cipher *Cipher) *Conn {
	return &Conn{
		Conn:   c,
		Cipher: cipher,
	}
}

func (c *Conn) Read(b []byte) (int, error) {
	// get iv in the first read
	if c.dec == nil {
		iv := make([]byte, c.info.ivLen)
		if _, err := io.ReadFull(c.Conn, iv); err != nil {
			return 0, err
		}
		if err := c.initDecrypt(iv); err != nil {
			return 0, err
		}

		if len(c.iv) == 0 {
			c.iv = iv
		}
	}
	cipherData := make([]byte, len(b))

	n, err := c.Conn.Read(cipherData)
	if n > 0 {
		c.Decrypt(b[0:n], cipherData[0:n])
	}
	return n, err
}

func (c *Conn) Write(b []byte) (int, error) {
	var iv []byte
	if c.enc == nil {
		if err := c.initEncrypt(); err != nil {
			return 0, err
		}
		if len(c.iv) == 0 {
			return 0, errors.New("get iv error")
		}
		iv = c.iv
	}

	cipherData := make([]byte, len(iv)+len(b))
	if len(iv) > 0 {
		copy(cipherData, iv)
	}

	c.Encrypt(cipherData[len(iv):], b)
	n, err := c.Conn.Write(cipherData)
	return n, err
}
