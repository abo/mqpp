package mqpp

import "encoding/binary"

// Unsuback - unsubscribe acknowledgement
type Unsuback struct {
	src []byte
}

func NewUnsuback(data []byte) (*Unsuback, error) {
	if len(data) != 4 || data[0] != (UNSUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Unsuback{src: data}, nil
}

func (s *Unsuback) Length() uint32 { return uint32(len(s.src)) }

func (s *Unsuback) Type() byte { return s.src[0] >> 4 }

func (s *Unsuback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(s.src[2:])
}
