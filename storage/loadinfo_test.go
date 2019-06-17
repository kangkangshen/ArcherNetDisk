package storage

import (
	"runtime"
	"testing"
)

func TestGetLoadInfo(t *testing.T) {
	t.Log(runtime.MemStats{})
}
