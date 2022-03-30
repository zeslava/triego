package triego

// Trie is interface for prefix tree
// Trie performs both tree and node roles, because node == sub-tree
type Trie[ValueType any] interface {
	// Insert inserting value by key
	Insert(key []byte, value ValueType) Trie[ValueType]
	// Delete remove subtree by key
	Delete(key []byte) bool
	// Find finds a subtree by key
	Find(key []byte) (Trie[ValueType], bool)
	// Sub returns subtree in next level if exists
	Sub(k byte) (Trie[ValueType], bool)
	// Value returns a value stored in root of tree
	Value() ValueType
}

// trie implements Trie interface
type trie[ValueType any] struct {
	key      byte
	val      ValueType
	bmap     bmap
	children []*trie[ValueType]
}

// New creates a default Trie
func New[ValueType any]() *trie[ValueType] {
	return newTrieWithKey[ValueType](0)
}

func newTrieWithKey[ValueType any](key byte) *trie[ValueType] {
	return &trie[ValueType]{key: key, children: make([]*trie[ValueType], 0)}
}

func (t *trie[ValueType]) Insert(key []byte, value ValueType) Trie[ValueType] {
	cur := t
	for _, k := range key {
		cur = cur.insChild(k)
	}
	cur.val = value
	return cur
}

func (t *trie[ValueType]) Delete(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	ch := t.take(key[:len(key)-1])
	if ch == nil {
		return false
	}
	return ch.delChild(key[len(key)-1])
}

func (t *trie[ValueType]) Find(key []byte) (Trie[ValueType], bool) {
	if len(key) == 0 {
		return nil, false
	}
	v := t.take(key)
	return v, v != nil
}

func (t *trie[ValueType]) Sub(key byte) (Trie[ValueType], bool) {
	child := t.getChild(key)
	return child, child != nil
}

func (t *trie[ValueType]) Value() ValueType {
	return t.val
}

// take returns subtree by key
func (t *trie[ValueType]) take(key []byte) *trie[ValueType] {
	cur := t
	for i := 0; i < len(key); i++ {
		cur = cur.getChild(key[i])
		if cur == nil {
			return nil
		}
	}
	return cur
}

// insChild inserts new child with the key, or return child if already exists
func (t *trie[ValueType]) insChild(key byte) *trie[ValueType] {
	child := t.getChild(key)
	if child != nil {
		return child
	}

	child = newTrieWithKey[ValueType](key)
	t.bmap.set(key)
	i := t.bmap.index(key)
	if i == len(t.children) {
		t.children = append(t.children, child)
	} else {
		t.children = append(t.children, nil)
		copy(t.children[i+1:], t.children[i:])
		t.children[i] = child
	}
	return child
}

// getChild returns child with the key
func (t *trie[ValueType]) getChild(key byte) *trie[ValueType] {
	i := t.bmap.index(key)
	if i == -1 {
		return nil
	}
	return t.children[i]
}

// delChild removes child with the key and returns flag is removed or not
func (t *trie[ValueType]) delChild(key byte) bool {
	i := t.bmap.index(key)
	if i == -1 {
		return false
	}
	copy(t.children[i:], t.children[i+1:])
	t.children[len(t.children)-1] = nil
	t.children = t.children[:len(t.children)-1]
	t.bmap.unset(key)
	return true
}
