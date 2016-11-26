package mqpp_test

import (
	"testing"

	. "github.com/abo/mqpp"
)

func TestConnack(t *testing.T) {
	connack := MakeConnack(true, Accepted)
	if connack.Type() != CONNACK {
		t.Fatalf("expect connack, actual %d", connack.Type())
	}

	//TODO
}
