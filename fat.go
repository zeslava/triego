package triego

type fat[ValueType any] struct {
	parent   *fat[ValueType]
	val      ValueType
	children [256]*fat[ValueType]
}

// NewFat creates Trie implementation, which more "fat"(uses more memory) than default
func NewFat[ValueType any, T *fat[ValueType]]() T {
	return newFatWithParent[ValueType]((T)(nil))
}

func newFatWithParent[ValueType any, T *fat[ValueType]](parent T) T {
	return &fat[ValueType]{
		parent:   parent,
		children: [256]*fat[ValueType]{},
	}
}

func (f *fat[ValueType]) Insert(key []byte, value ValueType) Trie[ValueType] {
	cur := f
	for _, k := range key {
		if cur.children[k] == nil {
			cur.children[k] = newFatWithParent(f)
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
	ch := f.get(key[:len(key)-1])
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
	v := f.get(key)
	return v, v != nil
}

func (f *fat[ValueType]) Sub(b byte) (Trie[ValueType], bool) {
	child := f.children[b]
	return child, child != nil
}

func (f *fat[ValueType]) Value() ValueType {
	return f.val
}

func (f *fat[ValueType]) get(key []byte) *fat[ValueType] {
	cur := f
	i := 0
begin:
	cur = cur.children[key[i]]
	if cur == nil {
		return nil
	}
	i++
	if i < len(key) {
		goto begin
	}

	return cur
}
