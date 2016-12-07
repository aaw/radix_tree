package radix_tree

import (
	"errors"
	"fmt"
)

type RadixTree struct {
	root node
}

type node interface {
	isTerminal() bool
}

type inode struct {
	lc       node
	rc       node
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
	n := t.root
	var i uint32
	for {
		switch x := n.(type) {
		case inode:
			i = x.critbyte
			if key[i]&x.critmask == 0 {
				n = x.lc
			} else {
				n = x.rc
			}
		case tnode:
			if len(key) != len(x.key) {
				return "", errors.New("Not found")
			}
			for j := int(i); j < len(key); j++ {
				if x.key[j] != key[j] {
					return "", errors.New("Not found")
				}
			}
			return x.val, nil
		}
	}
}

func (t *RadixTree) Set(key string, val string) {
	n := &t.root
	var i uint32
	for {
		if *n == nil {
			*n = &tnode{key: key, val: val}
			return
		} else {
			switch x := (*n).(type) {
			case inode:
				i = x.critbyte
				if key[i]&x.critmask == 0 {
					n = &x.lc
				} else {
					n = &x.rc
				}
			case tnode:
				// TODO: can maybe speed this comp up by starting from i
				equal := true
				j := i
				// TODO: figure out how to treat case where one string is
				// a prefix of another. Always branch right on longer string?
				minLen := len(key)
				if len(x.key) > minLen {
					minLen = len(x.key)
				}
				for ; equal && j < minLen; j++ {
					if x.key[j] != key[j] {
						equal = false
					}
				}
				if equal {
					x.val = val
					return
				} else {
					// TODO: figure out differing bit, set vals in inode
					// accordingly
					in := inode{lc: nil, rc: nil, critbyte: 0, critmask: 0}
					*n = &in
					n = &in.lc
				}
			}
		}
	}
}

func (t RadixTree) DumpContents() {
	dumpContents(t.root, 0)
}

func dumpContents(n node, indent int) {
	switch x := n.(type) {
	case inode:
		dumpContents(x.lc, indent+2)
		dumpContents(x.rc, indent+2)
	case tnode:
		fmt.Println("%s -> %s", x.key, x.val)
	}
}
