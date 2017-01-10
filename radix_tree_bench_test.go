package radix_tree

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

var data []string
var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var words []string

func randString(n int) string {
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(runes)
}

func ensureData(n int) {
	if n <= len(data) {
		return
	}
	for i := len(data); i < n; i++ {
		data = append(data, randString(rand.Intn(100)))
	}
}

func ensureWords() {
	if len(words) > 0 {
		return
	}
	filename := "/usr/share/dict/words"
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("%v: %v", filename, err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		words = append(words, word)
	}
}

func benchmarkSuggest(d int, b *testing.B) {
	ensureWords()
	r := NewTree()
	for _, word := range words {
		r.Set(word, word)
	}
	b.ResetTimer()
	data = []string{
		"acetonylacetone",
		"barbaralalia",
		"calcic",
		"dark",
		"using",
		"volt",
		"wrenchingly",
		"xenos",
		"yore",
		"zymosis",
	}
	for i := 0; i < b.N; i++ {
		r.Suggest(data[i%len(data)], d)
	}
}

func BenchmarkSuggest1(b *testing.B) {
	benchmarkSuggest(1, b)
}

func BenchmarkSuggest2(b *testing.B) {
	benchmarkSuggest(2, b)
}

func BenchmarkSuggest3(b *testing.B) {
	benchmarkSuggest(3, b)
}

func BenchmarkSuggest4(b *testing.B) {
	benchmarkSuggest(4, b)
}

func BenchmarkSuggest5(b *testing.B) {
	benchmarkSuggest(5, b)
}

func BenchmarkRadixTreeSet(b *testing.B) {
	ensureData(b.N)
	b.ResetTimer()
	r := NewTree()
	for i := 0; i < b.N; i++ {
		r.Set(data[i], data[i])
	}
}

func BenchmarkMapSet(b *testing.B) {
	ensureData(b.N)
	b.ResetTimer()
	m := make(map[string]string)
	for i := 0; i < b.N; i++ {
		m[data[i]] = data[i]
	}
}

func BenchmarkRadixTreeGet(b *testing.B) {
	ensureData(b.N)
	r := NewTree()
	for i := 0; i < b.N; i++ {
		r.Set(data[i], data[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Get(data[i])
	}
}

func BenchmarkMapGet(b *testing.B) {
	ensureData(b.N)
	m := make(map[string]string)
	for i := 0; i < b.N; i++ {
		m[data[i]] = data[i]
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[data[i]]
	}
}

func BenchmarkRadixTreeDelete(b *testing.B) {
	ensureData(b.N)
	r := NewTree()
	for i := 0; i < b.N; i++ {
		r.Set(data[i], data[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Delete(data[i])
	}
}

func BenchmarkMapDelete(b *testing.B) {
	ensureData(b.N)
	m := make(map[string]string)
	for i := 0; i < b.N; i++ {
		m[data[i]] = data[i]
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, data[i])
	}
}
