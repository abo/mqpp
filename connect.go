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

// Connect mqtt client requests a connection to server, structure:
// fixed header
// variable header: Protocol Name, Protocol Level, Connect Flags, Keep Alive
// payload: Client Identifier, Will Topic, Will Message, User Name, Password
type Connect struct {
	endecBytes

	protocolNamePos  int
	protocolLevelPos int
	connectFlagsPos  int
	keepalivePos     int
	clientIDPos      int
	willTopicPos     int
	willMessagePos   int
	usernamePos      int
	passwordPos      int
}

func newConnect(data []byte) (*Connect, error) {
	if len(data) < 1 || data[0] != (TCONNECT<<4) {
		return nil, ErrProtocolViolation
	}

	pkt := &Connect{endecBytes: data}

	_, offset := pkt.byte(0)             // 1)packet type, reserved
	remlen, offset := pkt.remlen(offset) // 2)remaining length

	if offset <= 1 {
		return nil, ErrMalformedRemLen
	}
	pktLen := offset + int(remlen)
	if len(data) < pktLen {
		return nil, ErrProtocolViolation
	}
	pkt.protocolNamePos = offset
	_, pkt.protocolLevelPos = pkt.string(offset)            // 3)protocol name
	_, pkt.connectFlagsPos = pkt.byte(pkt.protocolLevelPos) // 4)protocol level
	usernameFlag := pkt.bit(pkt.connectFlagsPos, 7)
	passwordFlag := pkt.bit(pkt.connectFlagsPos, 6)
	willFlag := pkt.bit(pkt.connectFlagsPos, 2)
	_, pkt.keepalivePos = pkt.byte(pkt.connectFlagsPos) // 5)connect flags
	_, pkt.clientIDPos = pkt.uint16(pkt.keepalivePos)   // 6)keep alive
	_, offset = pkt.string(pkt.clientIDPos)             // 7)clientid

	if willFlag {
		pkt.willTopicPos = offset
		_, pkt.willMessagePos = pkt.string(pkt.willTopicPos) // 8)will topic
		_, offset = pkt.string(pkt.willMessagePos)           // 9)will message
	}
	if usernameFlag {
		pkt.usernamePos = offset
		_, offset = pkt.string(pkt.usernamePos) // 10)user name
	}
	if passwordFlag {
		pkt.passwordPos = offset
		// _, _ = pkt.string(pkt.passwordPos) // 11)password
	}

	return pkt, nil
}

// MakeConnect create a mqtt connect packet with fields
func MakeConnect(protocolName string, protocolLevel byte, willRetain bool, willQoS byte, cleanSession bool, keepAlive uint16, clientIdentifier string, willTopic string, willMessage []byte, username string, password []byte) Connect {
	p := Connect{}
	remlen := p.calc(protocolName, protocolLevel, willQoS, keepAlive, clientIdentifier)
	willFlag, usernameFlag, passwordFlag := len(willTopic) > 0, len(username) > 0, len(password) > 0
	if willFlag {
		remlen += p.calc(willTopic, string(willMessage))
	}
	if usernameFlag {
		remlen += p.calc(username)
	}
	if passwordFlag {
		remlen += p.calc(string(password))
	}
	pktLen := 1 + p.calc(uint32(remlen)) + remlen

	p.endecBytes = make([]byte, pktLen)
	p.protocolNamePos = p.fill(0, TCONNECT<<4, uint32(remlen))
	p.protocolLevelPos = p.fill(p.protocolNamePos, protocolName)
	p.connectFlagsPos = p.fill(p.protocolLevelPos, protocolLevel)
	p.keepalivePos = p.fill(p.connectFlagsPos, willQoS<<3)
	p.set(p.connectFlagsPos, 7, usernameFlag)
	p.set(p.connectFlagsPos, 6, passwordFlag)
	p.set(p.connectFlagsPos, 5, willRetain)
	p.set(p.connectFlagsPos, 2, willFlag)
	p.set(p.connectFlagsPos, 1, cleanSession)
	p.clientIDPos = p.fill(p.keepalivePos, keepAlive)
	offset := p.fill(p.clientIDPos, clientIdentifier)
	if willFlag {
		p.willTopicPos = offset
		p.willMessagePos = p.fill(p.willTopicPos, willTopic)
		offset = p.fill(p.willMessagePos, string(willMessage))
	}
	if usernameFlag {
		p.usernamePos = offset
		offset = p.fill(p.usernamePos, username)
	}
	if passwordFlag {
		p.passwordPos = offset
		offset = p.fill(p.passwordPos, string(password))
	}

	return p
}

// ProtocolName return protocol name, "MQTT" in 3.1.1
func (c *Connect) ProtocolName() string {
	protoName, _ := c.string(c.protocolNamePos)
	return protoName
}

// ProtocolLevel return Protocol Level, 4 in 3.1.1
func (c *Connect) ProtocolLevel() byte {
	protoLevel, _ := c.byte(c.protocolLevelPos)
	return protoLevel
}

// UsernameFlag return is username present in the payload
func (c *Connect) UsernameFlag() bool {
	return c.bit(c.connectFlagsPos, 7)
}

// PasswordFlag return is password present in the payload
func (c *Connect) PasswordFlag() bool {
	return c.bit(c.connectFlagsPos, 6)
}

// WillRetain return is server should publish will message
func (c *Connect) WillRetain() bool {
	return c.bit(c.connectFlagsPos, 5)
}

// WillQoS return the QoS level to be used when publishing the Will Message.
func (c *Connect) WillQoS() byte {
	qos, _ := c.byte(c.connectFlagsPos)
	return qos << 3 >> 6
}

// WillFlag return is will message present int the payload
func (c *Connect) WillFlag() bool {
	return c.bit(c.connectFlagsPos, 2)
}

// CleanSession return is server should clean session when disconnect
func (c *Connect) CleanSession() bool {
	return c.bit(c.connectFlagsPos, 1)
}

// KeepAlive return  maximum time interval between client packets transmitting
func (c *Connect) KeepAlive() uint16 {
	keepalive, _ := c.uint16(c.keepalivePos)
	return keepalive
}

// ClientIdentifier return client id
func (c *Connect) ClientIdentifier() string {
	cid, _ := c.string(c.clientIDPos)
	return cid
}

// WillTopic return will topic if willflag is set, or "" when willflag not set
func (c *Connect) WillTopic() string {
	if !c.WillFlag() {
		return ""
	}
	topic, _ := c.string(c.willTopicPos)
	return topic
}

// WillMessage return will message if willflag is set, or []byte{} when willflag not set
func (c *Connect) WillMessage() []byte {
	if c.WillFlag() {
		return []byte{}
	}
	msg, _ := c.string(c.willMessagePos)
	return []byte(msg)
}

// Username return username when usernameFlag set, or "" when it not set
func (c *Connect) Username() string {
	if !c.UsernameFlag() {
		return ""
	}

	uname, _ := c.string(c.usernamePos)
	return uname
}

// Password return password when passwordFlag set, or []byte{} when it not set
func (c *Connect) Password() []byte {
	if !c.PasswordFlag() {
		return []byte{}
	}

	pwd, _ := c.string(c.passwordPos)
	return []byte(pwd)
}
