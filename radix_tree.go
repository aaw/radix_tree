package radix_tree

import (
	"unicode/utf8"
)

type RadixTree struct {
	root *node
}

type KV struct {
	key   string
	value string
}

type node struct {
	child map[rune]*node
	data  *KV
}

func newNode() *node {
	return &node{child: make(map[rune]*node)}
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
	for _, r := range runes {
		if n, ok = n.child[r]; !ok {
			return "", false
		}
	}
	if n.data != nil {
		return n.data.value, true
	} else {
		return "", false
	}
}

func (t *RadixTree) Set(key string, val string) {
	n := t.root
	runes, _ := stringToRunes(key, len(key))
	for _, r := range runes {
		if x, ok := n.child[r]; !ok {
			z := newNode()
			n.child[r] = z
			n = z
		} else {
			n = x
		}

	}
	n.data = &KV{key: key, value: val}
}

// TODO: clean up unused paths here?
func (t *RadixTree) Delete(key string) {
	n := t.root
	var ok bool
	runes, _ := stringToRunes(key, len(key))
	for _, r := range runes {
		if n, ok = n.child[r]; !ok {
			return
		}
	}
	n.data = nil
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
	n *node  // child: map[rune]*node, vals: map[rune]string
	s *state // offset, arr
}

func (t RadixTree) Suggest(key string, d int8, n int) []KV {
	results := []KV{}
	process := func(nd *node) bool {
		results = append(results, *nd.data)
		return len(results) < n
	}
	suggest(process, t.root, key, d, n)
	return results
}

func (t RadixTree) SuggestSuffixesAfterExactPrefix(key string, np int, d int8, n int) []KV {
	return []KV{}
}

func (t RadixTree) SuggestSuffixes(key string, d int8, n int) []KV {
	return []KV{}
}

func (t RadixTree) SuggestAfterExactPrefix(key string, np int, d int8, n int) []KV {
	runes, s := stringToRunes(key, np)
	var ok bool
	curr := t.root
	results := []KV{}
	for _, r := range runes {
		if curr, ok = curr.child[r]; !ok {
			return results
		}
	}
	process := func(nd *node) bool {
		results = append(results, *nd.data)
		return len(results) < n
	}
	suggest(process, curr, s, d, n)
	return results
}

func suggest(process func(*node) bool, root *node, key string, d int8, n int) {
	runes, _ := stringToRunes(key, len(key))
	initial := newState(d, int(-2*d))
	initial.arr[2*d] = 0
	stack := []frame{frame{n: root, s: initial}}
	for len(stack) > 0 {
		var f frame
		f, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if f.n.data != nil && f.s.isAccepting(len(key), d) {
			if !process(f.n) {
				return
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				stack = append(stack, frame{n: node, s: ns})
			}
		}
	}
}
