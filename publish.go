package mqpp

import "encoding/binary"

// Publish - publish message
type Publish struct {
	packetBytes
	remainingLengthBytes int
	topicNameBytes       int
}

func newPublish(data []byte) (*Publish, error) {
	if len(data) < 1 || data[0]>>4 != PUBLISH {
		return nil, ErrProtocolViolation
	}
	offset := 1
	remlen, remlenLen := decRemLen(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	packetLen := offset + int(remlen)
	if len(data) < packetLen {
		return nil, ErrProtocolViolation
	}

	topicLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	// offset += (2 + topicLen)

	// pid := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	return &Publish{
		packetBytes:          data[0:packetLen],
		remainingLengthBytes: remlenLen,
		topicNameBytes:       topicLen,
	}, nil
}

// MakePublish create a mqtt publish packet
func MakePublish(dup bool, qos byte, retain bool, topicName string, packetIdentifier uint16, payload []byte) Publish {
	topicLen := len(topicName)
	remlen := (2 + topicLen) + 2 + len(payload)
	remlenLen := lenRemLen(uint32(remlen))

	pb := make([]byte, 1+remlenLen+remlen)
	offset := fill(pb, (PUBLISH<<4)|set(3, dup)|(qos<<1)|set(0, retain))
	offset += fill(pb[offset:], uint32(remlen), topicName, packetIdentifier, payload)

	return Publish{
		packetBytes:          pb,
		remainingLengthBytes: remlenLen,
		topicNameBytes:       topicLen,
	}
}

// Dup return is dup
func (p *Publish) Dup() bool {
	return bit(p.packetBytes[0], 3)
}

// QoS return qos
func (p *Publish) QoS() byte {
	return p.packetBytes[0] << 5 >> 6
}

// Retain return is retain set
func (p *Publish) Retain() bool {
	return bit(p.packetBytes[0], 0)
}

// TopicName return topic name
func (p *Publish) TopicName() string {
	return string(p.variableHeader()[2 : 2+p.topicNameBytes])
}

// PacketIdentifier return packet id if qos > 0, or zero when qos = 0
func (p *Publish) PacketIdentifier() uint16 {
	if p.QoS() > QosAtMostOnce {
		return binary.BigEndian.Uint16(p.variableHeader()[2+p.topicNameBytes:])
	}

	return 0
}

func (p *Publish) variableHeader() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2 + p.topicNameBytes
	if p.QoS() > QosAtMostOnce {
		variableHeaderLen += 2
	}
	return p.packetBytes[fixedHeaderLen : fixedHeaderLen+variableHeaderLen]
}

// Payload return publish content
func (p *Publish) Payload() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2 + p.topicNameBytes
	if p.QoS() > QosAtMostOnce {
		variableHeaderLen += 2
	}
	return p.packetBytes[fixedHeaderLen+variableHeaderLen:]
}
