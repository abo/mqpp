package mqpp

import "encoding/binary"

// Pubrel - publish release(qos 2 publish received, part 2)
type Pubrel struct {
	src []byte
}

func NewPubrel(data []byte) (*Pubrel, error) {
	if len(data) != 4 || data[0] != (PUBREL<<4|0x02) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrel{src: data}, nil
}

func (p *Pubrel) Length() uint32 { return uint32(len(p.src)) }

func (p *Pubrel) Type() byte { return p.src[0] >> 4 }

func (p *Pubrel) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.src[2:])
}
