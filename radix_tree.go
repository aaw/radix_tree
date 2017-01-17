package radix_tree

import (
	"unicode/utf8"
)

type RadixTree struct {
	root *node
}

type node struct {
	child map[rune]*node
	vals  map[rune]string
}

func newNode() *node {
	return &node{child: make(map[rune]*node), vals: make(map[rune]string)}
}

func NewTree() RadixTree {
	return RadixTree{root: newNode()}
}

// Read at most the first n runes from the string, return those
// runes in an array and the remaining string.
func stringToRunes(s string, n int) ([]rune, string) {
	rs := []rune{}
	for i, w := 0, 0; i < len(s); i += w {
		if len(rs) >= n {
			return rs, s[i:]
		}
		r, width := utf8.DecodeRuneInString(s[i:])
		rs = append(rs, r)
		w = width
	}
	return rs, ""
}

func (t *RadixTree) Get(key string) (string, bool) {
	n := t.root
	var ok bool
	runes, _ := stringToRunes(key, len(key))
	for i, r := range runes {
		if i < len(runes)-1 {
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
	if len(key) == 0 {
		panic("Empty key not allowed.")
	}
	n := t.root
	runes, _ := stringToRunes(key, len(key))
	for i, r := range runes {
		if i < len(runes)-1 {
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
	runes, _ := stringToRunes(key, len(key))
	for i, r := range runes {
		if i < len(runes)-1 {
			if n, ok = n.child[r]; !ok {
				return
			}
		} else {
			delete(n.vals, r)
		}
	}
}

type state struct {
	offset int
	arr    []int8
}

func (s state) isAccepting(wl int, d int8) bool {
	for i, x := range s.arr {
		dist := int8(wl - s.offset - i)
		if dist <= d && dist >= x {
			//fmt.Printf("[%v](%v,%v) with %v, %v [Accepting]\n", s, i, x, wl, d)
			return true
		}
	}
	//fmt.Printf("[%v] with %v, %v [Not accepting]\n", s, wl, d)
	return false
}

func newState(d int8, offset int) *state {
	arr := make([]int8, 2*d+1)
	for i := range arr {
		arr[i] = int8(d + 1)
	}
	return &state{offset: offset, arr: arr}
}

// diagonal in the actual NFA is offset + d. index into arr is d for
// the main diagonal, so initial state is offset: -d, arr = [d+1, ... , d+1, 0, d+1, ..., d+1]
func (s state) transition(w []rune, r rune, d int8) (*state, bool) {
	ns := newState(d, s.offset+1)
	isValid := false
	for j := range ns.arr {
		// Compute carry right
		val := d + 1
		for k := s.arr[j]; k < d+1; k++ {
			if j+s.offset+int(k) < len(w) && w[j+s.offset+int(k)] == r {
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

type frame struct {
	n  *node  // child: map[rune]*node, vals: map[rune]string
	s  *state // offset, arr
	rp *runePath
}

type runePath struct {
	parent *runePath
	value  rune
}

func walkRunePath(node *runePath) []rune {
	runes := []rune{}
	for node != nil {
		runes = append([]rune{node.value}, runes...)
		node = node.parent
	}
	return runes
}

type KV struct {
	key   string
	value string
}

func (t RadixTree) Suggest(key string, d int8, n int) []KV {
	return suggest(t.root, nil, key, d, n)
}

func (t RadixTree) SuggestSuffixesAfterExactPrefix(key string, np int, d int8, n int) []KV {
	return []KV{}
}

func (t RadixTree) SuggestSuffixes(key string, d int8, n int) []KV {
	return []KV{}
}

func (t RadixTree) SuggestAfterExactPrefix(key string, np int, d int8, n int) []KV {
	runes, s := stringToRunes(key, np)
	var rp *runePath
	node := t.root
	var ok bool
	for _, r := range runes {
		if node, ok = node.child[r]; !ok {
			return []KV{}
		}
		rp = &runePath{parent: rp, value: r}
	}
	return suggest(node, rp, s, d, n)
}

func suggest(root *node, rp *runePath, key string, d int8, n int) []KV {
	runes, _ := stringToRunes(key, len(key))
	results := []KV{}
	initial := newState(d, int(-2*d))
	initial.arr[2*d] = 0
	stack := []frame{frame{n: root, s: initial, rp: rp}}
	for len(stack) > 0 {
		var f frame
		f, stack = stack[len(stack)-1], stack[:len(stack)-1]
		// fmt.Printf("Stack size: %v, current frame: %v\n", len(stack), f)
		for r, val := range f.n.vals {
			// fmt.Printf("\nConsidering %v:%v...\n", string(f.rs), string(r))
			// TODO: if isAccepting depends on len(key) like i think it does, push len(key) into state
			if ns, ok := f.s.transition(runes, r, d); ok && ns.isAccepting(len(key), d) {
				//fmt.Printf("Accepting: %v:%v\n", string(f.rs), string(r))
				results = append(results, KV{key: string(append(walkRunePath(f.rp), r)), value: val})
				if len(results) >= n {
					return results
				}
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				//fmt.Printf("Transition: %v:%v\n", string(f.rs), string(r))
				stack = append(stack, frame{n: node, s: ns, rp: &runePath{parent: f.rp, value: r}})
			}
		}
	}
	return results
}
