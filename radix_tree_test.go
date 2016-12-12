package radix_tree

import (
	"testing"
)

// TODO: replace most assertions below with expectGet
func expectGet(t *testing.T, r RadixTree, key string, val string) {
	if aval, err := r.Get(key); err != nil || aval != val {
		t.Errorf("Want err == nil, val == \"%v\". Got err == %v, val == %v", val, err, aval)
	}
}

func TestGetEmpty(t *testing.T) {
	r := RadixTree{}
	_, err := r.Get("foo")
	if err == nil {
		t.Error("Want err != nil, got err == nil")
	}
}

func TestSetAndGetBasic(t *testing.T) {
	r := RadixTree{}
	r.Set("foo", "bar")
	r.DumpContents()
	expectGet(t, r, "foo", "bar")
}

func TestGetUnsuccessful(t *testing.T) {
	r := RadixTree{}
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	expectGet(t, r, "fooey", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foozle", "barc")
}

func TestSetAndGetCommonPrefix(t *testing.T) {
	r := RadixTree{}
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	if _, err := r.Get("foo"); err == nil {
		t.Errorf("Want err != nil, got err == %v\n", err)
	}
	expectGet(t, r, "fooey", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foozle", "barc")
}

func TestSetAndGetSubstrings(t *testing.T) {
	r := RadixTree{}
	r.Set("fooingly", "bara")
	r.Set("fooing", "barb")
	r.Set("foo", "barc")
	expectGet(t, r, "fooingly", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foo", "barc")
}

func TestSetAndGetExhaustive(t *testing.T) {
	//var b [3]byte
	var b [2]byte
	r := RadixTree{}
	keys := make([]string, 0)
	//for i := 97; i < 123; i++ {
	for j := 97; j < 123; j++ {
		for k := 97; k < 123; k++ {
			//b[0], b[1], b[2] = byte(i), byte(j), byte(k)
			b[0], b[1] = byte(j), byte(k)
			key := string(b[:])
			keys = append(keys, key)
		}
	}
	//}
	for _, key := range keys {
		r.Set(key, key)
	}
	r.DumpContents()
	for _, key := range keys {
		expectGet(t, r, key, key)
	}
}
