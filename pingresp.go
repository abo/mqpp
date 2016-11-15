package mqpp

// Pingresp - ping response
type Pingresp struct {
	src []byte
}

func NewPingresp(data []byte) (*Pingresp, error) {
	if len(data) != 2 || data[0] != (PINGRESP<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Pingresp{src: data}, nil
}

func (p *Pingresp) Length() uint32 { return uint32(len(p.src)) }

func (p *Pingresp) Type() byte { return p.src[0] >> 4 }
