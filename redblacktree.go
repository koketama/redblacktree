package redblacktree

import (
	"sync"

	"github.com/koketama/redblacktree/internal/pkg"

	"github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"
)

type Value = pkg.Value

var _ Tree = (*tree)(nil)

type Tree interface {
	Put(key interface{}, value Value)
	Get(key interface{}) (values []Value, found bool)
	Remove(key interface{})
	Empty() bool
	Size() int
	Min() (key interface{}, values []Value)
	Max() (key interface{}, values []Value)
	String() string
}

type tree struct {
	sync.RWMutex
	rbt *pkg.Tree
}

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
		return nil, false
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

func (t *tree) Max() (key interface{}, values []Value) {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.Right()
}

func (t *tree) String() string {
	t.RLock()
	defer t.RUnlock()

	return t.rbt.String()
}
