package tlv

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertOpen(t *testing.T, tlv *TLV, tag TLVTag) {
	assert.Equal(t, tag, tlv.Tag())
	assert.Equal(t, TLVState_Open, tlv.state)
}

func assertAllocated(t *testing.T, tlv *TLV, tag TLVTag, l int) {
	assert.Equal(t, tag, tlv.Tag())
	assert.Equal(t, l, tlv.Len())
	assert.Equal(t, TAG_SIZE+LEN_SIZE+l, len(tlv.body))

	assert.Equal(t, TLVState_Allocated, tlv.state)
}

func assertClosed(t *testing.T, tlv *TLV, tag TLVTag, l int, v []byte) {
	assert.Equal(t, tag, tlv.Tag())
	assert.Equal(t, l, tlv.Len())
	assert.Equal(t, TAG_SIZE+LEN_SIZE+l, len(tlv.body))

	assert.Equal(t, TLVState_Closed, tlv.state)
	assert.Equal(t, l, len(tlv.Val()))
	assert.Equal(t, VAL_OFFSET+l, tlv.curr)
	if len(v) == 0 {
		assert.Equal(t, 0, len(tlv.Val()))
	}
	if len(v) > 0 {
		assert.True(t, bytes.Equal(v, tlv.Val()))
	}
}

func Test_Build_1(t *testing.T) {
	buf := make([]byte, 64)

	var tlv *TLV
	var err error
	tag := RawTag(234)

	//empty value
	tlv, err = Build(buf, nil, tag, nil)
	assert.Nil(t, err)
	assertClosed(t, tlv, tag, 0, nil)

	tagbinary := make([]byte, 4)
	lenbinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(tagbinary, uint32(tag))
	binary.LittleEndian.PutUint32(lenbinary, 0)
	assert.True(t, bytes.Equal(buf[0:4], tagbinary))
	assert.True(t, bytes.Equal(buf[4:8], lenbinary))

	buf = buf[32:]
	tag++

	tlv, err = Build(buf, nil, tag, []byte(""))
	assert.Nil(t, err)
	assertClosed(t, tlv, tag, 0, []byte(""))

	tagbinary = make([]byte, 4)
	lenbinary = make([]byte, 4)
	binary.LittleEndian.PutUint32(tagbinary, uint32(tag))
	binary.LittleEndian.PutUint32(lenbinary, 0)
	assert.True(t, bytes.Equal(buf[0:4], tagbinary))
	assert.True(t, bytes.Equal(buf[4:8], lenbinary))
}

func Test_Build_2(t *testing.T) {
	buf := make([]byte, 21)

	var tlv *TLV
	var err error
	tag := RawTag(235)

	//1. create would overflow
	tlv, err = Build(buf, nil, tag, make([]byte, 14)) //14+TAG_SIZE+LEN_SIZE > 21
	assert.Equal(t, Err_WouldOverflow, err)
	assert.Nil(t, tlv)

	//2. create
	tlv, err = Build(buf, nil, tag, []byte("AABBCCDDEEFFG"))
	assert.Nil(t, err)
	assertClosed(t, tlv, tag, 13, []byte("AABBCCDDEEFFG"))

	bin := make([]byte, 21)
	binary.LittleEndian.PutUint32(bin[0:4], uint32(tag))
	binary.LittleEndian.PutUint32(bin[4:8], 13)
	copy(bin[8:], []byte("AABBCCDDEEFFG"))
	assert.True(t, bytes.Equal(buf, bin))
}

func Test_Allocate_1(t *testing.T) {
	buf := make([]byte, 64)

	var tlv *TLV
	var err error
	tag := RawTag(123)

	//empty TLV is valid
	tlv, err = Allocate(buf, nil, tag, 0)
	assert.Nil(t, err)
	assertAllocated(t, tlv, tag, 0)

	err = tlv.FillRaw([]byte(""))
	assert.Nil(t, err)

	err = tlv.FillRaw([]byte("A"))
	assert.Equal(t, Err_WouldOverflow, err)

	err = tlv.Close()
	assert.Nil(t, err)
	assertClosed(t, tlv, tag, 0, nil)

	tagbinary := make([]byte, 4)
	lenbinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(tagbinary, uint32(tag))
	binary.LittleEndian.PutUint32(lenbinary, 0)
	assert.True(t, bytes.Equal(buf[0:4], tagbinary))
	assert.True(t, bytes.Equal(buf[4:8], lenbinary))
}

