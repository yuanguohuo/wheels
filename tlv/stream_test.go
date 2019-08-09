package tlv

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleHandler struct {
	root *TLV
}

func (h *simpleHandler) OnFailure(err error) {
	log.Panicf("OnFailure. error:%s", err)
}
func (h *simpleHandler) OnSuccess() {
	log.Printf("OnSuccess")
}
func (h *simpleHandler) OnBegin(tlv *TLV) {
	log.Printf("OnBegin")
	if h.root == nil {
		h.root = tlv
	}
}
func (h *simpleHandler) OnEnd(tlv *TLV) {
	log.Printf("OnBegin")
}

func Test_ParseStream(t *testing.T) {
	buf := make([]byte, 1024)
	var err error

	rootTag := NestingTag(0)

	son1Tag := NestingTag(1)
	son11Tag := NestingTag(11)
	son111Tag := RawTag(111)
	son112Tag := RawTag(112)
	son113Tag := RawTag(113)
	son12Tag := RawTag(12)

	son2Tag := NestingTag(2)

	son21Tag := RawTag(21)

	son22Tag := NestingTag(22)
	son221Tag := RawTag(221)
	son222Tag := RawTag(222)
	son223Tag := RawTag(223)

	son23Tag := NestingTag(23)
	son231Tag := NestingTag(231)
	son2311Tag := RawTag(2311)

	//
	//                                          Root 156+8
	//                                           |
	//                       +-------------------+-----------------------+
	//                       |                                           |
	//                     Son1 72                                     Son2  84
	//                       |                                           |
	//             +---------+----------+                +---------------+-----------------+
	//             |                    |                |               |                 |
	//           Son11 44              Son12 20        Son21 8         Son22 44          Son23  24
	//             |                                                     |                 |
	//   +---------+-------+                                    +--------+---------+       |
	//   |         |       |                                    |        |         |       |
	// Son111   Son112  Son113                               Son221   Son222     Son223  Son231  16
	//   12       16      8                                     8        12        16      |
	//                                                                                   Son2311 8

	root, err := Open(buf, nil, rootTag)
	assert.Nil(t, err)

	son1, err := root.OpenChild(son1Tag)
	assert.Nil(t, err)

	son11, err := son1.AllocateChild(son11Tag, 12+16+8)
	assert.Nil(t, err)

	son111, err := son11.BuildChildInt32(son111Tag, 111)
	assert.Nil(t, err)
	assert.Equal(t, int32(111), son111.ValAsInt32())
	son112, err := son11.BuildChildInt64(son112Tag, 112)
	assert.Nil(t, err)
	assert.Equal(t, int64(112), son112.ValAsInt64())
	son113, err := son11.BuildChild(son113Tag, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(son113.Val()))

	assert.Nil(t, son11.Close())

	son12, err := son1.OpenChild(son12Tag)
	assert.Nil(t, son12.FillRawInt32(121))
	assert.Nil(t, son12.FillRawInt64(122))
	assert.Nil(t, son12.Close())
	assert.Equal(t, 12, len(son12.Val()))
	assert.Equal(t, 20, len(son12.body))

	assert.Nil(t, son1.Close())

	son2, err := root.OpenChild(son2Tag)
	assert.Nil(t, err)

	son21, err := son2.BuildChild(son21Tag, []byte(""))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(son21.Val()))
	assert.Equal(t, 8, len(son21.body))

	son22, err := son2.OpenChild(son22Tag)

	son221, err := son22.OpenChild(son221Tag)
	assert.Nil(t, son221.FillRaw([]byte("")))
	assert.Nil(t, son221.Close())

	son222, err := son22.BuildChildInt32(son222Tag, 222)
	assert.Nil(t, err)
	assert.Equal(t, int32(222), son222.ValAsInt32())

	son223, err := son22.AllocateChild(son223Tag, 8)
	assert.Nil(t, err)
	assert.Nil(t, son223.FillRawInt64(223))
	assert.Nil(t, son223.Close())

	assert.Nil(t, son22.Close())

	son23, err := son2.OpenChild(son23Tag)
	assert.Nil(t, err)

	son231, err := son23.OpenChild(son231Tag)
	assert.Nil(t, err)

	son2311, err := son231.OpenChild(son2311Tag)
	assert.Nil(t, err)

	assert.Nil(t, son2311.Close())
	assert.Nil(t, son231.Close())
	assert.Nil(t, son23.Close())

	assert.Nil(t, son2.Close())
	assert.Nil(t, root.Close())

	//////////////////////////
	//    hand built        //
	//////////////////////////

	var rootValLen int = 156
	handBuilt := make([]byte, rootValLen+8)
	var i int = 0

	encoder.PutUint32(handBuilt[i:], uint32(rootTag))
	i += 4
	encoder.PutUint32(handBuilt[i:], uint32(rootValLen))
	i += 4

	var son1ValLen int = 64
	encoder.PutUint32(handBuilt[i:], uint32(son1Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], uint32(son1ValLen))
	i += 4

	var son11ValLen int = 12 + 16 + 8
	encoder.PutUint32(handBuilt[i:], uint32(son11Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], uint32(son11ValLen))
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son111Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 4)
	i += 4
	encoder.PutUint32(handBuilt[i:], 111)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son112Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 8)
	i += 4
	encoder.PutUint64(handBuilt[i:], 112)
	i += 8

	encoder.PutUint32(handBuilt[i:], uint32(son113Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 0)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son12Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 12)
	i += 4
	encoder.PutUint32(handBuilt[i:], 121)
	i += 4
	encoder.PutUint64(handBuilt[i:], 122)
	i += 8

	var son2ValLen int = 76
	encoder.PutUint32(handBuilt[i:], uint32(son2Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], uint32(son2ValLen))
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son21Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 0)
	i += 4

	var son22ValLen int = 8 + 12 + 16
	encoder.PutUint32(handBuilt[i:], uint32(son22Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], uint32(son22ValLen))
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son221Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 0)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son222Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 4)
	i += 4
	encoder.PutUint32(handBuilt[i:], 222)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son223Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 8)
	i += 4
	encoder.PutUint64(handBuilt[i:], 223)
	i += 8

	encoder.PutUint32(handBuilt[i:], uint32(son23Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 16)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son231Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 8)
	i += 4

	encoder.PutUint32(handBuilt[i:], uint32(son2311Tag))
	i += 4
	encoder.PutUint32(handBuilt[i:], 0)
	i += 4

	assert.Equal(t, 164, i)

	assert.True(t, bytes.Equal(handBuilt, root.body))

	//////////////////////////
	//    parse stream      //
	//////////////////////////
	newBuf := make([]byte, 1024)
	reader := bytes.NewReader(handBuilt)
	cb := new(simpleHandler)

	assert.Nil(t, ParseStream(newBuf, reader, cb))
	parsed := cb.root
	assert.True(t, bytes.Equal(parsed.body, root.body))
}
