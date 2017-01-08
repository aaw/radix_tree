package radix_tree

import (
	"fmt"
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
		dist := wl - s.offset - i
		if dist <= d && dist >= x {
			fmt.Printf("[%v](%v,%v) with %v, %v [Accepting]\n", s, i, x, wl, d)
			return true
		}
	}
	fmt.Printf("[%v] with %v, %v [Not accepting]\n", s, wl, d)
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
	for j := range ns.arr {
		cr := d + 1
		for k := s.arr[j]; k < d+1; k++ {
			fmt.Printf("considering carry\n")
			if j+s.offset+k < len(w) && w[j+s.offset+k] == r /* TODO: right comp here? */ {
				if cr > k {
					fmt.Printf("carrying %v right\n", k)
					cr = k
				}
			}
		}
		x := d + 1
		if j < len(s.arr)-1 {
			x = s.arr[j+1] + 1
		}
		cu := d + 1
		if j < len(s.arr)-2 {
			cu = s.arr[j+2] + 1
		}
		fmt.Printf("[%v] x: %v, cr: %v, cu: %v\n", j, x+1, cr, cu)
		carry := min(x, cr, cu)
		if carry < d+1 {
			ns.arr[j], isValid = carry, true
		}
	}
	fmt.Printf("transition from %v to %v (%v)\n", s, ns, isValid)
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
	initial := newState(d, -2*d)
	initial.arr[2*d] = 0
	stack := []frame{frame{n: t.root, s: initial, rs: []rune{}}}
	for len(stack) > 0 {
		var f frame
		f, stack = stack[len(stack)-1], stack[:len(stack)-1]
		fmt.Printf("Stack size: %v, current frame: %v\n", len(stack), f)
		for r, _ := range f.n.vals {
			fmt.Printf("\nConsidering %v:%v...\n", string(f.rs), string(r))
			// TODO: if isAccepting depends on len(key) like i think it does, push len(key) into state
			if ns, ok := f.s.transition(runes, r, d); ok && ns.isAccepting(len(key), d) {
				fmt.Printf("Accepting: %v:%v\n", string(f.rs), string(r))
				nrs := make([]rune, len(f.rs))
				copy(nrs, f.rs)
				nrs = append(nrs, r)
				results = append(results, string(nrs))
			} else {
				fmt.Printf("-> %v Not accepting\n", string(r))
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				fmt.Printf("Transition: %v:%v\n", string(f.rs), string(r))
				nrs := make([]rune, len(f.rs))
				copy(nrs, f.rs)
				nrs = append(nrs, r)
				stack = append(stack, frame{n: node, s: ns, rs: nrs})
			}
		}
	}
	return results
}
