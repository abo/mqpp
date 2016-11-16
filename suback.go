package mqpp

import "encoding/binary"

// Suback - subscribe acknowledgement
type Suback struct {
	packetBytes
	remainingLengthBytes int
}

func newSuback(data []byte) (*Suback, error) {
	if data[0] != (SUBACK << 4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	_, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}

	return &Suback{
		packetBytes:          data,
		remainingLengthBytes: remlenLen,
	}, nil
}

func (p *Suback) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	return binary.BigEndian.Uint16(p.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

func (p *Suback) ReturnCodes() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2
	return p.packetBytes[fixedHeaderLen+variableHeaderLen:]
}
