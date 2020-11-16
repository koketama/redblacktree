package redblacktree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/emirpasic/gods/utils"
)

type value struct {
	v int
}

func (v *value) ID() string {
	return fmt.Sprintf("%d", v)
}

func (v *value) String() string {
	return fmt.Sprintf("%d", v)
}

func TestRedBlackTree(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for k := 0; k < 10; k++ {
		seeds := random.Perm(1000000)

		tree, _ := NewWith(utils.IntComparator)
		for _, seed := range seeds {
			tree.Put(seed, &value{v: seed})
		}

		t.Log(tree.Size(), tree.Left(), tree.Right())
	}
}
