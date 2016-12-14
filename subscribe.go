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

// Subscribe mqtt subscribe to topics, sturcture:
// fixed header:
// variable header: Packet Identifier
// payload: (Topic Filter, Requested QoS)s
type Subscribe struct {
	endecBytes
	packetIDPos     int
	topicFilterPoss []int
}

// Subscription - topic filter and requested qos
type Subscription struct {
	TopicFilter  string
	RequestedQoS byte
}

func newSubscribe(data []byte) (*Subscribe, error) {
	if data[0] != (TSUBSCRIBE<<4 | 0x02) {
		return nil, ErrProtocolViolation
	}

	p := &Subscribe{endecBytes: data}
	remlen, offset := p.remlen(1)
	if offset <= 1 {
		return nil, ErrMalformedRemLen
	}
	pktLen := offset + int(remlen)
	if len(data) < pktLen {
		return nil, ErrProtocolViolation
	}
	p.packetIDPos = offset
	_, offset = p.uint16(p.packetIDPos)
	p.topicFilterPoss = []int{}
	for offset < pktLen {
		p.topicFilterPoss = append(p.topicFilterPoss, offset)
		_, offset = p.string(offset)
		_, offset = p.byte(offset)
	}

	return p, nil
}

// MakeSubscribe create a mqtt subscribe packet
func MakeSubscribe(packetIdentifier uint16, payload []Subscription) Subscribe {
	p := Subscribe{}
	remlen := p.calc(packetIdentifier)
	for _, s := range payload {
		remlen += p.calc(s.TopicFilter, s.RequestedQoS)
	}
	pktLen := 1 + p.calc(uint32(remlen)) + remlen

	p.endecBytes = make([]byte, pktLen)
	p.packetIDPos = p.fill(0, (TSUBSCRIBE<<4 | 0x02), uint32(remlen))
	offset := p.fill(p.packetIDPos, packetIdentifier)
	p.topicFilterPoss = make([]int, len(payload))
	for i, s := range payload {
		p.topicFilterPoss[i] = offset
		offset = p.fill(offset, s.TopicFilter, s.RequestedQoS)
	}
	return p
}

// PacketIdentifier return packet id
func (s *Subscribe) PacketIdentifier() uint16 {
	pid, _ := s.uint16(s.packetIDPos)
	return pid
}

// Payload return topicfilters and requested qoss
func (s *Subscribe) Payload() []Subscription {
	subs := make([]Subscription, len(s.topicFilterPoss))
	for i, offset := range s.topicFilterPoss {
		filter, pos := s.string(offset)
		qos, _ := s.byte(pos)
		subs[i] = Subscription{
			TopicFilter:  filter,
			RequestedQoS: qos,
		}
	}
	return subs
}
