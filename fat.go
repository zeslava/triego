package triego

type fat struct {
	key      byte
	val      interface{}
	children *[256]*fat
}

// NewFat creates Trie implementation, which more "fat"(uses more memory) than default
func NewFat() *fat {
	return newFatWithKey(0)
}

func newFatWithKey(key byte) *fat {
	return &fat{key: key, children: &[256]*fat{}}
}

func (f *fat) Insert(key []byte, value interface{}) Trie {
	cur := f
	for _, k := range key {
		if cur.children[k] == nil {
			cur.children[k] = newFatWithKey(k)
		}
		cur = cur.children[k]
	}
	cur.val = value
	return cur
}

func (f *fat) Delete(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	ch := f.take(key[:len(key)-1])
	if ch == nil {
		return false
	}
	k := key[len(key)-1]
	has := f.children[k] != nil
	f.children[k] = nil
	return has
}

func (f *fat) Find(key []byte) (Trie, bool) {
	if len(key) == 0 {
		return nil, false
	}
	v := f.take(key)
	return v, v != nil
}

func (f *fat) Sub(b byte) (Trie, bool) {
	child := f.children[b]
	return child, child != nil
}

func (f *fat) Value() interface{} {
	return f.val
}

func (f *fat) take(key []byte) *fat {
	cur := f
	for i := 0; i < len(key); i++ {
		cur = cur.children[key[i]]
		if cur == nil {
			return nil
		}
	}

	return cur
}
