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

// Disconnect mqtt disconnect notification, structure:
// fixed header
type Disconnect struct {
	endecBytes
}

func newDisconnect(data []byte) (*Disconnect, error) {
	// check packet length, packet type, remaining length
	if len(data) < 2 || data[0] != (TDISCONNECT<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Disconnect{endecBytes: data[0:2]}, nil
}

// MakeDisconnect create a mqtt disconnect packet
func MakeDisconnect() Disconnect {
	p := Disconnect{endecBytes: make([]byte, 2)}
	p.fill(0, TDISCONNECT<<4, uint32(0))
	return p
}
