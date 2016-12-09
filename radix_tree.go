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
	critbyte int
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
	var i int
	if n == nil {
		return "", errors.New("Not found")
	}
	for {
		fmt.Printf("loop: %T", n)
		switch x := n.(type) {
		case *inode:
			fmt.Printf("Going down inode %v\n", x)
			i = x.critbyte
			if i >= len(key) {
				return "", errors.New("Not found")
			}
			if key[i]&x.critmask == 0 {
				n = x.lc
			} else {
				n = x.rc
			}
		case *tnode:
			fmt.Printf("Going down tnode %v\n", x)
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

func msb_mask(b byte) byte {
	b |= b >> 1
	b |= b >> 2
	b |= b >> 4
	return b & ^(b >> 1)
}

func (t *RadixTree) Set(key string, val string) {
	n := &t.root
	var i int
	for {
		if *n == nil {
			fmt.Println("Setting from nil")
			*n = &tnode{key: key, val: val}
			return
		} else {
			switch x := (*n).(type) {
			case *inode:
				fmt.Printf("Going down inode %v\n", x)
				i = x.critbyte
				if i >= len(key) {
					// what do
				}
				if key[i]&x.critmask == 0 {
					n = &x.lc
				} else {
					n = &x.rc
				}
			case *tnode:
				fmt.Printf("Going down tnode %v\n", x)
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
					mask := msb_mask(x.key[j] ^ key[j])
					in := inode{lc: nil, rc: nil, critbyte: j, critmask: mask}
					*n = &in
					if mask&key[j] == 0 {
						in.rc = x
						n = &in.lc
					} else {
						in.lc = x
						n = &in.rc
					}
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
	case *inode:
		dumpContents(x.lc, indent+2)
		dumpContents(x.rc, indent+2)
	case *tnode:
		for i := 0; i < indent; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("%v -> %v\n", x.key, x.val)
	}
}
