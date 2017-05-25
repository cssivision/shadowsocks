package shadowsocks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func newStream(block cipher.Block, iv []byte, isEncrypt bool) (cipher.Stream, error) {
	if isEncrypt {
		return cipher.NewCFBEncrypter(block, iv), nil
	} else {
		return cipher.NewCFBDecrypter(block, iv), nil
	}
}

func newAESCFBStream(key, iv []byte, isEncrypt bool) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return newStream(block, iv, isEncrypt)
}

type Cipher struct {
	enc  cipher.Stream
	dec  cipher.Stream
	key  []byte
	iv   []byte
	info *CipherInfo
}

type CipherInfo struct {
	keyLen    int
	ivLen     int
	newStream func([]byte, []byte, bool) (cipher.Stream, error)
}

var cipherMethods = map[string]*CipherInfo{
	"aes-128-cfb": &CipherInfo{16, 16, newAESCFBStream},
	"aes-256-cfb": &CipherInfo{32, 16, newAESCFBStream},
}

func NewCipher(method, password string) (*Cipher, error) {
	c := new(Cipher)
	info, ok := cipherMethods[method]
	if !ok {
		return nil, errors.New(fmt.Sprintf("unsupport encrypt method: %v", method))
	}
	c.info = info
	key := generateKey(password, info.keyLen)
	c.key = key

	return c, nil
}

func (c *Cipher) initDecrypt(iv []byte) error {
	var err error
	c.dec, err = c.info.newStream(c.key, iv, false)
	return err
}

func (c *Cipher) initEncrypt() error {
	if c.iv == nil {
		ivLen := c.info.ivLen
		iv := make([]byte, ivLen)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}
		c.iv = iv
	}
	var err error
	c.enc, err = c.info.newStream(c.key, c.iv, true)
	return err
}

func (c *Cipher) Encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

func (c *Cipher) Decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}

func (c *Cipher) Clone() *Cipher {
	nc := *c
	nc.dec = nil
	nc.enc = nil

	return &nc
}
