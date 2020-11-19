package pkg

import (
	"fmt"
	"testing"

	"github.com/emirpasic/gods/utils"
)

type Entity string

func (e Entity) ID() string {
	return string(e)
}

func TestRBT(t *testing.T) {
	tree := NewWith(utils.IntComparator)

	tree.Put(1, Entity("A"))
	tree.Put(1, Entity("B"))
	tree.Put(1, Entity("C"))

	fmt.Printf("size:%d\n%s", tree.Size(), tree.String())
	fmt.Println(tree.Get(1))
	fmt.Println("----------------------")

	tree.Put(2, Entity("D"))
	tree.Put(2, Entity("E"))

	fmt.Printf("size:%d\n%s", tree.Size(), tree.String())
	fmt.Println(tree.Get(2))
	fmt.Println("----------------------")

	tree.Put(3, Entity("F"))
	tree.Put(3, Entity("G"))
	tree.Put(3, Entity("H"))
	tree.Put(3, Entity("I"))

	fmt.Printf("size:%d\n%s", tree.Size(), tree.String())
	fmt.Println(tree.Get(3))
	fmt.Println("----------------------")

	tree.Put(4, Entity("J"))

	fmt.Printf("size:%d\n%s", tree.Size(), tree.String())
	fmt.Println(tree.Get(4))
	fmt.Println("----------------------")

	for !tree.Empty() {
		key, values := tree.Left()
		fmt.Println(tree.Size(), key, values)
		tree.Remove(key)
	}
	fmt.Println("----------------------")
}
