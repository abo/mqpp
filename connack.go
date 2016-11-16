package mqpp

// Connack - acknowledge connection request
type Connack struct {
	packetBytes
}

// newConnack create Connack from byte slice,
func newConnack(data []byte) (*Connack, error) {
	// check packet length, packet type, remaining length, conack flags, return code
	if len(data) < 4 || data[0] != (CONNACK<<4) || data[1] != 2 || (data[2]>>1) != 0 || uint8(data[3]) > 5 {
		return nil, ErrProtocolViolation
	}

	return &Connack{packetBytes: data[0:4]}, nil
}

// func NewConnack() *Connack {
// 	return &Connack{
// 		src: []byte{CONNACK << 4, byte(2), 0x00, 0x00},
// 	}
// }

func (p *Connack) SetSessionPresent(present bool) {
	if present {
		p.packetBytes[2] = 0x01
	} else {
		p.packetBytes[2] = 0x00
	}
}

// SessionPresent return is session present
func (p *Connack) SessionPresent() bool {
	return p.packetBytes[2]&0x01 == 0x01
}

func (p *Connack) SetReturnCode(code byte) {
	p.packetBytes[3] = code
}

// ReturnCode return connect return code
func (p *Connack) ReturnCode() byte {
	return p.packetBytes[3]
}
