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

// func NewDisconnect() *Disconnect {
// 	return &Disconnect{
// 		src: []byte{DISCONNECT << 4, 0x00},
// 	}
// }
