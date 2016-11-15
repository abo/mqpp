package mqpp_test

import (
	"fmt"
	"testing"

	. "github.com/abo/mqpp"
)

func TestDecodePublish(t *testing.T) {
	data := []byte{0x31, 0x0a, 0x00, 0x08, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x41, 0x2f, 0x43}
	pub, err := NewPublish(data)
	if err != nil {
		t.Fatalf("failed to decode %v, Err:%v", data, err)
	}
	if pub.Dup() {
		t.Fatal("expect not dup")
	}
	if !pub.Retain() {
		t.Fatal("expect retain")
	}
	if pub.Length() != uint32(len(data)) {
		t.Fatal("length changed")
	}
	if pub.QoS() != 0 {
		t.Fatal("expect at most once")
	}
	if pub.PacketIdentifier() != 0 {
		t.Fatal("expect no packet identifier")
	}

	fmt.Println("payload:", pub.Payload())
	fmt.Println("topic:", pub.TopicName())
	fmt.Println("type:", pub.Type())

	// data2 := []byte{0x31, 0x9, 0x0, 0x7, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x43}
	// pub2, err2 := NewPublish(data2)
	// if err2 != nil {
	// 	t.Fatal(err.Error())
	// }
	// fmt.Println("dup:", pub2.Dup())
	// fmt.Println("len:", pub2.Length())

	// fmt.Println("packetid:", pub2.PacketIdentifier())

	// fmt.Println("payload:", pub2.Payload())
	// fmt.Println("qos:", pub2.QoS())
	// fmt.Println("retain:", pub2.Retain())
	// fmt.Println("topic:", pub2.TopicName())
	// fmt.Println("type:", pub2.Type())
}
