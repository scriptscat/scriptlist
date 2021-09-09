package utils

import (
	"testing"
)

func BenchmarkGetRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRandomString(12)
	}
}
