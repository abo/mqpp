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

// Unsuback mqtt unsubscribe acknowledgement, structure:
// fixed header:
// variable header: Packet Identifier
type Unsuback struct {
	endecBytes
}

func newUnsuback(data []byte) (*Unsuback, error) {
	if len(data) < 4 || data[0] != (TUNSUBACK<<4) || data[1] != 2 {
		return nil, ErrProtocolViolation
	}
	return &Unsuback{endecBytes: data[0:4]}, nil
}

// MakeUnsuback create a mqtt unsuback packet
func MakeUnsuback(packetIdentifier uint16) Unsuback {
	p := Unsuback{endecBytes: make([]byte, 4)}
	p.fill(0, TUNSUBACK<<4, uint32(2), packetIdentifier)
	return p
}

// PacketIdentifier return packet id
func (s *Unsuback) PacketIdentifier() uint16 {
	pid, _ := s.uint16(2)
	return pid
}
