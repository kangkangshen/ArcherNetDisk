package bufs

import "testing"

func BenchmarkNewBufferPool(b *testing.B) {
	b.N = 1000000
	bp := NewBufferPool(1024, 1024)

	for i := 0; i < b.N; i++ {
		b := bp.Get()
		bp.Put(b)
	}
}
