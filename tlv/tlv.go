package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//var encoder = binary.BigEndian
var encoder = binary.LittleEndian

var (
	Err_WouldOverflow   error = errors.New("buffer would overflow")
	Err_Corrupted       error = errors.New("tlv is corrupted")
	Err_UnexpectedTag   error = errors.New("unexpected tag")
	Err_UnexpectedState error = errors.New("unexpected state")
	Err_HasOpenChild    error = errors.New("already has open child")
)

const (
	TAG_OFFSET = 0
	TAG_SIZE   = 4

	LEN_OFFSET = 4
	LEN_SIZE   = 4

	VAL_OFFSET = 8
)

type TLVStateEnum int

const (
	TLVState_Open      TLVStateEnum = 1
	TLVState_Allocated TLVStateEnum = 2
	TLVState_Closed    TLVStateEnum = 3
)

type TLV struct {
	parent       *TLV
	body         []byte
	curr         int
	state        TLVStateEnum
	hasOpenChild bool
	children     map[TLVTag][]*TLV
}

func Open(buf []byte, p *TLV, t TLVTag) (*TLV, error) {
	if TAG_SIZE+LEN_SIZE > len(buf) {
		return nil, Err_WouldOverflow
	}
	tlv := &TLV{
		parent:       p,
		body:         buf,
		curr:         VAL_OFFSET,
		state:        TLVState_Open,
		hasOpenChild: false,
	}
	if IsNestingTag(t) {
		tlv.children = make(map[TLVTag][]*TLV)
	}
	if p != nil {
		p.mustAddChild(t, tlv)
	}
	tlv.mustSetTag(t)
	return tlv, nil
}

func Allocate(buf []byte, p *TLV, t TLVTag, l int) (*TLV, error) {
	if TAG_SIZE+LEN_SIZE+l > len(buf) {
		return nil, Err_WouldOverflow
	}
	tlv, err := Open(buf, p, t)
	if err != nil {
		return nil, err
	}
	if err := tlv.allocate(l); err != nil {
		return nil, err
	}
	return tlv, nil
}

func Build(buf []byte, p *TLV, t TLVTag, v []byte) (*TLV, error) {
	tlv, err := Allocate(buf, p, t, len(v))
	if err != nil {
		return nil, err
	}
	if len(v) > 0 {
		tlv.mustFillRaw(v)
	}
	if err := tlv.close(len(v)); err != nil {
		return nil, err
	}
	return tlv, nil
}

func (tlv *TLV) Close() error {
	if tlv.hasOpenChild {
		return Err_HasOpenChild
	}

	valLen := tlv.curr - VAL_OFFSET
	if tlv.state == TLVState_Open {
		if err := tlv.allocate(valLen); err != nil {
			return err
		}
	}
	if err := tlv.close(valLen); err != nil {
		return err
	}
	return nil
}

//fill in raw bytes as VALUE of 'this';
//this.curr is increased by len(v) correspondently;
func (tlv *TLV) FillRaw(v []byte) error {
	if tlv.state != TLVState_Allocated && tlv.state != TLVState_Open {
		return Err_UnexpectedState
	}
	if !IsRawTag(tlv.Tag()) {
		return Err_UnexpectedTag
	}
	if len(v) == 0 {
		return nil
	}
	if tlv.curr+len(v) > len(tlv.body) {
		return Err_WouldOverflow
	}
	tlv.mustFillRaw(v)
	return nil
}

func (tlv *TLV) FillRawInt32(v int32) error {
	b := make([]byte, 4)
	encoder.PutUint32(b, uint32(v))
	return tlv.FillRaw(b)
}

func (tlv *TLV) FillRawInt64(v int64) error {
	b := make([]byte, 8)
	encoder.PutUint64(b, uint64(v))
	return tlv.FillRaw(b)
}

//open a nested TLV (child) as VALUE of 'this';
//this.curr will not be increased until child.Close (see allocateForOpenChild)
func (tlv *TLV) OpenChild(t TLVTag) (*TLV, error) {
	if tlv.state != TLVState_Allocated && tlv.state != TLVState_Open {
		return nil, Err_UnexpectedState
	}
	if !IsNestingTag(tlv.Tag()) {
		return nil, Err_UnexpectedTag
	}
	if tlv.hasOpenChild {
		return nil, Err_HasOpenChild
	}
	child, err := Open(tlv.body[tlv.curr:], tlv, t)
	if err != nil {
		return nil, err
	}
	tlv.hasOpenChild = true
	return child, nil
}

//allocate a nested TLV (child) as VALUE of 'this';
//this.curr is increased by the given l;
func (tlv *TLV) AllocateChild(t TLVTag, l int) (*TLV, error) {
	if tlv.state != TLVState_Allocated && tlv.state != TLVState_Open {
		return nil, Err_UnexpectedState
	}
	if !IsNestingTag(tlv.Tag()) {
		return nil, Err_UnexpectedTag
	}
	if tlv.hasOpenChild {
		return nil, Err_HasOpenChild
	}
	child, err := Allocate(tlv.body[tlv.curr:], tlv, t, l)
	if err != nil {
		return nil, err
	}
	return child, nil
}

//build a nested TLV (child) as VALUE of 'this';
//this.curr is increased by the len of the given v;
func (tlv *TLV) BuildChild(t TLVTag, v []byte) (*TLV, error) {
	if tlv.state != TLVState_Allocated && tlv.state != TLVState_Open {
		return nil, Err_UnexpectedState
	}
	if !IsNestingTag(tlv.Tag()) {
		return nil, Err_UnexpectedTag
	}
	if tlv.hasOpenChild {
		return nil, Err_HasOpenChild
	}
	child, err := Build(tlv.body[tlv.curr:], tlv, t, v)
	if err != nil {
		return nil, err
	}
	return child, nil
}

