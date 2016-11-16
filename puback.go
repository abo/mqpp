package mqpp

import "encoding/binary"

//Puback - publish acknowledgement
type Puback struct {
	packetBytes
}

func newPuback(data []byte) (*Puback, error) {
	if len(data) != 4 || data[0] != (PUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Puback{packetBytes: data}, nil
}

func (p *Puback) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
