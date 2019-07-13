package triego

import "sort"

// Trie is interface for radix tree
// Trie performs both tree and node roles, because node == sub-tree
type Trie interface {
	Insert([]byte, interface{}) Trie
	Delete([]byte) bool
	Find([]byte) (Trie, bool)
	Sub(byte) (Trie, bool)
	Value() interface{}
}

// trie implements Trie interface
type trie struct {
	key byte
	val interface{}
	sub []*trie
}

// New creates a Trie
func New() Trie {
	return &trie{}
}

// Insert inserting value by path(key)
func (t *trie) Insert(path []byte, value interface{}) Trie {
	cur := t
	for _, k := range path {
		cur = cur.insChild(k)
	}
	cur.val = value
	return cur
}

// Delete remove sub-tree by path(key)
func (t *trie) Delete(path []byte) bool {
	ch := t.take(path[:len(path)-1])
	if ch == nil {
		return false
	}
	return ch.delChild(path[len(path)-1])
}

// Find finds a sub-tree by path(key)
func (t *trie) Find(path []byte) (Trie, bool) {
	v := t.take(path)
	return v, v != nil
}

// Sub returns sub-tree in next level if exists
func (t *trie) Sub(key byte) (Trie, bool) {
	v := t.getChild(key)
	return v, v != nil
}

// Value returns a value stored in root of tree
func (t *trie) Value() interface{} {
	return t.val
}

// take returns sub-tree by path(key)
func (t *trie) take(path []byte) *trie {
	cur := t
	for _, k := range path {
		if child := cur.getChild(k); child != nil {
			cur = child
		} else {
			return cur
		}
	}

	return cur
}

// insChild inserts new child with the key, or return child if already exists
func (t *trie) insChild(key byte) *trie {
	if t.sub == nil {
		v := &trie{key: key}
		t.sub = []*trie{v}
		return v
	}

	v := t.getChild(key)
	if v != nil {
		return v
	}

	v = &trie{key: key}
	i := sort.Search(len(t.sub), func(i int) bool { return (t.sub)[i].key >= key })
	if i == len(t.sub) {
		t.sub = append(t.sub, v)
	} else {
		t.sub = append(t.sub, nil)
		copy(t.sub[i+1:], t.sub[i:])
		(t.sub)[i] = v
	}
	return v
}

// getChild returns child with the key
func (t *trie) getChild(key byte) *trie {
	i := sort.Search(len(t.sub), func(i int) bool { return (t.sub)[i].key >= key })
	if i < len(t.sub) && (t.sub)[i].key == key {
		return t.sub[i]
	}
	return nil
}

// delChild removes child with the key and returns flag is removed or not
func (t *trie) delChild(key byte) bool {
	num := len(t.sub)
	i := sort.Search(num, func(i int) bool { return t.sub[i].key >= key })
	if i < len(t.sub) && (t.sub)[i].key == key {
		copy(t.sub[i:], t.sub[i+1:])
		t.sub[len(t.sub)-1] = nil
		t.sub = t.sub[:len(t.sub)-1]
		return true
	}
	return false
}
