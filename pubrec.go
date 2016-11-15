package mqpp

import "encoding/binary"

// Pubrec - publish received(qos 2 publish received, part 1)
type Pubrec struct {
	src []byte
}

func NewPubrec(data []byte) (*Pubrec, error) {
	if len(data) != 4 || data[0] != (PUBREC<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrec{src: data}, nil
}

func (p *Pubrec) Length() uint32 { return uint32(len(p.src)) }

func (p *Pubrec) Type() byte { return p.src[0] >> 4 }

func (p *Pubrec) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.src[2:])
}
