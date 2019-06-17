package bufs

import "bytes"

// Provides leaky buffer, based on the example in Effective Go.

type BufferPool struct {
	bufSize  int // size of each buffer
	freeList chan *Buffer
}

const (
	BUF_SIZE      = 1024
	BUF_POOL_SIZE = 512
)

type Buffer struct {
	*bytes.Buffer
	refCount int
}

func (bp *BufferPool) newBuffer() *Buffer {
	buf := Buffer{bytes.NewBuffer(make([]byte, bp.bufSize)), 0}
	return &buf
}

func NewBufferPool(n, bufSize int) *BufferPool {
	if n <= 0 || bufSize <= 0 {
		return &BufferPool{bufSize: BUF_SIZE, freeList: make(chan *Buffer, BUF_POOL_SIZE)}
	}
	return &BufferPool{bufSize: bufSize, freeList: make(chan *Buffer, n)}
}

func (bp *BufferPool) Get() (b *Buffer) {
	select {
	case b = <-bp.freeList:
		b.refCount++
	default:
		b = bp.newBuffer()
	}
	return
}

func (bp *BufferPool) Put(b *Buffer) {
	b.refCount--
	b.Reset()
	if b.refCount == 0 {
		bp.freeList <- b
	}
}

/*
//Default parameter ， Globally managed
const leakyBufSize = 20<<1 // data.len(2) + hmacsha1(10) + data(4096)
const maxNBuf = 2048

// NewLeakyBuf creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewLeakyBuf(n, bufSize int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan *bytes.Buffer, n),
	}
}


// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *LeakyBuf) Get() (b *bytes.Buffer ){
	select {
	case b = <-lb.freeList:
	default:
		b=
	}
	return
}

// Put add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func (lb *LeakyBuf) Put(b []byte) {
	if len(b) != lb.bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	select {
	case lb.freeList <- b:
	default:
	}
	return
}
*/
