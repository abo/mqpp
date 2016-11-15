package mqpp

import "encoding/binary"

// Pubcomp - publish complete (qos 2 publish received, part 3)
type Pubcomp struct {
	src []byte
}

func NewPubcomp(data []byte) (*Pubcomp, error) {
	if len(data) != 4 || data[0] != (PUBCOMP<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubcomp{src: data}, nil
}

func (p *Pubcomp) Length() uint32 { return uint32(len(p.src)) }

func (p *Pubcomp) Type() byte { return p.src[0] >> 4 }

func (p *Pubcomp) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.src[2:])
}
