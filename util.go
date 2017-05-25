package shadowsocks

import (
	"crypto/md5"
	"io"
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

func CopyBuffer(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
