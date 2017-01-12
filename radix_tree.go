package radix_tree

import (
	"unicode/utf8"
)

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
			//fmt.Printf("[%v](%v,%v) with %v, %v [Accepting]\n", s, i, x, wl, d)
			return true
		}
	}
	//fmt.Printf("[%v] with %v, %v [Not accepting]\n", s, wl, d)
	return false
}

func newState(d int, offset int) *state {
	arr := make([]int, 2*d+1)
	for i := range arr {
		arr[i] = d + 1
	}
	return &state{offset: offset, arr: arr}
}

// diagonal in the actual NFA is offset + d. index into arr is d for
// the main diagonal, so initial state is offset: -d, arr = [d+1, ... , d+1, 0, d+1, ..., d+1]
func (s state) transition(w []rune, r rune, d int) (*state, bool) {
	ns := newState(d, s.offset+1)
	isValid := false
	for j := range ns.arr {
		// Compute carry right
		val := d + 1
		for k := s.arr[j]; k < d+1; k++ {
			if j+s.offset+k < len(w) && w[j+s.offset+k] == r {
				if val > k {
					//fmt.Printf("carrying %v right\n", k)
					val = k
				}
			}
		}
		// Compute diagonal contribution
		if j < len(s.arr)-1 && s.arr[j+1]+1 < val {
			val = s.arr[j+1] + 1
		}
		// Compute carry up
		if j < len(s.arr)-2 && s.arr[j+2]+1 < val {
			val = s.arr[j+2] + 1
		}
		if val < d+1 {
			ns.arr[j], isValid = val, true
		}
	}
	//fmt.Printf("transition from %v to %v (%v)\n", s, ns, isValid)
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

func pushRune(rs []rune, r rune) []rune {
	newrs := make([]rune, len(rs))
	copy(newrs, rs)
	return append(newrs, r)
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
		// fmt.Printf("Stack size: %v, current frame: %v\n", len(stack), f)
		for r, _ := range f.n.vals {
			// fmt.Printf("\nConsidering %v:%v...\n", string(f.rs), string(r))
			// TODO: if isAccepting depends on len(key) like i think it does, push len(key) into state
			if ns, ok := f.s.transition(runes, r, d); ok && ns.isAccepting(len(key), d) {
				//fmt.Printf("Accepting: %v:%v\n", string(f.rs), string(r))
				results = append(results, string(pushRune(f.rs, r)))
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				//fmt.Printf("Transition: %v:%v\n", string(f.rs), string(r))
				stack = append(stack, frame{n: node, s: ns, rs: pushRune(f.rs, r)})
			}
		}
	}
	// TODO: add a limit param, also investigate returning items on a channel.
	// TODO: return key, value pairs here.
	return results
}
