package mqpp

import "encoding/binary"

//Puback - publish acknowledgement
type Puback struct {
	src []byte
}

func NewPuback(data []byte) (*Puback, error) {
	if len(data) != 4 || data[0] != (PUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Puback{src: data}, nil
}

func (p *Puback) Length() uint32 { return uint32(len(p.src)) }

func (p *Puback) Type() byte { return p.src[0] >> 4 }

func (p *Puback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.src[2:])
}
