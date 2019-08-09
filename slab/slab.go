package slab

import (
	"sync"
	"sync/atomic"
)

type Deallocator func(*[]byte)

type Slab struct {
	size      int
	inflights int64

	pool    sync.Pool
	dealloc Deallocator
}

func NewSlab(size int) *Slab {
	s := &Slab{
		size:      size,
		inflights: 0,
	}
	s.pool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, s.size)
			return &buf
		},
	}
	s.dealloc = func(buf *[]byte) {
		*buf = (*buf)[0:size]
		s.pool.Put(buf)
		atomic.AddInt64(&s.inflights, -1)
	}
	return s
}

func (s *Slab) Allocate() (*[]byte, Deallocator) {
	atomic.AddInt64(&s.inflights, 1)
	return s.pool.Get().(*[]byte), s.dealloc
}

func (s *Slab) Inflights() int64 {
	return atomic.LoadInt64(&s.inflights)
}
