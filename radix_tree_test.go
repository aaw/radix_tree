package radix_tree

import (
	"math/rand"
	"sort"
	"strings"
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

func keystr(x []KV) string {
	z := []string{}
	for _, y := range x {
		z = append(z, y.key)
	}
	sort.Strings(z)
	return strings.Join(z, " ")
}

func TestSuggest(t *testing.T) {
	data := []string{
		"f",
		"x",
		"fo",
		"fx",
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
	r := NewTree()
	var got, want string
	unlimited := len(data) + 1
	for _, key := range data {
		r.Set(key, key)
	}
	got = keystr(r.Suggest("foo", 0, unlimited))
	want = "foo"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("foo", 1, unlimited))
	want = "fo foo fooY fooZ fooa foob fooc"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("foo", 2, unlimited))
	want = "f fo foo fooY fooZ fooa fooaa fooab foob fooc fx"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("foo", 3, unlimited))
	want = "f fo foo fooY fooZ fooa fooaa fooaaZ fooaaa fooab foob fooc fx x"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("fooaaa", 3, unlimited))
	want = "foo fooY fooZ fooa fooaa fooaaZ fooaaa fooaaaa fooaaaaY fooaaaaa fooaaaaaa fooaaac fooab foob fooc"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("foobbb", 3, unlimited))
	want = "foo fooY fooZ fooa fooaa fooaaZ fooaaa fooab foob fooc"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.Suggest("foobbb", 4, unlimited))
	want = "fo foo fooY fooZ fooa fooaa fooaaZ fooaaa fooaaaa fooaaac fooab foob fooc"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
}

func TestSuggestWithLimit(t *testing.T) {
	// TODO: fill this out, add *WithLimit for all forms of suggest*
}

func TestSuggestAfterExactPrefix(t *testing.T) {
	data := []string{
		"a",
		"aa",
		"aaafoo",
		"aaf",
		"aafo",
		"aafoo",
		"aafoox",
		"aafooxx",
		"aafooxxx",
		"aafox",
		"aafx",
		"aafxx",
		"abfoo",
		"abfooxx",
		"b",
		"bbfoo",
		"foo",
	}
	r := NewTree()
	var got, want string
	unlimited := len(data) + 1
	for _, key := range data {
		r.Set(key, key)
	}
	got = keystr(r.SuggestAfterExactPrefix("aafoo", 2, 0, unlimited))
	want = "aafoo"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestAfterExactPrefix("aafoo", 2, 1, unlimited))
	want = "aaafoo aafo aafoo aafoox aafox"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestAfterExactPrefix("aafoo", 2, 2, unlimited))
	want = "aaafoo aaf aafo aafoo aafoox aafooxx aafox aafx aafxx"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestAfterExactPrefix("aafoo", 2, 3, unlimited))
	want = "aa aaafoo aaf aafo aafoo aafoox aafooxx aafooxxx aafox aafx aafxx"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
}

func TestSuggestSuffixes(t *testing.T) {
	data := []string{
		"", "afoo", "f", "fo", "foo", "fooey", "fooeyz", "fooeyzz", "foox",
		"fooxx", "fooxxx", "fooxxxaaaaa", "fooz", "fox", "fx", "fxx", "gog",
		"gogx", "gogy", "gogyy", "gogyyy",
	}
	r := NewTree()
	var got, want string
	unlimited := len(data) + 1
	for _, key := range data {
		r.Set(key, key)
	}
	got = keystr(r.SuggestSuffixes("foo", 0, unlimited))
	want = "foo fooey fooeyz fooeyzz foox fooxx fooxxx fooxxxaaaaa fooz"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestSuffixes("foo", 1, unlimited))
	want = "afoo fo foo fooey fooeyz fooeyzz foox fooxx fooxxx fooxxxaaaaa fooz fox"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestSuffixes("foo", 2, unlimited))
	want = "afoo f fo foo fooey fooeyz fooeyzz foox fooxx fooxxx fooxxxaaaaa fooz fox fx fxx gog gogx gogy gogyy gogyyy"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
	got = keystr(r.SuggestSuffixes("foo", 3, unlimited))
	want = " afoo f fo foo fooey fooeyz fooeyzz foox fooxx fooxxx fooxxxaaaaa fooz fox fx fxx gog gogx gogy gogyy gogyyy"
	if got != want {
		t.Errorf("Want '%v', got '%v'\n", want, got)
	}
}

func TestSuggestSuffixesAfterExactPrefix(t *testing.T) {
	// TODO: fill this out
}

func editDistance(s string, t string) int {
	rs, _ := stringToRunes(s, len(s))
	rt, _ := stringToRunes(t, len(s))
	return editDistanceHelper(rs, rt)
}

func editDistanceHelper(s []rune, t []rune) int {
	if len(s) == 0 {
		return len(t)
	} else if len(t) == 0 {
		return len(s)
	} else if s[len(s)-1] == t[len(t)-1] {
		return editDistanceHelper(s[:len(s)-1], t[:len(t)-1])
	} else {
		x := editDistanceHelper(s, t[:len(t)-1])
		y := editDistanceHelper(s[:len(s)-1], t)
		z := editDistanceHelper(s[:len(s)-1], t[:len(t)-1])
		d := x
		if y < d {
			d = y
		}
		if z < d {
			d = z
		}
		return 1 + d
	}
}

/*
func allStringsWithinEditDistance(needle string, haystack []string, d int) {
	for _, key := range haystack {

	}
}


func TestSuggestExhaustive(t *testing.T) {
	rand.Seed(0)
	for i := 0; i < 100; i++ {
		r := NewTree()
		keys := map[string]bool{}
		for j := 0; j < 10000; j++ {
			s := randString(10)
			keys[s] = true
			r.Set(s, s)
		}
		haystack := []string{}
		for k := range keys {
			haystack := append(haystack, k)
		}
	}
	// TODO: exhaustive edit distance test using editDistance above...
}*/

// TODO: benchmark that loads words file
