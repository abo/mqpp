package mqpp

// Pingresp - ping response
type Pingresp struct {
	packetBytes
}

func newPingresp(data []byte) (*Pingresp, error) {
	if len(data) < 2 || data[0] != (PINGRESP<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Pingresp{packetBytes: data[0:2]}, nil
}

// MakePingresp create a mqtt pingresp packet
func MakePingresp() Pingresp {
	pb := make([]byte, 2)
	fill(pb, PINGRESP<<4, uint32(0))
	return Pingresp{packetBytes: pb}
}
