package radix_tree

import (
	"errors"
	"strings"
)

type RadixTree struct {
	root node
}

type node interface {
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

func shouldDescendLeft(s string, x inode) bool {
	if x.critbyte < len(s) {
		return s[x.critbyte]&x.critmask == 0
	}
	return true
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
	// Find first differing byte between k and key
	minLen := len(key)
	if len(k) < minLen {
		minLen = len(k)
	}
	diffByte := 0
	var keyByte, kByte byte
	for i := 0; i < minLen; i++ {
		if key[i] != k[i] {
			keyByte, kByte = key[i], k[i]
			break
		}
		diffByte++
	}
	if keyByte == 0 && kByte == 0 {
		if diffByte < len(key) {
			keyByte = key[diffByte]
		}
		if diffByte < len(k) {
			kByte = k[diffByte]
		}
	}
	// mask is the most significant differing bit set in key[i] and k[i]
	mask := keyByte ^ kByte
	mask |= mask >> 1
	mask |= mask >> 2
	mask |= mask >> 4
	mask &= ^(mask >> 1)
	t.set(key, val, diffByte, mask, keyByte)
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
			if shouldDescendLeft(key, *x) {
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

func (t RadixTree) PrefixMatch(prefix string, limit int) []string {
	if t.root == nil {
		return []string{}
	}
	n := t.root
	results := make([]string, 0, limit)
	stack := make([]node, 0)
	for {
		switch x := n.(type) {
		case *inode:
			if x.critbyte < len(prefix) {
				if prefix[x.critbyte]&x.critmask == 0 {
					n = x.lc
				} else {
					n = x.rc
				}
			} else {
				n = x.lc
				stack = append(stack, x.rc)
			}
		case *tnode:
			// TODO: return pairs here
			if strings.HasPrefix(x.key, prefix) {
				results = append(results, x.key)
			}
			if len(stack) == 0 || len(results) >= limit {
				return results
			} else {
				n, stack = stack[len(stack)-1], stack[:len(stack)-1]
			}
		}
	}
}

func (t RadixTree) get(key string) (string, string) {
	if t.root == nil {
		return "", ""
	}
	n := t.root
	for {
		switch x := n.(type) {
		case *inode:
			if shouldDescendLeft(key, *x) {
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
				if critmask&keyByte == 0 {
					i := inode{lc: nil, rc: *n, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.lc
				} else {
					i := inode{lc: *n, rc: nil, critbyte: critbyte, critmask: critmask}
					*n = &i
					n = &i.rc
				}
			} else if shouldDescendLeft(key, *x) {
				n = &x.lc
			} else {
				n = &x.rc
			}
		case *tnode:
			if len(x.key) == critbyte && len(key) == critbyte {
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
