package mqpp

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// MQTT Protocol Name and Level
const (
	ProtocolName  = "MQTT"
	ProtocolLevel = 4
)

// MQTT Control Packet types
const (
	RESERVED byte = iota
	CONNECT
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
)

// QoS definitions
const (
	QosAtMostOnce byte = iota
	QosAtLeastOnce
	QosExactlyOnce
)

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

// ControlPacket is interface of basic MQTT packet
type ControlPacket interface {
	Type() byte
	Length() uint32
	// writeTo?
}

var (
	// ErrMalformedRemLen - can not decoding packet's remaining length, or remaining
	// length larger than 268,435,455(max length according MQTT specification)
	ErrMalformedRemLen = errors.New("mqpp: Malformed Remaining Length")
	// ErrIncompletePacket - there are no more data for remaining length
	ErrIncompletePacket = errors.New("mqpp: Incomplete Packet")
	// ErrProtocolViolation - protocol violation according MQTT specification
	ErrProtocolViolation = errors.New("mqpp: Protocol Violation")
)

// ScanPackets is a split function for a bufio.Scanner that returns each
// MQTT packet as a token. it just cut by length, the token maybe protocol
// violation, so check it yourself
func ScanPackets(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	l, n := binary.Uvarint(data[1:])
	if n == 0 {
		if atEOF {
			return 1, nil, ErrIncompletePacket
		}

		return 0, nil, nil
	}

	if l >= (1 << 28) {
		return n, nil, ErrMalformedRemLen
	}

	packetLen := int(l) + n + 1
	if len(data) >= packetLen {
		return packetLen, data[0:packetLen], nil
	} else if atEOF {
		return len(data), nil, ErrIncompletePacket
	} else {
		return 0, nil, nil
	}
}

func NewPacket(data []byte) (ControlPacket, error) {
	switch data[0] >> 4 {
	case CONNECT:
		return NewConnect(data)
	case CONNACK:
		return NewConnack(data)
	case PUBLISH:
		return NewPublish(data)
	case PUBACK:
		return NewPuback(data)
	case PUBREC:
		return NewPubrec(data)
	case PUBREL:
		return NewPubrel(data)
	case PUBCOMP:
		return NewPubcomp(data)
	case SUBSCRIBE:
		return NewSubscribe(data)
	case SUBACK:
		return NewSuback(data)
	case UNSUBSCRIBE:
		return NewUnsubscribe(data)
	case UNSUBACK:
		return NewUnsuback(data)
	case PINGREQ:
		return NewPingreq(data)
	case PINGRESP:
		return NewPingresp(data)
	case DISCONNECT:
		return NewDisconnect(data)
	default:
		return nil, errors.New("Reserved Packet Type")
	}
}

// is bit position set or not
func bit(b byte, n uint8) bool {
	return b&(1<<n) != 0
}

// decode remaining length, and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 268,435,455 (overflow)
//              and -n is the number of bytes read
//
func remainingLength(data []byte) (uint32, int) {
	val, n := binary.Uvarint(data)
	if n > 0 && val >= (1<<28) {
		return 0, -n
	}

	return uint32(val), n
}
