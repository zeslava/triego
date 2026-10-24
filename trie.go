package triego

import "math/bits"

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
	val      ValueType
	bmap     bmap
	children []*trie[ValueType]
}

// New creates a default Trie
func New[ValueType any, T *trie[ValueType]]() T {
	return newTrie[ValueType]()
}

func newTrie[ValueType any]() *trie[ValueType] {
	return &trie[ValueType]{
		children: make([]*trie[ValueType], 1), // [0] is a nil sentinel for absent keys, see bmap.index
	}
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

// Cursor walks the trie byte by byte, e.g. over a stream read in chunks.
// Unlike Sub it is a concrete type, so Step is inlined and nothing is boxed into an interface
type Cursor[ValueType any] struct {
	n *trie[ValueType]
}

// Cursor returns a cursor at the root of t
func (t *trie[ValueType]) Cursor() Cursor[ValueType] {
	return Cursor[ValueType]{t}
}

// Step moves the cursor to the child by key. If there is no such child,
// it returns false and the cursor stays in place
func (c *Cursor[ValueType]) Step(key byte) bool {
	// getChild is written out by hand, through it Step costs 73 instead of 60 and is slower
	b := &c.n.bmap
	w := b.f[key>>shift]
	if w>>(key&mask)&1 == 0 {
		return false
	}
	c.n = c.n.children[int(b.p[key>>shift])+bits.OnesCount64(w&(1<<(key&mask)-1))+1]
	return true
}

// Value returns the value at the cursor, the zero value if no key ends here
func (c *Cursor[ValueType]) Value() ValueType {
	return c.n.val
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

	child = newTrie[ValueType]()
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
	// bmap.index is written out: an early return for absent keys is cheaper
	// than its branchless 0, which is not needed with the branch here
	w := t.bmap.f[key>>shift]
	if w>>(key&mask)&1 == 0 {
		return nil
	}
	return t.children[int(t.bmap.p[key>>shift])+bits.OnesCount64(w&(1<<(key&mask)-1))+1]
}

// delChild removes child with the key and returns flag is removed or not
func (t *trie[ValueType]) delChild(key byte) bool {
	i := t.bmap.index(key)
	if i == 0 {
		return false
	}
	copy(t.children[i:], t.children[i+1:])
	t.children[len(t.children)-1] = nil
	t.children = t.children[:len(t.children)-1]
	t.bmap.unset(key)
	return true
}
