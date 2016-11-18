package mqpp

import "encoding/binary"

// Pubrel - publish release(qos 2 publish received, part 2)
type Pubrel struct {
	packetBytes
}

func newPubrel(data []byte) (*Pubrel, error) {
	if len(data) < 4 || data[0] != (PUBREL<<4|0x02) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrel{packetBytes: data[0:4]}, nil
}

// MakePubrel create a mqtt pubrel packet
func MakePubrel(packetIdentifier uint16) Pubrel {
	pb := make([]byte, 4)
	fill(pb, PUBREL<<4|0x02, uint32(2), packetIdentifier)
	return Pubrel{packetBytes: pb}
}

// PacketIdentifier return packet id
func (p *Pubrel) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