func Test_Allocate_2(t *testing.T) {
	buf := make([]byte, 64)

	var tlv *TLV
	var err error
	tag := RawTag(456)

	//1. create would overflow
	tlv, err = Allocate(buf, nil, tag, 57) //57+TAG_SIZE+LEN_SIZE > 64
	assert.Equal(t, Err_WouldOverflow, err)
	assert.Nil(t, tlv)

	//2. create
	tlv, err = Allocate(buf, nil, tag, 16)
	assert.Nil(t, err)
	assertAllocated(t, tlv, tag, 16)
	assert.Equal(t, VAL_OFFSET, tlv.curr)
	assert.Equal(t, tlv.curr+tlv.Len(), len(tlv.body))

	//3. append

	//3.1 append child
	childTag := RawTag(457)
	child, err := tlv.AllocateChild(childTag, 1)
	assert.Nil(t, child)
	assert.Equal(t, Err_UnexpectedTag, err)

	//3.1 overflow
	err = tlv.FillRaw(make([]byte, 17))
	assert.Equal(t, Err_WouldOverflow, err)
	assert.Equal(t, VAL_OFFSET, tlv.curr)
	assert.Equal(t, tlv.curr+tlv.Len(), len(tlv.body))

	//3.2 append 3 bytes, 3 bytes now
	err = tlv.FillRaw([]byte("AAA"))
	assert.Nil(t, err)
	assert.Equal(t, VAL_OFFSET+3, tlv.curr)

	//3.3 append another 9 bytes, 12 bytes now
	err = tlv.FillRaw([]byte("BBBBBBBBB"))
	assert.Nil(t, err)
	assert.Equal(t, VAL_OFFSET+3+9, tlv.curr)

	//3.4 call Close early: the room is 16 bytes, we filled only 12 bytes, so Close will fail
	err = tlv.Close()
	assert.Equal(t, Err_Corrupted, err)

	//3.5 append another 5 bytes , overflow
	err = tlv.FillRaw(make([]byte, 5))
	assert.Equal(t, Err_WouldOverflow, err)
	assert.Equal(t, 16, tlv.Len())
	assert.Equal(t, VAL_OFFSET+3+9, tlv.curr)

	//3.6 append empty
	err = tlv.FillRaw([]byte(""))
	assert.Nil(t, err)

	//3.7 append nil
	err = tlv.FillRaw(nil)
	assert.Nil(t, err)

	//3.8 append another 4 bytes, full now
	err = tlv.FillRaw([]byte("CCCC"))
	assert.Nil(t, err)
	assert.Equal(t, VAL_OFFSET+3+9+4, tlv.curr)

	//3.9 append another 1 byte , overflow
	err = tlv.FillRaw(make([]byte, 1))
	assert.Equal(t, Err_WouldOverflow, err)

	//3.10 append empty
	err = tlv.FillRaw([]byte(""))
	assert.Nil(t, err)

	//3.11 append nil
	err = tlv.FillRaw(nil)
	assert.Nil(t, err)

	//3.12 Close
	err = tlv.Close()
	assert.Nil(t, err)
	assertClosed(t, tlv, tag, 16, []byte("AAABBBBBBBBBCCCC"))

	//4 check bytes
	val := tlv.Val()
	assert.Equal(t, "AAABBBBBBBBBCCCC", string(val))

	tagbinary := make([]byte, 4)
	lenbinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(tagbinary, uint32(tag))
	binary.LittleEndian.PutUint32(lenbinary, 16)
	assert.True(t, bytes.Equal(buf[0:4], tagbinary))
	assert.True(t, bytes.Equal(buf[4:8], lenbinary))
	assert.True(t, bytes.Equal(buf[8:24], []byte("AAABBBBBBBBBCCCC")))
}

