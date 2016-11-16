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

// PacketIdentifier return packet id
func (s *Unsuback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(s.packetBytes[2:])
}
