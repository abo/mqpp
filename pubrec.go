package mqpp

import "encoding/binary"

// Pubrec - publish received(qos 2 publish received, part 1)
type Pubrec struct {
	packetBytes
}

func newPubrec(data []byte) (*Pubrec, error) {
	if len(data) < 4 || data[0] != (PUBREC<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrec{packetBytes: data[0:4]}, nil
}

// MakePubrec create a mqtt pubrec packet
func MakePubrec(packetIdentifier uint16) Pubrec {
	pb := make([]byte, 4)
	fill(pb, PUBREC<<4, uint32(2), packetIdentifier)
	return Pubrec{packetBytes: pb}
}

// PacketIdentifier return packet id
func (p *Pubrec) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
