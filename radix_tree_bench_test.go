package radix_tree

import (
	"math/rand"
	"testing"
)

var data []string
var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
