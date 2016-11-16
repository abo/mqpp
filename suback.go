package mqpp

import "encoding/binary"

// Suback - subscribe acknowledgement
type Suback struct {
	packetBytes
	remainingLengthBytes int
}

func newSuback(data []byte) (*Suback, error) {
	if len(data) < 1 || data[0] != (SUBACK<<4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	packetLen := 1 + remlenLen + int(remlen)
	if len(data) < packetLen {
		return nil, ErrProtocolViolation
	}

	return &Suback{
		packetBytes:          data[0:packetLen],
		remainingLengthBytes: remlenLen,
	}, nil
}

// PacketIdentifier return packet id
func (p *Suback) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	return binary.BigEndian.Uint16(p.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

// ReturnCodes return sub return codes
func (p *Suback) ReturnCodes() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2
	return p.packetBytes[fixedHeaderLen+variableHeaderLen:]
}
