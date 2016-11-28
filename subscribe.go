// Copyright (c) 2016 The MQPP Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mqpp

import "encoding/binary"

// Subscribe mqtt subscribe to topics, sturcture:
// fixed header:
// variable header: Packet Identifier
// payload: (Topic Filter, Requested QoS)s
type Subscribe struct {
	packetBytes
	remainingLengthBytes int
	topicFiltersBytes    []int
}

// Subscription - topic filter and requested qos
type Subscription struct {
	TopicFilter  string
	RequestedQoS byte
}

func newSubscribe(data []byte) (*Subscribe, error) {
	if data[0] != (SUBSCRIBE<<4 | 0x02) {
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

	offset += 2 // variable header

	filterLens := []int{}
	for offset < 1+remlenLen+int(remlen) {
		filterLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + filterLen + 1)
		filterLens = append(filterLens, filterLen)
	}
	return &Subscribe{
		packetBytes:          data[0:packetLen],
		remainingLengthBytes: remlenLen,
		topicFiltersBytes:    filterLens,
	}, nil
}

// MakeSubscribe create a mqtt subscribe packet
func MakeSubscribe(packetIdentifier uint16, payload []Subscription) Subscribe {
	remlen := 0
	remlen += 2
	filterLens := []int{}
	for _, s := range payload {
		filterLens = append(filterLens, len(s.TopicFilter))
		remlen += (2 + len([]byte(s.TopicFilter)) + 1)
	}

	pb := make([]byte, 1+lenRemLen(uint32(remlen))+remlen)
	offset := fill(pb, (SUBSCRIBE<<4 | 0x02), uint32(remlen), packetIdentifier)
	for _, s := range payload {
		offset += fill(pb[offset:], s.TopicFilter, s.RequestedQoS)
	}
	return Subscribe{
		packetBytes:          pb,
		remainingLengthBytes: lenRemLen(uint32(remlen)),
		topicFiltersBytes:    filterLens,
	}
}

// PacketIdentifier return packet id
func (s *Subscribe) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + s.remainingLengthBytes
	return binary.BigEndian.Uint16(s.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

// Payload return topicfilters and requested qoss
func (s *Subscribe) Payload() []Subscription {
	subs := make([]Subscription, len(s.topicFiltersBytes))
	offset := 1 + s.remainingLengthBytes + 2
	for i, l := range s.topicFiltersBytes {
		filter := string(s.packetBytes[offset+2 : offset+2+l])
		qos := s.packetBytes[offset+2+l]
		subs[i] = Subscription{
			TopicFilter:  filter,
			RequestedQoS: qos,
		}
		offset += (2 + l + 1)
	}
	return subs
}
