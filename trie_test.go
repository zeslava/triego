package triego

import (
	"math/rand"
	"testing"
)

var (
	testTrie    = New[string]()
	testFatTrie = NewFat[string]()
)

func init() {
	for _, s := range words {
		s := s
		testTrie.Insert([]byte(s), s)
		testFatTrie.Insert([]byte(s), s)
	}
}

func Test_trie(t *testing.T) {
	for i := range byteWords {
		i := i
		t.Run(words[i], func(t *testing.T) {
			if s := tr(byteWords[i]); s != words[i] {
				t.Fatalf("bad word: got %s, want %s", s, words[i])
			}
		})
	}
}

func Benchmark_map(b *testing.B) {
	for _, word := range words {
		word := word
		b.Run(word, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if s := mapWords[word]; s != word {
					b.Fatalf("bad word %#v", s)
				}
			}
		})
	}
}

func Benchmark_trie(b *testing.B) {
	for i := range byteWords {
		i := i
		b.Run(words[i], func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				if s := tr(byteWords[i]); s != words[i] {
					b.Fatalf("bad word %#v", s)
				}
			}
		})
	}
}

func tr(b []byte) string {
	t, ok := testTrie.Find(b)
	if !ok {
		return ""
	}
	v := t.Value()
	return v
}

var impls = map[string]func() Trie[string]{
	"trie": func() Trie[string] { return New[string]() },
	"fat":  func() Trie[string] { return NewFat[string]() },
}

func Test_findMissing(t *testing.T) {
	for name, newTrie := range impls {
		t.Run(name, func(t *testing.T) {
			tr := newTrie()
			tr.Insert([]byte("a"), "a")
			tr.Insert([]byte("c"), "c")
			tr.Insert([]byte{200}, "200")
			for _, k := range [][]byte{[]byte("b"), []byte("d"), []byte("ab"), {199}, {201}, {0}, {255}} {
				if n, ok := tr.Find(k); ok {
					t.Errorf("Find(%q): got %q, want not found", k, n.Value())
				}
			}
		})
	}
}

func Test_delete(t *testing.T) {
	for name, newTrie := range impls {
		t.Run(name, func(t *testing.T) {
			tr := newTrie()
			tr.Insert([]byte("abc"), "abc")
			tr.Insert([]byte("abd"), "abd")
			if !tr.Delete([]byte("abc")) {
				t.Fatal("Delete(abc): got false, want true")
			}
			if _, ok := tr.Find([]byte("abc")); ok {
				t.Error("Find(abc) after delete: got found")
			}
			if n, ok := tr.Find([]byte("abd")); !ok || n.Value() != "abd" {
				t.Error("Find(abd) after deleting sibling: got not found")
			}
			if tr.Delete([]byte("abc")) {
				t.Error("second Delete(abc): got true, want false")
			}
		})
	}
}

func Test_cursor(t *testing.T) {
	tr := New[string]()
	tr.Insert([]byte("ab"), "ab")
	tr.Insert([]byte("abcd"), "abcd")

	c := tr.Cursor()
	for _, step := range []struct {
		key   byte
		moved bool
		val   string
	}{
		{'a', true, ""},
		{'b', true, "ab"},
		{'x', false, "ab"}, // no child, the cursor stays at "ab"
		{'c', true, ""},
		{'d', true, "abcd"},
		{'e', false, "abcd"},
	} {
		if moved := c.Step(step.key); moved != step.moved {
			t.Fatalf("Step(%q) = %v", step.key, moved)
		}
		if val := c.Value(); val != step.val {
			t.Fatalf("after %q: Value() = %q", step.key, val)
		}
	}
}

// Test_random checks both implementations against a map oracle on keys using the full byte range
func Test_random(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	genKey := func() []byte {
		k := make([]byte, 1+rnd.Intn(8))
		rnd.Read(k)
		return k
	}

	for name, newTrie := range impls {
		t.Run(name, func(t *testing.T) {
			tr := newTrie()
			paths := map[string]bool{}
			values := map[string]string{}
			for i := 0; i < 2000; i++ {
				k := genKey()
				tr.Insert(k, string(k))
				values[string(k)] = string(k)
				for j := 1; j <= len(k); j++ {
					paths[string(k[:j])] = true
				}
			}

			for k, v := range values {
				if n, ok := tr.Find([]byte(k)); !ok || n.Value() != v {
					t.Fatalf("Find(%q): inserted key not found", k)
				}
			}
			for i := 0; i < 10000; i++ {
				k := genKey()
				if _, ok := tr.Find(k); ok != paths[string(k)] {
					t.Fatalf("Find(%q): got %v, want %v", k, ok, paths[string(k)])
				}
			}
		})
	}
}

