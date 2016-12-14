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
	"errors"
	"fmt"
	"io"
)

// MQTT Protocol Name and Level
const (
	ProtocolName  string = "MQTT"
	ProtocolLevel byte   = 4
)

// MQTT Control Packet types
const (
	TCONNECT byte = iota + 1
	TCONNACK
	TPUBLISH
	TPUBACK
	TPUBREC
	TPUBREL
	TPUBCOMP
	TSUBSCRIBE
	TSUBACK
	TUNSUBSCRIBE
	TUNSUBACK
	TPINGREQ
	TPINGRESP
	TDISCONNECT
)

// QoS definitions
const (
	QosAtMostOnce byte = iota
	QosAtLeastOnce
	QosExactlyOnce
)

// SubackFailure suback return code - failed
const SubackFailure byte = 0x80

// Connect Return Code
const (
	Accepted byte = iota
	RefusedProtocolVersion
	RefusedInvalidIdentifier
	RefusedServerUnavailable
	RefusedBadCredentials
	RefusedUnauthorized
)

// ConnectReturnCodeResponses - Connect Return code descriptions
var ConnectReturnCodeResponses = map[byte]string{
	Accepted:                 fmt.Sprintf("%#02x Connection Accepted", Accepted),
	RefusedProtocolVersion:   fmt.Sprintf("%#02x Connection Refused, unacceptable protocol version", RefusedProtocolVersion),
	RefusedInvalidIdentifier: fmt.Sprintf("%#02x Connection Refused, identifier rejected", RefusedInvalidIdentifier),
	RefusedServerUnavailable: fmt.Sprintf("%#02x Connection Refused, server unavailable", RefusedServerUnavailable),
	RefusedBadCredentials:    fmt.Sprintf("%#02x Connection Refused, bad user name or password", RefusedBadCredentials),
	RefusedUnauthorized:      fmt.Sprintf("%#02x Connection Refused, not authorized", RefusedUnauthorized),
}

var (
	// ErrMalformedRemLen - can not decoding packet's remaining length, or remaining
	// length larger than 268,435,455(max length according MQTT specification)
	ErrMalformedRemLen = errors.New("mqpp: Malformed Remaining Length")
	// ErrIncompletePacket - there are no more data for remaining length
	ErrIncompletePacket = errors.New("mqpp: Incomplete Packet")
	// ErrProtocolViolation - protocol violation according MQTT specification
	ErrProtocolViolation = errors.New("mqpp: Protocol Violation")
	// ErrReservedPacketType - unknown packet type
	ErrReservedPacketType = errors.New("mqpp: Reserved Packet Type")
)

// ControlPacket is interface of basic MQTT packet
type ControlPacket interface {
	Type() byte
	Length() uint32
	Bytes() []byte
	io.WriterTo
}
