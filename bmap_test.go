package triego

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func Test_bmap(t *testing.T) {
	rand.New(rand.NewSource(time.Now().Unix()))

	genbs := func() []byte {
		n := rand.Intn(150)
		x := rand.Intn(97)
		bs := make([]byte, n)
		for i := 0; i < n; i++ {
			bs[i] = byte(x + i)
		}
		return bs
	}

	t.Run("empty", func(t *testing.T) {
		bm := newBmap()
		for i := 0; i < 255; i++ {
			if bm.index(byte(i)) > -1 {
				t.Fail()
			}
		}
	})

	for j := 0; j < 10; j++ {
		j := j
		t.Run(strconv.Itoa(j), func(t *testing.T) {
			bs := genbs()

			bm := newBmap()
			for i, b := range bs {
				i := i
				b := b
				bm.set(b)
				index := bm.index(b)
				if index != i {
					t.Errorf("%d: %d", i, b)
				}
			}

			for _, b := range bs {
				if !bm.has(b) {
					t.Fail()
				}
			}

			for _, b := range bs {
				bm.unset(b)
			}

			for i := 0; i < 255; i++ {
				if bm.has(byte(i)) {
					t.Fail()
				}
			}
		})
	}
}

func Benchmark_bmap_index(b *testing.B) {
	bm := newBmap()

	for index, v := range []byte{10, 70, 150, 200, 250} {
		bm.set(v)
		index := index
		v := v
		b.Run(strconv.Itoa(int(v)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if j := bm.index(v); j != index {
					b.Fatal("wrong index", j)
				}
			}
		})
	}
}
