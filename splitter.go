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

import (
	"bufio"
	"io"
)

// Splitter wrap bufio.Scanner with SplitFunc which split a file into mqtt packets
type Splitter struct {
	bufio.Scanner
}

// Packet returns the most recent token generated by a call to Scan as a mqtt packet holding its bytes.
func (s *Splitter) Packet() (ControlPacket, error) {
	data := s.Bytes()
	switch data[0] >> 4 {
	case TCONNECT:
		return newConnect(data)
	case TCONNACK:
		return newConnack(data)
	case TPUBLISH:
		return newPublish(data)
	case TPUBACK:
		return newPuback(data)
	case TPUBREC:
		return newPubrec(data)
	case TPUBREL:
		return newPubrel(data)
	case TPUBCOMP:
		return newPubcomp(data)
	case TSUBSCRIBE:
		return newSubscribe(data)
	case TSUBACK:
		return newSuback(data)
	case TUNSUBSCRIBE:
		return newUnsubscribe(data)
	case TUNSUBACK:
		return newUnsuback(data)
	case TPINGREQ:
		return newPingreq(data)
	case TPINGRESP:
		return newPingresp(data)
	case TDISCONNECT:
		return newDisconnect(data)
	default:
		return nil, ErrReservedPacketType
	}
}

// NextPacket advances the Splitter to the next packet, and return it.
// it return any error that
// occurred during scanning and parsing, except that if it was io.EOF, Err
// will return nil.
func (s *Splitter) NextPacket() (ControlPacket, error) {
	if !s.Scan() {
		return nil, s.Err()
	}
	return s.Packet()
}

// NewSplitter returns a new Splitter to read from r, with The split function splitPackets.
func NewSplitter(r io.Reader) *Splitter {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitPackets)
	return &Splitter{
		Scanner: *scanner,
	}
}

// splitPackets is a split function for a bufio.Scanner that returns each
// MQTT packet as a token. it just cut by length, the token maybe protocol
// violation, so check it yourself
func splitPackets(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	l, n := endecBytes(data).remlen(1)

	if n == 0 {
		if atEOF {
			return 1, nil, ErrIncompletePacket
		}

		return 0, nil, nil
	}

	if n < 0 {
		return n, nil, ErrMalformedRemLen
	}

	packetLen := int(l) + n
	if len(data) >= packetLen {
		return packetLen, data[0:packetLen], nil
	} else if atEOF {
		return len(data), nil, ErrIncompletePacket
	} else {
		return 0, nil, nil
	}
}
