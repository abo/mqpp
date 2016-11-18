package mqpp

import "encoding/binary"

// Unsubscribe - unsubscribe from topics
type Unsubscribe struct {
	packetBytes
	remainingLengthBytes int
	topicFiltersBytes    []int
}

func newUnsubscribe(data []byte) (*Unsubscribe, error) {
	if data[0] != (UNSUBSCRIBE<<4 | 0x02) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := decRemLen(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen
	packetLen := offset + int(remlen)
	if len(data) < packetLen {
		return nil, ErrProtocolViolation
	}

	offset += 2 //variable header

	filterLens := []int{}
	for offset < 1+remlenLen+int(remlen) {
		filterLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + filterLen)
		filterLens = append(filterLens, filterLen)
	}

	return &Unsubscribe{
		packetBytes:          data[0:packetLen],
		remainingLengthBytes: remlenLen,
		topicFiltersBytes:    filterLens,
	}, nil
}

// MakeUnsubscribe create a mqtt unsubscribe packet
func MakeUnsubscribe(packetIdentifier uint16, payload []string) Unsubscribe {
	remlen := 0
	remlen += 2
	filterLens := []int{}
	for _, filter := range payload {
		remlen += (2 + len(filter))
		filterLens = append(filterLens, len(filter))
	}
	pb := make([]byte, 1+lenRemLen(uint32(remlen))+remlen)
	offset := fill(pb, UNSUBSCRIBE<<4|0x02, uint32(remlen), packetIdentifier)
	for _, filter := range payload {
		offset += fill(pb[offset:], filter)
	}
	return Unsubscribe{
		packetBytes:          pb,
		remainingLengthBytes: lenRemLen(uint32(remlen)),
		topicFiltersBytes:    filterLens,
	}
}

// PacketIdentifier return packet id
func (u *Unsubscribe) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + u.remainingLengthBytes
	return binary.BigEndian.Uint16(u.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

// Payload return topic filters
func (u *Unsubscribe) Payload() []string {
	filters := make([]string, len(u.topicFiltersBytes))
	offset := 1 + u.remainingLengthBytes + 2
	for i, l := range u.topicFiltersBytes {
		filters[i] = string(u.packetBytes[offset+2 : offset+2+l])
		offset += (2 + l)
	}
	return filters
}