func (tlv *TLV) BuildChildInt32(t TLVTag, v int32) (*TLV, error) {
	b := make([]byte, 4)
	encoder.PutUint32(b, uint32(v))
	return tlv.BuildChild(t, b)
}

func (tlv *TLV) BuildChildInt64(t TLVTag, v int64) (*TLV, error) {
	b := make([]byte, 8)
	encoder.PutUint64(b, uint64(v))
	return tlv.BuildChild(t, b)
}

func (tlv *TLV) Tag() TLVTag {
	return TLVTag(encoder.Uint32(tlv.body[TAG_OFFSET:]))
}

func (tlv *TLV) Len() int {
	return int(encoder.Uint32(tlv.body[LEN_OFFSET:]))
}

func (tlv *TLV) Val() []byte {
	return tlv.body[VAL_OFFSET:]
}

func (tlv *TLV) ValAsInt32() int32 {
	return int32(encoder.Uint32(tlv.body[VAL_OFFSET:]))
}

func (tlv *TLV) ValAsInt64() int64 {
	return int64(encoder.Uint64(tlv.body[VAL_OFFSET:]))
}

func (tlv *TLV) Equal(other *TLV) bool {
	return bytes.Equal(tlv.body, other.body)
}

func (tlv *TLV) ChildWithTag(t TLVTag) *TLV {
	c := tlv.children[t]
	if len(c) == 0 {
		return nil
	}
	return c[0]
}

func (tlv *TLV) ChildrenWithTag(t TLVTag) []*TLV {
	return tlv.children[t]
}

func Marshal(tlv *TLV) ([]byte, error) {
	if tlv.state != TLVState_Closed {
		return nil, Err_UnexpectedState
	}
	return tlv.body, nil
}

func Unmarshal(bin []byte, p *TLV) (map[TLVTag][]*TLV, error) {
	var err error
	tlvs := make(map[TLVTag][]*TLV)
	for len(bin) > 0 {
		if len(bin) < TAG_SIZE+LEN_SIZE {
			return nil, Err_Corrupted
		}
		t := TLVTag(encoder.Uint32(bin[TAG_OFFSET:]))
		l := int(encoder.Uint32(bin[LEN_OFFSET:]))

		if len(bin) < VAL_OFFSET+l {
			return nil, Err_Corrupted
		}
		body := bin[0 : VAL_OFFSET+l]
		bin = bin[VAL_OFFSET+l:]

		tlv := &TLV{
			parent:       p,
			body:         body,
			curr:         VAL_OFFSET + l,
			state:        TLVState_Closed,
			hasOpenChild: false,
		}
		tlvs[t] = append(tlvs[t], tlv)

		if IsRawTag(t) {
			continue
		}
		if IsNestingTag(t) {
			tlv.children, err = Unmarshal(body[VAL_OFFSET:], tlv)
			if err != nil {
				return nil, err
			}
			continue
		}
		return nil, Err_Corrupted //invalid tag, corrupted
	}
	return tlvs, nil
}

func (tlv *TLV) allocate(valLen int) error {
	if tlv.state == TLVState_Allocated {
		return nil
	}
	if tlv.state != TLVState_Open {
		return Err_UnexpectedState
	}
	bodyLen := TAG_SIZE + LEN_SIZE + valLen
	if tlv.parent != nil {
		if err := tlv.parent.allocateForOpenChild(bodyLen); err != nil {
			return err
		}
	}
	tlv.mustSetLen(valLen)
	tlv.body = tlv.body[:bodyLen]
	tlv.state = TLVState_Allocated
	return nil
}

func (tlv *TLV) close(valLen int) error {
	if tlv.state == TLVState_Closed {
		return nil
	}
	if tlv.state != TLVState_Allocated {
		return Err_UnexpectedState
	}
	if tlv.Len() != valLen || tlv.curr != VAL_OFFSET+valLen || tlv.curr != len(tlv.body) {
		return Err_Corrupted
	}
	tlv.state = TLVState_Closed
	return nil
}

func (tlv *TLV) allocateForOpenChild(childLen int) error {
	if tlv.curr+childLen > len(tlv.body) {
		return Err_WouldOverflow
	}
	tlv.curr += childLen
	tlv.hasOpenChild = false
	return nil
}

func (tlv *TLV) mustSetTag(t TLVTag) {
	if len(tlv.body[TAG_OFFSET:]) < TAG_SIZE {
		panic("buffer overflow")
	}
	encoder.PutUint32(tlv.body[TAG_OFFSET:], uint32(t))
}

func (tlv *TLV) mustSetLen(l int) {
	if len(tlv.body[LEN_OFFSET:]) < LEN_SIZE {
		panic("buffer overflow")
	}
	encoder.PutUint32(tlv.body[LEN_OFFSET:], uint32(l))
}

func (tlv *TLV) mustFillRaw(v []byte) {
	if len(tlv.body[tlv.curr:]) < len(v) {
		panic("buffer overflow")
	}
	copy(tlv.body[tlv.curr:], v)
	tlv.curr += len(v)
}

func (tlv *TLV) mustAddChild(childTag TLVTag, child *TLV) {
	if tlv.children == nil {
		panic("add child to non-nesting tlv")
	}
	tlv.children[childTag] = append(tlv.children[childTag], child)
}