func Test_NestedTLV_1(t *testing.T) {
	TOTAL := 100
	buf := make([]byte, TOTAL)

	var root *TLV
	var err error
	rootTag := NestingTag(111)

	root, err = Open(buf, nil, rootTag)
	assert.Nil(t, err)
	assertOpen(t, root, rootTag)
	assert.Equal(t, VAL_OFFSET, root.curr)
	assert.Equal(t, TOTAL, len(root.body)) //root is open sized now (taking all size of buf); will be fixed when root.Close is called

	//1. append son1: fixed len
	SON1_START_POS_IN_BUF := root.curr
	assert.Equal(t, 8, SON1_START_POS_IN_BUF)

	son1Tag := RawTag(222)
	son1, err := root.AllocateChild(son1Tag, 7)
	assert.Nil(t, err)
	assertAllocated(t, son1, son1Tag, 7)
	assert.Equal(t, VAL_OFFSET, son1.curr)
	//root has increased its curr
	assert.Equal(t, VAL_OFFSET+len(son1.body), root.curr)

	//2. append son2: var len
	SON2_START_POS_IN_BUF := root.curr
	assert.Equal(t, 8+15, SON2_START_POS_IN_BUF) //15 is space of son1

	son2Tag := RawTag(333)
	son2, err := root.OpenChild(son2Tag)
	assert.Nil(t, err)
	assertOpen(t, son2, son2Tag)

	//son2 is open sized now (taking all size from parent's cur); will be fixed when root.Close is called
	assert.Equal(t, TOTAL-SON2_START_POS_IN_BUF, len(son2.body))
	//root has NOT increased its curr
	assert.Equal(t, VAL_OFFSET+len(son1.body), root.curr)

	err = son2.FillRaw([]byte("AAAA"))
	assert.Nil(t, err)
	assert.Equal(t, VAL_OFFSET+4, son2.curr)
	err = son2.FillRaw([]byte("BBBB"))
	assert.Nil(t, err)
	assert.Equal(t, VAL_OFFSET+4+4, son2.curr)

	err = son2.Close()
	assert.Nil(t, err)
	assertClosed(t, son2, son2Tag, 8, []byte("AAAABBBB"))

	//son2 is NOT open sized now
	assert.Equal(t, TAG_SIZE+LEN_SIZE+8, len(son2.body))
	//root has increased its curr
	assert.Equal(t, VAL_OFFSET+len(son1.body)+len(son2.body), root.curr) //parent has increased its curr

	//  +----------------------+
	//  |    root TAG          | 4
	//  |    root LEN          | 4
	//  +----------------------+     <-------------  SON1_START_POS_IN_BUF = 8
	//  |    son1 TAG          | 4
	//  |    son1 LEN          | 4
	//  +----------------------+
	//  |                      |
	//  |    HHHHHHH           | 7
	//  |                      |
	//  +----------------------+     <-------------  SON2_START_POS_IN_BUF = 23
	//  |    son2 TAG          | 4
	//  |    son2 LEN          | 4
	//  +----------------------+
	//  |    AAAA              | 8
	//  |    BBBB              |
	//  +----------------------+     <-------------  SON3_START_POS_IN_BUF = 39
	//  |    son3 TAG          | 4
	//  |    son3 LEN          | 4
	//  +----------------------+
	//  |  son1 of son3 TAG    | 4
	//  |  son1 of son3 LEN    | 4
	//  +----------------------+
	//  |   A                  | 1
	//  +----------------------+
	//  |  son2 of son3 TAG    | 4
	//  |  son2 of son3 LEN    | 4
	//  +----------------------+
	//  |   CCCC               | 6
	//  |   DD                 |
	//  +----------------------+
	//  |  son3 of son3 TAG    | 4
	//  |  son3 of son3 LEN    | 4
	//  +----------------------+
	//  |   X                  | 1
	//  +----------------------+
	//  |  son4 of son3 TAG    | 4
	//  |  son4 of son3 LEN    | 4   val is empty
	//  +----------------------+     <------------- SON4_START_POS_IN_BUF = 87
	//  |    son4 TAG          | 4
	//  |    son4 LEN          | 4
	//  +----------------------+
	//  |                      |
	//  |   ABCDE              | 5
	//  |                      |
	//  +----------------------+     <------------- 100

	//3. append son3: var len
	SON3_START_POS_IN_BUF := root.curr
	assert.Equal(t, 8+15+16, SON3_START_POS_IN_BUF) //15 is space of son1; 16 is space of son2;

	son3Tag := NestingTag(444)
	son3, err := root.OpenChild(son3Tag)
	assert.Nil(t, err)
	assertOpen(t, son3, son3Tag)

	//3.1 append son31 to son3
	son31Tag := RawTag(4441)
	son31, err := son3.OpenChild(son31Tag)
	assert.Nil(t, err)
	assertOpen(t, son31, son31Tag)

	//son3 has NOT increased its curr
	assert.Equal(t, VAL_OFFSET, son3.curr)
	//root has NOT increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF, root.curr)

	err = son31.FillRaw([]byte("A"))
	assert.Nil(t, err)

	err = son31.Close()
	assert.Nil(t, err)
	assertClosed(t, son31, son31Tag, 1, []byte("A"))

	//son3 has increased its curr
	assert.Equal(t, VAL_OFFSET+9, son3.curr) //9 is the space of son31
	//root has NOT increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF, root.curr)

	//3.2 append son32 to son3
	son32Tag := RawTag(4442)
	son32, err := son3.BuildChild(son32Tag, []byte("CCCCDD"))
	assert.Nil(t, err)
	assertClosed(t, son32, son32Tag, 6, []byte("CCCCDD"))

	//son3 has increased its curr
	assert.Equal(t, VAL_OFFSET+9+14, son3.curr) //9 is the space of son31; 14 is the space of son32
	//root has NOT increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF, root.curr)

	//3.3 append son33 to son3
	son33Tag := RawTag(4443)
	son33, err := son3.AllocateChild(son33Tag, 1)
	assert.Nil(t, err)
	assertAllocated(t, son33, son33Tag, 1)

	err = son33.FillRaw([]byte("X"))
	assert.Nil(t, err)

	err = son33.Close()
	assert.Nil(t, err)
	assertClosed(t, son33, son33Tag, 1, []byte("X"))

	//son3 has increased its curr
	assert.Equal(t, VAL_OFFSET+9+14+9, son3.curr) //9 is the space of son31; 14 is the space of son32; 9 is the space of son33
	//root has NOT increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF, root.curr)

	//3.4 append son34 to son3
	son34Tag := RawTag(4444)
	son34, err := son3.AllocateChild(son34Tag, 0)
	assert.Nil(t, err)
	assertAllocated(t, son34, son34Tag, 0)

	err = son34.Close()
	assert.Nil(t, err)
	assertClosed(t, son34, son34Tag, 0, nil)

	//son3 has increased its curr
	assert.Equal(t, VAL_OFFSET+9+14+9+8, son3.curr) //9 is the space of son31; 14 is the space of son32; 9 is the space of son33; 8 is the space of son34
	//root has NOT increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF, root.curr)

	//3.5 son3 Close
	err = son3.Close()
	assert.Nil(t, err)

	son3ValBin := make([]byte, 40)
	var i int = 0
	binary.LittleEndian.PutUint32(son3ValBin[i:], uint32(son31Tag))
	i += 4
	binary.LittleEndian.PutUint32(son3ValBin[i:], 1)
	i += 4
	copy(son3ValBin[i:], []byte("A"))
	i += 1

	binary.LittleEndian.PutUint32(son3ValBin[i:], uint32(son32Tag))
	i += 4
	binary.LittleEndian.PutUint32(son3ValBin[i:], 6)
	i += 4
	copy(son3ValBin[i:], []byte("CCCCDD"))
	i += 6

	binary.LittleEndian.PutUint32(son3ValBin[i:], uint32(son33Tag))
	i += 4
	binary.LittleEndian.PutUint32(son3ValBin[i:], 1)
	i += 4
	copy(son3ValBin[i:], []byte("X"))
	i += 1

	binary.LittleEndian.PutUint32(son3ValBin[i:], uint32(son34Tag))
	i += 4
	binary.LittleEndian.PutUint32(son3ValBin[i:], 0)
	i += 4

	assertClosed(t, son3, son3Tag, 40, son3ValBin)
	//root has increased its curr
	assert.Equal(t, SON3_START_POS_IN_BUF+48, root.curr) //48 is the total space of son3

	//4. try to append son4
	SON4_START_POS_IN_BUF := root.curr
	SPACE_LEFT_IN_BUF := TOTAL - SON4_START_POS_IN_BUF

	son4Tag := RawTag(555)
	son4, err := root.AllocateChild(son4Tag, SPACE_LEFT_IN_BUF-TAG_SIZE-LEN_SIZE+1)
	assert.Equal(t, Err_WouldOverflow, err)
	assert.Nil(t, son4)

	son4, err = root.AllocateChild(son4Tag, SPACE_LEFT_IN_BUF-TAG_SIZE-LEN_SIZE)
	assert.Nil(t, err)
	assertAllocated(t, son4, son4Tag, SPACE_LEFT_IN_BUF-TAG_SIZE-LEN_SIZE)

	//root has increased its curr
	assert.Equal(t, SON4_START_POS_IN_BUF+13, root.curr) //13 is the total space of son4
	assert.Equal(t, TOTAL, root.curr)                    //13 is the total space of son4

	err = son4.FillRaw(make([]byte, SPACE_LEFT_IN_BUF-TAG_SIZE-LEN_SIZE+1))
	assert.Equal(t, Err_WouldOverflow, err)

	err = son4.FillRaw([]byte("ABCD"))
	assert.Nil(t, err)

	err = son4.Close()
	assert.Equal(t, Err_Corrupted, err)

	err = son4.FillRaw([]byte("E"))
	assert.Nil(t, err)

	err = son4.Close()
	assert.Nil(t, err)
	assertClosed(t, son4, son4Tag, 5, []byte("ABCDE"))

	//5. son1 has not filled with value yet
	assertAllocated(t, son1, son1Tag, 7)
	err = son1.Close()
	assert.Equal(t, Err_Corrupted, err)
	err = son1.FillRaw([]byte("HHHHHHH"))
	assert.Nil(t, err)
	err = son1.Close()
	assert.Nil(t, err)
	assertClosed(t, son1, son1Tag, 7, []byte("HHHHHHH"))

	//6. root Close
	assertOpen(t, root, rootTag)
	err = root.Close()
	assert.Nil(t, err)
	assert.Equal(t, TLVState_Closed, root.state)
	assert.Equal(t, TOTAL, len(root.body))
	assert.Equal(t, TOTAL, root.curr)

	//7. check children
	assert.Equal(t, 4, len(root.children))

	s1, ok := root.children[son1Tag]
	assert.True(t, ok)
	assert.True(t, son1.Equal(s1[0]))

	s2, ok := root.children[son2Tag]
	assert.True(t, ok)
	assert.True(t, son2.Equal(s2[0]))

	s3, ok := root.children[son3Tag]
	assert.True(t, ok)
	assert.True(t, son3.Equal(s3[0]))

	s31, ok := son3.children[son31Tag]
	assert.True(t, ok)
	assert.True(t, son31.Equal(s31[0]))

	s32, ok := son3.children[son32Tag]
	assert.True(t, ok)
	assert.True(t, son32.Equal(s32[0]))

	s33, ok := son3.children[son33Tag]
	assert.True(t, ok)
	assert.True(t, son33.Equal(s33[0]))

	s34, ok := son3.children[son34Tag]
	assert.True(t, ok)
	assert.True(t, son34.Equal(s34[0]))

	s4, ok := root.children[son4Tag]
	assert.True(t, ok)
	assert.True(t, son4.Equal(s4[0]))

	//8. final checks
	allTLVs := []*TLV{root, son1, son2, son3, son31, son32, son33, son34, son4}
	for _, tx := range allTLVs {
		assert.Equal(t, TLVState_Closed, tx.state)
	}

	bin := make([]byte, TOTAL)
	i = 0
	binary.LittleEndian.PutUint32(bin[i:], uint32(rootTag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], uint32(TOTAL-8)) //VAL of root is TOTAL-8 bytes
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], uint32(son1Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 7) //VAL os son1 is 7 bytes
	i += 4
	copy(bin[i:], []byte("HHHHHHH"))
	i += 7

	binary.LittleEndian.PutUint32(bin[i:], uint32(son2Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 8) //VAL of son2 is 8 bytes;
	i += 4
	copy(bin[i:], []byte("AAAABBBB"))
	i += 8

	binary.LittleEndian.PutUint32(bin[i:], uint32(son3Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 40) //VAL of son3 is 40 bytes;
	i += 4

	binary.LittleEndian.PutUint32(bin[i:], uint32(son31Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 1) //VAL of son31 is 1 byte;
	i += 4
	copy(bin[i:], []byte("A"))
	i += 1

	binary.LittleEndian.PutUint32(bin[i:], uint32(son32Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 6) //VAL of son32 is 6 bytes;
	i += 4
	copy(bin[i:], []byte("CCCCDD"))
	i += 6

	binary.LittleEndian.PutUint32(bin[i:], uint32(son33Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 1) //VAL of son33 is 1 byte;
	i += 4
	copy(bin[i:], []byte("X"))
	i += 1

	binary.LittleEndian.PutUint32(bin[i:], uint32(son34Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 0) //VAL of son34 is 0 byte;
	i += 4

	binary.LittleEndian.PutUint32(bin[i:], uint32(son4Tag))
	i += 4
	binary.LittleEndian.PutUint32(bin[i:], 5) //VAL of son34 is 5 bytes;
	i += 4
	copy(bin[i:], []byte("ABCDE"))
	i += 5

	assert.Equal(t, TOTAL, i)
	assert.True(t, bytes.Equal(bin, buf))
	assert.True(t, bytes.Equal(bin, root.body))

	//9. marshal and unmarsha
	flat, err := Marshal(root)
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(bin, flat))

	tlvs, err := Unmarshal(bin, nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tlvs[rootTag]))

	root_m := tlvs[rootTag][0]
	son1_m := root_m.children[son1Tag][0]
	son2_m := root_m.children[son2Tag][0]

	son3_m := root_m.children[son3Tag][0]
	son31_m := son3_m.children[son31Tag][0]
	son32_m := son3_m.children[son32Tag][0]
	son33_m := son3_m.children[son33Tag][0]
	son34_m := son3_m.children[son34Tag][0]

	son4_m := root_m.children[son4Tag][0]

	assert.Equal(t, root.Tag(), root_m.Tag())
	assert.Equal(t, root.Len(), root_m.Len())
	assert.True(t, bytes.Equal(root.Val(), root_m.Val()))

	assert.Equal(t, son1.Tag(), son1_m.Tag())
	assert.Equal(t, son1.Len(), son1_m.Len())
	assert.True(t, bytes.Equal(son1.Val(), son1_m.Val()))

	assert.Equal(t, son2.Tag(), son2_m.Tag())
	assert.Equal(t, son2.Len(), son2_m.Len())
	assert.True(t, bytes.Equal(son2.Val(), son2_m.Val()))

	assert.Equal(t, son3.Tag(), son3_m.Tag())
	assert.Equal(t, son3.Len(), son3_m.Len())
	assert.True(t, bytes.Equal(son3.Val(), son3_m.Val()))

	assert.Equal(t, son31.Tag(), son31_m.Tag())
	assert.Equal(t, son31.Len(), son31_m.Len())
	assert.True(t, bytes.Equal(son31.Val(), son31_m.Val()))

	assert.Equal(t, son32.Tag(), son32_m.Tag())
	assert.Equal(t, son32.Len(), son32_m.Len())
	assert.True(t, bytes.Equal(son32.Val(), son32_m.Val()))

	assert.Equal(t, son33.Tag(), son33_m.Tag())
	assert.Equal(t, son33.Len(), son33_m.Len())
	assert.True(t, bytes.Equal(son33.Val(), son33_m.Val()))

	assert.Equal(t, son34.Tag(), son34_m.Tag())
	assert.Equal(t, son34.Len(), son34_m.Len())
	assert.True(t, bytes.Equal(son34.Val(), son34_m.Val()))

	assert.Equal(t, son4.Tag(), son4_m.Tag())
	assert.Equal(t, son4.Len(), son4_m.Len())
	assert.True(t, bytes.Equal(son4.Val(), son4_m.Val()))
}
