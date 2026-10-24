package triego

import (
	"math/bits"
)

// bmap is a 256-bit set of child keys, one bit per byte value.
// It is split into 4 words of 64 bits: a key x lives in word x>>shift
// at bit x&mask, e.g. 'a' (97) is word 1, bit 33.
const (
	shift = 6    // x>>6 == x/64: which of the 4 words holds x
	mask  = 0x3f // x&63 == x%64: which bit inside that word
)

type bmap struct {
	f [4]uint64 // fields: bit (x&mask) of f[x>>shift] is set if key x is present
	// p[i] is the number of set bits in f[0..i-1], so p[0] is always 0.
	// The word f[i] itself is not included, its bits are counted by popcount in index,
	// so the max is p[3] = 3*64 = 192 (f[0..2] full), which fits uint8.
	// The total count of 256 would not fit, but it is never stored.
	// uint8 instead of int shrinks bmap from 64 to 36 bytes, so more nodes fit in cache
	p [4]uint8
}

// set marks x as present. Must not be called for an already present x,
// otherwise prefix counts in p become wrong
func (b *bmap) set(x byte) {
	// 1<<(x&mask) is a word with only x's bit set; OR turns it on
	b.f[x>>shift] |= 1 << (x & mask)
	// x is now counted before every following word
	for i := x>>shift + 1; i < 4; i++ {
		b.p[i]++
	}
}

// unset marks x as absent. Must not be called for an absent x
func (b *bmap) unset(x byte) {
	// &^ is AND NOT: clears x's bit, keeps the others
	b.f[x>>shift] &^= 1 << (x & mask)
	for i := x>>shift + 1; i < 4; i++ {
		b.p[i]--
	}
}

func (b *bmap) has(x byte) bool {
	// AND with x's bit is non-zero only if that bit is set
	return b.f[x>>shift]&(1<<(x&mask)) != 0
}

// index returns 1-based position of x among present keys in ascending order,
// or 0 if x is absent. Children are stored in the same order, so the result
// is an index into children, where [0] is a nil sentinel
func (b *bmap) index(x byte) int {
	w := b.f[x>>shift]

	// shift x's bit down to position 0 and drop the rest: 1 if present, 0 if absent
	bit := int(w >> (x & mask) & 1)

	// 1<<(x&mask)-1 is a word with all bits below x's set (e.g. x&mask=3 -> 0b0111),
	// so the popcount counts present keys in x's word that are less than x.
	// Adding p gives the count of all present keys less than x
	rank := int(b.p[x>>shift]) + bits.OnesCount64(w&(1<<(x&mask)-1))

	// -bit is 0 for absent and all ones (-1) for present x,
	// so AND gives either 0 or rank+1 without a branch
	return (rank + 1) & -bit
}
