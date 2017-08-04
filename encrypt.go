package shadowsocks

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	"io"
)

// Cipher struct
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
	"rc4-md5":     &CipherInfo{16, 16, newRC4MD5Stream},
}

func newRC4MD5Stream(key, iv []byte, _ bool) (cipher.Stream, error) {
	rc4key := md5sum(append(key, iv...))
	return rc4.NewCipher(rc4key)
}

func newAESCFBStream(key, iv []byte, isEncrypt bool) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if isEncrypt {
		return cipher.NewCFBEncrypter(block, iv), nil
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

// NewCipher ...
func NewCipher(method, password string) (*Cipher, error) {
	c := new(Cipher)
	info, ok := cipherMethods[method]
	if !ok {
		return nil, fmt.Errorf("unsupport encrypt method: %v", method)
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

// Encrypt encrypt data
func (c *Cipher) Encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

// Decrypt decrypt data
func (c *Cipher) Decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}

// Clone cipher clone
func (c *Cipher) Clone() *Cipher {
	nc := *c
	nc.dec = nil
	nc.enc = nil

	return &nc
}
