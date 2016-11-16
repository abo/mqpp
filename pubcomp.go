package mqpp

import "encoding/binary"

// Pubcomp - publish complete (qos 2 publish received, part 3)
type Pubcomp struct {
	packetBytes
}

func newPubcomp(data []byte) (*Pubcomp, error) {
	if len(data) < 4 || data[0] != (PUBCOMP<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubcomp{packetBytes: data[0:4]}, nil
}

// PacketIdentifier return packet id
func (p *Pubcomp) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
