package mqpp

// Pingreq - ping request
type Pingreq struct {
	src []byte
}

func NewPingreq(data []byte) (*Pingreq, error) {
	if len(data) != 2 || data[0] != (PINGREQ<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Pingreq{src: data}, nil
}

func (p *Pingreq) Length() uint32 { return uint32(len(p.src)) }

func (p *Pingreq) Type() byte { return p.src[0] >> 4 }
