package radix_tree

import (
	"testing"
)

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
	r.DumpContents()
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
	val, ok := r.Get("fooey")
	if ok != nil {
		t.Errorf("Want ok == nil, got ok = %v", ok)
	}
	if val != "bara" {
		t.Errorf("Want val == 'bara', got val = %v", val)
	}
}

func TestSetAndGetSubstrings(t *testing.T) {

}
