package triego

// Trie is interface for radix tree
// Trie performs both tree and node roles, because node == sub-tree
type Trie interface {
	// Insert inserting value by key
	Insert(key []byte, value interface{}) Trie
	// Delete remove sub-tree by key
	Delete(key []byte) bool
	// Find finds a sub-tree by key
	Find(key []byte) (Trie, bool)
	// Sub returns sub-tree in next level if exists
	Sub(k byte) (Trie, bool)
	// Value returns a value stored in root of tree
	Value() interface{}
}

// trie implements Trie interface
type trie struct {
	key      byte
	val      interface{}
	bmap     bmap
	children []*trie
}

// New creates a default Trie
func New() *trie {
	return newTrieWithKey(0)
}

func newTrieWithKey(key byte) *trie {
	return &trie{key: key, children: []*trie{}}
}

func (t *trie) Insert(key []byte, value interface{}) Trie {
	cur := t
	for _, k := range key {
		cur = cur.insChild(k)
	}
	cur.val = value
	return cur
}

func (t *trie) Delete(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	ch := t.take(key[:len(key)-1])
	if ch == nil {
		return false
	}
	return ch.delChild(key[len(key)-1])
}

func (t *trie) Find(key []byte) (Trie, bool) {
	if len(key) == 0 {
		return nil, false
	}
	v := t.take(key)
	return v, v != nil
}

func (t *trie) Sub(key byte) (Trie, bool) {
	child := t.getChild(key)
	return child, child != nil
}

func (t *trie) Value() interface{} {
	return t.val
}

// take returns sub-tree by key
func (t *trie) take(key []byte) *trie {
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
func (t *trie) insChild(key byte) *trie {
	child := t.getChild(key)
	if child != nil {
		return child
	}

	child = newTrieWithKey(key)
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
func (t *trie) getChild(key byte) *trie {
	i := t.bmap.index(key)
	if i == -1 {
		return nil
	}
	return t.children[i]
}

// delChild removes child with the key and returns flag is removed or not
func (t *trie) delChild(key byte) bool {
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
