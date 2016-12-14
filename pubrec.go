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

// Pubrec mqtt publish received(qos 2 publish received, part 1), structure:
// fixed header
// variable header: Packet Identifier
type Pubrec struct {
	endecBytes
}

func newPubrec(data []byte) (*Pubrec, error) {
	if len(data) < 4 || data[0] != (TPUBREC<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrec{endecBytes: data[0:4]}, nil
}

// MakePubrec create a mqtt pubrec packet
func MakePubrec(packetIdentifier uint16) Pubrec {
	p := Pubrec{endecBytes: make([]byte, 4)}
	p.fill(0, TPUBREC<<4, uint32(2), packetIdentifier)
	return p
}

// PacketIdentifier return packet id
func (p *Pubrec) PacketIdentifier() uint16 {
	pid, _ := p.uint16(2)
	return pid
}
