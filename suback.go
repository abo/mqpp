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

// Suback mqtt subscribe acknowledgement, structure:
// fixed header:
// variable header: Packet Identifier
// payload: Return Codes
type Suback struct {
	packetBytes
	remainingLengthBytes int
}

func newSuback(data []byte) (*Suback, error) {
	if len(data) < 1 || data[0] != (SUBACK<<4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := decRemLen(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	packetLen := 1 + remlenLen + int(remlen)
	if len(data) < packetLen {
		return nil, ErrProtocolViolation
	}

	return &Suback{
		packetBytes:          data[0:packetLen],
		remainingLengthBytes: remlenLen,
	}, nil
}

// MakeSuback create a mqtt suback packet
func MakeSuback(packetIdentifier uint16, returnCodes []byte) Suback {
	remlen := 2 + len(returnCodes)
	remlenLen := lenRemLen(uint32(remlen))
	pb := make([]byte, 1+remlenLen+remlen)

	fill(pb, SUBACK<<4, uint32(remlen), packetIdentifier, returnCodes)

	return Suback{
		packetBytes:          pb,
		remainingLengthBytes: remlenLen,
	}
}

// PacketIdentifier return packet id
func (p *Suback) PacketIdentifier() uint16 {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	return binary.BigEndian.Uint16(p.packetBytes[fixedHeaderLen : fixedHeaderLen+2])
}

// ReturnCodes return sub return codes
func (p *Suback) ReturnCodes() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2
	return p.packetBytes[fixedHeaderLen+variableHeaderLen:]
}
