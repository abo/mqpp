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

// Publish mqtt publish message, structure:
// fixed header
// variable header: Topic Name, Packet Identifier
// payload: content
type Publish struct {
	endecBytes
	topicNamePos int
	packetIDPos  int
	payloadPos   int
}

func newPublish(data []byte) (*Publish, error) {
	if len(data) < 1 || data[0]>>4 != TPUBLISH {
		return nil, ErrProtocolViolation
	}
	p := &Publish{endecBytes: data}
	fixedHeader, offset := p.byte(0)
	qos := fixedHeader << 5 >> 6
	remlen, offset := p.remlen(offset)
	if offset <= 1 {
		return nil, ErrMalformedRemLen
	}
	pktLen := offset + int(remlen)
	if len(data) < pktLen {
		return nil, ErrProtocolViolation
	}
	p.topicNamePos = offset
	_, offset = p.string(p.topicNamePos)
	if qos > QosAtMostOnce {
		p.packetIDPos = offset
		_, offset = p.uint16(offset)
	}
	p.payloadPos = offset

	return p, nil
}

// MakePublish create a mqtt publish packet
func MakePublish(dup bool, qos byte, retain bool, topicName string, packetIdentifier uint16, payload []byte) Publish {
	p := Publish{}
	remlen := p.calc(topicName, payload)
	if qos > QosAtMostOnce {
		remlen += p.calc(packetIdentifier)
	}
	pktLen := 1 + p.calc(uint32(remlen)) + remlen

	p.endecBytes = make([]byte, pktLen)
	offset := p.fill(0, (TPUBLISH<<4)|(qos<<1))
	p.set(0, 3, dup).set(0, 0, retain)
	p.topicNamePos = p.fill(offset, uint32(remlen))
	offset = p.fill(p.topicNamePos, topicName)
	if qos > QosAtMostOnce {
		p.packetIDPos = offset
		offset = p.fill(p.packetIDPos, packetIdentifier)
	}
	p.payloadPos = offset
	p.fill(p.payloadPos, payload)
	return p
}

// Dup return is dup
func (p *Publish) Dup() bool {
	return p.bit(0, 3)
}

// QoS return qos
func (p *Publish) QoS() byte {
	qos, _ := p.byte(0)
	return qos << 5 >> 6
}

// Retain return is retain set
func (p *Publish) Retain() bool {
	return p.bit(0, 0)
}

// TopicName return topic name
func (p *Publish) TopicName() string {
	topic, _ := p.string(p.topicNamePos)
	return topic
}

// PacketIdentifier return packet id if qos > 0, or zero when qos = 0
func (p *Publish) PacketIdentifier() uint16 {
	if p.QoS() > QosAtMostOnce {
		pid, _ := p.uint16(p.packetIDPos)
		return pid
	}

	return 0
}

// Payload return publish content
func (p *Publish) Payload() []byte {
	return p.bytes(p.payloadPos)
}
