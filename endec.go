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
	"encoding/binary"
	"io"
)

// A endecBytes is a slice of bytes with Read and Write methods for MQTT fields.The zero value for endecBytes is an calc for fields length.
type endecBytes []byte

// Length returns how many bytes this packet
func (bs endecBytes) Length() uint32 {
	return uint32(len(bs))
}

// Type returns packet type
func (bs endecBytes) Type() byte {
	return bs[0] >> 4
}

// WriteTo write packet to w
func (bs endecBytes) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(bs)
	return int64(n), err
}

// Bytes returns
func (bs endecBytes) Bytes() []byte {
	return []byte(bs)
}

// returns how many bytes will be writen
func (bs *endecBytes) calc(fields ...interface{}) int {
	total := 0
	for _, val := range fields {
		switch val.(type) {
		case byte:
			total++
		case uint16:
			total += 2
		case uint32: // remaingLength
			remlen := val.(uint32)
			if remlen <= 127 {
				total++
			} else if remlen <= 16383 {
				total += 2
			} else if remlen <= 2097151 {
				total += 3
			} else if remlen <= 268435455 {
				total += 4
			}
		case string:
			total += (2 + len([]byte(val.(string))))
		case []string:
			for _, s := range val.([]string) {
				total += (2 + len([]byte(s)))
			}
		case []byte:
			total += len(val.([]byte))
		default: // unknown type
			return -total
		}
	}
	return total
}

// fill byte array, return new offset(old offset + n bytes writen)
// if new offset < 0, meanings unsupport type, n bytes writen(n = -new offset-old offset)
func (bs endecBytes) fill(offset int, fields ...interface{}) int {
	for _, val := range fields {
		switch val.(type) {
		case byte:
			bs[offset] = val.(byte)
			offset++
		case uint16:
			binary.BigEndian.PutUint16(bs[offset:], val.(uint16))
			offset += 2
		case uint32: // remaingLength
			n := binary.PutUvarint(bs[offset:], uint64(val.(uint32)))
			offset += n
		case string:
			str := val.(string)
			binary.BigEndian.PutUint16(bs[offset:], uint16(len(str)))
			offset += 2
			offset += copy(bs[offset:], str)
		case []byte:
			offset += copy(bs[offset:], val.([]byte))
		default: // unknown type
			return -offset
		}
	}
	return offset
}

func (bs endecBytes) set(offset int, pos uint8, v bool) endecBytes {
	b := byte(1) << pos
	if v {
		bs[offset] = bs[offset] | b
	} else {
		bs[offset] = bs[offset] & (^b)
	}
	return bs
}

func (bs endecBytes) bit(offset int, pos uint8) bool {
	return bs[offset]&(1<<pos) != 0
}

func (bs endecBytes) byte(offset int) (byte, int) {
	return bs[offset], offset + 1
}

func (bs endecBytes) uint16(offset int) (uint16, int) {
	return binary.BigEndian.Uint16(bs[offset : offset+2]), offset + 2
}

func (bs endecBytes) string(offset int) (string, int) {
	l, start := bs.uint16(offset)
	return string(bs[start : start+int(l)]), start + int(l)
}

func (bs endecBytes) remlen(offset int) (uint32, int) {
	val, n := binary.Uvarint(bs[offset:])

	if val > 268435455 {
		return 0, offset - n
	}
	return uint32(val), offset + n
}

func (bs endecBytes) bytes(offset int) []byte {
	return bs[offset:]
}
