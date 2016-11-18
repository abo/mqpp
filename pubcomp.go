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

// MakePubcomp create a mqtt pubcomp packet
func MakePubcomp(packetIdentifier uint16) Pubcomp {
	pb := make([]byte, 4)
	fill(pb, PUBCOMP<<4, uint32(2), packetIdentifier)
	return Pubcomp{packetBytes: pb}
}

// PacketIdentifier return packet id
func (p *Pubcomp) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
