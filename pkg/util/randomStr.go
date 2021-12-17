package util

import (
	"math/rand"
	"time"
)

const letterBytes = "ABCDEF0123456789"
const (
	letterIdxBits = 4                    // 4 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src1 = rand.NewSource(time.Now().UnixNano())
var src2 = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandStringBytesMaskImprSrc returns a random hexadecimal string of length n.
func RandStringBytesMaskImprSrc1(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src1.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src1.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
