package mqpp

import "encoding/binary"

// Puback - publish acknowledgement
type Puback struct {
	packetBytes
}

func newPuback(data []byte) (*Puback, error) {
	if len(data) < 4 || data[0] != (PUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Puback{packetBytes: data[0:4]}, nil
}

// MakePuback create a mqtt puback Packet
func MakePuback(packetIdentifier uint16) Puback {
	pb := make([]byte, 4)
	fill(pb, PUBACK<<4, uint32(2), packetIdentifier)
	return Puback{packetBytes: pb}
}

// PacketIdentifier return packet id
func (p *Puback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
