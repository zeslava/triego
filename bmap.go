package triego

import (
	"math/bits"
)

const (
	shift = 6
	mask  = 0x3f
)

type bmap struct {
	f [4]uint64 // fields
	c [4]int    // counts
}

func newBmap() *bmap {
	return &bmap{
		f: [4]uint64{},
		c: [4]int{},
	}
}

func (b *bmap) set(x byte) {
	i := x >> shift
	b.f[i] |= 1 << uint(x&mask)
	b.c[i] += 1
}

func (b *bmap) unset(x byte) {
	i := x >> shift
	b.f[i] &^= 1 << uint(x&mask)
	b.c[i] -= 1
}

func (b *bmap) has(x byte) bool {
	return b.f[x>>shift]&(1<<uint(x&mask)) != 0
}

func (b *bmap) index(x byte) int {
	i := x >> shift
	if b.f[i] == 0 {
		return -1
	}

	if x < 64 {
		return bits.OnesCount64(b.f[i]&ones[bits.Len64(1<<uint64(x&mask))]) - 1
	}
	if x < 128 {
		return bits.OnesCount64(b.f[i]&ones[bits.Len64(1<<uint64(x&mask))]) - 1 + b.c[0]
	}
	if x < 192 {
		return bits.OnesCount64(b.f[i]&ones[bits.Len64(1<<uint64(x&mask))]) - 1 + b.c[0] + b.c[1]
	}

	return bits.OnesCount64(b.f[i]&ones[bits.Len64(1<<uint64(x&mask))]) - 1 + b.c[0] + b.c[1] + b.c[2]
}

var ones = [65]uint64{
	0x0, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f, 0xff,
	0x1ff, 0x3ff, 0x7ff, 0xfff, 0x1fff, 0x3fff, 0x7fff, 0xffff, 0x1ffff,
	0x3ffff, 0x7ffff, 0xfffff, 0x1fffff, 0x3fffff, 0x7fffff, 0xffffff, 0x1ffffff, 0x3ffffff,
	0x7ffffff, 0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, 0xffffffff, 0x1ffffffff, 0x3ffffffff, 0x7ffffffff,
	0xfffffffff, 0x1fffffffff, 0x3fffffffff, 0x7fffffffff, 0xffffffffff, 0x1ffffffffff, 0x3ffffffffff, 0x7ffffffffff, 0xfffffffffff,
	0x1fffffffffff, 0x3fffffffffff, 0x7fffffffffff, 0xffffffffffff, 0x1ffffffffffff, 0x3ffffffffffff, 0x7ffffffffffff, 0xfffffffffffff, 0x1fffffffffffff,
	0x3fffffffffffff, 0x7fffffffffffff, 0xffffffffffffff, 0x1ffffffffffffff, 0x3ffffffffffffff, 0x7ffffffffffffff, 0xfffffffffffffff, 0x1fffffffffffffff, 0x3fffffffffffffff,
	0x7fffffffffffffff, 0xffffffffffffffff,
}
