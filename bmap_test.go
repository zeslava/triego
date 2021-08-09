package triego

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_bmap(t *testing.T) {
	rand.Seed(time.Now().Unix())

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
		bm := bmap{}
		for i := 0; i < 255; i++ {
			assert.Equal(t, -1, bm.index(byte(i)))
		}
	})

	for j := 0; j < 10; j++ {
		t.Run(strconv.Itoa(j), func(t *testing.T) {
			bs := genbs()

			bm := bmap{}
			for i, b := range bs {
				i := i
				b := b
				bm.set(b)
				index := bm.index(b)
				assert.Equal(t, i, index, fmt.Sprintf("%d: %d", i, b))
				if i != index {
					index = bm.index(b)
				}
			}

			for _, b := range bs {
				assert.True(t, bm.has(b))
			}

			for _, b := range bs {
				bm.unset(b)
			}

			for i := 0; i < 255; i++ {
				assert.False(t, bm.has(byte(i)))
			}
		})
	}
}

func Benchmark_bmap_index(b *testing.B) {
	bm := bmap{}

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
