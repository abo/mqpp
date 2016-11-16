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
	if len(data) < 1 || data[0] != (CONNECT<<4) {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	packetLen := offset + int(remlen)
	if len(data) < packetLen {
		return nil, ErrProtocolViolation
	}

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
		packetBytes:          data[0:packetLen],
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

// ProtocolName return protocol name, "MQTT" in 3.1.1
func (c *Connect) ProtocolName() string {
	return string(c.variableHeader()[2 : 2+c.protocolNameBytes])
}

// ProtocolLevel return Protocol Level, 4 in 3.1.1
func (c *Connect) ProtocolLevel() byte {
	return c.variableHeader()[2+c.protocolNameBytes]
}

// UsernameFlag return is username present in the payload
func (c *Connect) UsernameFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 7)
}

// PasswordFlag return is password present in the payload
func (c *Connect) PasswordFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 6)
}

// WillRetain return is server should publish will message
func (c *Connect) WillRetain() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 5)
}

// WillQoS return the QoS level to be used when publishing the Will Message.
func (c *Connect) WillQoS() byte {
	return c.variableHeader()[2+c.protocolNameBytes+1] << 3 >> 6
}

// WillFlag return is will message present int the payload
func (c *Connect) WillFlag() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 2)
}

// CleanSession return is server should clean session when disconnect
func (c *Connect) CleanSession() bool {
	return bit(c.variableHeader()[2+c.protocolNameBytes+1], 1)
}

// KeepAlive return  maximum time interval between client packets transmitting
func (c *Connect) KeepAlive() uint16 {
	return binary.BigEndian.Uint16(c.variableHeader()[2+c.protocolNameBytes+1+1:])
}

// ClientIdentifier return client id
func (c *Connect) ClientIdentifier() string {
	return string(c.payload()[2 : 2+c.clientIDBytes])
}

// WillTopic return will topic if willflag is set, or "" when willflag not set
func (c *Connect) WillTopic() string {
	if !c.WillFlag() {
		return ""
	}
	willTopicOffset := 2 + c.clientIDBytes + 2
	return string(c.payload()[willTopicOffset : willTopicOffset+c.willTopicBytes])
}

// WillMessage return will message if willflag is set, or []byte{} when willflag not set
func (c *Connect) WillMessage() []byte {
	if c.WillFlag() {
		return []byte{}
	}
	willMsgOffset := 2 + c.clientIDBytes + 2 + c.willTopicBytes + 2
	return c.payload()[willMsgOffset : willMsgOffset+c.willMessageBytes]
}

// Username return username when usernameFlag set, or "" when it not set
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

// Password return password when passwordFlag set, or []byte{} when it not set
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
