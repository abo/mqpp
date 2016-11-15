package mqpp

// Disconnect - disconnect notification
type Disconnect struct {
	src []byte
}

func NewDisconnect(data []byte) (*Disconnect, error) {
	// check packet length, packet type, remaining length
	if len(data) != 2 || data[0] != (DISCONNECT<<4) || data[1] != 0 {
		return nil, ErrProtocolViolation
	}
	return &Disconnect{src: data}, nil
}

// func NewDisconnect() *Disconnect {
// 	return &Disconnect{
// 		src: []byte{DISCONNECT << 4, 0x00},
// 	}
// }
func (d *Disconnect) Length() uint32 { return uint32(len(d.src)) }

func (d *Disconnect) Type() byte { return d.src[0] >> 4 }
