package levtrie

import (
	"unicode/utf8"
)

type Trie struct {
	root *node
}

type KV struct {
	key   string
	value string
}

// A Trie node.
type node struct {
	child map[rune]*node
	data  *KV
}

func newNode() *node {
	return &node{child: make(map[rune]*node)}
}

func New() Trie {
	return Trie{root: newNode()}
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

func (t *Trie) Get(key string) (string, bool) {
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

func (t *Trie) Set(key string, val string) {
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
func (t *Trie) Delete(key string) {
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
	n *node  // child: map[rune]*node, data: *KV
	s *state // offset, arr
}

type strategy interface {
	processAcceptingNode(n *node, results *[]KV, limit int)
	keepGoingAfterAccept() bool
}

type doNotExpandSuffixes struct{}
type expandSuffixes struct{}

func (x doNotExpandSuffixes) processAcceptingNode(n *node, results *[]KV, limit int) {
	if n.data != nil {
		*results = append(*results, *n.data)
	}
}

func (x expandSuffixes) processAcceptingNode(n *node, results *[]KV, limit int) {
	stack := []*node{n}
	for len(stack) > 0 {
		var x *node
		x, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if x.data != nil {
			*results = append(*results, *x.data)
			if len(*results) >= limit {
				return
			}
		}
		for _, child := range x.child {
			stack = append(stack, child)
		}
	}
}

func (x doNotExpandSuffixes) keepGoingAfterAccept() bool { return true }

func (x expandSuffixes) keepGoingAfterAccept() bool { return false }

func (t Trie) Suggest(key string, d int8, n int) []KV {
	return suggest(doNotExpandSuffixes{}, t.root, key, d, n)
}

func (t Trie) SuggestSuffixes(key string, d int8, n int) []KV {
	return suggest(expandSuffixes{}, t.root, key, d, n)
}

func (t Trie) SuggestAfterExactPrefix(key string, np int, d int8, n int) []KV {
	runes, s := stringToRunes(key, np)
	var ok bool
	curr := t.root
	for _, r := range runes {
		if curr, ok = curr.child[r]; !ok {
			return []KV{}
		}
	}
	return suggest(doNotExpandSuffixes{}, curr, s, d, n)
}

func (t Trie) SuggestSuffixesAfterExactPrefix(key string, np int, d int8, n int) []KV {
	runes, s := stringToRunes(key, np)
	var ok bool
	curr := t.root
	for _, r := range runes {
		if curr, ok = curr.child[r]; !ok {
			return []KV{}
		}
	}
	return suggest(expandSuffixes{}, curr, s, d, n)
}

func suggest(s strategy, root *node, key string, d int8, n int) []KV {
	runes, _ := stringToRunes(key, len(key))
	initial := newState(d, int(-2*d))
	initial.arr[2*d] = 0
	stack := []frame{frame{n: root, s: initial}}
	results := []KV{}
	for len(stack) > 0 {
		var f frame
		f, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if f.s.isAccepting(len(runes), d) {
			s.processAcceptingNode(f.n, &results, n)
			if len(results) >= n {
				break
			}
			if !s.keepGoingAfterAccept() {
				continue
			}
		}
		for r, node := range f.n.child {
			if ns, ok := f.s.transition(runes, r, d); ok {
				stack = append(stack, frame{n: node, s: ns})
			}
		}
	}
	return results
}
