package mqpp

// Connack - acknowledge connection request
type Connack struct {
	src []byte
}

// NewConnack create Connack from byte slice,
func NewConnack(data []byte) (*Connack, error) {
	// check packet length, packet type, remaining length, conack flags, return code
	if len(data) != 4 || data[0] != (CONNACK<<4) || data[1] != 2 || (data[2]>>1) != 0 || uint8(data[3]) > 5 {
		return nil, ErrProtocolViolation
	}

	return &Connack{src: data}, nil
}

// func NewConnack() *Connack {
// 	return &Connack{
// 		src: []byte{CONNACK << 4, byte(2), 0x00, 0x00},
// 	}
// }

func (p *Connack) Type() byte {
	return p.src[0] >> 4
}

func (p *Connack) Length() uint32 {
	return uint32(len(p.src))
}

func (p *Connack) SetSessionPresent(present bool) {
	if present {
		p.src[2] = 0x01
	} else {
		p.src[2] = 0x00
	}
}

func (p *Connack) SessionPresent() bool {
	return p.src[2]&0x01 == 0x01
}

func (p *Connack) SetReturnCode(code byte) {
	p.src[3] = code
}

// ReturnCode
func (p *Connack) ReturnCode() byte {
	return p.src[3]
}
