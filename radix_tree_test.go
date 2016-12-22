package radix_tree

import (
	"math/rand"
	"testing"
)

func expectGet(t *testing.T, r RadixTree, key string, val string) {
	if actual, ok := r.Get(key); ok && actual != val {
		t.Errorf("Want val == \"%v\", ok. Got val == %v, ok == %v", val, actual, ok)
	}
}

func expectNotGet(t *testing.T, r RadixTree, key string) {
	if actual, ok := r.Get(key); ok {
		t.Errorf("Want !ok. Got val == %v, ok == %v", actual, ok)
	}
}

func TestGetEmpty(t *testing.T) {
	r := NewTree()
	if _, ok := r.Get("foo"); ok {
		t.Error("Want !ok, got ok")
	}
}

func TestSetGet(t *testing.T) {
	r := NewTree()
	r.Set("foo", "bar")
	expectGet(t, r, "foo", "bar")
}

func TestSetDelete(t *testing.T) {
	r := NewTree()
	r.Set("foo", "bar")
	r.Delete("foo")
	expectNotGet(t, r, "foo")
}

func TestSetSetDeleteDelete(t *testing.T) {
	r := NewTree()
	r.Set("foo", "bar")
	r.Set("bar", "foo")
	r.Delete("foo")
	expectNotGet(t, r, "foo")
	expectGet(t, r, "bar", "foo")
	r.Delete("bar")
	expectNotGet(t, r, "foo")
	expectNotGet(t, r, "bar")
}

func TestSetSetSetDeleteDeleteDelete(t *testing.T) {
	r := NewTree()
	r.Set("foo", "bar")
	r.Set("bar", "foo")
	r.Set("baz", "biz")
	r.Delete("foo")
	expectNotGet(t, r, "foo")
	expectGet(t, r, "bar", "foo")
	expectGet(t, r, "baz", "biz")
	r.Delete("bar")
	expectNotGet(t, r, "foo")
	expectNotGet(t, r, "bar")
	expectGet(t, r, "baz", "biz")
	r.Delete("baz")
	expectNotGet(t, r, "foo")
	expectNotGet(t, r, "bar")
	expectNotGet(t, r, "baz")
}

func TestGetUnsuccessful(t *testing.T) {
	r := NewTree()
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	expectGet(t, r, "fooey", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foozle", "barc")
}

func TestDeleteUnsuccessful(t *testing.T) {
	r := NewTree()
	r.Delete("foo")
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	r.Delete("foo")
	r.Delete("fooe")
	r.Delete("fooeyy")
	expectGet(t, r, "fooey", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foozle", "barc")
}

func TestSetAndGetCommonPrefix(t *testing.T) {
	r := NewTree()
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	expectNotGet(t, r, "foo")
	expectGet(t, r, "fooey", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foozle", "barc")
}

func TestSetAndGetSubstrings(t *testing.T) {
	r := NewTree()
	r.Set("fooingly", "bara")
	r.Set("fooing", "barb")
	r.Set("foo", "barc")
	expectGet(t, r, "fooingly", "bara")
	expectGet(t, r, "fooing", "barb")
	expectGet(t, r, "foo", "barc")
}

func TestSetGetDeleteMixedOrder(t *testing.T) {
	rand.Seed(0)
	data := []string{
		"foo",
		"fooa",
		"foob",
		"fooc",
		"fooY",
		"fooZ",
		"fooaa",
		"fooab",
		"fooaaa",
		"fooaaZ",
		"fooaaaa",
		"fooaaac",
		"fooaaaaa",
		"fooaaaaY",
		"fooaaaaaa",
		"fooaaaaaaa",
		"fooaaaaaaaa",
	}
	for i := 0; i < 1000; i++ {
		r := NewTree()
		for j := 0; j < 10; j++ {
			for _, k := range rand.Perm(len(data)) {
				expectNotGet(t, r, data[k])
				r.Set(data[k], data[k])
			}
			for _, key := range data {
				expectGet(t, r, key, key)
			}
			for _, k := range rand.Perm(len(data)) {
				r.Delete(data[k])
			}
		}
	}
}

func TestSetAndGetExhaustive3ByteLowercaseEnglish(t *testing.T) {
	var b [3]byte
	r := NewTree()
	keys := make([]string, 0)
	for i := 97; i < 123; i++ {
		for j := 97; j < 123; j++ {
			for k := 97; k < 123; k++ {
				b[0], b[1], b[2] = byte(i), byte(j), byte(k)
				key := string(b[:])
				keys = append(keys, key)
			}
		}
	}
	for _, key := range keys {
		r.Set(key, key)
	}
	for _, key := range keys {
		expectGet(t, r, key, key)
	}
	for _, key := range keys {
		r.Delete(key)
		expectNotGet(t, r, key)
	}
}

/*
func contains(x []string, s string) bool {
	for _, y := range x {
		if s == y {
			return true
		}
	}
	return false
}

func TestPrefixMatch(t *testing.T) {
	r := NewTree()
	matches := []string{
		"foo",
		"fooa",
		"foob",
		"food",
		"fooaa",
		"fooab",
		"fooba",
		"foobb",
		"foobaa",
		"fooaab",
	}
	non_matches := []string{
		"foa",
		"xxx",
		"fox",
		"foaaa",
		"foxaaa",
		"foxaaaa",
	}
	for _, key := range matches {
		r.Set(key, key)
	}
	for _, key := range non_matches {
		r.Set(key, key)
	}
	actual := r.PrefixMatch("foo", 100)
	for _, key := range matches {
		if !contains(actual, key) {
			t.Errorf("Want %v in list of prefix matches, but it wasn't there\n", key)
		}
	}
	for _, key := range non_matches {
		if contains(actual, key) {
			t.Errorf("Want %v not in list of prefix matches, but it was there\n", key)
		}
	}
}

func TestPrefixMatchUnsuccessful(t *testing.T) {
	r := NewTree()
	matches := []string{
		"axxx",
		"bxxx",
		"cxxx",
		"dxxx",
		"exxx",
	}
	for _, key := range matches {
		r.Set(key, key)
	}
	for _, prefix := range []string{"f", "g", "h", "Z"} {
		if x := r.PrefixMatch(prefix, 100); len(x) > 0 {
			t.Errorf("Want no prefixes to match \"%v\", got %v", prefix, x)
		}
	}
}
*/
