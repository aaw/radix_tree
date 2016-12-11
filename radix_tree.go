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

func msbMask(b byte) byte {
	b |= b >> 1
	b |= b >> 2
	b |= b >> 4
	return b & ^(b >> 1)
}

func firstDifferingIndex(s string, t string) int {
	minLen, tl := len(s), len(t)
	if tl < minLen {
		minLen = tl
	}
	for i := 0; i < minLen; i++ {
		if s[i] != t[i] {
			return i
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

func (t RadixTree) Get(key string) (string, error) {
	k, v := t.get(key)
	if k == key {
		return v, nil
	} else {
		return "", errors.New("Not found")
	}
}

func (t *RadixTree) Set(key string, val string) {
	k, _ := t.get(key)
	fmt.Printf("** key is [% x]\n", key)
	fmt.Printf("** k is [% x]\n", k)
	i := firstDifferingIndex(k, key)
	fmt.Printf("** fdi is [%v]\n", i)
	fmt.Printf("** gioz(key,%v) = % X, gioz(k,%v) = % X\n", i, getItemOrZero(key, i), i, getItemOrZero(k, i))
	mask := msbMask(getItemOrZero(key, i) ^ getItemOrZero(k, i))
	fmt.Printf("** mask is [%x]. mask & k = %x, mask & key = %v\n", mask, mask&getItemOrZero(k, i), mask&getItemOrZero(key, i))
	t.set(key, val, i, mask)
}

// In a RadixTree t with t.root != nil, get the ...
func (t RadixTree) get(key string) (string, string) {
	if t.root == nil {
		return "", ""
	}
	n := t.root
	for {
		switch x := n.(type) {
		case *inode:
			if getItemOrZero(key, x.critbyte)&x.critmask == 0 {
				n = x.lc
			} else {
				n = x.rc
			}
		case *tnode:
			return x.key, x.val
		}
	}
}

func (t *RadixTree) set(key string, val string, critbyte int, critmask byte) {
	n := &t.root
	for {
		if *n == nil {
			fmt.Printf("setting %v = %v\n", key, val)
			*n = &tnode{key: key, val: val}
			return
		}
		switch x := (*n).(type) {
		case *inode:
			if x.critbyte >= critbyte {
				// An internal node already discriminates at a byte index
				// that's greater than critbyte. Insert a new internal node here
				// that discriminates on byte index critbyte instead. To do that,
				// we have to walk down the tree and find a terminal node so that
				// we know the first bit set in the len(key)-th byte (which is
				// the same for any string in this subtree, since they are all
				// equal up to byte i and i > len(key).
				fmt.Printf("Replacing internal node %v\n", x)
				if critmask&key[critbyte] == 0 {
					i := inode{lc: nil, rc: *n, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.lc
				} else {
					i := inode{lc: *n, rc: nil, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.rc
				}
			} else if key[x.critbyte]&x.critmask == 0 {
				n = &x.lc
			} else {
				n = &x.rc
			}
		case *tnode:
			if len(x.key) == critbyte && len(key) == critbyte {
				// We're at a terminal node with the same key that we're
				// trying to insert, so just overwrite the value.
				fmt.Printf("Replacing %v's value with %v (was %v)\n", x.key, val, x.val)
				x.val = val
				return
			} else {
				in := inode{lc: nil, rc: nil, critbyte: critbyte, critmask: critmask}
				*n = &in
				if critmask&key[critbyte] == 0 {
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
