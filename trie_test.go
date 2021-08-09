package triego

import (
	"testing"
)

var testTrie = New()
var testFatTrie = NewFat()

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

func tr(b []byte) string {
	t, ok := testTrie.Find(b)
	if !ok {
		return ""
	}
	v, _ := t.Value().(string)
	return v
}

func trfat(b []byte) string {
	t, ok := testFatTrie.Find(b)
	if !ok {
		return ""
	}
	v, _ := t.Value().(string)
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
