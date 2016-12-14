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

// Unsubscribe mqtt unsubscribe from topics, structure:
// fixed header:
// variable header: Packet Identifier
// payload: Topic Filters
type Unsubscribe struct {
	endecBytes
	packetIDPos     int
	topicFilterPoss []int
}

func newUnsubscribe(data []byte) (*Unsubscribe, error) {
	if data[0] != (TUNSUBSCRIBE<<4 | 0x02) {
		return nil, ErrProtocolViolation
	}

	p := &Unsubscribe{endecBytes: data}
	// 1) packet type
	remlen, offset := p.remlen(1) // 2) remaining length
	p.packetIDPos = offset
	pktLen := offset + int(remlen)
	_, offset = p.uint16(p.packetIDPos) // 3) packet identifier
	p.topicFilterPoss = []int{}
	for offset < pktLen {
		p.topicFilterPoss = append(p.topicFilterPoss, offset)
		_, offset = p.string(offset) // 4~N) topic filter
	}

	return p, nil
}

// MakeUnsubscribe create a mqtt unsubscribe packet
func MakeUnsubscribe(packetIdentifier uint16, payload []string) Unsubscribe {
	p := Unsubscribe{}
	remlen := p.calc(packetIdentifier, payload)
	pktLen := 1 + p.calc(uint32(remlen)) + remlen

	p.endecBytes = make([]byte, pktLen)
	p.packetIDPos = p.fill(0, TUNSUBSCRIBE<<4|0x02, uint32(remlen))
	offset := p.fill(p.packetIDPos, packetIdentifier)
	p.topicFilterPoss = make([]int, len(payload))
	for i, filter := range payload {
		p.topicFilterPoss[i] = offset
		offset = p.fill(offset, filter)
	}

	return p
}

// PacketIdentifier return packet id
func (u *Unsubscribe) PacketIdentifier() uint16 {
	pid, _ := u.uint16(u.packetIDPos)
	return pid
}

// Payload return topic filters
func (u *Unsubscribe) Payload() []string {
	filters := make([]string, len(u.topicFilterPoss))
	for i, l := range u.topicFilterPoss {
		filters[i], _ = u.string(l)
	}
	return filters
}
