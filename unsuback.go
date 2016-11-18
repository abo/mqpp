package mqpp

import "encoding/binary"

// Unsuback - unsubscribe acknowledgement
type Unsuback struct {
	packetBytes
}

func newUnsuback(data []byte) (*Unsuback, error) {
	if len(data) < 4 || data[0] != (UNSUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Unsuback{packetBytes: data[0:4]}, nil
}

// MakeUnsuback create a mqtt unsuback packet
func MakeUnsuback(packetIdentifier uint16) Unsuback {
	pb := make([]byte, 4)
	fill(pb, UNSUBACK<<4, byte(2), packetIdentifier)
	return Unsuback{packetBytes: pb}
}

// PacketIdentifier return packet id
func (s *Unsuback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(s.packetBytes[2:])
}
