package radix_tree

import (
	"fmt"
)

type RadixTree struct {
	root *node
}

type node interface {
	isTerminal() bool
}

type inode struct {
	lc       *node
	rc       *node
	critbyte uint32
	critmask uint8
}

type tnode struct {
	key string
	val string
}

func (n tnode) isTerminal() bool { return true }
func (n inode) isTerminal() bool { return false }

func (t RadixTree) Get(key string) (val string, err error) {
	return "", nil
}

func (t *RadixTree) Set(key string, val string) {

}

func (t RadixTree) DumpContents() {
	dumpContents(t.root, 0)
}

func dumpContents(n *node, indent int) {
	switch x := (*n).(type) {
	case inode:
		dumpContents(x.lc, indent+2)
		dumpContents(x.rc, indent+2)
	case tnode:
		fmt.Println("%s -> %s", x.key, x.val)
	}
}
