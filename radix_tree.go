package radix_tree

type RadixTree struct {
	root node
}

type node struct {
	child map[rune]node
	vals  map[rune]string
}

func newNode() node {
	return node{child: make(map[rune]node), vals: make(map[rune]string)}
}

func NewTree() RadixTree {
	return RadixTree{root: newNode()}
}

func (t RadixTree) Get(key string) (string, bool) {
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
