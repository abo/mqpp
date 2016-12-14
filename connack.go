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

// Connack mqtt acknowledge connection request, structure:
// fixed header
// variable header: Connect Acknowledge Flags(1 byte) + Connect Return code(1 byte)
type Connack struct {
	endecBytes
}

// newConnack parse Connack from byte slice
func newConnack(data []byte) (*Connack, error) {
	// check packet length, packet type, remaining length, conack flags, return code
	if len(data) < 4 || data[0] != (TCONNACK<<4) || data[1] != 2 || (data[2]>>1) != 0 || uint8(data[3]) > 5 {
		return nil, ErrProtocolViolation
	}

	return &Connack{endecBytes: data[0:4]}, nil
}

// MakeConnack create a mqtt connack packet with SessionPresent and ReturnCode
func MakeConnack(sessionPresent bool, returnCode byte) Connack {
	p := Connack{endecBytes: make([]byte, 4)}
	offset := p.fill(0, TCONNACK<<4, uint32(2))
	p.set(offset, 0, sessionPresent)
	p.fill(offset+1, returnCode)
	return p
}

// SetSessionPresent set is the session present or not
func (p *Connack) SetSessionPresent(sessionPresent bool) {
	p.set(2, 0, sessionPresent)
}

// SessionPresent returns is the session present or not
func (p *Connack) SessionPresent() bool {
	return p.bit(2, 0)
}

// SetReturnCode set the return code
func (p *Connack) SetReturnCode(returnCode byte) {
	p.fill(3, returnCode)
}

// ReturnCode return the return code
func (p *Connack) ReturnCode() byte {
	code, _ := p.byte(3)
	return code
}
