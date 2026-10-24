package triego

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"os"
	"sync"
	"testing"
)

// Benchmarks on a real dictionary (testdata/words.txt, ~104k words, ~238k trie nodes).
// Queries go in shuffled order, so the working set does not fit in L1/L2
// and cache misses are part of the measurement, unlike the small benchmarks in trie_test.go

var dict struct {
	once    sync.Once
	words   [][]byte
	queries [][]byte // words in shuffled order
	missing [][]byte // queries with the last byte changed, mostly absent keys with a long common path
}

func loadDict(b *testing.B) {
	b.Helper()
	dict.once.Do(func() {
		data, err := os.ReadFile("testdata/words.txt")
		if err != nil {
			panic(err)
		}
		for _, w := range bytes.Split(data, []byte("\n")) {
			if len(w) > 0 {
				dict.words = append(dict.words, w)
			}
		}

		rnd := rand.New(rand.NewSource(1))
		dict.queries = append([][]byte(nil), dict.words...)
		rnd.Shuffle(len(dict.queries), func(i, j int) {
			dict.queries[i], dict.queries[j] = dict.queries[j], dict.queries[i]
		})

		for _, q := range dict.queries {
			m := append([]byte(nil), q...)
			m[len(m)-1] = '~'
			dict.missing = append(dict.missing, m)
		}
	})
}

type dictImpl struct {
	name  string
	build func([][]byte) func([]byte) bool // returns lookup
}

var dictImpls = []dictImpl{
	{"map", func(words [][]byte) func([]byte) bool {
		m := make(map[string]string)
		for _, w := range words {
			m[string(w)] = string(w)
		}
		return func(k []byte) bool { _, ok := m[string(k)]; return ok }
	}},
	{"trie", func(words [][]byte) func([]byte) bool {
		return buildTrie(New[string](), words)
	}},
	{"fat", func(words [][]byte) func([]byte) bool {
		return buildTrie(NewFat[string](), words)
	}},
}

func buildTrie(tr Trie[string], words [][]byte) func([]byte) bool {
	for _, w := range words {
		tr.Insert(w, string(w))
	}
	return func(k []byte) bool { _, ok := tr.Find(k); return ok }
}

func Benchmark_dictFind(b *testing.B) {
	loadDict(b)
	for _, impl := range dictImpls {
		find := impl.build(dict.words)
		b.Run(impl.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if !find(dict.queries[i%len(dict.queries)]) {
					b.Fatal("word not found")
				}
			}
		})
	}
}

func Benchmark_dictFindMissing(b *testing.B) {
	loadDict(b)
	for _, impl := range dictImpls {
		find := impl.build(dict.words)
		b.Run(impl.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				find(dict.missing[i%len(dict.missing)])
			}
		})
	}
}

// File signatures from github.com/zeslava/sigo (testdata/sigs.txt, hex, one per line).
// sigo detects a file type from a stream, walking the tree byte by byte with Sub,
// so this is the access pattern to measure. The tree is small (~2.4k nodes) and stays in cache
var sigs struct {
	once   sync.Once
	keys   [][]byte
	inputs [][]byte // shuffled file headers: a signature followed by random bytes
}

func loadSigs(b *testing.B) {
	b.Helper()
	sigs.once.Do(func() {
		data, err := os.ReadFile("testdata/sigs.txt")
		if err != nil {
			panic(err)
		}
		rnd := rand.New(rand.NewSource(1))
		for _, line := range bytes.Fields(data) {
			k, err := hex.DecodeString(string(line))
			if err != nil {
				panic(err)
			}
			sigs.keys = append(sigs.keys, k)

			tail := make([]byte, 16)
			rnd.Read(tail)
			sigs.inputs = append(sigs.inputs, append(append([]byte(nil), k...), tail...))
		}
		rnd.Shuffle(len(sigs.inputs), func(i, j int) {
			sigs.inputs[i], sigs.inputs[j] = sigs.inputs[j], sigs.inputs[i]
		})
	})
}

// Benchmark_sigsSub detects one header per op the way sigo does: Sub until there is no child,
// remembering the last node with a value. ns/step is the cost of a single Sub
func Benchmark_sigsSub(b *testing.B) {
	loadSigs(b)
	impls := []struct {
		name string
		tr   Trie[string]
	}{
		{"trie", New[string]()},
		{"fat", NewFat[string]()},
	}
	for _, impl := range impls {
		for _, k := range sigs.keys {
			impl.tr.Insert(k, string(k))
		}
		tr := impl.tr
		b.Run(impl.name, func(b *testing.B) {
			steps := 0
			for i := 0; i < b.N; i++ {
				cur, last := tr, ""
				for _, c := range sigs.inputs[i%len(sigs.inputs)] {
					next, ok := cur.Sub(c)
					if !ok {
						break
					}
					steps++
					cur = next
					if v := cur.Value(); v != "" {
						last = v
					}
				}
				if last == "" {
					b.Fatal("signature not found")
				}
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(steps), "ns/step")
		})
	}

	tr := New[string]()
	for _, k := range sigs.keys {
		tr.Insert(k, string(k))
	}
	b.Run("cursor", func(b *testing.B) {
		steps := 0
		for i := 0; i < b.N; i++ {
			c, last := tr.Cursor(), ""
			for _, k := range sigs.inputs[i%len(sigs.inputs)] {
				if !c.Step(k) {
					break
				}
				steps++
				if v := c.Value(); v != "" {
					last = v
				}
			}
			if last == "" {
				b.Fatal("signature not found")
			}
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(steps), "ns/step")
	})
}

// Benchmark_dictInsert builds the whole dictionary per op, B/op is the memory footprint of the structure
func Benchmark_dictInsert(b *testing.B) {
	loadDict(b)
	for _, impl := range dictImpls {
		b.Run(impl.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				impl.build(dict.words)
			}
		})
	}
}
