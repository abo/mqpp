package mqpp

// Disconnect - disconnect notification
type Disconnect struct {
	packetBytes
}

func newDisconnect(data []byte) (*Disconnect, error) {
	// check packet length, packet type, remaining length
	if len(data) < 2 || data[0] != (DISCONNECT<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Disconnect{packetBytes: data[0:2]}, nil
}

// MakeDisconnect create a mqtt disconnect packet
func MakeDisconnect() Disconnect {
	pb := make([]byte, 2)
	fill(pb, DISCONNECT<<4, uint32(0))
	return Disconnect{packetBytes: pb}
}
