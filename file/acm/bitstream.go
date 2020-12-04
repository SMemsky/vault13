package acm

import (
	"io"
)

type bitstream struct {
	r      io.Reader
	cache  uint8
	cached uint8
}

func newBitstream(r io.Reader) *bitstream {
	return &bitstream{
		r: r,
	}
}

// Error from Read is ignored, this is a feature
func (b *bitstream) bit() uint8 {
	if b.cached == 0 {
		var tmp [1]uint8
		b.r.Read(tmp[:])
		b.cache = tmp[0]
		b.cached = 8
	}
	r := b.cache & 0x1
	b.cache >>= 1
	b.cached--
	return r
}

func (b *bitstream) bits(n int) uint64 {
	var r uint64
	for i := 0; i < n; i++ {
		r |= uint64(b.bit()) << i
	}
	return r
}
