package thermal

import (
	"bytes"
	"testing"
)

func TestBasicPacketEncodeDecode(t *testing.T) {

	originaljson := "{\"color\": \"green\"}"
	originalbody := []byte("This is a test body")

	encPacket, err := encodePacket(originaljson, originalbody)
	if err != nil {
		t.Log("Error encoding packet (err: %s)\n", err)
		t.Fail()
	}

	packet, err := decodePacket(encPacket)
	if err != nil {
		t.Log("Error decoding packet (err: %s)\n", err)
		t.Fail()
	}

	t.Logf("packet.json: %s\n", packet.json)
	t.Logf("originaljson: %s\n", originaljson)
	t.Logf("packet.body: %s\n", packet.body)
	t.Logf("originalbody: %s\n", originalbody)

	if packet.json != originaljson {
		t.Log("json was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

	if bytes.Compare(packet.body, originalbody) != 0 {
		t.Log("body was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

}
