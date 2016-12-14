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
	i := firstDifferingIndex(k, key)
	keyByte := getItemOrZero(key, i)
	mask := msbMask(keyByte ^ getItemOrZero(k, i))
	t.set(key, val, i, mask, keyByte)
}

func (t *RadixTree) Delete(key string) (string, error) {
	if t.root == nil {
		return "", errors.New("Not found")
	}
	var oc *node
	n := &t.root
	parent := n
	for {
		switch x := (*n).(type) {
		case *inode:
			parent = n
			if getItemOrZero(key, x.critbyte)&x.critmask == 0 {
				oc, n = &x.rc, &x.lc
			} else {
				oc, n = &x.lc, &x.rc
			}
		case *tnode:
			if x.key != key {
				return "", errors.New("Not found")
			}
			if oc == nil {
				*parent = nil
			} else {
				*parent = *oc
			}
			return x.val, nil
		}
	}
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

func (t *RadixTree) set(key string, val string, critbyte int, critmask byte, keyByte byte) {
	n := &t.root
	for {
		if *n == nil {
			*n = &tnode{key: key, val: val}
			return
		}
		switch x := (*n).(type) {
		case *inode:
			if x.critbyte > critbyte || (x.critbyte == critbyte && x.critmask < critmask) {
				// An internal node already discriminates at a byte index
				// that's greater than critbyte. Insert a new internal node here
				// that discriminates on byte index critbyte instead. To do that,
				// we have to walk down the tree and find a terminal node so that
				// we know the first bit set in the len(key)-th byte (which is
				// the same for any string in this subtree, since they are all
				// equal up to byte i and i > len(key).
				if critmask&keyByte == 0 {
					i := inode{lc: nil, rc: *n, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.lc
				} else {
					i := inode{lc: *n, rc: nil, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.rc
				}
			} else if getItemOrZero(key, x.critbyte)&x.critmask == 0 {
				n = &x.lc
			} else {
				n = &x.rc
			}
		case *tnode:
			if len(x.key) == critbyte && len(key) == critbyte {
				// We're at a terminal node with the same key that we're
				// trying to insert, so just overwrite the value.
				x.val = val
				return
			} else {
				in := inode{lc: nil, rc: nil, critbyte: critbyte, critmask: critmask}
				*n = &in
				if critmask&keyByte == 0 {
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
