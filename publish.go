package mqpp

import "encoding/binary"

// Publish - publish message
type Publish struct {
	src                  []byte
	remainingLengthBytes int
	topicNameBytes       int
}

func NewPublish(data []byte) (*Publish, error) {
	if data[0]>>4 != PUBLISH {
		return nil, ErrProtocolViolation
	}
	offset := 1
	_, remlenLen := remainingLength(data[offset:])
	if remlenLen <= 0 {
		return nil, ErrMalformedRemLen
	}
	offset += remlenLen

	topicLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	// offset += (2 + topicLen)

	// pid := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	return &Publish{
		src:                  data,
		remainingLengthBytes: remlenLen,
		topicNameBytes:       topicLen,
	}, nil
}

func (p *Publish) Length() uint32 { return uint32(len(p.src)) }

func (p *Publish) Type() byte { return p.src[0] >> 4 }

func (p *Publish) Dup() bool {
	return bit(p.src[0], 3)
}

func (p *Publish) QoS() byte {
	return p.src[0] << 5 >> 6
}

func (p *Publish) Retain() bool {
	return bit(p.src[0], 0)
}

func (p *Publish) TopicName() string {
	return string(p.variableHeader()[2 : 2+p.topicNameBytes])
}

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
	return p.src[fixedHeaderLen : fixedHeaderLen+variableHeaderLen]
}

func (p *Publish) Payload() []byte {
	fixedHeaderLen := 1 + p.remainingLengthBytes
	variableHeaderLen := 2 + p.topicNameBytes
	if p.QoS() > QosAtMostOnce {
		variableHeaderLen += 2
	}
	return p.src[fixedHeaderLen+variableHeaderLen:]
}
