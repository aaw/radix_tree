package radix_tree

import (
	"fmt"
	"testing"
)

func TestRadixTree(t *testing.T) {
	const nihongo = "日本語"
	for i := 0; i < len(nihongo); i++ {
		fmt.Printf("%v: %v\n", i, nihongo[i])
	}
}
