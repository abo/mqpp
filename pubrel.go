package mqpp

import "encoding/binary"

// Pubrel - publish release(qos 2 publish received, part 2)
type Pubrel struct {
	packetBytes
}

func newPubrel(data []byte) (*Pubrel, error) {
	if len(data) != 4 || data[0] != (PUBREL<<4|0x02) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrel{packetBytes: data}, nil
}

func (p *Pubrel) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