func Benchmark_findMissing(b *testing.B) {
	keys := [][]byte{[]byte("abatiz"), []byte("aaz"), []byte("abandonwarez"), []byte("zzz")}
	for name, tr := range map[string]Trie[string]{"trie": testTrie, "fat": testFatTrie} {
		for _, k := range keys {
			b.Run(name+"/"+string(k), func(b *testing.B) {
				for j := 0; j < b.N; j++ {
					tr.Find(k)
				}
			})
		}
	}
}

func trfat(b []byte) string {
	t, ok := testFatTrie.Find(b)
	if !ok {
		return ""
	}
	v := t.Value()
	return v
}

var words = []string{
	"abatis",
	"abatises",
	"abator",
	"abandonwares",
	"abands",
	"abapical",
	"aasvogels",
	"ab",
	"aal",
	"aalii",
	"",
}

var byteWords = [][]byte{
	[]byte("abatis"),
	[]byte("abatises"),
	[]byte("abator"),
	[]byte("abandonwares"),
	[]byte("abands"),
	[]byte("abapical"),
	[]byte("aasvogels"),
	[]byte("ab"),
	[]byte("aal"),
	[]byte("aalii"),
	[]byte(""),
}

var mapWords = map[string]string{
	"aa":           "aa",
	"aah":          "aah",
	"aahed":        "aahed",
	"aahing":       "aahing",
	"aahs":         "aahs",
	"aal":          "aal",
	"aalii":        "aalii",
	"aaliis":       "aaliis",
	"aals":         "aals",
	"aardvark":     "aardvark",
	"aardvarks":    "aardvarks",
	"aardwolf":     "aardwolf",
	"aardwolves":   "aardwolves",
	"aargh":        "aargh",
	"aarrgh":       "aarrgh",
	"aarrghh":      "aarrghh",
	"aarti":        "aarti",
	"aartis":       "aartis",
	"aas":          "aas",
	"aasvogel":     "aasvogel",
	"aasvogels":    "aasvogels",
	"ab":           "ab",
	"aba":          "aba",
	"abac":         "abac",
	"abaca":        "abaca",
	"abacas":       "abacas",
	"abaci":        "abaci",
	"aback":        "aback",
	"abacs":        "abacs",
	"abacterial":   "abacterial",
	"abactinal":    "abactinal",
	"abactinally":  "abactinally",
	"abactor":      "abactor",
	"abactors":     "abactors",
	"abacus":       "abacus",
	"abacuses":     "abacuses",
	"abaft":        "abaft",
	"abaka":        "abaka",
	"abakas":       "abakas",
	"abalone":      "abalone",
	"abalones":     "abalones",
	"abamp":        "abamp",
	"abampere":     "abampere",
	"abamperes":    "abamperes",
	"abamps":       "abamps",
	"aband":        "aband",
	"abanded":      "abanded",
	"abanding":     "abanding",
	"abandon":      "abandon",
	"abandoned":    "abandoned",
	"abandonedly":  "abandonedly",
	"abandonee":    "abandonee",
	"abandonees":   "abandonees",
	"abandoner":    "abandoner",
	"abandoners":   "abandoners",
	"abandoning":   "abandoning",
	"abandonment":  "abandonment",
	"abandonments": "abandonments",
	"abandons":     "abandons",
	"abandonware":  "abandonware",
	"abandonwares": "abandonwares",
	"abands":       "abands",
	"abapical":     "abapical",
	"abas":         "abas",
	"abase":        "abase",
	"abased":       "abased",
	"abasedly":     "abasedly",
	"abasement":    "abasement",
	"abasements":   "abasements",
	"abaser":       "abaser",
	"abasers":      "abasers",
	"abases":       "abases",
	"abash":        "abash",
	"abashed":      "abashed",
	"abashedly":    "abashedly",
	"abashes":      "abashes",
	"abashing":     "abashing",
	"abashless":    "abashless",
	"abashment":    "abashment",
	"abashments":   "abashments",
	"abasia":       "abasia",
	"abasias":      "abasias",
	"abasing":      "abasing",
	"abask":        "abask",
	"abatable":     "abatable",
	"abate":        "abate",
	"abated":       "abated",
	"abatement":    "abatement",
	"abatements":   "abatements",
	"abater":       "abater",
	"abaters":      "abaters",
	"abates":       "abates",
	"abating":      "abating",
	"abatis":       "abatis",
	"abatises":     "abatises",
	"abator":       "abator",
	"abators":      "abators",
	"abattis":      "abattis",
	"abattises":    "abattises",
	"abattoir":     "abattoir",
}
