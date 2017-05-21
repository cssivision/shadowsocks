package shadowsocks

import (
	"crypto/md5"
	"math"
)

func md5sum(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func generateKey(password string, keyLen int) []byte {
	md5Len := 16
	count := int(math.Ceil(float64(keyLen) / float64(md5Len)))
	r := make([]byte, count*md5Len)
	copy(r, md5sum([]byte(password)))

	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < count; i++ {
		start += md5Len
		copy(d[:md5Len], r[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(r[start:start+md5Len], md5sum(d))
	}

	return r[:keyLen]
}
