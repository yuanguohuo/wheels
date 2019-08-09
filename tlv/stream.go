package tlv

import "io"

type StreamHandler interface {
	OnFailure(err error)
	OnSuccess()

	OnBegin(tlv *TLV)
	OnEnd(tlv *TLV)
}

func Stream(tlv *TLV, writer io.Writer) error {
	body, err := Marshal(tlv)
	if err != nil {
		return err
	}
	if _, err := writer.Write(body); err != nil {
		return err
	}
	return nil
}

func ParseStream(buffer []byte, reader io.Reader, cb StreamHandler) error {
	var t TLVTag
	var l int
	var num int
	var err error
	var nearestParent *TLV = nil

	for {
		num, err = readn(buffer[0:TAG_SIZE+LEN_SIZE], reader)
		if num < TAG_SIZE+LEN_SIZE {
			if num == 0 && err == io.EOF {
				cb.OnSuccess()
				break
			}
			// 0<num<TAG_SIZE+LEN_SIZE  || err != io.EOF
			cb.OnFailure(Err_Corrupted)
			return Err_Corrupted
		}

		//assert: num == TAG_SIZE+LEN_SIZE
		t = TLVTag(encoder.Uint32(buffer[0:TAG_SIZE]))
		l = int(encoder.Uint32(buffer[TAG_SIZE : TAG_SIZE+LEN_SIZE]))
		if len(buffer) < VAL_OFFSET+l {
			cb.OnFailure(Err_WouldOverflow)
			return Err_WouldOverflow
		}

		switch {
		case IsRawTag(t):
			num, err = readn(buffer[VAL_OFFSET:VAL_OFFSET+l], reader)
			if num < l {
				cb.OnFailure(err)
				return err
			}

			raw := takeOverRaw(buffer[0:VAL_OFFSET+l], nearestParent)
			buffer = buffer[VAL_OFFSET+l:]
			cb.OnBegin(raw)
			cb.OnEnd(raw)

			c := raw
			for {
				if nearestParent == nil {
					break
				}
				nearestParent.curr += len(c.body)
				if nearestParent.curr < len(nearestParent.body) { //nearest parent is not full
					break
				}
				if nearestParent.curr == len(nearestParent.body) { //nearest parent is full
					nearestParent.state = TLVState_Closed
					cb.OnEnd(nearestParent)
					c = nearestParent
					nearestParent = c.parent
					continue
				}
				//never: nearestParent.curr > len(nearestParent.body)
				cb.OnFailure(Err_WouldOverflow)
				return Err_WouldOverflow
			}

		case IsNestingTag(t):
			nesting := takeOverNesting(buffer[0:VAL_OFFSET+l], nearestParent)
			nearestParent = nesting
			buffer = buffer[TAG_SIZE+LEN_SIZE:]
			cb.OnBegin(nesting)

		default:
			cb.OnFailure(Err_UnexpectedTag)
			return Err_UnexpectedTag
		}
	}

	if nearestParent != nil {
		cb.OnFailure(Err_Corrupted)
		return Err_Corrupted
	}

	return nil
}

func takeOverRaw(body []byte, p *TLV) *TLV {
	tlv := &TLV{
		parent:       p,
		body:         body,
		curr:         len(body),
		state:        TLVState_Closed,
		hasOpenChild: false,
		children:     nil,
	}
	if p != nil {
		p.mustAddChild(tlv.Tag(), tlv)
	}
	return tlv
}

func takeOverNesting(body []byte, p *TLV) *TLV {
	tlv := &TLV{
		parent:       p,
		body:         body,
		curr:         VAL_OFFSET,
		state:        TLVState_Allocated,
		hasOpenChild: false,
		children:     make(map[TLVTag][]*TLV),
	}
	if p != nil {
		p.mustAddChild(tlv.Tag(), tlv)
	}
	return tlv
}

func readn(b []byte, reader io.Reader) (int, error) {
	var n int = 0
	for len(b) > 0 {
		m, e := reader.Read(b)
		n += m
		if e != nil {
			return n, e
		}
		b = b[m:]
	}
	return n, nil
}
