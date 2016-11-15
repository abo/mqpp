package mqpp

import "encoding/binary"

// Subscribe - subscribe to topics
type Subscribe struct {
	src                  []byte
	remainingLengthBytes int
	topicFiltersBytes    []int
}

// Subscription - topic filter and requested qos
type Subscription struct {
	TopicFilter  string
	RequestedQoS byte
}

func NewSubscribe(data []byte) (*Subscribe, error) {
	if data[0] != (SUBSCRIBE<<4 | 0x02) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	offset += 2 // variable header

	filterLens := []int{}
	for offset < 1+remlenLen+int(remlen) {
		filterLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + filterLen + 1)
		filterLens = append(filterLens, filterLen)
	}
	return &Subscribe{
		src:                  data,
		remainingLengthBytes: remlenLen,
		topicFiltersBytes:    filterLens,
	}, nil
}

func (s *Subscribe) Length() uint32 { return uint32(len(s.src)) }

func (s *Subscribe) Type() byte { return s.src[0] >> 4 }

func (s *Subscribe) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + s.remainingLengthBytes
	return binary.BigEndian.Uint16(s.src[fixedHeaderLen : fixedHeaderLen+2])
}

func (s *Subscribe) Payload() []Subscription {
	subs := make([]Subscription, len(s.topicFiltersBytes))
	offset := 1 + s.remainingLengthBytes + 2
	for i, l := range s.topicFiltersBytes {
		filter := string(s.src[offset+2 : offset+2+l])
		qos := s.src[offset+2+l]
		subs[i] = Subscription{
			TopicFilter:  filter,
			RequestedQoS: qos,
		}
		offset += (2 + l + 1)
	}
	return subs
}
