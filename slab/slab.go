package slab

import (
	"sync"
	"sync/atomic"
)

type Deallocator func(*[]byte)

type Slab struct {
	size     int
	inflight int64

	pool    sync.Pool
	dealloc Deallocator
}

func NewSlab(size int) *Slab {
	s := &Slab{
		size:     size,
		inflight: 0,
	}
	s.pool = sync.Pool{
		New: func() interface{} {
			atomic.AddInt64(&s.inflight, 1)
			buf := make([]byte, s.size)
			return &buf
		},
	}
	s.dealloc = func(buf *[]byte) {
		*buf = (*buf)[0:0]
		s.pool.Put(buf)
		atomic.AddInt64(&s.inflight, -1)
	}
	return s
}

func (s *Slab) Allocate() (*[]byte, Deallocator) {
	return s.pool.Get().(*[]byte), s.dealloc
}

func (s *Slab) Inflight() int64 {
	return atomic.LoadInt64(&s.inflight)
}
