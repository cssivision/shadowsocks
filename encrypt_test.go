package shadowsocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCipher(t *testing.T) {
	t.Run("should return error", func(t *testing.T) {
		cipher, err := NewCipher("not support method", "password")
		assert.Nil(t, cipher)
		assert.NotNil(t, err)
	})

	t.Run("normal new cipher", func(t *testing.T) {
		cipher, err := NewCipher("aes-128-cfb", "password")
		assert.Nil(t, err)
		assert.NotNil(t, cipher)
		assert.Equal(t, 16, cipher.info.keyLen)
		assert.Equal(t, 16, cipher.info.ivLen)
		assert.Equal(t, generateKey("password", 16), cipher.key)
		assert.Nil(t, cipher.enc)
		assert.Nil(t, cipher.dec)
	})
}
