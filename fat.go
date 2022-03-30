package triego

type fat[ValueType any] struct {
	key      byte
	val      ValueType
	children *[256]*fat[ValueType]
}

// NewFat creates Trie implementation, which more "fat"(uses more memory) than default
func NewFat[ValueType any]() *fat[ValueType] {
	return newFatWithKey[ValueType](0)
}

func newFatWithKey[ValueType any](key byte) *fat[ValueType] {
	return &fat[ValueType]{key: key, children: &[256]*fat[ValueType]{}}
}

func (f *fat[ValueType]) Insert(key []byte, value ValueType) Trie[ValueType] {
	cur := f
	for _, k := range key {
		if cur.children[k] == nil {
			cur.children[k] = newFatWithKey[ValueType](k)
		}
		cur = cur.children[k]
	}
	cur.val = value
	return cur
}

func (f *fat[ValueType]) Delete(key []byte) bool {
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

func (f *fat[ValueType]) Find(key []byte) (Trie[ValueType], bool) {
	if len(key) == 0 {
		return nil, false
	}
	v := f.take(key)
	return v, v != nil
}

func (f *fat[ValueType]) Sub(b byte) (Trie[ValueType], bool) {
	child := f.children[b]
	return child, child != nil
}

func (f *fat[ValueType]) Value() ValueType {
	return f.val
}

func (f *fat[ValueType]) take(key []byte) *fat[ValueType] {
	cur := f
	for i := 0; i < len(key); i++ {
		cur = cur.children[key[i]]
		if cur == nil {
			return nil
		}
	}

	return cur
}
