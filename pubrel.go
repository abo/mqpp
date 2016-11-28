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

// Pubrel mqtt publish release(qos 2 publish received, part 2), structure:
// fixed header
// variable header: Packet Identifier
type Pubrel struct {
	packetBytes
}

func newPubrel(data []byte) (*Pubrel, error) {
	if len(data) < 4 || data[0] != (PUBREL<<4|0x02) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Pubrel{packetBytes: data[0:4]}, nil
}

// MakePubrel create a mqtt pubrel packet
func MakePubrel(packetIdentifier uint16) Pubrel {
	pb := make([]byte, 4)
	fill(pb, PUBREL<<4|0x02, uint32(2), packetIdentifier)
	return Pubrel{packetBytes: pb}
}

// PacketIdentifier return packet id
func (p *Pubrel) PacketIdentifier() uint16 {
	return binary.BigEndian.Uint16(p.packetBytes[2:])
}
