package mqpp

import "encoding/binary"

//Connect - client requests a connection to server
type Connect struct {
	packetBytes

	// 1 byte) packet type, reserved
	// 1~4 bytes) remaining length
	remainingLengthBytes int

	// 2+4=6 bytes) protocol name
	protocolNameBytes int
	// 1 byte) protocol level
	// 1 byte) connect flags (usernameFlag, passwordFlag, willRetain, willQoS(2 bit), willFlag, cleanSession, reserved)
	// 2 bytes) keep alive

	// 2+n bytes) clientid
	clientIDBytes int
	// 2+n bytes) will topic
	willTopicBytes int
	// 2+n bytes) will message
	willMessageBytes int
	// 2+n bytes) user name
	usernameBytes int
	// 2+n bytes) password
	passwordBytes int
}

func newConnect(data []byte) (*Connect, error) {
	if data[0] != (CONNECT << 4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	_, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	pnameLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	offset += (2 + pnameLen)

	// plevel
	offset++

	cflags := data[offset]
	offset++

	// keepalive
	offset += 2

	cidLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	offset += (2 + cidLen)

	wtopicLen, wmsgLen := 0, 0
	if bit(cflags, 2) {
		wtopicLen = int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + wtopicLen)

		wmsgLen = int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + wmsgLen)
	}

	unameLen := 0
	if bit(cflags, 7) {
		unameLen = int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += (2 + unameLen)
	}

	pwdLen := 0
	if bit(cflags, 6) {
		pwdLen = int(binary.BigEndian.Uint16(data[offset : offset+2]))
	}

	return &Connect{
		packetBytes:          data,
		remainingLengthBytes: remlenLen,
		protocolNameBytes:    pnameLen,
		clientIDBytes:        cidLen,
		willTopicBytes:       wtopicLen,
		willMessageBytes:     wmsgLen,
		usernameBytes:        unameLen,
		passwordBytes:        pwdLen,
	}, nil
}

// Packet Type + Reserved (1 byte)
// Remaining Length (1~4 bytes)
func (c *Connect) fixedHeader() []byte {
	fixedHeaderLen := 1 + c.remainingLengthBytes
	return c.packetBytes[0:fixedHeaderLen]
}

// Protocol Name (2 + n bytes)
// Protocol Level (1 byte)
// Connect Flags (1 byte)
// Keep Alive (2 bytes)
func (c *Connect) variableHeader() []byte {
	fixedHeaderLen := 1 + c.remainingLengthBytes
	variableHeaderLen := 2 + c.protocolNameBytes + 1 + 1 + 2
	return c.packetBytes[fixedHeaderLen : fixedHeaderLen+variableHeaderLen]
}

// Client Identifier (2 + n bytes)
// Will Topic (2 + n bytes)
// Will Message
// User Name
// Password
func (c *Connect) payload() []byte {
	fixedHeaderLen := 1 + c.remainingLengthBytes
	variableHeaderLen := 2 + c.protocolNameBytes + 1 + 1 + 2
	return c.packetBytes[fixedHeaderLen+variableHeaderLen:]
}

// ProtocolName decode Protocol Name, 6 byte in variable header
// A UTF-8 encoded string that represents the protocol name “MQTT”
// according the specification
func (c *Connect) ProtocolName() string {
	return string(c.variableHeader()[2 : 2+c.protocolNameBytes])
}

// ProtocolLevel decode Protocol Level, 1 byte in variable header
func (c *Connect) ProtocolLevel() byte {
	return c.variableHeader()[2+c.protocolNameBytes]
}

func (c *Connect) UsernameFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 7)
}

func (c *Connect) PasswordFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 6)
}

func (c *Connect) WillRetain() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 5)
}

func (c *Connect) WillQoS() byte {
	return c.variableHeader()[2+c.protocolNameBytes+1] << 3 >> 6
}

func (c *Connect) WillFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 2)
}

func (c *Connect) CleanSession() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 1)
}

func (c *Connect) KeepAlive() uint16 {
	return binary.BigEndian.Uint16(c.variableHeader()[2+c.protocolNameBytes+1+1:])
}

func (c *Connect) ClientIdentifier() string {
	return string(c.payload()[2 : 2+c.clientIDBytes])
}

func (c *Connect) WillTopic() string {
	if !c.WillFlag() {
		return ""
	}
	willTopicOffset := 2 + c.clientIDBytes + 2
	return string(c.payload()[willTopicOffset : willTopicOffset+c.willTopicBytes])
}

func (c *Connect) WillMessage() []byte {
	if c.WillFlag() {
		return []byte{}
	}
	willMsgOffset := 2 + c.clientIDBytes + 2 + c.willTopicBytes + 2
	return c.payload()[willMsgOffset : willMsgOffset+c.willMessageBytes]
}
func (c *Connect) Username() string {
	if !c.UsernameFlag() {
		return ""
	}

	usernameOffset := 2 + c.clientIDBytes + 2
	if c.WillFlag() {
		usernameOffset += (2 + c.willTopicBytes + 2 + c.willMessageBytes)
	}
	return string(c.payload()[usernameOffset : usernameOffset+c.usernameBytes])
}

func (c *Connect) Password() []byte {
	if !c.PasswordFlag() {
		return []byte{}
	}

	passwordOffset := 2 + c.clientIDBytes + 2
	if c.WillFlag() {
		passwordOffset += (2 + c.willTopicBytes + 2 + c.willMessageBytes)
	}
	if c.UsernameFlag() {
		passwordOffset += (2 + c.usernameBytes)
	}
	return c.payload()[passwordOffset:]
}
