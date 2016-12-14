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

// Suback mqtt subscribe acknowledgement, structure:
// fixed header:
// variable header: Packet Identifier
// payload: Return Codes
type Suback struct {
	endecBytes
	packetIDPos int
}

func newSuback(data []byte) (*Suback, error) {
	if len(data) < 1 || data[0] != (TSUBACK<<4) {
		return nil, ErrProtocolViolation
	}

	p := &Suback{endecBytes: data}
	remlen, offset := p.remlen(1)
	if offset <= 1 {
		return nil, ErrMalformedRemLen
	}
	pktLen := offset + int(remlen)
	if len(data) < pktLen {
		return nil, ErrProtocolViolation
	}
	p.packetIDPos = offset
	return p, nil
}

// MakeSuback create a mqtt suback packet
func MakeSuback(packetIdentifier uint16, returnCodes []byte) Suback {
	p := Suback{}
	remlen := p.calc(packetIdentifier, returnCodes)
	pktLen := 1 + p.calc(uint32(remlen)) + remlen
	p.endecBytes = make([]byte, pktLen)
	p.packetIDPos = p.fill(0, TSUBACK<<4, uint32(remlen))
	p.fill(p.packetIDPos, packetIdentifier, returnCodes)

	return p
}

// PacketIdentifier return packet id
func (p *Suback) PacketIdentifier() uint16 {
	pid, _ := p.uint16(p.packetIDPos)
	return pid
}

// ReturnCodes return sub return codes
func (p *Suback) ReturnCodes() []byte {
	_, offset := p.uint16(p.packetIDPos)
	return p.bytes(offset)
}
