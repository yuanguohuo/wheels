package tlv

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
