package thermal

import (
	"bytes"
	"testing"
)

func TestBasicPacketEncodeDecode(t *testing.T) {

	originalhead := []byte("{\"color\": \"green\"}")
	originalbody := []byte("This is a test body")

	encPacket, err := encodePacket(originalhead, originalbody)
	if err != nil {
		t.Log("Error encoding packet (err: %s)\n", err)
		t.Fail()
	}

	packet, err := decodePacket(encPacket)
	if err != nil {
		t.Log("Error decoding packet (err: %s)\n", err)
		t.Fail()
	}

	t.Logf("packet.head: %s\n", packet.head)
	t.Logf("originalhead: %s\n", originalhead)
	t.Logf("packet.body: %s\n", packet.body)
	t.Logf("originalbody: %s\n", originalbody)

	if bytes.Compare(packet.head, originalhead) != 0 {
		t.Log("head was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

	if bytes.Compare(packet.body, originalbody) != 0 {
		t.Log("body was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

}
