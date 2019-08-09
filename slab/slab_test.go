package slab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_allocate(t *testing.T) {
	var s *Slab = NewSlab(1024 * 1024 * 4)

	defer func() {
		assert.Equal(t, int64(0), s.Inflights())
	}()

	buf1, dealloc1 := s.Allocate()
	defer dealloc1(buf1)

	assert.Equal(t, 1024*1024*4, len(*buf1))
	assert.Equal(t, 1024*1024*4, cap(*buf1))
	assert.Equal(t, int64(1), s.Inflights())

	buf2, dealloc2 := s.Allocate()
	defer dealloc2(buf2)

	assert.Equal(t, 1024*1024*4, len(*buf2))
	assert.Equal(t, 1024*1024*4, cap(*buf2))
	assert.Equal(t, int64(2), s.Inflights())

	buf3, dealloc3 := s.Allocate()
	assert.Equal(t, 1024*1024*4, len(*buf3))
	assert.Equal(t, 1024*1024*4, cap(*buf3))
	assert.Equal(t, int64(3), s.Inflights())

	dealloc3(buf3)
	assert.Equal(t, int64(2), s.Inflights())
}
