package mqpp

import "encoding/binary"

// Suback - subscribe acknowledgement
type Suback struct {
	src                  []byte
	remainingLengthBytes int
}

func NewSuback(data []byte) (*Suback, error) {
	if data[0] != (SUBACK << 4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	_, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}

	return &Suback{
		src:                  data,
		remainingLengthBytes: remlenLen,
	}, nil
}

func (p *Suback) Length() uint32 { return uint32(len(p.src)) }

func (p *Suback) Type() byte { return p.src[0] >> 4 }

func (p *Suback) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	return binary.BigEndian.Uint16(p.src[fixedHeaderLen : fixedHeaderLen+2])
}

func (p *Suback) ReturnCodes() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2
	return p.src[fixedHeaderLen+variableHeaderLen:]
}
