package triego

import "testing"

func Test_fattrie(t *testing.T) {
	for i := range byteWords {
		i := i
		t.Run(words[i], func(t *testing.T) {
			if s := trfat(byteWords[i]); s != words[i] {
				t.Fatalf("bad word: got %s, want %s", s, words[i])
			}
		})
	}
}

func Benchmark_fattrie(b *testing.B) {
	for i := range byteWords {
		i := i
		b.Run(words[i], func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				if s := trfat(byteWords[i]); s != words[i] {
					b.Fatalf("bad word %#v", s)
				}
			}
		})
	}
}
