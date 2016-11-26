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

// MakeConnack return a mqtt connack packet with SessionPresent and ReturnCode
func MakeConnack(sessionPresent bool, returnCode byte) Connack {
	pb := make([]byte, 4)
	fill(pb, CONNACK<<4, uint32(2), set(0, sessionPresent), returnCode)
	return Connack{packetBytes: pb}
}

// SetSessionPresent set is session present
func (p *Connack) SetSessionPresent(sessionPresent bool) {
	p.packetBytes[2] = set(0, sessionPresent)
}

// SessionPresent return is session present
func (p *Connack) SessionPresent() bool {
	return p.packetBytes[2]&0x01 == 0x01
}

// SetReturnCode set return code
func (p *Connack) SetReturnCode(returnCode byte) {
	p.packetBytes[3] = returnCode
}

// ReturnCode return connect return code
func (p *Connack) ReturnCode() byte {
	return p.packetBytes[3]
}
