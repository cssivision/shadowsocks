package shadowsocks

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMd5sum(t *testing.T) {
	sum := "5d41402abc4b2a76b9719d911017c592"
	assert.Equal(t, fmt.Sprintf("%x", md5sum([]byte("hello"))), sum)
}

func TestGenerateKey(t *testing.T) {
	t.Run("key length 16", func(t *testing.T) {
		sum := "5d41402abc4b2a76b9719d911017c592"
		assert.Equal(t, fmt.Sprintf("%x", generateKey("hello", 16)), sum)
	})

	t.Run("key length 32", func(t *testing.T) {
		sum := "5d41402abc4b2a76b9719d911017c59228b46ed3c111e85102909b1cfb50ea0f"
		assert.Equal(t, fmt.Sprintf("%x", generateKey("hello", 32)), sum)
	})
}

func TestCopyBuffer(t *testing.T) {
	sourceDate := bytes.Repeat([]byte("Hello World!"), 1<<10)
	src := bytes.NewBuffer(nil)
	dst := bytes.NewBuffer(nil)
	src.Write(sourceDate)
	CopyBuffer(dst, src)
	assert.Equal(t, dst.Bytes(), sourceDate)
}

func TestWriteRandomData(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	require.Nil(t, WriteRandomData(buf))
}
