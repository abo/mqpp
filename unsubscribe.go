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
	remlen, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	offset += 2 //variable header

	filterLens := []int{}
	for offset < 1+remlenLen+int(remlen) {
		filterLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + filterLen)
		filterLens = append(filterLens, filterLen)
	}

	return &Unsubscribe{
		packetBytes:          data,
		remainingLengthBytes: remlenLen,
		topicFiltersBytes:    filterLens,
	}, nil
}

func (u *Unsubscribe) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + u.remainingLengthBytes
	return binary.BigEndian.Uint16(u.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

func (u *Unsubscribe) Payload() []string {
	filters := make([]string, len(u.topicFiltersBytes))
	offset := 1 + u.remainingLengthBytes + 2
	for i, l := range u.topicFiltersBytes {
		filters[i] = string(u.packetBytes[offset+2 : offset+2+l])
		offset += (2 + l)
	}
	return filters
}
