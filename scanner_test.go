package mqpp_test

import "testing"
import "bytes"

import . "github.com/abo/mqpp"

func TestScan(t *testing.T) {
	data := []byte{0x31, 0x0a, 0x00, 0x08, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x41, 0x2f, 0x43, 0x31, 0x09, 0x00, 0x07, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x43}
	r := bytes.NewReader(data)
	s := NewScanner(r)
	if !s.Scan() {
		t.Fatalf("scan first packet failed,%v", s.Err())
	}
	if bytes.Compare(s.Bytes(), data[0:12]) != 0 {
		t.Fatalf("first packet is not expected,%v", s.Bytes())
	}
	if !s.Scan() {
		t.Fatalf("scan second packet failed,%v", s.Err())
	}
	if bytes.Compare(s.Bytes(), data[12:]) != 0 {
		t.Fatalf("second packet expect:%v, actual:%v", s.Bytes(), data[12:])
	}
}
