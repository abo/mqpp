package mqpp_test

import (
	"bufio"
	"fmt"
	"testing"
)

import . "github.com/abo/mqpp"

func TestScan(t *testing.T) {
	scanner := bufio.NewScanner(nil)

	scanner.Split(ScanPackets)
	for scanner.Scan() {
		p, err := NewPacket(scanner.Bytes())
		fmt.Println(p.Type())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
