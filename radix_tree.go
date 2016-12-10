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
	lc       node // lc === critbit 0
	rc       node // rc === critbit 1
	critbyte int
	critmask uint8
}

type tnode struct {
	key string
	val string
}

func (n tnode) isTerminal() bool { return true }
func (n inode) isTerminal() bool { return false }

func msb_mask(b byte) byte {
	b |= b >> 1
	b |= b >> 2
	b |= b >> 4
	return b & ^(b >> 1)
}

func firstDifferingIndex(s string, t string, i int) int {
	minLen, tl := len(s), len(t)
	if tl < minLen {
		minLen = tl
	}
	for j := i; j < minLen; j++ {
		if s[j] != t[j] {
			return j
		}
	}
	return minLen
}

func getItemOrZero(s string, i int) byte {
	if i < len(s) {
		return s[i]
	}
	return 0
}

func (t RadixTree) Get(key string) (val string, err error) {
	n := t.root
	var i int
	if n == nil {
		return "", errors.New("Not found")
	}
	for {
		switch x := n.(type) {
		case *inode:
			i = x.critbyte
			if getItemOrZero(key, i)&x.critmask == 0 {
				n = x.lc
			} else {
				n = x.rc
			}
		case *tnode:
			j := firstDifferingIndex(key, x.key, i)
			if j < len(key) || len(key) != len(x.key) {
				return "", errors.New("Not found")
			}
			return x.val, nil
		}
	}
}

func (t *RadixTree) Set(key string, val string) {
	n := &t.root
	kl := len(key)
	var i int
	for {
		if *n == nil {
			*n = &tnode{key: key, val: val}
			return
		} else {
			switch x := (*n).(type) {
			case *inode:
				i = x.critbyte
				if i >= kl {
					// Some internal node already discriminates at a byte index i
					// that's greater than len(key). Insert a new internal node here
					// that discriminates on byte index len(key) instead. To do that,
					// we have to walk down the tree and find a terminal node so that
					// we know the first bit set in the len(key)-th byte (which is
					// the same for any string in this subtree, since they are all
					// equal up to byte i and i > len(key).
					t := x.lc
					for !t.isTerminal() {
						t = t.(inode).lc
					}
					mask := msb_mask(t.(*tnode).key[len(key)])
					in := inode{lc: nil, rc: *n, critbyte: len(key), critmask: mask}
					*n = &in
					n = &in.lc
				} else if key[i]&x.critmask == 0 {
					n = &x.lc
				} else {
					n = &x.rc
				}
			case *tnode:
				j := firstDifferingIndex(key, x.key, i)
				if kl == len(x.key) && kl == j {
					// We're at a terminal node with the same key that we're
					// trying to insert, so just overwrite the value.
					x.val = val
					return
				} else {
					kb := getItemOrZero(key, j)
					xb := getItemOrZero(x.key, j)
					mask := msb_mask(kb ^ xb)
					in := inode{lc: nil, rc: nil, critbyte: j, critmask: mask}
					*n = &in
					if mask&kb == 0 {
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
	fmt.Printf("[%v]\n", n)
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
