package etag

import "unsafe"

// b2s converts byte slice to a string without memory allocation
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// appendUint appends n to dst and returns the extended dst
func appendUint(dst []byte, n uint32) []byte {
	var b [20]byte
	buf := b[:]
	i := len(buf)
	var q uint32
	for n >= 10 {
		i--
		q = n / 10
		buf[i] = '0' + byte(n-q*10)
		n = q
	}
	i--
	buf[i] = '0' + byte(n)
	dst = append(dst, buf[i:]...)
	return dst
}
