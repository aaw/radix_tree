package radix_tree

import (
	"unicode/utf8"
)

// (6)---->(7)---->((8))
//  ^     ↗ ^      ↗ ^
//  |    /  |     /  |
//  |  /    |   /    |
//  |/      | /      |
// (3)---->(4)---->((5))
//  ^     ↗ ^      ↗ ^
//  |    /  |     /  |
//  |  /    |   /    |
//  |/      | /      |
// (0)---->(1)---->((2))

type RadixTree struct {
	root *node
}

type node struct {
	child map[rune]*node
	// TODO: store vals as []rune, run DecodeRuneInString at insertion?
	vals map[rune]string
}

func newNode() *node {
	return &node{child: make(map[rune]*node), vals: make(map[rune]string)}
}

func NewTree() RadixTree {
	return RadixTree{root: newNode()}
}

func (t *RadixTree) Get(key string) (string, bool) {
	n := t.root
	var ok bool
	for i, r := range key {
		if i < len(key)-1 {
			if n, ok = n.child[r]; !ok {
				return "", false
			}
		} else {
			val, ok := n.vals[r]
			return val, ok
		}
	}
	return "", false
}

func (t *RadixTree) Set(key string, val string) {
	// TODO: disallow key == "", since that will break this code.
	n := t.root
	for i, r := range key {
		if i < len(key)-1 {
			if x, ok := n.child[r]; !ok {
				z := newNode()
				n.child[r] = z
				n = z
			} else {
				n = x
			}
		} else {
			n.vals[r] = val
		}
	}
}

func (t *RadixTree) Delete(key string) {
	n := t.root
	var ok bool
	for i, r := range key {
		if i < len(key)-1 {
			if n, ok = n.child[r]; !ok {
				return
			}
		} else {
			delete(n.vals, r)
		}
	}
}

func (t RadixTree) PrefixMatch(prefix string, limit int) []string {
	return []string{}
}

// TODO: pack arr into a 64 bit int for d <= 5?
type state struct {
	offset int
	arr    []int
}

func (s state) isAccepting(wl int, d int) bool {
	for i, x := range s.arr {
		if x < d+1 && s.offset+i+x == wl {
			return true
		}
	}
	return false
}

func newState(d int, offset int) *state {
	arr := make([]int, 2*d+1)
	for i := range arr {
		arr[i] = d + 1
	}
	return &state{offset: offset, arr: arr}
}

func min(x int, y int, z int) int {
	a := x
	if y < a {
		a = y
	}
	if z < a {
		a = z
	}
	return a
}

// diagonal in the actual NFA is offset + d. index into arr is d for
// the main diagonal, so initial state is offset: -d, arr = [d+1, ... , d+1, 0, d+1, ..., d+1]
func (s state) transition(w []rune, r rune, d int) (*state, bool) {
	// cr == carry right, right transition
	// cu == carry up, up transition from diagonal to the right
	ns := newState(d, s.offset+1)
	isValid := false
	cr := d + 1
	for j, x := range s.arr {
		// Calculate carry up from j+1st diagonal
		cu := d + 1
		if j < len(s.arr)-1 {
			cu = s.arr[j+1]
		}
		carry := min(x, cr, cu)
		if carry < d+1 {
			isValid = true
		}
		ns.arr[j] = carry
		for k := x; k < d+1; k++ {
			if w[j+s.offset+k] == r /* TODO: right comp here? */ {
				cr = k
			}
		}
	}
	return ns, isValid
}

func stringToRunes(s string) []rune {
	rs := []rune{}
	for i, w := 0, 0; i < len(s); i += w {
		r, width := utf8.DecodeRuneInString(s[i:])
		rs = append(rs, r)
		w = width
	}
	return rs
}

type frame struct {
	n  *node  // child: map[rune]*node, vals: map[rune]string
	s  *state // offset, arr
	rs []rune
}

func (t RadixTree) Suggest(key string, d int) []string {
	runes := stringToRunes(key)
	results := []string{}
	initial := newState(d, -d)
	initial.arr[d] = 0
	stack := []frame{frame{n: t.root, s: initial, rs: []rune{}}}
	for len(stack) > 0 {
		f, stack := stack[len(stack)-1], stack[:len(stack)-1]
		for r, _ := range f.n.vals {
			if ns, ok := f.s.transition(runes, r, d); ok && ns.isAccepting(len(f.rs)+1, d) {
				nrs := make([]rune, len(f.rs)+1)
				copy(nrs, f.rs)
				nrs = append(nrs, r)
				results = append(results, string(nrs))
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				nrs := make([]rune, len(f.rs)+1)
				copy(nrs, f.rs)
				nrs = append(nrs, r)
				stack = append(stack, frame{n: node, s: ns, rs: nrs})
			}
		}
	}
	return results
}
