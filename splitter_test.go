// Copyright (c) 2016 The MQPP Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mqpp

import (
	"bytes"
	"testing"
)

func TestScan(t *testing.T) {
	pkts := []ControlPacket{
		MakeConnack(false, Accepted),
		MakeConnect(ProtocolName, ProtocolLevel, true, QosExactlyOnce, false, 128, "clientIdentifier", "willTopic", []byte("willMessage"), "username", []byte("password")),
		MakeDisconnect(),
		MakePingreq(),
		MakePingresp(),
		MakePuback(123),
		MakePubcomp(124),
		MakePublish(false, QosAtMostOnce, true, "topicName", 125, []byte("payload")),
		MakePubrec(126),
		MakePubrel(127),
		MakeSuback(65530, []byte{QosAtLeastOnce, QosAtMostOnce, QosExactlyOnce, SubackFailure}),
		MakeSubscribe(2, []Subscription{{TopicFilter: "/topic/filter", RequestedQoS: QosExactlyOnce}, {TopicFilter: "/topic/#", RequestedQoS: QosAtMostOnce}}),
		MakeUnsuback(1),
		MakeUnsubscribe(65535, []string{"#", "/topic/a/aa"}),
	}

	var buf bytes.Buffer
	for _, pkt := range pkts {
		buf.Write(pkt.Bytes())
	}

	// data := []byte{0x31, 0x0a, 0x00, 0x08, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x41, 0x2f, 0x43, 0x31, 0x09, 0x00, 0x07, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x43}
	data := buf.Bytes()
	r := bytes.NewReader(data)
	s := NewSplitter(r)

	for i, origin := range pkts {
		p, err := s.NextPacket()
		if err != nil || 0 != bytes.Compare(origin.Bytes(), p.Bytes()) {
			t.Fatalf("no.%d : expect %v, actual %v, with err:%v", i, origin, p, err)
		}
	}
	p, err := s.NextPacket()
	if p != nil || err != nil {
		t.Fatal("expect EOF")
	}

	// data := []byte{0x31, 0x0a, 0x00, 0x08, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x41, 0x2f, 0x43}
	// data2 := []byte{0x31, 0x9, 0x0, 0x7, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x43}
}
