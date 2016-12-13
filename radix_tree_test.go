package radix_tree

import (
	"math/rand"
	"testing"
)

func expectGet(t *testing.T, r RadixTree, key string, val string) {
	if actual, err := r.Get(key); err != nil || actual != val {
		t.Errorf("Want err == nil, val == \"%v\". Got err == %v, val == %v", val, err, actual)
	}
}

func expectNotGet(t *testing.T, r RadixTree, key string) {
	if actual, err := r.Get(key); err == nil || actual != "" {
		t.Errorf("Want err != nil, val == \"\". Got err == %v, val == %v", err, actual)
	}
}

func TestGetEmpty(t *testing.T) {
	r := RadixTree{}
	_, err := r.Get("foo")
	if err == nil {
		t.Error("Want err != nil, got err == nil")
	}
}

func TestSetGet(t *testing.T) {
	r := RadixTree{}
	r.Set("foo", "bar")
	expectGet(t, r, "foo", "bar")
}

func TestSetDelete(t *testing.T) {
	r := RadixTree{}
	r.Set("foo", "bar")
	r.Delete("foo")
	expectNotGet(t, r, "foo")
}

func TestSetSetDeleteDelete(t *testing.T) {
	r := RadixTree{}
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
	r := RadixTree{}
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

func TestSetAndGetMixedOrder(t *testing.T) {
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
		r := RadixTree{}
		for _, j := range rand.Perm(len(data)) {
			r.Set(data[j], data[j])
		}
		for _, key := range data {
			expectGet(t, r, key, key)
		}
	}
}

func TestSetAndGetExhaustive3ByteLowercaseEnglish(t *testing.T) {
	var b [3]byte
	r := RadixTree{}
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
}
