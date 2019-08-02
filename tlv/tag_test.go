package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tag(t *testing.T) {
	tag1 := NestingTag(123)
	tag2 := RawTag(456)

	assert.True(t, IsNestingTag(tag1))
	assert.False(t, IsRawTag(tag1))

	assert.False(t, IsNestingTag(tag2))
	assert.True(t, IsRawTag(tag2))
}
