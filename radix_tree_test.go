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
