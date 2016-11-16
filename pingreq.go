package mqpp

// Pingreq - ping request
type Pingreq struct {
	packetBytes
}

func newPingreq(data []byte) (*Pingreq, error) {
	if len(data) < 2 || data[0] != (PINGREQ<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Pingreq{packetBytes: data[0:2]}, nil
}
