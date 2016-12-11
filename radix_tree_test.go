package radix_tree

import (
	"fmt"
	"testing"
)

// TODO: replace most assertions below with expectGet
func expectGet(t *testing.T, r RadixTree, key string, val string) {
	if aval, ok := r.Get(key); ok != nil || aval != val {
		t.Errorf("Want ok == nil, val == \"%v\". Got ok == %v, val == %v", val, ok, aval)
	}
}

func TestGetEmpty(t *testing.T) {
	r := RadixTree{}
	_, ok := r.Get("foo")
	if ok == nil {
		t.Error("Want ok != nil, got ok = nil")
	}
}

func TestSetAndGetBasic(t *testing.T) {
	r := RadixTree{}
	r.Set("foo", "bar")
	val, ok := r.Get("foo")
	if ok != nil {
		t.Errorf("Want ok == nil, got ok = %v", ok)
	}
	if val != "bar" {
		t.Errorf("Want val == 'bar', got val = %v", val)
	}
}

func TestGetUnsuccessful(t *testing.T) {
	r := RadixTree{}
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	if _, ok := r.Get("foo"); ok == nil {
		t.Error("Want ok != nil, got ok == nil")
	}
	if _, ok := r.Get("fooe"); ok == nil {
		t.Error("Want ok != nil, got ok == nil")
	}
	if _, ok := r.Get("fooeyz"); ok == nil {
		t.Error("Want ok != nil, got ok == nil")
	}
}

func TestSetAndGetCommonPrefix(t *testing.T) {
	r := RadixTree{}
	r.Set("fooey", "bara")
	r.Set("fooing", "barb")
	r.Set("foozle", "barc")
	if _, ok := r.Get("foo"); ok == nil {
		t.Error("Want ok != nil, got ok == nil")
	}
	if val, ok := r.Get("fooey"); ok != nil || val != "bara" {
		t.Errorf("Want ok == nil, val == \"bara\". Got ok == %v, val == %v", ok, val)
	}
	if val, ok := r.Get("fooing"); ok != nil || val != "barb" {
		t.Errorf("Want ok == nil, val == \"barb\". Got ok == %v, val == %v", ok, val)
	}
	if val, ok := r.Get("foozle"); ok != nil || val != "barc" {
		t.Errorf("Want ok == nil, val == \"barc\". Got ok == %v, val == %v", ok, val)
	}
}

func TestSetAndGetSubstrings(t *testing.T) {
	r := RadixTree{}
	r.Set("fooingly", "bara")
	r.Set("fooing", "barb")
	r.Set("foo", "barc")
	if val, ok := r.Get("fooingly"); ok != nil || val != "bara" {
		t.Errorf("Want ok == nil, val == \"bara\". Got ok == %v, val == %v", ok, val)
	}
	if val, ok := r.Get("fooing"); ok != nil || val != "barb" {
		t.Errorf("Want ok == nil, val == \"barb\". Got ok == %v, val == %v", ok, val)
	}
	if val, ok := r.Get("foo"); ok != nil || val != "barc" {
		t.Errorf("Want ok == nil, val == \"barc\". Got ok == %v, val == %v", ok, val)
	}
}

func TestFooBar(t *testing.T) {
	r := RadixTree{}
	fmt.Printf("Setting cb: z\n")
	r.Set("cb", "z")
	r.DumpContents()
	fmt.Println("")
	fmt.Printf("Setting ca: zz\n")
	r.Set("ca", "zz")
	r.DumpContents()
	fmt.Println("")
	fmt.Printf("Setting bb: y\n")
	r.Set("bb", "y")
	r.DumpContents()
	fmt.Println("")
	fmt.Printf("Setting ab: x\n")
	r.Set("ab", "x")
	r.DumpContents()
	expectGet(t, r, "bb", "y")
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
	for _, key := range keys {
		if val, ok := r.Get(key); ok != nil || val != key {
			t.Errorf("Want ok == nil, val == \"%v\". Got ok == %v, val == %v", key, ok, val)
		}
	}
}
