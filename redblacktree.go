package redblacktree

import (
	"sync"

	"github.com/koketama/redblacktree/internal/pkg"

	"github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"
)

// Value if the value.id is duplicated, it will be ignored.
type Value = pkg.Value

var _ Tree = (*tree)(nil)

// Tree supported methods on red-black-tree
type Tree interface {
	Put(key interface{}, value Value)
	Get(key interface{}) (values []Value, found bool)
	Remove(key interface{})
	Empty() bool
	Size() int
	Min() (key interface{}, values []Value)
	PopMin() (key interface{}, values []Value)
	Max() (key interface{}, values []Value)
	PopMax() (key interface{}, values []Value)
	Topology() string
	Iterator() Iterator
}

// Iterator a stateful iterator whose elements are key/value pairs.
type Iterator interface {
	Next() bool
	Key() interface{}
	Values() []Value
}

type tree struct {
	sync.RWMutex
	rbt *pkg.Tree
}

// New create a thread safe red-black-tree based on github.com/emirpasic/gods/trees/redblacktree
func New(comparator utils.Comparator) (Tree, error) {
	if comparator == nil {
		return nil, errors.New("comparator required")
	}

	return &tree{rbt: pkg.NewWith(comparator)}, nil
}

func (t *tree) Put(key interface{}, value Value) {
	if key == nil || value == nil {
		return
	}

	t.Lock()
	defer t.Unlock()

	t.rbt.Put(key, value)
}

func (t *tree) Get(key interface{}) (values []Value, found bool) {
	if key == nil {
		return
	}

	t.RLock()
	defer t.RUnlock()

	return t.rbt.Get(key)
}

func (t *tree) Remove(key interface{}) {
	if key == nil {
		return
	}

	t.Lock()
	defer t.Unlock()

	t.rbt.Remove(key)
}

func (t *tree) Empty() bool {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.Empty()
}

func (t *tree) Size() int {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.Size()
}

func (t *tree) Min() (key interface{}, values []Value) {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.Left()
}

func (t *tree) PopMin() (key interface{}, values []Value) {
	t.Lock()
	defer t.Unlock()

	return t.rbt.PopLeft()
}

func (t *tree) Max() (key interface{}, values []Value) {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.Right()
}

func (t *tree) PopMax() (key interface{}, values []Value) {
	t.Lock()
	defer t.Unlock()

	return t.rbt.PopRight()
}

func (t *tree) Topology() string {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.String()
}

func (t *tree) Iterator() Iterator {
	iterator := t.rbt.Iterator()
	return &iterator
}
