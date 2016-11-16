package mqpp

import "encoding/binary"

// Pubrec - publish received(qos 2 publish received, part 1)
type Pubrec struct {
	packetBytes
}

func newPubrec(data []byte) (*Pubrec, error) {
	if len(data) != 4 || data[0] != (PUBREC<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrec{packetBytes: data}, nil
}

func (p *Pubrec) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
